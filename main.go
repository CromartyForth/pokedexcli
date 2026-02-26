package main

import (
	"fmt"
	"bufio"
	"os"
	"github.com/CromartyForth/pokedexcli/internal/pokecache"
)

var locationsEP = "https://pokeapi.co/api/v2/location-area/"
var cache_duration = 7
var poke_cache = pokecache.NewCache(cache_duration)

func main(){

	scanner := bufio.NewScanner(os.Stdin)

	pConfig := Config{}
	for ;; {
		fmt.Printf("Pokedex > ")
		if scanner.Scan() == false {
			break
		}
		query := scanner.Text()
		value, ok := getCommands()[query]
		if ok == false {
			fmt.Println("Unknown command")
		} else {
			value.Callback(&pConfig)
		}
	}
}





