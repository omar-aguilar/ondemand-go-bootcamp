package transport

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

const (
	formFileCSV        = "csv"
	paramID            = "id"
	paramPage          = "page"
	queryOutputFormat  = "outputFormat"
	queryStorageFormat = "inputFormat"
)

func HTTPLoadCharactersCSV(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile(formFileCSV)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	err = interactor.LoadAndStore(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, "file loaded correctly")
}

func HTTPGetCharacterById(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.Atoi(chi.URLParam(r, paramID))
	outputFormat := r.URL.Query().Get(queryOutputFormat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	character, err := interactor.GetById(ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeFormattedResponse(w, character, outputFormat)
}

func HTTPGetCharactersFromAPI(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(chi.URLParam(r, paramPage))
	storageFormat := r.URL.Query().Get(queryStorageFormat)
	outputFormat := r.URL.Query().Get(queryOutputFormat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	characterList, err := interactor.StoreCharactersByPageFromAPI(page, storageFormat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusFailedDependency)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeFormattedResponse(w, characterList, outputFormat)
}

func HTTPGetCharactersStoredFromAPI(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(chi.URLParam(r, paramPage))
	storageFormat := r.URL.Query().Get(queryStorageFormat)
	outputFormat := r.URL.Query().Get(queryOutputFormat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	characterList, err := interactor.GetCharactersStoredByPageFromAPI(page, storageFormat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusFailedDependency)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeFormattedResponse(w, characterList, outputFormat)
}
