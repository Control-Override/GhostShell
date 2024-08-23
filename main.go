package main

import (
	"bufio"
	"fmt"
	"os"
	"plugin"
	"strings"
	"github.com/Control-Override/GhostCommandInterface"
)


var commands = map[string] command_interface.Command{}

func loadPlugin(path string) error {
	p, err := plugin.Open(path)
	if err != nil {
		return err
	}

	newCmdSymbol, err := p.Lookup("New")
	if err != nil {
		return err
	}


	newCmd, ok := newCmdSymbol.(func() command_interface.Command)
	if !ok {
		return fmt.Errorf("invalid plugin signature")
	}

	command := newCmd()

	commandName := command.Alias()
	commands[commandName] = command
	return nil
}

func main() {
	// Some config stuff we will move to a better place later
	projectName := "Ghost Shell"
	shellPrompt := "GhostShell> "
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
	fmt.Println(projectName+" with Plugins - Type 'exit' to quit.")
	for {
		fmt.Print(shellPrompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)

		
		
	
		args := strings.Fields(input)
		if len(args) == 0 {
			continue
		}

		// Embedded Commmands
		if args[0] == "exit" {
			fmt.Println("Exiting "+projectName+".....")
			break
		} else if args[0] == "help" {
			if len(args) == 1 {
				// Help with no args should call help but pass in the "short" parameter
				fmt.Println("Commands:")
				// for through all commands and display short help
				for _, cmd := range commands {
					fmt.Println(cmd.Alias()," - ",cmd.Help(true,args[1:]))
					//fmt.Println("Detailed Help:", cmd.Help(false, nil))
				}
				//fmt.Println(command.Help(true))
			} else {
				// Help with an arg should just call that modules help with no parameter
				// Process Plugin Commands
				commandName := args[1]
				command, exists := commands[commandName]
				if exists {
					if len(args)>=2 {
						fmt.Println(command.Help(false, args[2:]))
					} else {
						
					}
				} else {
					fmt.Println("Unknown plugin:", commandName)
				}
			}
		} else {
			// Process Plugin Commands
			commandName := args[0]
			command, exists := commands[commandName]
			if exists {
				command.Execute(args[1:])
			} else {
				fmt.Println("Unknown command:", commandName)
			}
		}

		
	}
}
