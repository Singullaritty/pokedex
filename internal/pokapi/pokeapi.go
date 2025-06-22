package pokapi

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Singullaritty/pokedexcli/internal/pokecache"
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
	cache := pokecache.NewCache(5 * time.Second)

	if cacheData, ok := cache.Get(url); !ok {
		response, err := http.Get(url)
		if err != nil {
			return Locations{}, err
		}
		defer response.Body.Close()
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return Locations{}, err
		}
		errs := json.Unmarshal(responseData, &locs)
		if errs != nil {
			return Locations{}, err
		}
		cache.Add(url, responseData)
	} else {
		err := json.Unmarshal(cacheData, &locs)
		if err != nil {
			return Locations{}, err
		}

	}
	return locs, nil
}
