package transport

import (
	"io"

	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/rickandmorty"
	"github.com/stretchr/testify/mock"
)

type interactorMock struct {
	mock.Mock
}

func (i *interactorMock) LoadAndStore(file io.Reader) error {
	args := i.Called(file)
	return args.Error(0)
}

func (i *interactorMock) GetById(ID int) (rickandmorty.Character, error) {
	args := i.Called(ID)
	return args.Get(0).(rickandmorty.Character), args.Error(1)
}

func (i *interactorMock) StoreCharactersByPageFromAPI(page int, storageFormat string) (rickandmorty.CharacterList, error) {
	args := i.Called(page, storageFormat)
	return args.Get(0).(rickandmorty.CharacterList), args.Error(1)
}

func (i *interactorMock) GetCharactersStoredByPageFromAPI(page int, storageFormat string) (rickandmorty.CharacterList, error) {
	args := i.Called(page, storageFormat)
	return args.Get(0).(rickandmorty.CharacterList), args.Error(1)
}
