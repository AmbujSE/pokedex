package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

// Structs for API Response
type Pokemon struct {
	Name           string     `json:"name"`
	BaseExperience int        `json:"base_experience"`
	Height         int        `json:"height"`
	Weight         int        `json:"weight"`
	Stats          []Stat     `json:"stats"`
	Types          []PokeType `json:"types"`
}

type Stat struct {
	BaseStat int     `json:"base_stat"`
	Stat     StatObj `json:"stat"`
}

type StatObj struct {
	Name string `json:"name"`
}

type PokeType struct {
	Type TypeObj `json:"type"`
}

type TypeObj struct {
	Name string `json:"name"`
}

type config struct {
	nextURL     *string
	previousURL *string
	cache       map[string][]string // Cache for explored locations
	pokedex     map[string]Pokemon  // Caught Pokémon storage
}

var commands map[string]cliCommand

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandPokedex(cfg *config, args []string) error {
	if len(cfg.pokedex) == 0 {
		fmt.Println("Your Pokédex is empty!")
		return nil
	}

	fmt.Println("Your Pokédex:")
	for name := range cfg.pokedex {
		fmt.Println(" -", name)
	}
	return nil
}

func fetchPokemonData(pokemonName string) (*Pokemon, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", pokemonName)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Pokémon data: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result Pokemon
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &result, nil
}

func commandCatch(cfg *config, args []string) error {
	if len(args) < 1 {
		fmt.Println("Usage: catch <pokemon-name>")
		return nil
	}

	pokemonName := strings.ToLower(args[0])
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	// Check if Pokémon is already caught
	if _, found := cfg.pokedex[pokemonName]; found {
		fmt.Printf("%s is already in your Pokedex!\n", pokemonName)
		return nil
	}

	pokemon, err := fetchPokemonData(pokemonName)
	if err != nil {
		fmt.Println("Error fetching Pokémon:", err)
		return err
	}

	// Calculate catch probability
	rand.Seed(time.Now().UnixNano())
	catchChance := 1.0 - (float64(pokemon.BaseExperience) / 500.0)
	if catchChance < 0.1 {
		catchChance = 0.1 // Minimum catch chance of 10%
	}

	if rand.Float64() < catchChance {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		cfg.pokedex[pokemon.Name] = *pokemon // Add to Pokédex
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

func commandInspect(cfg *config, args []string) error {
	if len(args) < 1 {
		fmt.Println("Usage: inspect <pokemon-name>")
		return nil
	}

	pokemonName := strings.ToLower(args[0])

	// Check if the Pokémon is in the user's Pokédex
	pokemon, found := cfg.pokedex[pokemonName]
	if !found {
		fmt.Println("You have not caught that Pokémon")
		return nil
	}

	// Print Pokémon details
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, pokeType := range pokemon.Types {
		fmt.Printf("  - %s\n", pokeType.Type.Name)
	}

	return nil
}

func init() {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List caught Pokémon",
			callback:    commandPokedex,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a Pokémon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a caught Pokémon’s details",
			callback:    commandInspect,
		},
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(strings.TrimSpace(text)))
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{
		cache:   make(map[string][]string),
		pokedex: make(map[string]Pokemon),
	}

	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		words := cleanInput(input)
		if len(words) > 0 {
			command, exists := commands[words[0]]
			if exists {
				if err := command.callback(cfg, words[1:]); err != nil {
					fmt.Println("Error:", err)
				}
			} else {
				fmt.Println("Unknown command")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}
