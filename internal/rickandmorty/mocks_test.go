package rickandmorty

import (
	"io"

	"github.com/stretchr/testify/mock"
)

type characterRepositoryMock struct {
	mock.Mock
}

func (c *characterRepositoryMock) GetById(ID int) (Character, error) {
	args := c.Called(ID)
	return args.Get(0).(Character), args.Error(1)
}

func (c *characterRepositoryMock) Load(file io.Reader) (CharacterList, error) {
	args := c.Called(file)
	return args.Get(0).(CharacterList), args.Error(1)
}

func (c *characterRepositoryMock) ReadConcurrent(file io.Reader, params ReadConcurrentParams, lines chan<- string) error {
	args := c.Called(file, params, lines)
	return args.Error(0)
}

type apiGetterMock struct {
	mock.Mock
}

func (a *apiGetterMock) GetCharactersByPage(page int) (CharacterList, error) {
	args := a.Called(page)
	return args.Get(0).(CharacterList), args.Error(1)
}

type characterStorerMock struct {
	mock.Mock
}

func (s *characterStorerMock) Write(filename string, data CharacterList, storageFormat string) error {
	args := s.Called(filename, data, storageFormat)
	return args.Error(0)
}

func (s *characterStorerMock) Read(filename string, data *CharacterList, storageFormat string) error {
	args := s.Called(filename, data, storageFormat)
	return args.Error(0)
}
