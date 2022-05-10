package datasource

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/config"
	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/rickandmorty"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	type testCase struct {
		name          string
		input         io.Reader
		expectedError error
	}

	mockCharacterDB := rickandmorty.CharacterList{
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

	var buff bytes.Buffer
	rickandmorty.NewCSVCharacterCodec().Encode(bufio.NewWriter(&buff), mockCharacterDB)
	fileMock := bufio.NewReader(&buff)

	csvWithErrorStr := `1,"`
	buffWithError := strings.NewReader(csvWithErrorStr)
	errDecode := rickandmorty.NewCSVCharacterCodec().Decode(buffWithError, &rickandmorty.CharacterList{})

	testCases := []testCase{
		{
			name:  "successfully reads an input file",
			input: fileMock,
		},
		{
			name:          "fails when cannot decode db",
			input:         strings.NewReader(csvWithErrorStr),
			expectedError: errDecode,
		},
	}

	storeFolder := os.TempDir()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := config.Config{
				StoreFolder: storeFolder,
			}
			memoryDS := NewMemoryDS()
			csvDS := NewCSVDS(cfg, memoryDS)
			_, err := csvDS.Load(tc.input)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetByID(t *testing.T) {
	type testCase struct {
		name          string
		ID            int
		expectedError error
	}

	mockCharacterDB := rickandmorty.CharacterList{
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

	testCases := []testCase{
		{
			name: "successfully reads an input file",
			ID:   1,
		},
		{
			name:          "returns error when id is not found",
			ID:            2,
			expectedError: rickandmorty.ErrCharacterNotFound,
		},
	}

	storeFolder := os.TempDir()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := config.Config{
				StoreFolder: storeFolder,
			}
			memoryDS := NewMemoryDS()
			csvDS := NewCSVDS(cfg, memoryDS)

			var buff bytes.Buffer
			rickandmorty.NewCSVCharacterCodec().Encode(bufio.NewWriter(&buff), mockCharacterDB)
			dbMock := bufio.NewReader(&buff)
			csvDS.Load(dbMock)
			_, err := csvDS.GetById(tc.ID)
			assert.Equal(t, tc.expectedError, err)
		})
	}

}
