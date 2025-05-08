package shell

import (
	"fmt"
	"os"
	"strings"
)

// parseCommandLine splits the input into command name and arguments
// while handling quotes and escapes
func parseCommandLine(input string) (string, []string) {
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