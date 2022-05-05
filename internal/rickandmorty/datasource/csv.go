package datasource

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/config"
	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/rickandmorty"
)

var ErrIncompatibleCSV = errors.New("incompatible csv file")

type csvDS struct {
	config      config.Config
	memoryStore MemoryDS
}

func NewCSVDS(config config.Config, memoryStore MemoryDS) rickandmorty.CharacterRepository {
	ds := csvDS{
		config,
		memoryStore,
	}
	ds.init()
	return ds
}

func (d csvDS) init() {
	csvFile, err := os.OpenFile(d.config.GetDBPath(), os.O_RDWR, 0444)
	if err != nil {
		log.Println("empty csv, please make sure to load a csv first")
		return
	}
	_, err = d.Load(csvFile)
	if err != nil {
		log.Println("error loading csv file", err.Error())
		return
	}
	log.Println("successfully loaded csv file")
}

func (d csvDS) GetById(ID int) (rickandmorty.Character, error) {
	return d.memoryStore.GetById(ID)
}

func (d csvDS) Load(file io.Reader) (rickandmorty.CharacterList, error) {
	csvCodec := rickandmorty.NewCSVCharacterCodec()
	characterList := rickandmorty.CharacterList{}
	err := csvCodec.Decode(file, &characterList)
	if err != nil {
		return rickandmorty.CharacterList{}, err
	}
	d.memoryStore.UpsertDB(characterList)
	return characterList, err
}
