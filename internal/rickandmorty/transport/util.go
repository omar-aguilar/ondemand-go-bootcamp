package transport

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/rickandmorty"
)

func writeFormattedCharacter(channel io.Writer, character rickandmorty.Character, format string) {
	switch format {
	case "csv":
		fmt.Fprintln(channel, character.ToCSVEntry())
	default:
		json.NewEncoder(channel).Encode(character)
	}
}
