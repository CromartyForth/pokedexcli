package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"io"
)

type Config struct {
	Next string
	Previous string
}

type Location struct {
	Count    int `json:"count"`
	Next string    `json:"next"`
	Previous string `json:"previous"`
	Results []struct {
		Name string `json:"name"`
		Url string `json:"url"`
	}
}

type CliCommand struct {
	Name        string
	Description string
	Callback    func(*Config) error
}

func getCommands() map[string]CliCommand {
	var mapCommands = map[string]CliCommand{
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"help": {
			Name: "help",
			Description: "Displays a help message",
			Callback: commandHelp,
		},
		"map": {
			Name: "map",
			Description: "Displays the next 20 locations",
			Callback: commandMap,
		},
		"mapb": {
			Name: "mapb",
			Description: "Desplays the previous 20 locations",
			Callback: commandMapb,
		},
	}
	return mapCommands
}

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n")
	for _, command := range(getCommands()) {
		// Welcome to the Pokedex!
		fmt.Printf("%v: %v\n",command.Name, command.Description)
	}
	return nil
}
func commandMapb(ptrConfig *Config) error {
	// is on first page
	if ptrConfig.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	} 
	url := ptrConfig.Previous

	// make api call
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Network Error: %v\n", err)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return fmt.Errorf("HTTP Status Code: %v\n", res.StatusCode)
	}

	// add res to cache


	// read response body
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Reading Error: %v\n", err)
	}

	// convert Json to go struct
	pokeLocations := Location{}
	err = json.Unmarshal(data, &pokeLocations)
	if err != nil {
		return fmt.Errorf("JSON Error: %v\n", err)
	}

	// update config with previous and next urls
	ptrConfig.Next = pokeLocations.Next
	ptrConfig.Previous = pokeLocations.Previous

	// print out locations
	for _, location := range(pokeLocations.Results){
		fmt.Println(location.Name)
	}

	return nil
}

func commandMap(ptrConfig *Config) error {
	
	// is repeated call?
	url := locationsEP
	if ptrConfig.Next != "" {
		url = ptrConfig.Next
	}

	// is data already in cache
	data, ok := poke_cache.Get(url)
	if ok == false {
		fmt.Println("Making API call")
		// make api call
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Network Error: %v\n", err)
		}
		defer res.Body.Close()

		if res.StatusCode > 299 {
			return fmt.Errorf("HTTP Status Code: %v\n", res.StatusCode)
		}

		// read response body returns []byte
		data, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Reading Error: %v\n", err)
		}
		
		// put data into cache
		poke_cache.Add(url, data)
	} else {
		fmt.Printf("Read from Cache...\n")
	}

	fmt.Println(string(data))
	// convert Json to go struct
	pokeLocations := Location{}
	err := json.Unmarshal(data, &pokeLocations)
	if err != nil {
		return fmt.Errorf("JSON Error: %v\n", err)
	}

	// update config with previous and next urls
	ptrConfig.Next = pokeLocations.Next
	ptrConfig.Previous = pokeLocations.Previous

	// print out locations
	for _, location := range(pokeLocations.Results){
		fmt.Println(location.Name)
	}

	return nil
}