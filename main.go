package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func cleanInput(input string) []string {
	loweredInput := strings.ToLower(input)
	output := strings.Fields(loweredInput)
	return output
}

func main() {
	// initialize commandline reader
	fmt.Print("Welcome to SCRY CLI.\n")
	scanner := bufio.NewScanner(os.Stdin)

	// initialize config
	cfg := &config{
		client:   &http.Client{},
		commands: getCommands(),
	}

	// CLI loop
	for {
		fmt.Print("SCRY > ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			fmt.Fprint(os.Stderr, "readig standard input:", err)
		}

		// clean input to list of lower strings
		input := cleanInput(scanner.Text())
		if len(input) == 0 {
			continue
		}
		commandName := input[0]
		var args []string
		if len(input) > 1 {
			args = input[1:]
		}
		// fmt.Printf("Command: %s\nArguments: %s\n", commandName, args)

		_, ok := cfg.commands[commandName]
		if ok {
			cfg.commands[commandName].Command(cfg, args)
		} else {
			fmt.Printf("command is not valid: %s\n", commandName)
			fmt.Println("call 'help' for information")
		}
	}
}
