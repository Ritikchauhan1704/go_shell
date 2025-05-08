package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// runExternalCommand executes an external program
func runExternalCommand(cmd string, args []string) {
	path, err := exec.LookPath(cmd)
	if err != nil {
		fmt.Printf("%s: not found\n", cmd)
		return
	}
	
	// Create the command
	ec := exec.Command(path, args...)
	ec.Args[0] = filepath.Base(path)
	
	// Connect standard I/O
	ec.Stdin = os.Stdin
	ec.Stdout = os.Stdout
	ec.Stderr = os.Stderr
	
	// Execute the command
	err = ec.Run()
	if err != nil {
		// We don't print the error since most commands handle their own error messages
		// and this would result in duplicate messages
		return
	}
}