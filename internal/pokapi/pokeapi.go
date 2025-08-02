package pokapi

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Singullaritty/pokedexcli/internal/pokecache"
)

func GetApiData(url string, cache *pokecache.Cache) ([]byte, error) {
	responseData, ok := cache.Get(url)
	if !ok {
		response, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error: received status code %d", response.StatusCode)
		}
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		cache.Add(url, responseData)
		return responseData, nil
	}
	return responseData, nil
}
