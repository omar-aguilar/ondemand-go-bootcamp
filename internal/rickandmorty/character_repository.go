package rickandmorty

import (
	"errors"
	"io"
)

var (
	ErrCharacterNotFound = errors.New("character not found")
	ErrInvalidID         = errors.New("id should be greater than 0")
	ErrInvalidPage       = errors.New("page should be greater than 0")
	ErrInvalidFormat     = errors.New("invalid format")
)

const FormatCSV = "csv"
const FormatJSON = "json"

type CharacterGetter interface {
	GetById(ID int) (Character, error)
}

type CharacterLoader interface {
	Load(file io.Reader) (CharacterList, error)
}

type APIGetter interface {
	GetCharactersByPage(page int) (CharacterList, error)
}

type CharacterStorer interface {
	Write(filename string, data CharacterList, storageFormat string) error
	Read(filename string, data *CharacterList, storageFormat string) error
}

type CharacterRepository interface {
	CharacterGetter
	CharacterLoader
}
