package main

import (
	"fmt"
	"net/http"
	"os"
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
			Description: "Displays a help message.",
			Command:     handlerHelp,
		},
		"search": {
			Name:        "search",
			Description: "search one or more cards",
			Command:     handlerSearch,
		},
		"exit": {
			Name:        "exit",
			Description: "exit ScryCLI",
			Command:     handlerExit,
		},
	}
	return commands
}

// set request Header information for API requirements
func setRequest(req *http.Request) {
	// set request header to API requirements
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
	// #TODO
	fmt.Println("here comes help")
	return nil
}

// call search function with given arguments to find card information
func handlerSearch(cfg *config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("no args given")
	}
	results, err := search(cfg.client, args)
	if err != nil {
		return err
	}
	for i, result := range results {
		fmt.Printf("%2d. %s\n", i+1, result)
	}
	return nil
}
