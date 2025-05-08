package shell

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Command represents a shell command with its name, description, and handler
type Command struct {
	Name        string
	Description string
	Handler     func(args []string)
}

// initCommands initializes the built-in shell commands
func initCommands() map[string]Command {
	return map[string]Command{
		"exit": {
			Name:        "exit",
			Description: "Exit the shell",
			Handler:     exitCommand,
		},
		"echo": {
			Name:        "echo",
			Description: "Display a line of text",
			Handler:     echoCommand,
		},
		"type": {
			Name:        "type",
			Description: "Describe a command",
			Handler:     typeCommand,
		},
		"pwd": {
			Name:        "pwd",
			Description: "Print the working directory",
			Handler:     pwdCommand,
		},
		"cd": {
			Name:        "cd",
			Description: "Change the directory",
			Handler:     cdCommand,
		},
	}
}

// exitCommand implements the exit command
func exitCommand(args []string) {
	if len(args) == 1 && args[0] == "0" {
		os.Exit(0)
	}
	os.Exit(0)
}

// echoCommand implements the echo command
func echoCommand(args []string) {
	fmt.Println(strings.Join(args, " "))
}

// pwdCommand implements the pwd command
func pwdCommand(args []string) {
	wd, _ := os.Getwd()
	fmt.Println(wd)
}

// cdCommand implements the cd command
func cdCommand(args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: cd <directory>")
		return
	}
	dir := args[0]
	if dir == "~" {
		dir = os.Getenv("HOME")
	}
	err := os.Chdir(dir)
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", dir)
	}
}

// typeCommand implements the type command
func typeCommand(args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: type <command>")
		return
	}
	
	name := args[0]
	
	// Check if it's a built-in command
	commands := initCommands()
	if _, ok := commands[name]; ok {
		fmt.Printf("%s is a shell builtin\n", name)
		return
	}
	
	// Use exec.LookPath for external commands
	path, err := exec.LookPath(name)
	if err != nil {
		fmt.Printf("%s: not found\n", name)
	} else {
		fmt.Printf("%s is %s\n", name, path)
	}
}