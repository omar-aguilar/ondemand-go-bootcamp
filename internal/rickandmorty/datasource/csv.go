package datasource

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/rickandmorty"
)

var ErrIncompatibleCSV = errors.New("incompatible csv file")

type csvDS struct {
	csvSource   string
	memoryStore MemoryDS
}

func NewCSVDS(csvSource string, memoryStore MemoryDS) rickandmorty.CharacterRepository {
	ds := csvDS{
		csvSource,
		memoryStore,
	}
	ds.init()
	return ds
}

func (d csvDS) init() {
	csvFile, err := os.OpenFile(d.csvSource, os.O_RDWR, 0444)
	if err != nil {
		log.Println("empty csv, please make sure to load a csv first")
		return
	}
	err = d.Load(csvFile)
	if err != nil {
		log.Println("error loading csv file", err.Error())
		return
	}
	log.Println("successfully loaded csv file")
}

func getCharacterFromRow(row []string) *rickandmorty.Character {
	id, err := strconv.Atoi(row[0])
	if err != nil {
		return nil
	}

	return &rickandmorty.Character{
		ID:      id,
		Name:    row[1],
		Species: row[2],
		Type:    row[3],
		Gender:  row[4],
		Image:   row[5],
		Url:     row[6],
		Created: row[7],
	}
}

func (d csvDS) saveDB(header []string, content [][]string) error {
	expectedHeader := []string{"ID", "Name", "Species", "Type", "Gender", "Image", "Url", "Created"}
	if !reflect.DeepEqual(expectedHeader, header) {
		return ErrIncompatibleCSV
	}

	csvFile, err := os.OpenFile(d.csvSource, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0444)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	writer.Write(header)
	writer.WriteAll(content)
	return nil
}

func (d csvDS) fillMemoryStore(content [][]string) {
	characterList := []rickandmorty.Character{}
	for _, row := range content {
		character := getCharacterFromRow(row)
		if character == nil {
			continue
		}
		characterList = append(characterList, *character)
	}
	d.memoryStore.UpsertDB(characterList)
}

func (d csvDS) GetById(ID int) (rickandmorty.Character, error) {
	return d.memoryStore.GetById(ID)
}

func (d csvDS) Load(file io.Reader) error {
	csvReader := csv.NewReader(file)
	data, err := csvReader.ReadAll()
	if err != nil {
		return err
	}
	header := data[0]
	content := data[1:]
	err = d.saveDB(header, content)
	if err != nil {
		return err
	}
	d.fillMemoryStore(content)
	return err
}
