package transport

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/rickandmorty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCLILoadCSVFromFileName(t *testing.T) {
	type testCase struct {
		name          string
		loadError     error
		expectedError error
		filename      string
	}

	tempFile, _ := ioutil.TempFile(os.TempDir(), "")
	tempFile.Write([]byte("test"))
	tempFilename := tempFile.Name()
	tempFile.Close()

	tempEmptyFile, _ := ioutil.TempFile(os.TempDir(), "")
	tempEmptyFilename := tempEmptyFile.Name()
	tempEmptyFile.Close()

	unexistentFilename := "unexistent-test-file.abc"
	_, errUnexistentFile := os.Open(unexistentFilename)

	errLoad := errors.New("error in load and store")

	testCases := []testCase{
		{
			name:     "loads file successfully",
			filename: tempFilename,
		},
		{
			name:          "fails when file is empty",
			filename:      tempEmptyFilename,
			expectedError: ErrEmptyFile,
		},
		{
			name:          "fails to load when the file is invalid",
			filename:      unexistentFilename,
			expectedError: errUnexistentFile,
		},
		{
			name:          "fails when load and store sends an error",
			filename:      tempFilename,
			loadError:     errLoad,
			expectedError: errLoad,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			iMock := interactorMock{}
			interactor = &iMock
			iMock.On("LoadAndStore", mock.Anything).Return(tc.loadError)
			err := CLILoadCSVFromFileName(tc.filename)
			assert.Equal(t, tc.expectedError, err)
		})
	}

	os.Remove(tempFilename)
	os.Remove(tempEmptyFilename)
}

func TestCLILoadCSVFromStdin(t *testing.T) {
	type testCase struct {
		name          string
		loadError     error
		expectedError error
		stdinInput    *os.File
	}

	tempFile, _ := ioutil.TempFile(os.TempDir(), "")
	tempFile.Write([]byte("test"))
	tempFilename := tempFile.Name()
	defer tempFile.Close()

	tempEmptyFile, _ := ioutil.TempFile(os.TempDir(), "")
	tempEmptyFilename := tempEmptyFile.Name()
	defer tempEmptyFile.Close()
	errLoad := errors.New("error in load and store")

	testCases := []testCase{
		{
			name:       "loads file successfully",
			stdinInput: tempFile,
		},
		{
			name:          "fails when file is empty",
			stdinInput:    tempEmptyFile,
			expectedError: ErrEmptyFile,
		},
		{
			name:          "fails when file is nil",
			expectedError: ErrFileStats,
		},
		{
			name:          "fails when load and store sends an error",
			stdinInput:    tempFile,
			loadError:     errLoad,
			expectedError: errLoad,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			iMock := interactorMock{}
			interactor = &iMock
			iMock.On("LoadAndStore", mock.Anything).Return(tc.loadError)
			err := CLILoadCSVFromStdin(tc.stdinInput)
			assert.Equal(t, tc.expectedError, err)
		})
	}

	os.Remove(tempFilename)
	os.Remove(tempEmptyFilename)
}

func TestCLIGetCharacterById(t *testing.T) {
	type testCase struct {
		name          string
		expectedError error
		ID            int
	}

	errGet := errors.New("error getting character")

	testCases := []testCase{
		{
			name: "gets character successfully",
			ID:   1,
		},
		{
			name:          "fails when get by if fails",
			ID:            0,
			expectedError: errGet,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			iMock := interactorMock{}
			interactor = &iMock
			iMock.On("GetById", mock.Anything).Return(rickandmorty.Character{}, tc.expectedError)
			err := CLIGetCharacterById(tc.ID, rickandmorty.FormatJSON)
			assert.Equal(t, tc.expectedError, err)
		})
	}

}
