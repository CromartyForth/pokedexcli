package main

import (
	"fmt"
	"bufio"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for ;; {
		fmt.Printf("Pokedex > ")
		if scanner.Scan() == false {
			break
		}
		query := scanner.Text()
		if query == "" {
			break
		}
		splitWords := cleanInput(query)
		fmt.Printf("Your command was: %v\n", splitWords[0])

	}
}




