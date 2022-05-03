package rickandmorty

import (
	"io"
)

type Interactor interface {
	Load(file io.Reader) error
	GetById(ID int) (Character, error)
}

type interactor struct {
	ds CharacterRepository
}

func NewInteractor(ds CharacterRepository) interactor {
	return interactor{
		ds,
	}
}

func (i interactor) Load(file io.Reader) error {
	err := i.ds.Load(file)
	if err != nil {
		return err
	}
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
