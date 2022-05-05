package transport

import (
	"fmt"
	"os"
)

func cliLoadCSV(file *os.File) {
	stats, err := file.Stat()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	isEmptyFile := stats.Size() == 0
	if isEmptyFile {
		return
	}
	interactor.LoadAndStore(file)
}

func CLILoadCSVFromFileName(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	cliLoadCSV(file)
}

func CLILoadCSVFromStdin() {
	cliLoadCSV(os.Stdin)
}

func CLIGetCharacterById(ID int, format string) {
	character, err := interactor.GetById(ID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	writeFormattedResponse(os.Stdout, character, format)
}
