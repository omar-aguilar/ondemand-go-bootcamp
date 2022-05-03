package transport

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func HTTPLoadCSV(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("csv")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	err = interactor.Load(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, "file loaded correctly")
}

func HTTPGetCharacterById(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.Atoi(chi.URLParam(r, "id"))
	format := r.URL.Query().Get("format")
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
	writeFormattedCharacter(w, character, format)
}
