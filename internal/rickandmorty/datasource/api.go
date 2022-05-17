package datasource

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/omar-aguilar/ondemand-go-bootcamp/internal/rickandmorty"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type api struct {
	httpClient HTTPClient
}

func NewApiDS(client HTTPClient) rickandmorty.APIGetter {
	return api{
		httpClient: client,
	}
}

func (a api) GetCharactersByPage(page int) (rickandmorty.CharacterList, error) {
	baseURL := "https://rickandmortyapi.com/api/character"
	requestURL, _ := url.Parse(baseURL)
	params := url.Values{}
	params.Add("page", strconv.Itoa(page))
	requestURL.RawQuery = params.Encode()
	request, _ := http.NewRequest(http.MethodGet, requestURL.String(), nil)
	response, err := a.httpClient.Do(request)
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
