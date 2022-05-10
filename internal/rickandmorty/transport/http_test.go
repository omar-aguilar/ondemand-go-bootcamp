package transport

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/rickandmorty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHTTPLoadCharactersCSV(t *testing.T) {
	type testCase struct {
		name        string
		fileContent io.Reader
		loadError   error
		statusCode  int
	}

	errLoad := errors.New("load error")
	testCases := []testCase{
		{
			name:        "fails with bad request when there is an error in form file",
			fileContent: nil,
			statusCode:  400,
		},
		{
			name:        "fails with bad request when there is an error in load an store",
			fileContent: strings.NewReader("test"),
			loadError:   errLoad,
			statusCode:  400,
		},
		{
			name:        "sends ok when loads file correctly",
			fileContent: strings.NewReader("test"),
			statusCode:  200,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			iMock := interactorMock{}
			interactor = &iMock
			iMock.On("LoadAndStore", mock.Anything).Return(tc.loadError)
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			if tc.fileContent != nil {
				fmt.Println("tc.name", tc.name)
				part, _ := writer.CreateFormFile(formFileCSV, "test.csv")
				io.Copy(part, tc.fileContent)
				writer.Close()
			}
			r := httptest.NewRequest(http.MethodPost, "/rickandmorty/load-csv", body)
			r.Header.Set("Content-Type", writer.FormDataContentType())
			r.Header.Set("Content-Length", fmt.Sprintf("%d", body.Len()))
			w := httptest.NewRecorder()
			HTTPLoadCharactersCSV(w, r)
			assert.Equal(t, tc.statusCode, w.Code)
		})
	}
}
func TestHTTPGetCharacterById(t *testing.T) {
	type testCase struct {
		name       string
		statusCode int
		getError   error
		ID         int
	}

	errGet := errors.New("error getting item")
	testCases := []testCase{
		{
			name:       "sends bad request when ID param is not found in the url",
			statusCode: 400,
		},
		{
			name:       "sends not found when get by id fails",
			statusCode: 404,
			ID:         1,
			getError:   errGet,
		},
		{
			name:       "sends ok when everything works as expected",
			statusCode: 200,
			ID:         1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			iMock := interactorMock{}
			interactor = &iMock
			iMock.On("GetById", mock.Anything).Return(rickandmorty.Character{}, tc.getError)
			url := fmt.Sprintf("/rickandmorty/character/%d", tc.ID)
			r := httptest.NewRequest(http.MethodGet, url, nil)
			routerCtx := chi.NewRouteContext()
			if tc.ID != 0 {
				routerCtx.URLParams.Add(paramID, strconv.Itoa(tc.ID))
			}
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, routerCtx))
			w := httptest.NewRecorder()
			HTTPGetCharacterById(w, r)
			assert.Equal(t, tc.statusCode, w.Code)
		})
	}
}
func TestHTTPGetCharactersFromAPI(t *testing.T) {
	type testCase struct {
		name       string
		statusCode int
		storeError error
		ID         int
	}

	errStore := errors.New("error storing")
	testCases := []testCase{
		{
			name:       "sends bad request when ID param is not found in the url",
			statusCode: 400,
		},
		{
			name:       "sends failed dependency when store fails",
			statusCode: 424,
			ID:         1,
			storeError: errStore,
		},
		{
			name:       "sends ok when everything works as expected",
			statusCode: 200,
			ID:         1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			iMock := interactorMock{}
			interactor = &iMock
			iMock.On("StoreCharactersByPageFromAPI", mock.Anything, mock.Anything).Return(rickandmorty.CharacterList{}, tc.storeError)
			url := fmt.Sprintf("/rickandmorty/character/%d", tc.ID)
			r := httptest.NewRequest(http.MethodGet, url, nil)
			routerCtx := chi.NewRouteContext()
			if tc.ID != 0 {
				routerCtx.URLParams.Add(paramPage, strconv.Itoa(tc.ID))
			}
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, routerCtx))
			w := httptest.NewRecorder()
			HTTPGetCharactersFromAPI(w, r)
			assert.Equal(t, tc.statusCode, w.Code)
		})
	}
}
func TestHTTPGetCharactersStoredFromAPI(t *testing.T) {
	type testCase struct {
		name              string
		statusCode        int
		getFromStoreError error
		ID                int
	}

	errStore := errors.New("error storing")
	testCases := []testCase{
		{
			name:       "sends bad request when ID param is not found in the url",
			statusCode: 400,
		},
		{
			name:              "sends failed dependency when get from store fails",
			statusCode:        424,
			ID:                1,
			getFromStoreError: errStore,
		},
		{
			name:       "sends ok when everything works as expected",
			statusCode: 200,
			ID:         1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			iMock := interactorMock{}
			interactor = &iMock
			iMock.On("GetCharactersStoredByPageFromAPI", mock.Anything, mock.Anything).Return(rickandmorty.CharacterList{}, tc.getFromStoreError)
			url := fmt.Sprintf("/rickandmorty/character/%d?%s=%s", tc.ID, queryOutputFormat, rickandmorty.FormatCSV)
			r := httptest.NewRequest(http.MethodGet, url, nil)
			routerCtx := chi.NewRouteContext()
			if tc.ID != 0 {
				routerCtx.URLParams.Add(paramPage, strconv.Itoa(tc.ID))
			}
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, routerCtx))
			w := httptest.NewRecorder()
			HTTPGetCharactersStoredFromAPI(w, r)
			assert.Equal(t, tc.statusCode, w.Code)
		})
	}
}
