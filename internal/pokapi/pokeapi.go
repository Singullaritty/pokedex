package pokapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Singullaritty/pokedexcli/internal/pokecache"
)

func GetAreas(url string, cache *pokecache.Cache) (Locations, error) {
	locs := Locations{}
	cacheData, ok := cache.Get(url)
	if !ok {
		response, err := http.Get(url)
		if err != nil {
			return locs, err
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			return locs, fmt.Errorf("error: received status code %d", response.StatusCode)
		}
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return locs, err
		}
		errs := json.Unmarshal(responseData, &locs)
		if errs != nil {
			return locs, err
		}
		cache.Add(url, responseData)
		return locs, nil
	}
	err := json.Unmarshal(cacheData, &locs)
	if err != nil {
		return locs, err
	}
	return locs, nil
}

func ExploreArea(exploreUrl string, cache *pokecache.Cache) (LocationsArea, error) {
	locsArea := LocationsArea{}
	cacheData, ok := cache.Get(exploreUrl)
	if !ok {
		response, err := http.Get(exploreUrl)
		if err != nil {
			return locsArea, err
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			return locsArea, fmt.Errorf("received status code %d", response.StatusCode)
		}
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return locsArea, err
		}
		errs := json.Unmarshal(responseData, &locsArea)
		if errs != nil {
			return locsArea, err
		}
		cache.Add(exploreUrl, responseData)
		return locsArea, nil
	}
	err := json.Unmarshal(cacheData, &locsArea)
	if err != nil {
		return locsArea, err
	}
	return locsArea, nil
}

func GetPokemonInfo(pokeUrl string, cache *pokecache.Cache) (Pokemon, error) {
	pokemon := Pokemon{}
	cacheData, ok := cache.Get(pokeUrl)
	if !ok {
		response, err := http.Get(pokeUrl)
		if err != nil {
			return pokemon, err
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			return pokemon, fmt.Errorf("received status code %d", response.StatusCode)
		}
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return pokemon, err
		}
		errs := json.Unmarshal(responseData, &pokemon)
		if errs != nil {
			return pokemon, err
		}
		cache.Add(pokeUrl, responseData)
		return pokemon, nil
	}
	err := json.Unmarshal(cacheData, &pokemon)
	if err != nil {
		return pokemon, err
	}
	return pokemon, nil
}
