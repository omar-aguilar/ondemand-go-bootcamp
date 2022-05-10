package datasource

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/rickandmorty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockHTTPClient struct {
	mock.Mock
}

func (h mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := h.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestGetCharactersByPage(t *testing.T) {
	mockCharacterResponse := rickandmorty.API{
		Results: rickandmorty.CharacterList{
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
		},
	}
	mockCharacterResponseBytes, _ := json.Marshal(mockCharacterResponse)
	mockCharacterResponseString := string(mockCharacterResponseBytes)

	type testCase struct {
		name           string
		requestOutput  string
		requestError   error
		expectedError  error
		expectedOutput rickandmorty.CharacterList
	}

	errAPI := errors.New("error from api")
	errDecode := json.NewDecoder(strings.NewReader("test")).Decode(&testCase{})

	testCases := []testCase{
		{
			name:           "returns character list when no errors in request",
			requestOutput:  mockCharacterResponseString,
			expectedOutput: mockCharacterResponse.Results,
		},
		{
			name:           "returns error from api when it fails",
			requestOutput:  "test",
			requestError:   errAPI,
			expectedError:  errAPI,
			expectedOutput: rickandmorty.CharacterList{},
		},
		{
			name:           "returns error when response body is not a valid api response",
			requestOutput:  "test",
			expectedError:  errDecode,
			expectedOutput: rickandmorty.CharacterList{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{}
			api := NewApiDS(mockClient)
			responseMock := &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader(tc.requestOutput)),
			}
			mockClient.On("Do", mock.Anything).Return(responseMock, tc.requestError)
			output, err := api.GetCharactersByPage(1)
			assert.Equal(t, tc.expectedOutput, output)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
