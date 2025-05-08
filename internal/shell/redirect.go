package shell

import (
	"fmt"
	"os"
)

// parseRedirection processes redirection operators in the command arguments
// and returns clean arguments and redirection information
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