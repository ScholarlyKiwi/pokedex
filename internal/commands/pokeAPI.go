package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func callApi(url string, config *CmdConfig) ([]byte, error) {

	var bodyData []byte

	bodyData, exists := config.Cache.Get(url)

	if !exists {
		res, err := http.Get(url)
		if err != nil {
			return bodyData, err
		}

		body, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		if res.StatusCode > 299 {
			return bodyData, fmt.Errorf("PokeAPI returned status code: %d and \nbody: %s\n", res.StatusCode, body)
		}

		if err != nil {
			return bodyData, err
		}
		err = config.Cache.Add(url, body)
		if err != nil {
			return bodyData, err
		}

		bodyData = body
	}

	return bodyData, nil

}

func getPokeLocationAreas(url string, config *CmdConfig) (pokeResourceList, error) {

	var data pokeResourceList
	var bodyData []byte
	var err error

	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	}

	bodyData, err = callApi(url, config)

	err = json.Unmarshal(bodyData, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func getpokeLocationAreaByName(area string, config *CmdConfig) (pokeLocationArea, error) {
	var locationArea pokeLocationArea
	if len(area) == 0 {
		return locationArea, fmt.Errorf("Error: missing area name %v", area)
	}

	url := "https://pokeapi.co/api/v2/location-area/" + area + "/"

	bodyData, err := callApi(url, config)
	if err != nil {
		return locationArea, nil
	}

	err = json.Unmarshal(bodyData, &locationArea)
	if err != nil {
		return locationArea, err
	}
	return locationArea, nil

}

func getpokePokemon(name string, config *CmdConfig) (pokePokemon, error) {
	var pokemon pokePokemon
	if len(name) <= 0 {
		return pokemon, fmt.Errorf("Error: missing pokemon name %v", name)
	}

	url := "https://pokeapi.co/api/v2/pokemon/" + name + "/"

	bodyData, err := callApi(url, config)
	if err != nil {
		return pokemon, err
	}

	if err := json.Unmarshal(bodyData, &pokemon); err != nil {
		return pokemon, err
	}

	return pokemon, nil
}
