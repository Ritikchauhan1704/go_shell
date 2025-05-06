package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Command defines a shell command
type Command struct {
    Name        string
    Description string
    Handler     func(args []string)
}

var builtins map[string]Command

func init() {
    builtins = make(map[string]Command)
    register(Command{"exit", "Exit the shell", exit})
    register(Command{"echo", "Display a line of text", echo})
    register(Command{"pwd", "Print the working directory", pwd})
    register(Command{"cd", "Change the directory", cd})
    register(Command{"type", "Describe a command", shellType})
}

func register(cmd Command) {
    builtins[cmd.Name] = cmd
}

// Run dispatches to builtin or external
func Run(name string, args []string, input string) {
    if cmd, ok := builtins[name]; ok {
        cmd.Handler(args)
    } else {
        external(name, args)
    }
}

// HandleRedirection wraps execution with stdout/stderr redirection
func HandleRedirection(
    outFile string, outR, outAppend bool,
    errFile string, errR, errAppend bool,
    execFunc func(),
) {
    var origOut, origErr *os.File
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
            return
        }
        origOut = os.Stdout
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
            fmt.Fprintf(os.Stderr, "redirect error: %v\n", err)
            if outR { os.Stdout = origOut }
            return
        }
        origErr = os.Stderr
        os.Stderr = fErr
        defer fErr.Close()
    }

    execFunc()

    if outR {
        os.Stdout.Close()
        os.Stdout = origOut
    }
    if errR {
        os.Stderr.Close()
        os.Stderr = origErr
    }
}

func external(name string, args []string) {
    path, err := exec.LookPath(name)
    if err != nil {
        fmt.Printf("%s: not found\n", name)
        return
    }
    cmd := exec.Command(path, args...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Run()
}

// Builtin Handlers
func exit(args []string) {
    if len(args) == 1 && args[0] == "0" {
        os.Exit(0)
    }
}

func echo(args []string) {
    fmt.Println(strings.Join(args, " "))
}

func pwd(args []string) {
    wd, _ := os.Getwd()
    fmt.Println(wd)
}

func cd(args []string) {
    if len(args) != 1 {
        fmt.Println("Usage: cd <dir>")
        return
    }
    dir := args[0]
    if dir == "~" {
        dir = os.Getenv("HOME")
    }
    if err := os.Chdir(dir); err != nil {
        fmt.Printf("cd: %s: No such file or directory\n", dir)
    }
}

func shellType(args []string) {
    if len(args) != 1 {
        fmt.Println("Usage: type <command>")
        return
    }
    name := args[0]
    if _, ok := builtins[name]; ok {
        fmt.Printf("%s is a shell builtin\n", name)
        return
    }
    path, err := exec.LookPath(name)
    if err != nil {
        fmt.Printf("%s: not found\n", name)
    } else {
        fmt.Printf("%s is %s\n", name, path)
    }
}