package main

import (
	"fmt"
	"bufio"
	"os"
	"time"
	"github.com/CromartyForth/pokedexcli/internal/pokecache"
)

const locationsEP = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
const cache_duration = 20 * time.Second
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





