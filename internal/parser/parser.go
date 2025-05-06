package parser

import (
    "fmt"
    "os"
    "strings"
)

// Trim removes trailing newline and spaces
func Trim(input string) string {
    return strings.TrimSpace(input)
}

// Split separates command name and arguments, handling quotes and escapes
func Split(input string) (string, []string) {
    var args []string
    var curr strings.Builder
    inSingle, inDouble, esc := false, false, false

    for i := 0; i < len(input); i++ {
        ch := input[i]
        if ch == '\\' && !esc && !inSingle {
            esc = true
            continue
        }
        if ch == '\'' && !esc && !inDouble {
            inSingle = !inSingle
            continue
        }
        if ch == '"' && !esc && !inSingle {
            inDouble = !inDouble
            continue
        }
        if ch == ' ' && !inSingle && !inDouble && !esc {
            if curr.Len() > 0 {
                args = append(args, curr.String())
                curr.Reset()
            }
            continue
        }
        if esc {
            curr.WriteByte(ch)
            esc = false
        } else {
            curr.WriteByte(ch)
        }
    }
    if curr.Len() > 0 {
        args = append(args, curr.String())
    }
    if len(args) == 0 {
        return "", []string{}
    }
    return args[0], args[1:]
}

// ParseRedirection processes stdout/stderr redirection tokens
func ParseRedirection(args []string) (
    clean []string,
    outFile string, outR, outAppend bool,
    errFile string, errR, errAppend bool,
) {
    clean = make([]string, 0, len(args))
    for i := 0; i < len(args); i++ {
        tok := args[i]
        switch tok {
        case ">", "1>":
            if i+1 < len(args) {
                outFile, outR = args[i+1], true
                i++
                continue
            }
            fmt.Fprintln(os.Stderr, "redirect: no file specified")
            return args, "", false, false, "", false, false
        case ">>", "1>>":
            if i+1 < len(args) {
                outFile, outR, outAppend = args[i+1], true, true
                i++
                continue
            }
            fmt.Fprintln(os.Stderr, "redirect: no file specified")
            return args, "", false, false, "", false, false
        case "2>":
            if i+1 < len(args) {
                errFile, errR = args[i+1], true
                i++
                continue
            }
            fmt.Fprintln(os.Stderr, "redirect: no file specified")
            return args, "", false, false, "", false, false
        case "2>>":
            if i+1 < len(args) {
                errFile, errR, errAppend = args[i+1], true, true
                i++
                continue
            }
            fmt.Fprintln(os.Stderr, "redirect: no file specified")
            return args, "", false, false, "", false, false
        }
        clean = append(clean, tok)
    }
    return
}