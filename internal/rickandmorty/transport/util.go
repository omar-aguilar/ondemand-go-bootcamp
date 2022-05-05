package transport

import (
	"io"

	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/rickandmorty"
)

func writeFormattedResponse(channel io.Writer, characterData interface{}, format string) {
	switch format {
	case rickandmorty.FormatCSV:
		csvCodec := rickandmorty.NewCSVCharacterCodec()
		csvCodec.Encode(channel, characterData)
	default:
		csvCodec := rickandmorty.NewJSONCharacterCodec()
		csvCodec.Encode(channel, characterData)
	}
}
