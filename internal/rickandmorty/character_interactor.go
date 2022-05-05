package rickandmorty

import (
	"fmt"
	"io"

	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/config"
)

type Interactor interface {
	LoadAndStore(file io.Reader) error
	GetById(ID int) (Character, error)
	StoreCharactersByPageFromAPI(page int, storageFormat string) (CharacterList, error)
	GetCharactersStoredByPageFromAPI(page int, storageFormat string) (CharacterList, error)
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

	i.storeDS.Write(i.config.DBFile, characterList, FormatCSV)
	return nil
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
