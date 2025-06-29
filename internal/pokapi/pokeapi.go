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
	if cacheData, ok := cache.Get(url); !ok {
		response, err := http.Get(url)
		if err != nil {
			return locs, err
		}
		defer response.Body.Close()
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return locs, err
		}
		errs := json.Unmarshal(responseData, &locs)
		if errs != nil {
			return locs, err
		}
		cache.Add(url, responseData)
	} else {
		err := json.Unmarshal(cacheData, &locs)
		if err != nil {
			return locs, err
		}
	}
	return locs, nil
}

func ExploreArea(name string, cache *pokecache.Cache) (LocationsArea, error) {
	exploreUrl := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", name)
	locsArea := LocationsArea{}
	if cacheData, ok := cache.Get(exploreUrl); !ok {
		response, err := http.Get(exploreUrl)
		if err != nil {
			return locsArea, err
		}
		defer response.Body.Close()
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return locsArea, err
		}
		errs := json.Unmarshal(responseData, &locsArea)
		if errs != nil {
			return locsArea, err
		}
		cache.Add(exploreUrl, responseData)
	} else {
		err := json.Unmarshal(cacheData, &locsArea)
		if err != nil {
			return locsArea, err
		}
	}
	return locsArea, nil
}

func GetPokemonInfo(name string, cache *pokecache.Cache) (Pokemon, error) {
	pokemonUrl := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)
	pokemon := Pokemon{}
	if cacheData, ok := cache.Get(pokemonUrl); !ok {
		response, err := http.Get(pokemonUrl)
		if err != nil {
			return pokemon, err
		}
		defer response.Body.Close()
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return pokemon, err
		}
		errs := json.Unmarshal(responseData, &pokemon)
		if errs != nil {
			return pokemon, err
		}
		cache.Add(pokemonUrl, responseData)
	} else {
		err := json.Unmarshal(cacheData, &pokemon)
		if err != nil {
			return pokemon, err
		}
	}
	return pokemon, nil
}
