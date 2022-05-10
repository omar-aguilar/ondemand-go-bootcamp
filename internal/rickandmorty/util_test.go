package rickandmorty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getJSONString(data interface{}) string {
	var buffer bytes.Buffer
	json.NewEncoder(&buffer).Encode(data)
	return buffer.String()
}

func TestJSONEncode(t *testing.T) {
	type testCase struct {
		name           string
		data           interface{}
		expectedOutput string
		expectedErr    error
	}

	mockCharacterInput := Character{
		ID:      1,
		Name:    "Rick Sanchez",
		Species: "Human",
		Type:    "",
		Gender:  "Male",
		Image:   "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
		Url:     "https://rickandmortyapi.com/api/character/1",
		Created: "2017-11-04T18:48:46.250Z",
	}
	mockCharacterListInput := CharacterList{mockCharacterInput}
	mockCharacterOutput := getJSONString(mockCharacterInput)
	mockCharacterListOutput := getJSONString(mockCharacterListInput)

	testCases := []testCase{
		{
			name:           "encodes a character list",
			data:           mockCharacterListInput,
			expectedOutput: mockCharacterListOutput,
			expectedErr:    nil,
		},
		{
			name:           "encodes a single character",
			data:           mockCharacterInput,
			expectedOutput: mockCharacterOutput,
			expectedErr:    nil,
		},
		{
			name:           "fails with invalid data type",
			data:           []string{"test"},
			expectedErr:    ErrInvalidCodecFormat,
			expectedOutput: "",
		},
	}

	jsonCodec := NewJSONCharacterCodec()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			err := jsonCodec.Encode(&buffer, tc.data)
			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Equal(t, tc.expectedOutput, buffer.String())
		})
	}
}

func TestJSONDecode(t *testing.T) {
	type testCase struct {
		name           string
		reader         io.Reader
		dataReceiver   interface{}
		expectedOutput interface{}
		expectedErr    error
	}

	mockCharacterInput := `{
		"id": 1,
		"name": "Rick Sanchez",
		"species": "Human",
		"type": "",
		"gender": "Male",
		"image": "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
		"url": "https://rickandmortyapi.com/api/character/1",
		"created": "2017-11-04T18:48:46.250Z"
	}`
	mockCharacterListInput := fmt.Sprintf("[%s]", mockCharacterInput)

	mockCharacterOutput := Character{
		ID:      1,
		Name:    "Rick Sanchez",
		Species: "Human",
		Type:    "",
		Gender:  "Male",
		Image:   "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
		Url:     "https://rickandmortyapi.com/api/character/1",
		Created: "2017-11-04T18:48:46.250Z",
	}
	mockCharacterListOutput := CharacterList{mockCharacterOutput}

	testCases := []testCase{
		{
			name:           "decodes a character list",
			dataReceiver:   &CharacterList{},
			reader:         strings.NewReader(mockCharacterListInput),
			expectedErr:    nil,
			expectedOutput: &mockCharacterListOutput,
		},
		{
			name:           "decodes a single character",
			dataReceiver:   &Character{},
			reader:         strings.NewReader(mockCharacterInput),
			expectedErr:    nil,
			expectedOutput: &mockCharacterOutput,
		},
		{
			name:           "fails to decode a valid type",
			reader:         strings.NewReader("test"),
			dataReceiver:   &Character{},
			expectedErr:    ErrDecodeInput,
			expectedOutput: &Character{},
		},
		{
			name:           "fails with invalid data receiver",
			reader:         strings.NewReader("test"),
			dataReceiver:   &[]string{"test"},
			expectedErr:    ErrInvalidCodecFormat,
			expectedOutput: &[]string{"test"},
		},
	}

	jsonCodec := NewJSONCharacterCodec()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := jsonCodec.Decode(tc.reader, tc.dataReceiver)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedOutput, tc.dataReceiver)
		})
	}
}

func TestCSVEncode(t *testing.T) {
	type testCase struct {
		name           string
		data           interface{}
		expectedOutput string
		expectedErr    error
	}

	mockCharacterInput := Character{
		ID:      1,
		Name:    "Rick Sanchez",
		Species: "Human",
		Type:    "",
		Gender:  "Male",
		Image:   "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
		Url:     "https://rickandmortyapi.com/api/character/1",
		Created: "2017-11-04T18:48:46.250Z",
	}
	mockCharacterListInput := CharacterList{mockCharacterInput}
	mockCharacterOutput := "ID,Name,Species,Type,Gender,Image,Url,Created\n1,Rick Sanchez,Human,,Male,https://rickandmortyapi.com/api/character/avatar/1.jpeg,https://rickandmortyapi.com/api/character/1,2017-11-04T18:48:46.250Z\n"
	mockCharacterListOutput := mockCharacterOutput

	testCases := []testCase{
		{
			name:           "encodes a character list",
			data:           mockCharacterListInput,
			expectedErr:    nil,
			expectedOutput: mockCharacterListOutput,
		},
		{
			name:           "encodes a single character",
			data:           mockCharacterInput,
			expectedErr:    nil,
			expectedOutput: mockCharacterOutput,
		},
		{
			name:           "fails with invalid data type",
			data:           []string{"test"},
			expectedErr:    ErrInvalidCodecFormat,
			expectedOutput: "",
		},
	}

	csvCodec := NewCSVCharacterCodec()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			err := csvCodec.Encode(&buffer, tc.data)
			assert.ErrorIs(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedOutput, buffer.String())
		})
	}
}

func TestCSVDecode(t *testing.T) {
	type testCase struct {
		name           string
		reader         io.Reader
		dataReceiver   interface{}
		expectedOutput interface{}
		expectedErr    error
	}

	mockCharacterInput := "ID,Name,Species,Type,Gender,Image,Url,Created\n1,Rick Sanchez,Human,,Male,https://rickandmortyapi.com/api/character/avatar/1.jpeg,https://rickandmortyapi.com/api/character/1,2017-11-04T18:48:46.250Z\n"
	mockCharacterListInput := mockCharacterInput

	mockCharacterOutput := Character{
		ID:      1,
		Name:    "Rick Sanchez",
		Species: "Human",
		Type:    "",
		Gender:  "Male",
		Image:   "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
		Url:     "https://rickandmortyapi.com/api/character/1",
		Created: "2017-11-04T18:48:46.250Z",
	}
	mockCharacterListOutput := CharacterList{mockCharacterOutput}

	testCases := []testCase{
		{
			name:           "decodes a character list",
			dataReceiver:   &CharacterList{},
			reader:         strings.NewReader(mockCharacterListInput),
			expectedErr:    nil,
			expectedOutput: &mockCharacterListOutput,
		},
		{
			name:           "decodes a single character",
			dataReceiver:   &Character{},
			reader:         strings.NewReader(mockCharacterInput),
			expectedErr:    nil,
			expectedOutput: &mockCharacterOutput,
		},
		{
			name:           "fails with empty csv",
			reader:         strings.NewReader(""),
			dataReceiver:   &Character{},
			expectedErr:    ErrEmptyCSV,
			expectedOutput: &Character{},
		},
		{
			name:           "fails to decode a valid type",
			reader:         strings.NewReader(`1,"`),
			dataReceiver:   &Character{},
			expectedErr:    ErrDecodeInput,
			expectedOutput: &Character{},
		},
		{
			name:           "fails with invalid data receiver",
			reader:         strings.NewReader("test"),
			dataReceiver:   &[]string{"test"},
			expectedErr:    ErrInvalidCodecFormat,
			expectedOutput: &[]string{"test"},
		},
	}

	csvCodec := NewCSVCharacterCodec()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := csvCodec.Decode(tc.reader, tc.dataReceiver)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedOutput, tc.dataReceiver)
		})
	}
}
