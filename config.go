package main

import (
	"fmt"
	"net/http"
	"os"
	"slices"
)

type config struct {
	client   *http.Client
	commands map[string]cliCommand
}

type cliCommand struct {
	Name        string
	Description string
	Command     func(*config, []string) error
}

// get a list of all commands for application usage
func getCommands() map[string]cliCommand {
	commands := map[string]cliCommand{
		"help": {
			Name:        "help",
			Description: "Displays a help message, explaining the usage of ScryCLI. No arguments needed.",
			Command:     handlerHelp,
		},
		"search": {
			Name:        "search",
			Description: "Search for one or more cards, just attach the search term (e.g. 'search Sheoldred the').",
			Command:     handlerSearch,
		},
		"random": {
			Name:        "random",
			Description: "Get a random card.",
			Command:     handlerRandom,
		},
		"exit": {
			Name:        "exit",
			Description: "Exits the ScryCLI. No arguments needed.",
			Command:     handlerExit,
		},
	}
	return commands
}

// set request Header information for API requirements
func setRequest(req *http.Request) {
	req.Header.Set("User-Agent", "ScryCLI/0.1")
	req.Header.Set("Accept", "*/*")
}

// exit the application
func handlerExit(cfg *config, args []string) error {
	os.Exit(0)
	return nil
}

// display information to user how to use application
func handlerHelp(cfg *config, args []string) error {
	fmt.Print("Usage of the ScryCLI\nJust type one of the following commands and add the arguments if necessary:\n\n")
	for name, obj := range cfg.commands {
		fmt.Printf("%8s: %s\n", name, obj.Description)
	}
	fmt.Println()
	return nil
}

// call search function with given arguments to find card information
func handlerSearch(cfg *config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("no args given")
	}
	if slices.Contains(args, "--help") {
		searchHelp()
		return nil
	}
	results, err := search(cfg.client, args)
	if err != nil {
		return err
	}
	if len(results) > 10 {
		fmt.Printf("Your search results in %d cards.", len(results))
		userinput := cliContinue()
		if !userinput {
			return nil
		}
	}
	for i, result := range results {
		fmt.Printf("%3d. %-30s %-30s\n", i+1, result.Name, result.ManaCost)
	}
	return nil
}

func handlerRandom(cfg *config, args []string) error {
	if slices.Contains(args, "--help") {
		randomHelp()
		return nil
	}
	result, err := randomCard(cfg.client, args)
	if err != nil {
		return err
	}
	fmt.Printf("%-30s %-30s\n", result.Name, result.ManaCost)

	return nil
}
