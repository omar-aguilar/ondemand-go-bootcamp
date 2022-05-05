package datasource

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/rickandmorty"
)

type api struct{}

func NewApiDS() rickandmorty.APIGetter {
	return api{}
}

func (a api) GetCharactersByPage(page int) (rickandmorty.CharacterList, error) {
	baseURL := "https://rickandmortyapi.com/api/character"
	requestURL, _ := url.Parse(baseURL)
	params := url.Values{}
	params.Add("page", strconv.Itoa(page))
	requestURL.RawQuery = params.Encode()

	response, err := http.Get(requestURL.String())
	if err != nil {
		return rickandmorty.CharacterList{}, err
	}

	api := rickandmorty.API{}
	err = json.NewDecoder(response.Body).Decode(&api)
	if err != nil {
		return rickandmorty.CharacterList{}, err
	}
	return api.Results, nil
}
