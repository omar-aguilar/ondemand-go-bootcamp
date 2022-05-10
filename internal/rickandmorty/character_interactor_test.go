package rickandmorty

import (
	"errors"
	"strings"
	"testing"

	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoadAndStore(t *testing.T) {
	type testCase struct {
		name          string
		loadError     error
		writeError    error
		loadResult    CharacterList
		expectedError error
	}

	errLoad := errors.New("error loading file")
	errWrite := errors.New("error writing file")
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

	testCases := []testCase{
		{
			name:       "returns no error",
			loadResult: CharacterList{mockCharacterInput},
		},
		{
			name:          "returns an error when load fails",
			loadError:     errLoad,
			expectedError: errLoad,
			loadResult:    CharacterList{},
		},
		{
			name:          "returns error when write fails",
			writeError:    errWrite,
			expectedError: errWrite,
			loadResult:    CharacterList{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			configMock := config.Config{}
			characterDSMock := &characterRepositoryMock{}
			apiGetterDSMock := &apiGetterMock{}
			characterStorerDSMock := &characterStorerMock{}
			interactorMock := NewInteractor(configMock, characterDSMock, apiGetterDSMock, characterStorerDSMock)
			characterDSMock.On("Load", mock.Anything).Return(tc.loadResult, tc.loadError)
			characterStorerDSMock.On("Write", mock.Anything, mock.Anything, mock.Anything).Return(tc.writeError)
			err := interactorMock.LoadAndStore(strings.NewReader("test"))
			assert.Equal(t, tc.expectedError, err)
			if err == nil {
				characterStorerDSMock.AssertCalled(t, "Write", configMock.DBFile, tc.loadResult, FormatCSV)
			}
		})
	}
}

func TestGetById(t *testing.T) {
	type testCase struct {
		name           string
		ID             int
		expectedError  error
		expectedOutput Character
	}

	errGet := errors.New("error getting character")
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

	testCases := []testCase{
		{
			name:           "returns valid character",
			expectedOutput: mockCharacterInput,
			ID:             1,
		},
		{
			name:           "returns an error for ids lower than or equal to 0",
			ID:             0,
			expectedError:  ErrInvalidID,
			expectedOutput: Character{},
		},
		{
			name:           "returns an error when get by id fails",
			ID:             1,
			expectedError:  errGet,
			expectedOutput: Character{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			configMock := config.Config{}
			characterDSMock := &characterRepositoryMock{}
			apiGetterDSMock := &apiGetterMock{}
			characterStorerDSMock := &characterStorerMock{}
			interactorMock := NewInteractor(configMock, characterDSMock, apiGetterDSMock, characterStorerDSMock)
			characterDSMock.On("GetById", tc.ID).Return(tc.expectedOutput, tc.expectedError)
			output, err := interactorMock.GetById(tc.ID)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedOutput, output)
		})
	}
}

func TestStoreCharactersByPageFromAPI(t *testing.T) {
	type testCase struct {
		name           string
		page           int
		expectedError  error
		expectedOutput CharacterList
		getByPageError error
		writeError     error
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

	errApiCall := errors.New("api call failed")
	errWrite := errors.New("write error")

	testCases := []testCase{
		{
			name:           "returns stored character list",
			page:           1,
			expectedOutput: CharacterList{mockCharacterInput},
		},
		{
			name:           "returns an error for page lower than or equal to 0",
			page:           0,
			expectedOutput: CharacterList{},
			expectedError:  ErrInvalidPage,
		},
		{
			name:           "returns an error when api call fail",
			page:           1,
			expectedOutput: CharacterList{},
			getByPageError: errApiCall,
			expectedError:  errApiCall,
		},
		{
			name:           "returns an error when write fails",
			page:           1,
			expectedOutput: CharacterList{},
			writeError:     errWrite,
			expectedError:  errWrite,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			configMock := config.Config{}
			characterDSMock := &characterRepositoryMock{}
			apiGetterDSMock := &apiGetterMock{}
			characterStorerDSMock := &characterStorerMock{}
			interactorMock := NewInteractor(configMock, characterDSMock, apiGetterDSMock, characterStorerDSMock)
			apiGetterDSMock.On("GetCharactersByPage", mock.Anything).Return(tc.expectedOutput, tc.getByPageError)
			characterStorerDSMock.On("Write", mock.Anything, mock.Anything, mock.Anything).Return(tc.writeError)
			output, err := interactorMock.StoreCharactersByPageFromAPI(tc.page, FormatJSON)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedOutput, output)
			if err == nil {
				characterStorerDSMock.AssertCalled(t, "Write", getCharactersFilenameFromPage(tc.page, FormatJSON), tc.expectedOutput, FormatJSON)
			}
		})
	}
}

func TestGetCharactersStoredByPageFromAPI(t *testing.T) {
	type testCase struct {
		name           string
		page           int
		expectedError  error
		expectedOutput CharacterList
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

	errGet := errors.New("error getting character list")
	testCases := []testCase{
		{
			name:           "returns valid character list",
			expectedOutput: CharacterList{mockCharacterInput},
			page:           1,
		},
		{
			name:           "returns an error for ids lower than or equal to 0",
			page:           0,
			expectedError:  ErrInvalidPage,
			expectedOutput: CharacterList{},
		},
		{
			name:           "returns an error when get character list fails",
			page:           1,
			expectedError:  errGet,
			expectedOutput: CharacterList{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			configMock := config.Config{}
			characterDSMock := &characterRepositoryMock{}
			apiGetterDSMock := &apiGetterMock{}
			characterStorerDSMock := &characterStorerMock{}
			interactorMock := NewInteractor(configMock, characterDSMock, apiGetterDSMock, characterStorerDSMock)
			characterStorerDSMock.On("Read", mock.Anything, mock.Anything, mock.Anything).
				Return(tc.expectedError).
				Run(func(args mock.Arguments) {
					charListPointer := args.Get(1).(*CharacterList)
					*charListPointer = tc.expectedOutput
				})
			output, err := interactorMock.GetCharactersStoredByPageFromAPI(tc.page, FormatJSON)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedOutput, output)
		})
	}
}
