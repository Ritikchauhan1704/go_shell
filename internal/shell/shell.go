package shell

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/peterh/liner"
)

// Shell represents the shell application with its state
type Shell struct {
	line     *liner.State
	commands map[string]Command
	lastPrefix string
	tabCount   int
}

// New creates a new shell instance
func New(line *liner.State) *Shell {
	sh := &Shell{
		line:     line,
		commands: initCommands(),
	}

	// Set up auto-completion
	sh.setupCompleter()

	return sh
}

// Run starts the shell's read-eval-print loop
func (sh *Shell) Run() error {
	for {
		input, err := sh.line.Prompt("$ ")
		if err == liner.ErrPromptAborted {
			continue // Ctrl+C pressed
		} else if err == io.EOF {
			return err // Ctrl+D or EOF
		} else if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading line:", err)
			continue
		}

		trimmed := strings.TrimRight(input, "\n")
		if strings.TrimSpace(trimmed) == "" {
			continue
		}

		sh.line.AppendHistory(trimmed)
		
		// Parse and execute the command
		sh.executeInput(trimmed)
	}
}

// executeInput processes a command line input
func (sh *Shell) executeInput(input string) {
	cmdName, args := parseCommandLine(input)
	if cmdName == "" {
		return
	}
	
	cleanArgs, outFile, outR, outAppend, errFile, errR, errAppend := parseRedirection(args)
	
	if outR || errR {
		var origStdout, origStderr *os.File
		var fOut, fErr *os.File
		var err error

		// Handle stdout redirection
		if outR {
			if outAppend {
				fOut, err = os.OpenFile(outFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			} else {
				fOut, err = os.Create(outFile)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "redirect error: %v\n", err)
				return
			}
			origStdout = os.Stdout
			os.Stdout = fOut
			defer fOut.Close()
		}

		// Handle stderr redirection
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
				return
			}
			origStderr = os.Stderr
			os.Stderr = fErr
			defer fErr.Close()
		}

		sh.executeCommand(cmdName, cleanArgs, input)

		// Restore standard outputs
		if outR {
			os.Stdout.Close()
			os.Stdout = origStdout
		}
		if errR {
			os.Stderr.Close()
			os.Stderr = origStderr
		}
		return
	}
	
	sh.executeCommand(cmdName, cleanArgs, input)
}

// executeCommand runs the given command with arguments
func (sh *Shell) executeCommand(cmdName string, args []string, input string) {
	if cmd, ok := sh.commands[cmdName]; ok {
		cmd.Handler(args)
	} else {
		runExternalCommand(cmdName, args)
	}
}

// setupCompleter configures command auto-completion
func (sh *Shell) setupCompleter() {
	sh.line.SetCompleter(func(input string) (c []string) {
		matches := make(map[string]bool)
		results := []string{}

		// Built-in command matches
		for cmdName := range sh.commands {
			if strings.HasPrefix(cmdName, input) {
				matches[cmdName] = true
				results = append(results, cmdName)
			}
		}

		// External command matches
		for _, dir := range strings.Split(os.Getenv("PATH"), ":") {
			entries, err := os.ReadDir(dir)
			if err != nil {
				continue
			}
			for _, entry := range entries {
				name := entry.Name()
				if matches[name] || !strings.HasPrefix(name, input) || !entry.Type().IsRegular() {
					continue
				}
				fullPath := filepath.Join(dir, name)
				if info, err := os.Stat(fullPath); err == nil && info.Mode().Perm()&0111 != 0 {
					matches[name] = true
					results = append(results, name)
				}
			}
		}

		if len(results) == 0 {
			fmt.Print("\a")
			return nil
		}

		// Track prefix and tab count for multi-match behavior
		if input == sh.lastPrefix {
			sh.tabCount++
		} else {
			sh.lastPrefix = input
			sh.tabCount = 1
		}

		sort.Strings(results)

		if len(results) == 1 {
			return []string{results[0] + " "}
		}

		// Common prefix calculation
		common := longestCommonPrefix(results)
		if common != input {
			return []string{common}
		}

		if sh.tabCount == 1 {
			fmt.Print("\a")
			return nil
		}

		fmt.Print("\n" + strings.Join(results, "  ") + "\n")
		fmt.Print("$ " + input)
		sh.tabCount = 0
		return nil
	})
}

// longestCommonPrefix finds the longest common prefix in a slice of strings
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]
	for _, s := range strs[1:] {
		for !strings.HasPrefix(s, prefix) {
			prefix = prefix[:len(prefix)-1]
			if prefix == "" {
				return ""
			}
		}
	}
	return prefix
}