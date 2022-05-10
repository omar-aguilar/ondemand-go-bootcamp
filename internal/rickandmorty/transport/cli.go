package transport

import (
	"errors"
	"fmt"
	"os"
)

var ErrEmptyFile = errors.New("empty file")
var ErrFileStats = errors.New("stats error in file")

func cliLoadCSV(file *os.File) error {
	stats, err := file.Stat()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return ErrFileStats
	}
	isEmptyFile := stats.Size() == 0
	if isEmptyFile {
		return ErrEmptyFile
	}
	return interactor.LoadAndStore(file)
}

func CLILoadCSVFromFileName(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return err
	}
	return cliLoadCSV(file)
}

func CLILoadCSVFromStdin(stdin *os.File) error {
	return cliLoadCSV(stdin)
}

func CLIGetCharacterById(ID int, format string) error {
	character, err := interactor.GetById(ID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return err
	}
	writeFormattedResponse(os.Stdout, character, format)
	return nil
}
