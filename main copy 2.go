package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanIsnput(text string) []string {
	words := strings.Fields(strings.ToLower(strings.TrimSpace(text)))
	return words
}

func masin() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		words := cleanInput(input)
		if len(words) > 0 {
			fmt.Printf("Your command was: %s\n", words[0])
		}
	}
}
