package datasource

import (
	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/rickandmorty"
)

type MemoryDS interface {
	rickandmorty.CharacterGetter
	UpsertDB(entries rickandmorty.CharacterList)
}

type DB = map[rickandmorty.CharacterID]rickandmorty.Character

type dbImpl struct {
	db DB
}

func NewMemoryDS() MemoryDS {
	return &dbImpl{
		db: DB{},
	}
}

func (d *dbImpl) UpsertDB(entries []rickandmorty.Character) {
	for _, character := range entries {
		d.db[character.ID] = character
	}
}

func (d dbImpl) GetById(ID int) (rickandmorty.Character, error) {
	character, found := d.db[ID]
	if !found {
		return rickandmorty.Character{}, rickandmorty.ErrCharacterNotFound
	}
	return character, nil
}
