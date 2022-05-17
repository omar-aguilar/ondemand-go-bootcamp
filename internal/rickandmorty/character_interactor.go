package rickandmorty

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/config"
)

type Interactor interface {
	LoadAndStore(file io.Reader) error
	GetById(ID int) (Character, error)
	StoreCharactersByPageFromAPI(page int, storageFormat string) (CharacterList, error)
	GetCharactersStoredByPageFromAPI(page int, storageFormat string) (CharacterList, error)
	ReadConcurrent(file io.Reader, params ReadConcurrentParams) (CharacterList, error)
}

type interactor struct {
	config  config.Config
	ds      CharacterRepository
	apiDS   APIGetter
	storeDS CharacterStorer
}

func getCharactersFilenameFromPage(page int, extension string) string {
	return fmt.Sprintf("character_%03d.%s", page, extension)
}

func NewInteractor(config config.Config, ds CharacterRepository, apiDS APIGetter, storeDS CharacterStorer) Interactor {
	return interactor{
		config,
		ds,
		apiDS,
		storeDS,
	}
}

func (i interactor) LoadAndStore(file io.Reader) error {
	characterList, err := i.ds.Load(file)
	if err != nil {
		return err
	}

	return i.storeDS.Write(i.config.DBFile, characterList, FormatCSV)
}

func (i interactor) GetById(ID int) (Character, error) {
	if ID <= 0 {
		return Character{}, ErrInvalidID
	}
	character, err := i.ds.GetById(ID)
	if err != nil {
		return Character{}, err
	}
	return character, nil
}

func (i interactor) StoreCharactersByPageFromAPI(page int, storageFormat string) (CharacterList, error) {
	if page <= 0 {
		return CharacterList{}, ErrInvalidPage
	}
	characterList, err := i.apiDS.GetCharactersByPage(page)
	if err != nil {
		return CharacterList{}, err
	}
	filename := getCharactersFilenameFromPage(page, storageFormat)
	err = i.storeDS.Write(filename, characterList, storageFormat)
	if err != nil {
		return CharacterList{}, err
	}
	return characterList, err
}

func (i interactor) GetCharactersStoredByPageFromAPI(page int, storageFormat string) (CharacterList, error) {
	if page <= 0 {
		return CharacterList{}, ErrInvalidPage
	}
	filename := getCharactersFilenameFromPage(page, storageFormat)
	characterList := CharacterList{}
	err := i.storeDS.Read(filename, &characterList, storageFormat)
	return characterList, err
}

type lineChecker func(number int) bool

func isEven(number int) bool {
	return number%2 == 0
}

func isOdd(number int) bool {
	return !isEven(number)
}

func worker(ctx context.Context, id int, processed chan<- Character, lines <-chan string, params ReadConcurrentParams) {
	codec := NewCSVCharacterCodec()
	count := 0
	var isValidLine lineChecker = isEven
	if params.Type == "odd" {
		isValidLine = isOdd
	}
loop:
	for line := range lines {
		character := Character{}
		codec.Decode(strings.NewReader(line), &character)
		if count == params.ItemsPerWorker {
			break
		}
		if character.ID == 0 || !isValidLine(character.ID) {
			continue
		}
		select {
		case <-ctx.Done():
			break loop
		case processed <- character:
		}
		count++
	}
}

func consumer(cancel context.CancelFunc, results chan<- CharacterList, processed <-chan Character, params ReadConcurrentParams) {
	characterList := CharacterList{}
	for character := range processed {
		characterList = append(characterList, character)
		if params.Items > 0 && len(characterList) == params.Items {
			cancel()
			break
		}
	}
	results <- characterList
	close(results)
}

func (i interactor) ReadConcurrent(file io.Reader, params ReadConcurrentParams) (CharacterList, error) {
	if err := Validate(params); err != nil {
		return CharacterList{}, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	linesChannel := make(chan string)
	processedChannel := make(chan Character)
	resultsChannel := make(chan CharacterList)
	numOfWorkers := getNumberOfWorkers(params.NumberOfWorkers)

	go i.ds.ReadConcurrent(file, params, linesChannel)
	for i := 1; i <= numOfWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(ctx, id, processedChannel, linesChannel, params)
		}(i)
	}
	go consumer(cancel, resultsChannel, processedChannel, params)
	wg.Wait()
	close(processedChannel)
	list := <-resultsChannel
	return list, nil
}
