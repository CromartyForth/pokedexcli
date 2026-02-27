package main

import(
	"strings"
)

func cleanInput(text string) []string {
	
	lowered := strings.ToLower(text)
	cleaned := strings.Fields(lowered)
	return cleaned
}

