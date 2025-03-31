package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			text := scanner.Text()
			newText := cleanInput(text)
			if len(newText) == 0 {
				fmt.Println("Please input a command.")
			} else {
				if _, ok := commandMap[newText[0]]; !ok {
					fmt.Println("Please input a valid command.")
					continue
				}
				err := commandMap[newText[0]].callback()
				if err != nil {
					fmt.Println(fmt.Errorf("error: %w", err))
				}
			}
			// fmt.Printf("Your command was: %s\n", newText[0]) \\ for printing first word
		}
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commandMap map[string]cliCommand

func init() {
	commandMap = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Get more info",
			callback:    commandHelp,
		},
	}
}

func cleanInput(text string) []string {
	strippedString := strings.TrimSpace(strings.ToLower(text))
	substrings := strings.Fields(strippedString)
	return substrings
}

func commandExit() error { // does this need to return an error? When would exiting not work?
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp() error { // does this need to return an error? When would exiting not work?
	usageStr := ""
	for key, val := range commandMap {
		currStr := fmt.Sprintf("%s: %s\n", key, val.description)
		usageStr += currStr
	}
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n%s", usageStr)
	return nil
}
