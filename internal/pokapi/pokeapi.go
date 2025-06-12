package pokapi

import (
	"encoding/json"
	"io"
	"net/http"
)

type Locations struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetAreas(url string) (Locations, error) {
	locs := Locations{}

	response, err := http.Get(url)
	if err != nil {
		return Locations{}, err
	}

	responseData, err := io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return Locations{}, err
	}

	errs := json.Unmarshal(responseData, &locs)
	if errs != nil {
		return Locations{}, err
	}

	return locs, nil
}
