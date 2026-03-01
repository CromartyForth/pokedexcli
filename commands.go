package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"io"
	"math/rand"
	"math"
)

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
		"explore": {
			Name: "explore",
			Description: "Displays all the pokemon at a given location\nUsage > explore {location}",
			Callback: commandExplore,
		},
		"catch": {
			Name: "catch",
			Description: "Tries to catch a pokemon and add it to your pokedex\nUsage > catch {pokemon}",
			Callback: commandCatch,
		},
		"inspect":{
			Name: "Inspect",
			Description: "Displays the stats of any previously caught pokemon\nUsage > inspect {pokemon}",
			Callback: commandInspect,
		},
	}
	return mapCommands
}

func commandExit(config *Config, parameter string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config, parameter string) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n")
	for _, command := range(getCommands()) {
		// Welcome to the Pokedex!
		fmt.Printf("%v: %v\n",command.Name, command.Description)
	}
	return nil
}

func commandInspect(config *Config, parameter string) error {
	// check for pokemon user input
	if parameter == ""{
		fmt.Println(getCommands()["inspect"].Description)
		return nil
	}

	// is already caught
	data, ok := pokedex[parameter]
	if ok == false {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %v\nHeight: %v\nWeight: %v\nStats:\n", data.Name, data.Height, data.Weight)
	for _, value := range data.Stats {
		fmt.Printf("  -%v: %v\n", value.Stat.Name, value.BaseStat)
	}
	fmt.Println("Types:")
	for _, value := range data.Types {
		fmt.Printf("  - %v\n", value.Type.Name)
	}

	return nil
}

func commandCatch(config *Config, parameter string) error{
	// check for pokemon user input
	if parameter == ""{
		fmt.Println(getCommands()["catch"].Description)
		return nil
	}
	fmt.Printf("Throwing a Pokeball at %v...\n", parameter)
	url := pokemonEP + parameter

	// is already caught
	_, ok := pokedex[parameter]
	if ok == true {
		fmt.Printf("%v has already been caught!\n", parameter)
		return nil
	}

	// check cache and make api call
	data, ok := poke_cache.Get(url)
	if ok == false {
		//fmt.Println("making api call")
		// make api call
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Network Error: %v\n", err)
		}
		defer res.Body.Close()

		if res.StatusCode > 299 {
			fmt.Printf("'%v', not sure that was a real pokemon\n", parameter)
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
		//fmt.Printf("reading from cache\n")
	}

	// convert Json to go struct
	pokemon := Pokemon{}
	err := json.Unmarshal(data, &pokemon)
	if err != nil {
		return fmt.Errorf("JSON Error: %v\n", err)
	}

	// catch Pokemon?
	// pokemon base XP minus min pokemon xp
	base := float64(pokemon.BaseExperience) - minXP
	fmt.Printf("%v experience: %v\n", parameter, base)
	// scale by ratio of spans
	ratio := (maxChance - minChance) / (maxXP - minXP)
	scaledBase := base * ratio
	// add new minimum
	percentChanceToCatch := math.Round((scaledBase + minChance))
	fmt.Printf("must roll greater than %.0f to catch %v\n", percentChanceToCatch, parameter)


	// get random number in range 1 to 100
	result := rand.Intn(100)
	fmt.Printf("You rolled: %v! - ", result)

	if result < int(percentChanceToCatch) {
		fmt.Printf("%v escaped!\n", parameter)
		return nil
	}

	// add pokemon to pokedex
	pokedex[parameter] = pokemon
	fmt.Printf("%v was caught!\n", parameter)

	return nil
}

func commandExplore(config *Config, parameter string) error{
	
	if parameter == ""{
		fmt.Println(getCommands()["explore"].Description)
		return nil
	}
	fmt.Printf("Exploring %v...", parameter)
	url := locationsEP + parameter
	
	// is data already in cache
	data, ok := poke_cache.Get(url)
	if ok == false {
		fmt.Println("..")
		// make api call
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Network Error: %v\n", err)
		}
		defer res.Body.Close()

		if res.StatusCode > 299 {
			fmt.Printf("\nLocation not found, check your entry\n")
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
		fmt.Printf(".\n")
	}

	// convert Json to go struct
	pokeLocations := NamedLocation{}
	err := json.Unmarshal(data, &pokeLocations)
	if err != nil {
		return fmt.Errorf("JSON Error: %v\n", err)
	}

	// print out Pokemon
	fmt.Println("Found Pokemon:")
	for _, location := range(pokeLocations.PokemonEncounters){
		fmt.Printf(" - %v\n", location.Pokemon.Name)
	}

	return nil
}

func commandMapb(ptrConfig *Config, parameter string) error {
	// is on first page
	if ptrConfig.Previous == "" {
		fmt.Println("You're on the first page")
		return nil
	} 
	url := ptrConfig.Previous
	
	// check cache, make api call and print result
	err := map_helper(ptrConfig, url)
	if err != nil {
		return err
	}
	return nil
}

func commandMap(ptrConfig *Config, parameter string) error {
	
	// is repeated call?
	url := locationsEP + locationsQuery
	if ptrConfig.Next != "" {
		url = ptrConfig.Next
	}

	// check cache, make api call and print result
	err := map_helper(ptrConfig, url)
	if err != nil {
		return err
	}
	return nil
}

func map_helper(ptrConfig *Config, url string) error {
	
	// is data already in cache
	data, ok := poke_cache.Get(url)
	if ok == false {
		fmt.Println("..")
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
		fmt.Printf(".\n")
	}

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
