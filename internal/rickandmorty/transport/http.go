package transport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/rickandmorty"
)

const (
	formFileCSV          = "csv"
	paramID              = "id"
	paramPage            = "page"
	queryOutputFormat    = "outputFormat"
	queryStorageFormat   = "inputFormat"
	queryType            = "type"
	queryItems           = "items"
	queryItemsPerWorker  = "itemsPerWorker"
	queryNumberOfWorkers = "numberOfWorkers"
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

func HTTPReadConcurrent(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile(formFileCSV)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	outputFormat := r.URL.Query().Get(queryOutputFormat)
	typee := r.URL.Query().Get(queryType)
	items, _ := strconv.Atoi(r.URL.Query().Get(queryItems))
	itemsPerWorker, _ := strconv.Atoi(r.URL.Query().Get(queryItemsPerWorker))
	numberOfWorkers, _ := strconv.Atoi(r.URL.Query().Get(queryNumberOfWorkers))

	params := rickandmorty.ReadConcurrentParams{
		Type:            typee,
		Items:           items,
		ItemsPerWorker:  itemsPerWorker,
		NumberOfWorkers: numberOfWorkers,
	}
	characterList, err := interactor.ReadConcurrent(file, params)
	if err != nil {
		if errs := rickandmorty.Translate(err); errs != nil {
			w.WriteHeader(http.StatusFailedDependency)
			json.NewEncoder(w).Encode(errs)
			return
		}
		http.Error(w, err.Error(), http.StatusFailedDependency)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeFormattedResponse(w, characterList, outputFormat)
}
