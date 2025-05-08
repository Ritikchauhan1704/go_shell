package main

import (
	"fmt"
	"io"
	"os"

	"github.com/Ritikchauhan1704/go_shell/internal/shell"

	"github.com/peterh/liner"
)

func main() {
	// Initialize liner for enhanced line editing
	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)

	// Create shell instance
	sh := shell.New(line)

	// Run the shell
	err := sh.Run()
	if err != nil && err != io.EOF {
		fmt.Fprintln(os.Stderr, "Shell exited with error:", err)
		os.Exit(1)
	}
}