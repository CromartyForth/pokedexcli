package main

import (
	"fmt"
	"bufio"
	"os"
)

var locationsEP = "https://pokeapi.co/api/v2/location/"

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





