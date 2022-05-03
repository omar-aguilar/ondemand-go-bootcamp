package transport

import (
	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/config"
	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/rickandmorty"
	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/rickandmorty/datasource"
)

var interactor rickandmorty.Interactor

func init() {
	config := config.GetConfig()
	csvSource := config.CSVSource
	memoryStore := datasource.NewMemoryDS()
	datastore := datasource.NewCSVDS(csvSource, memoryStore)
	interactor = rickandmorty.NewInteractor(datastore)
}
