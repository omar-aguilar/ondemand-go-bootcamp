package transport

import (
	"net/http"

	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/config"
	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/rickandmorty"
	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/rickandmorty/datasource"
)

var interactor rickandmorty.Interactor

func init() {
	config := config.GetConfig()
	memoryStore := datasource.NewMemoryDS()
	characterDS := datasource.NewCSVDS(config, memoryStore)
	apiDS := datasource.NewApiDS(http.DefaultClient)
	fsDS := datasource.NewFileSystemDS(config)
	interactor = rickandmorty.NewInteractor(config, characterDS, apiDS, fsDS)
}
