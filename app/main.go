package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)


// List of supported built-in commands
var builtins = map[string]bool{
	"echo": true,
	"exit": true,
	"type": true,
	"pwd": true,
	"cd": true,
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Get PATH environment variable and split it into directories
	pathEnv := os.Getenv("PATH")
	pathDirs := strings.Split(pathEnv, ":")

	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Read user input
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		// Trim whitespace and check for empty input
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// Split input into command and arguments
		fields := strings.Fields(input)
		cmd, args := fields[0], fields[1:]

		// Handle built-in commands using switch-case
		switch cmd {
		case "exit":
			// Handle "exit 0" to terminate shell
			if len(args) == 1 && args[0] == "0" {
				os.Exit(0)
			}
		case "echo":
			// Print the arguments as a string
			fmt.Println(strings.Join(args, " "))
		case "type":
			// Handle the "type" command logic
			handleType(args, pathDirs)
		case "pwd":
			// Print the current working directory
			workingDir, _ := os.Getwd()
			fmt.Println(workingDir)
		case "cd":
			// Change current directory
			// todo
		default:
			// Try to run the command as an external program
			runExternal(cmd, args, pathDirs, input)
		}
	}
}

// handleType checks if a command is built-in or an executable in PATH
func handleType(args []string, dirs []string) {
	if len(args) != 1 {
		fmt.Println("Usage: type <command>")
		return
	}

	arg := args[0]
	if builtins[arg] {
		// Command is a shell builtin
		fmt.Printf("%s is a shell builtin\n", arg)
	} else if path, found := findExecutable(arg, dirs); found {
		// Command found in PATH
		fmt.Printf("%s is %s\n", arg, path)
	} else {
		// Command not found
		fmt.Printf("%s: not found\n", arg)
	}
}

// runExternal finds and executes an external program
func runExternal(cmd string, args []string, dirs []string, fullInput string) {
	if path, found := findExecutable(cmd, dirs); found {
		// Prepare the command with correct path
		externalCmd := exec.Command(cmd, args...)
		externalCmd.Path = path
		externalCmd.Stdin = os.Stdin
		externalCmd.Stdout = os.Stdout
		externalCmd.Stderr = os.Stderr

		// Run the external command
		if err := externalCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running command: %v\n", err)
		}
	} else {
		// If command not found, show error
		fmt.Printf("%s: command not found\n", fullInput)
	}
}

func isExecutable(filePath string) bool {
	//check if the file exists and retrieves metadata
	fileInfo, err := os.Stat(filePath)

	if err != nil {
		return false // File doesn't exist or error accessing it
	}
	if fileInfo.IsDir() {
		return false // It's a directory, not a file
	}
	if fileInfo.Mode()&0111 == 0 {
		return false // File is not executable by anyone
	}
	return true // File exists, is not a directory, and is executable
}
func findExecutable(command string, dirs [] string) (string, bool){
	for _, dir := range dirs {
		fullPath := filepath.Join(dir, command)
		if isExecutable(fullPath){
			return fullPath, true
		}
	}
	return "", false
}

