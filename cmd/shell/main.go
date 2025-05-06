package main

import (
    "bufio"
    "fmt"
    "os"

    "github.com/Ritikchauhan1704/go_shell/internal/commands"
    "github.com/Ritikchauhan1704/go_shell/internal/parser"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("$ ")
        input, err := reader.ReadString('\n')
        if err != nil {
            fmt.Fprintln(os.Stderr, "Error reading input:", err)
            os.Exit(1)
        }

        input = parser.Trim(input)
        if input == "" {
            continue
        }

        cmdName, args := parser.Split(input)

        // Handle redirections
        cleanArgs, outFile, outR, outAppend, errFile, errR, errAppend := parser.ParseRedirection(args)
        if outR || errR {
            commands.HandleRedirection(outFile, outR, outAppend, errFile, errR, errAppend, func() {
                commands.Run(cmdName, cleanArgs, input)
            })
            continue
        }

        commands.Run(cmdName, cleanArgs, input)
    }
}