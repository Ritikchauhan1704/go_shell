package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Command represents a shell command with its name, description, and handler
type Command struct {
	Name        string
	Description string
	Handler     func(args []string)
}

var commands = map[string]Command{}

func init() {
	commands = map[string]Command{
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

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		cmdName, args := splitCommandAndArguments(input)
		// fmt.Println(cmdName)
		// for _, arg := range args {
		// 	fmt.Println(arg)
		// }

		// handle redirections
		// In your main loop, replace the redirection block with this:
		cleanArgs, outFile, outR, outAppend, errFile, errR, errAppend := parseRedirection(args)
		if outR || errR {
			var origStdout, origStderr *os.File
			var fOut, fErr *os.File
			var err error

			if outR {
				if outAppend {
					fOut, err = os.OpenFile(outFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
				} else {
					fOut, err = os.Create(outFile)
				}
				if err != nil {
					fmt.Fprintf(os.Stderr, "redirect error: %v\n", err)
					continue
				}
				origStdout = os.Stdout
				os.Stdout = fOut
				defer fOut.Close()
			}

			if errR {
				if errAppend {
					fErr, err = os.OpenFile(errFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
				} else {
					fErr, err = os.Create(errFile)
				}
				if err != nil {
					if outR {
						os.Stdout = origStdout
					}
					fmt.Fprintf(os.Stderr, "redirect error: %v\n", err)
					continue
				}
				origStderr = os.Stderr
				os.Stderr = fErr
				defer fErr.Close()
			}

			runCommands(cmdName, cleanArgs, input)

			if outR {
				os.Stdout.Close()
				os.Stdout = origStdout
			}
			if errR {
				os.Stderr.Close()
				os.Stderr = origStderr
			}
			continue
		}

		runCommands(cmdName, cleanArgs, input)


	}
}

func runCommands(cmdName string, args []string, input string) {
	if cmd, ok := commands[cmdName]; ok {
		cmd.Handler(args)
	} else {
		externalCommand(cmdName, args, input)
	}
}
// parseRedirection now tracks both stdout and stderr redirects and appends.
func parseRedirection(args []string) (
    clean       []string,
    outFile     string, outRedirect, outAppend bool,
    errFile     string, errRedirect, errAppend bool,
) {
    clean = make([]string, 0, len(args))
    for i := 0; i < len(args); i++ {
        tok := args[i]
        switch tok {
        case ">", "1>":
            if i+1 < len(args) {
                outFile = args[i+1]
                outRedirect = true
                i++
                continue
            }
            fmt.Fprintln(os.Stderr, "redirect: no file specified")
            return args, "", false, false, "", false, false

        case ">>", "1>>":
            if i+1 < len(args) {
                outFile = args[i+1]
                outRedirect = true
                outAppend = true
                i++
                continue
            }
            fmt.Fprintln(os.Stderr, "redirect: no file specified")
            return args, "", false, false, "", false, false

        case "2>":
            if i+1 < len(args) {
                errFile = args[i+1]
                errRedirect = true
                i++
                continue
            }
            fmt.Fprintln(os.Stderr, "redirect: no file specified")
            return args, "", false, false, "", false, false

        case "2>>":
            if i+1 < len(args) {
                errFile = args[i+1]
                errRedirect = true
                errAppend = true
                i++
                continue
            }
            fmt.Fprintln(os.Stderr, "redirect: no file specified")
            return args, "", false, false, "", false, false
        }

        clean = append(clean, tok)
    }
    return clean, outFile, outRedirect, outAppend, errFile, errRedirect, errAppend
}


func splitCommandAndArguments(input string) (string, []string) {
	var args []string
	var current strings.Builder
	var i int

	// State variables
	inSingleQuote := false
	inDoubleQuote := false
	escapeNext := false
	
	for i < len(input) {
		ch := input[i]
		
		// Handle escape character
		if ch == '\\' && !escapeNext && !inSingleQuote {
			if inDoubleQuote {
				// In double quotes, only certain chars are escaped
				if i+1 < len(input) {
					nextCh := input[i+1]
					if nextCh == '"' || nextCh == '\\' || nextCh == '$' || nextCh == '`' || nextCh == '\n' {
						escapeNext = true
						i++
						continue
					}
					// Not a special character: output one backslash and skip
					current.WriteByte('\\')
					i++
					continue
				} else {
					// Trailing backslash: output it
					current.WriteByte('\\')
					i++
					continue
				}
			} else {
				// Outside quotes, escape always works
				escapeNext = true
				i++
				continue
			}
		}
		
		// Handle quotes
		if ch == '\'' && !escapeNext && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			i++
			continue
		}
		if ch == '"' && !escapeNext && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			i++
			continue
		}
		
		// Handle space as argument separator
		if ch == ' ' && !inSingleQuote && !inDoubleQuote && !escapeNext {
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
				i++
			continue
		}
		
		// Add character to current argument
		if escapeNext {
			current.WriteByte(ch)
			escapeNext = false
		} else {
			current.WriteByte(ch)
		}
			i++
	}
	
	// Add last argument if exists
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	
	// Warn on unclosed quotes
	if inSingleQuote || inDoubleQuote {
		fmt.Fprintln(os.Stderr, "Warning: Unclosed quotes in command")
	}
	
	if len(args) == 0 {
		return "", []string{}
	}
	return args[0], args[1:]
}
func exitCommand(args []string) {
	if len(args) == 1 && args[0] == "0" {
		os.Exit(0)
	}
}
func echoCommand(args []string) {
	fmt.Println(strings.Join(args, " "))
}
func pwdCommand(args []string) {
	wd, _ := os.Getwd()
	fmt.Println(wd)
}

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
func typeCommand(args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: type <command>")
		return
	}
	name := args[0]
	// Check builtins first
	if _, ok := commands[name]; ok {
		fmt.Printf("%s is a shell builtin\n", name)
		return
	}
	// Use exec.LookPath for external
	path, err := exec.LookPath(name)
	if err != nil {
		fmt.Printf("%s: not found\n", name)
	} else {
		fmt.Printf("%s is %s\n", name, path)
	}
}
// externalCommand runs an external program or reports not found
func externalCommand(cmd string, args []string, input string) {
	path, err := exec.LookPath(cmd)
	if err != nil {
		fmt.Printf("%s: not found\n", cmd)
		return
	}
	ec := exec.Command(path, args...)
	ec.Args[0] = filepath.Base(path)
	ec.Stdin = os.Stdin
	ec.Stdout = os.Stdout
	ec.Stderr = os.Stderr
	// if err := ec.Run(); err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error running %s: %v\n", cmd, err)
	// }
	ec.Run()
}
