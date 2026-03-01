package main

import (
	"fmt"
	"bufio"
	"os"
	"time"
	"github.com/CromartyForth/pokedexcli/internal/pokecache"
)

const locationsEP = "https://pokeapi.co/api/v2/location-area/"
const locationsQuery = "?offset=0&limit=20"
const pokemonEP = "https://pokeapi.co/api/v2/pokemon/"
const cache_duration = 100 * time.Second
var poke_cache = pokecache.NewCache(cache_duration)
const minXP = 20.0// Magikarp
const maxXP = 608.0 // Blissey
const minChance = 10.0
const maxChance = 90.0



// THE POKEDEX!!
var pokedex = make(map[string]Pokemon)

func main(){

	scanner := bufio.NewScanner(os.Stdin)

	pConfig := Config{}
	for ;; {
		fmt.Printf("Pokedex > ")
		if scanner.Scan() == false {
			break
		}
		text := scanner.Text()

		// clean input into slice
		query := cleanInput(text)

		// check first item in slice is a valid command
		value, ok := getCommands()[query[0]]
		if ok == false {
			fmt.Println("Unknown command")
		} else {
			if len(query) > 1 {
				value.Callback(&pConfig, query[1])
			} else {
				value.Callback(&pConfig, "")
			}
		}
	}
}





