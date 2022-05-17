package rickandmorty

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"strconv"
)

var (
	ErrInvalidCodecFormat = errors.New("invalid codec format")
	ErrEmptyCSV           = errors.New("empty csv")
	ErrDecodeInput        = errors.New("cannot decode input")
)

type CharacterCodec interface {
	Encode(writer io.Writer, data interface{}) error
	Decode(reader io.Reader, data interface{}) error
}

type jsonCodec struct{}

func NewJSONCharacterCodec() CharacterCodec {
	return jsonCodec{}
}
func (j jsonCodec) Encode(writer io.Writer, data interface{}) error {
	switch data.(type) {
	case Character, CharacterList:
		return json.NewEncoder(writer).Encode(data)
	default:
		return ErrInvalidCodecFormat
	}
}
func (j jsonCodec) Decode(reader io.Reader, data interface{}) error {
	switch data.(type) {
	case *Character, *CharacterList:
		err := json.NewDecoder(reader).Decode(data)
		if err != nil {
			return ErrDecodeInput
		}
		return nil
	default:
		return ErrInvalidCodecFormat
	}
}

type csvCodec struct {
	header []string
}

func NewCSVCharacterCodec() CharacterCodec {
	return csvCodec{
		header: []string{"ID", "Name", "Species", "Type", "Gender", "Image", "Url", "Created"},
	}
}
func (j csvCodec) getCSVRowFromCharacter(character Character) []string {
	return []string{
		strconv.Itoa(character.ID),
		character.Name,
		character.Species,
		character.Type,
		character.Gender,
		character.Image,
		character.Url,
		character.Created,
	}
}
func (j csvCodec) getCharacterFromCSVRow(characterRow []string) Character {
	id, _ := strconv.Atoi(characterRow[0])
	return Character{
		ID:      id,
		Name:    characterRow[1],
		Species: characterRow[2],
		Type:    characterRow[3],
		Gender:  characterRow[4],
		Image:   characterRow[5],
		Url:     characterRow[6],
		Created: characterRow[7],
	}
}
func (j csvCodec) Encode(writer io.Writer, data interface{}) error {
	csvHeader := j.header
	csvData := [][]string{}
	csvData = append(csvData, csvHeader)
	switch dataOfType := data.(type) {
	case CharacterList:
		for _, character := range dataOfType {
			row := j.getCSVRowFromCharacter(character)
			csvData = append(csvData, row)
		}
	case Character:
		row := j.getCSVRowFromCharacter(dataOfType)
		csvData = append(csvData, row)
	default:
		return ErrInvalidCodecFormat
	}
	csvWriter := csv.NewWriter(writer)
	return csvWriter.WriteAll(csvData)
}

func (j csvCodec) Decode(reader io.Reader, data interface{}) error {
	csvHeader := j.header
	csvReader := csv.NewReader(reader)
	csvData, err := csvReader.ReadAll()
	if err != nil {
		return ErrDecodeInput
	}

	if len(csvData) == 0 {
		return ErrEmptyCSV
	}
	hasHeader := len(csvData) > 1 && reflect.DeepEqual(csvHeader, csvData[0])

	switch dataOfType := data.(type) {
	case *Character:
		characterRow := csvData[0]
		if hasHeader {
			characterRow = csvData[1]
		}
		character := j.getCharacterFromCSVRow(characterRow)
		*dataOfType = character
	case *CharacterList:
		characterRowList := csvData
		if hasHeader {
			characterRowList = csvData[1:]
		}
		characterList := CharacterList{}
		for _, characterRow := range characterRowList {
			character := j.getCharacterFromCSVRow(characterRow)
			characterList = append(characterList, character)
		}
		*dataOfType = characterList
	default:
		return ErrInvalidCodecFormat
	}
	return nil
}

func getNumberOfWorkers(numOfWorkers int) int {
	defaultWorkers := 2
	if numOfWorkers <= 0 || numOfWorkers >= 4 {
		return defaultWorkers
	}
	return numOfWorkers
}
