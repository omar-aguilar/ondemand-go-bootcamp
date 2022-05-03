package rickandmorty

import (
	"errors"
	"io"
)

var (
	ErrCharacterNotFound = errors.New("character not found")
	ErrInvalidID         = errors.New("id should be greater than 0")
)

type CharacterGetter interface {
	GetById(ID int) (Character, error)
}

type CharacterRepository interface {
	CharacterGetter
	Load(file io.Reader) error
}
