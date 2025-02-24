package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(input string) []string {
	loweredInput := strings.ToLower(input)
	output := strings.Fields(loweredInput)
	return output
}

func main() {
	fmt.Print("Welcome to SCRY CLI.\n\n")

	// initialize commandline reader
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("SCRY > ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			fmt.Fprint(os.Stderr, "readig standard input:", err)
		}
		input := cleanInput(scanner.Text())
		if len(input) == 0 {
			continue
		}
		commandName := input[0]
		var args []string
		if len(input) > 1 {
			args = input[1:]
		}
		fmt.Printf("Command: %s\nArguments: %s\n", commandName, args)
		if commandName == "exit" {
			break
		}
		if commandName == "search" && len(args) > 0 {
			results, err := search(args)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			for i, result := range results {
				fmt.Printf("%2d. %s\n", i, result)
			}
		}
	}
}
