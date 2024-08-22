package main

import (
	"bufio"
	"fmt"
	"os"
	"plugin"
	"strings"
)

type Command interface {
	Execute(args []string)
}

var commands = map[string]Command{}

func loadPlugin(path string) error {
	p, err := plugin.Open(path)
	if err != nil {
		return err
	}

	newCmdSymbol, err := p.Lookup("New")
	if err != nil {
		return err
	}

	fmt.Printf("Type: %T\n", newCmdSymbol)

	newCmd, ok := newCmdSymbol.(func() Command)
	if !ok {
		return fmt.Errorf("invalid plugin signature")
	}

	command := newCmd()
	commandName := strings.TrimSuffix(path, ".so")
	commands[commandName] = command
	return nil
}

func main() {
	// Load plugins from a folder
	pluginFolder := "./plugins"
	pluginFiles, _ := os.ReadDir(pluginFolder)

	for _, file := range pluginFiles {
		if strings.HasSuffix(file.Name(), ".so") {
			pluginPath := fmt.Sprintf("%s/%s", pluginFolder, file.Name())
			err := loadPlugin(pluginPath)
			if err != nil {
				fmt.Printf("Failed to load plugin %s: %s\n", file.Name(), err)
			}
		}
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Go Shell with Plugins - Type 'exit' to quit.")

	for {
		fmt.Print("go-shell> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "exit" {
			fmt.Println("Exiting Go Shell...")
			break
		}

		args := strings.Fields(input)
		if len(args) == 0 {
			continue
		}

		commandName := args[0]
		command, exists := commands[commandName]
		if exists {
			command.Execute(args[1:])
		} else {
			fmt.Println("Unknown command:", commandName)
		}
	}
}
