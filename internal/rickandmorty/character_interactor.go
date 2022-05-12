package rickandmorty

import (
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

func worker(id int, output chan<- string, lines <-chan string, params ReadConcurrentParams) {
	codec := NewCSVCharacterCodec()
	var isValidLine lineChecker = isEven
	if params.Type == "odd" {
		isValidLine = isOdd
	}
	for line := range lines {
		character := Character{}
		codec.Decode(strings.NewReader(line), &character)
		if character.ID == 0 || !isValidLine(character.ID) {
			continue
		}
		output <- line
	}
}

func (i interactor) ReadConcurrent(file io.Reader, params ReadConcurrentParams) (CharacterList, error) {
	if err := Validate(params); err != nil {
		return CharacterList{}, err
	}

	var wg sync.WaitGroup
	linesChannel := make(chan string, 1)
	outputChannel := make(chan string, 1)
	go i.ds.ReadConcurrent(file, params, linesChannel)
	wg.Add(2)
	go func() {
		defer wg.Done()
		worker(1, outputChannel, linesChannel, params)
	}()
	go func() {
		defer wg.Done()
		worker(2, outputChannel, linesChannel, params)
	}()
	wg.Wait()

	for consumedLine := range outputChannel {
		fmt.Println(consumedLine)
	}
	fmt.Println("finished")
	return CharacterList{}, nil
}
