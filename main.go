package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/config"
	ram "github.com/omar-aguilar/ondemand-go-bootcamp/internal/rickandmorty/transport"
)

func startRouter() {
	config := config.GetConfig()
	router := chi.NewRouter()
	router.Route("/rickandmorty", func(r chi.Router) {
		r.Post("/load", ram.HTTPLoadCSV)
		r.Get("/{id}", ram.HTTPGetCharacterById)
	})
	addr := fmt.Sprintf("127.0.0.1:%d", config.Port)
	http.ListenAndServe(addr, router)
}

func startCLI() {
	cliMode := flag.Bool("cli-mode", false, "indicates if cli should be used")
	csvFilename := flag.String("csv-filename", "", "initial csv db")
	outputFormat := flag.String("output-format", "json", "the output format (csv|json)")
	charID := flag.Int("char-id", 0, "the id od the character to look for")
	flag.Parse()

	if !*cliMode {
		return
	}

	fileStat, err := os.Stdin.Stat()
	hasStdinInput := err == nil && fileStat.Size() > 0
	hasCsvFilename := *csvFilename != ""
	isCharacterRequest := *charID != 0

	switch {
	case hasStdinInput:
		ram.CLILoadCSVFromStdin()
	case hasCsvFilename:
		ram.CLILoadCSVFromFileName(*csvFilename)
	case isCharacterRequest:
		ram.CLIGetCharacterById(*charID, *outputFormat)
	}
}

func main() {
	startCLI()
	startRouter()
}
