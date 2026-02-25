package commands

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/ScholarlyKiwi/pokedex/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*CmdConfig) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"map": {
			name:        "map",
			description: "Displays the next 20 areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 areas",
			callback:    commandMapb,
		},
		"help": {
			name:        "help",
			description: "Displays a Help Message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"explore": {
			name:        "explore",
			description: "Explore an area - requires area name",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a Pokemon by name",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspects a previously caught Pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Returns a list of caught pokemon",
			callback:    commandPokedex,
		},
	}
}

func CommandScan() error {

	var config CmdConfig
	scan := bufio.NewScanner(os.Stdin)
	commandReg := getCommands()

	config.Map_next = "https://pokeapi.co/api/v2/location-area/"
	config.CaughtPokemon = make(map[string]any)

	newCache, err := pokecache.NewCache(time.Second * 5)
	if err != nil {
		fmt.Printf("Error creating cache: %v", err)
	}
	config.Cache = newCache

	for {
		fmt.Print("Pokedex > ")
		scan.Scan()

		config.Input = cleanInput(scan.Text())

		if command, ok := commandReg[config.Input[0]]; ok == true {
			call := command.callback
			if err := call(&config); err != nil {
				fmt.Printf("Error executing command %v: %v", config.Input[0], err)
			}
		} else {
			fmt.Println("Invalid Command")
			commandHelp(&config)
		}
	}
}

func commandExit(config *CmdConfig) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *CmdConfig) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")

	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	return nil
}

func getLocationArea(config *CmdConfig, url string) error {

	json_locationAreas, err := getPokeLocationAreas(url, config)
	if err != nil {
		fmt.Printf("Unable to retrieve Location Areas: %v\n", err)
		return nil
	}

	config.Map_next = json_locationAreas.Next
	config.Map_prev = json_locationAreas.Previous
	for _, value := range json_locationAreas.Results {
		fmt.Println(value.Name)
	}
	return nil

}

func commandMap(config *CmdConfig) error {
	if config.Map_next == "" {
		fmt.Printf("you're on the last page\n")
		return nil
	}

	return getLocationArea(config, config.Map_next)
}

func commandMapb(config *CmdConfig) error {
	if config.Map_prev == "" {
		fmt.Printf("you're on the first page\n")
		return nil
	}

	return getLocationArea(config, config.Map_prev)
}

func commandExplore(config *CmdConfig) error {
	if len(config.Input) != 2 {
		fmt.Println("Usage - explore {area name}")
		return nil
	}

	locationArea, err := getpokeLocationAreaByName(config.Input[1], config)
	if err != nil {
		fmt.Printf("Error retrieving Location Area by Name: %v\n", err)
	}

	for _, encounter := range locationArea.PokemonEncounters {
		fmt.Printf(" - %v\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(config *CmdConfig) error {
	if len(config.Input) <= 1 {
		fmt.Println("The catch command must be followed by a Pokemon name.")
		return nil
	}
	pokemonName := config.Input[1]

	if _, exists := config.CaughtPokemon[pokemonName]; exists {
		fmt.Printf("You already have %v.\n", pokemonName)
	} else {
		pokemonDetails, err := getpokePokemon(pokemonName, config)
		if err != nil {
			fmt.Printf("Error catching pokemon %v: API returned %v\n", pokemonName, err)
			return nil
		}
		if pokemonDetails.Name != pokemonName {
			fmt.Printf("Pokemon %v not found.\n", pokemonName)
			return nil
		}

		fmt.Printf("Throwing a Pokeball at %v...\n", pokemonName)
		catch_attempt := rand.Intn(800)
		if catch_attempt >= pokemonDetails.BaseExperience {
			fmt.Printf("%v was caught!\n", pokemonName)
			config.CaughtPokemon[pokemonName] = pokemonDetails
		} else {
			fmt.Printf("%v escaped!\n", pokemonName)
		}
	}
	return nil
}

func commandInspect(config *CmdConfig) error {
	if len(config.Input) <= 1 {
		fmt.Println("The inspect command must be followed by a Pokemon name.")
		return nil
	}
	pokemonName := config.Input[1]
	data, exists := config.CaughtPokemon[pokemonName]

	if !exists {
		fmt.Println("you have not caught that pokemon")
	} else {
		pokemon := data.(pokePokemon)
		fmt.Printf("Name: %v\n", pokemon.Name)
		fmt.Printf("Height: %v\n", pokemon.Height)
		fmt.Printf("Weight: %v\n", pokemon.Weight)
		fmt.Printf("Stats:\n")
		for _, value := range pokemon.Stats {
			fmt.Printf(" -%v: %v\n", value.Stat.Name, value.BaseStat)
		}
		fmt.Printf("Types:\n")
		for _, value := range pokemon.Types {
			fmt.Printf(" - %v\n", value.Type.Name)
		}
	}

	return nil
}

func commandPokedex(config *CmdConfig) error {

	fmt.Println("Your Pokedex:")

	for _, pokemon := range config.CaughtPokemon {
		fmt.Printf(" - %v\n", pokemon.(pokePokemon).Name)
	}

	return nil
}
