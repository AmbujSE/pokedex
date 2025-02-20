package main

import (
	"fmt"
	"strings"
)

func mawin() {
	fmt.Println("Hello, World!")
}

func cleanInwput(text string) []string {
	words := strings.Fields(strings.ToLower(strings.TrimSpace(text)))
	return words
}
