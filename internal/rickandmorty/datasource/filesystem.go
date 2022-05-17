package datasource

import (
	"fmt"
	"os"

	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/config"
	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/rickandmorty"
)

type fsImpl struct {
	config config.Config
}

func NewFileSystemDS(config config.Config) rickandmorty.CharacterStorer {
	return fsImpl{
		config,
	}
}

func (f fsImpl) getFilePath(filename string) string {
	return fmt.Sprintf("%s/%s", f.config.StoreFolder, filename)
}

func (f fsImpl) Write(filename string, characterList rickandmorty.CharacterList, format string) error {
	var characterCodec rickandmorty.CharacterCodec
	switch format {
	case rickandmorty.FormatCSV:
		characterCodec = rickandmorty.NewCSVCharacterCodec()
	case rickandmorty.FormatJSON:
		characterCodec = rickandmorty.NewJSONCharacterCodec()
	default:
		return rickandmorty.ErrInvalidFormat
	}

	filepath := f.getFilePath(filename)
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	return characterCodec.Encode(file, characterList)
}

func (f fsImpl) Read(filename string, characterList *rickandmorty.CharacterList, format string) error {
	filepath := f.getFilePath(filename)
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0222)
	if err != nil {
		return err
	}
	defer file.Close()

	switch format {
	case rickandmorty.FormatCSV:
		csvCodec := rickandmorty.NewCSVCharacterCodec()
		return csvCodec.Decode(file, characterList)
	case rickandmorty.FormatJSON:
		jsonCodec := rickandmorty.NewJSONCharacterCodec()
		return jsonCodec.Decode(file, characterList)
	default:
		return rickandmorty.ErrInvalidFormat
	}
}
