package datasource

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/config"
	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/rickandmorty"
	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	type testCase struct {
		name           string
		filename       string
		format         string
		expectedError  error
		expectedOutput rickandmorty.CharacterList
	}

	mockCharacterList := rickandmorty.CharacterList{
		{
			ID:      1,
			Name:    "Rick Sanchez",
			Species: "Human",
			Type:    "",
			Gender:  "Male",
			Image:   "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
			Url:     "https://rickandmortyapi.com/api/character/1",
			Created: "2017-11-04T18:48:46.250Z",
		},
	}

	storeFolder := os.TempDir()

	mockJSONFile, _ := ioutil.TempFile(storeFolder, "")
	mockJSONFileName := filepath.Base(mockJSONFile.Name())
	rickandmorty.NewJSONCharacterCodec().Encode(mockJSONFile, mockCharacterList)
	mockJSONFile.Close()
	defer os.Remove(mockJSONFile.Name())

	mockCSVFile, _ := ioutil.TempFile(storeFolder, "")
	mockCSVFileName := filepath.Base(mockCSVFile.Name())
	rickandmorty.NewCSVCharacterCodec().Encode(mockCSVFile, mockCharacterList)
	mockCSVFile.Close()
	defer os.Remove(mockCSVFile.Name())

	testCases := []testCase{
		{
			name:           "succeed to decode in json file",
			format:         rickandmorty.FormatJSON,
			filename:       mockJSONFileName,
			expectedOutput: mockCharacterList,
		},
		{
			name:           "succeed to decode in csv file",
			format:         rickandmorty.FormatCSV,
			filename:       mockCSVFileName,
			expectedOutput: mockCharacterList,
		},
		{
			name:           "fails on invalid format",
			format:         "some-unsupported-format",
			expectedError:  rickandmorty.ErrInvalidFormat,
			expectedOutput: rickandmorty.CharacterList{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := config.Config{
				StoreFolder: storeFolder,
			}
			fsDS := NewFileSystemDS(cfg)
			output := rickandmorty.CharacterList{}
			err := fsDS.Read(tc.filename, &output, tc.format)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedOutput, output)
		})
	}
}

func TestWrite(t *testing.T) {
	type testCase struct {
		name          string
		filename      string
		format        string
		expectedError error
	}

	mockCharacterList := rickandmorty.CharacterList{
		{
			ID:      1,
			Name:    "Rick Sanchez",
			Species: "Human",
			Type:    "",
			Gender:  "Male",
			Image:   "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
			Url:     "https://rickandmortyapi.com/api/character/1",
			Created: "2017-11-04T18:48:46.250Z",
		},
	}

	storeFolder := os.TempDir()

	mockJSONFile, _ := ioutil.TempFile(storeFolder, "")
	mockJSONFileName := filepath.Base(mockJSONFile.Name())
	mockJSONFile.Close()
	defer os.Remove(mockJSONFile.Name())

	mockCSVFile, _ := ioutil.TempFile(storeFolder, "")
	mockCSVFileName := filepath.Base(mockCSVFile.Name())
	mockCSVFile.Close()
	defer os.Remove(mockCSVFile.Name())

	testCases := []testCase{
		{
			name:     "succeed to decode in json file",
			format:   rickandmorty.FormatJSON,
			filename: mockJSONFileName,
		},
		{
			name:     "succeed to decode in csv file",
			format:   rickandmorty.FormatCSV,
			filename: mockCSVFileName,
		},
		{
			name:          "fails on invalid format",
			format:        "some-unsupported-format",
			expectedError: rickandmorty.ErrInvalidFormat,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := config.Config{
				StoreFolder: storeFolder,
			}
			fsDS := NewFileSystemDS(cfg)
			err := fsDS.Write(tc.filename, mockCharacterList, tc.format)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
