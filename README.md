🐚 Go Shell

A feature-rich, interactive shell implementation written in Go with advanced line editing capabilities, command completion, and I/O redirection support.

## ✨ Features

### 🎯 Core Functionality
- **Interactive Command Line Interface** with enhanced line editing via `liner`
- **Command History** with up/down arrow navigation
- **Tab Completion** for both built-in and external commands
- **I/O Redirection** support (`>`, `>>`, `2>`, `2>>`)
- **Quote Handling** for arguments with spaces
- **External Command Execution** with proper PATH resolution

### 🛠️ Built-in Commands
- **`exit`** - Exit the shell
- **`echo`** - Display text output
- **`pwd`** - Print current working directory
- **`cd`** - Change directory (supports `~` for home)
- **`type`** - Show command type and location

### 🎨 Advanced Features
- **Smart Tab Completion** with common prefix matching
- **Multi-match Display** showing all available options
- **Quote Parsing** supporting both single (`'`) and double (`"`) quotes
- **Escape Sequence Handling** for special characters
- **Error Handling** with informative messages

## 🚀 Getting Started

```bash
# Clone the repository
git clone https://github.com/Ritikchauhan1704/go_shell
cd go_shell

# Install dependencies
go mod tidy

# Build and run using Makefile
make build    # Builds to bin/main
make run      # Run directly with go run
make run-bin  # Build and run binary

# Or build manually
go build -o bin/main cmd/shell/main.go
./bin/main
```

### Dependencies
- [`github.com/peterh/liner`](https://github.com/peterh/liner) - Enhanced line editing and history

## 📖 Usage Examples

### Basic Commands
```bash
$ echo Hello, World!
Hello, World!

$ pwd
/home/user/go_shell

$ cd /tmp
$ pwd
/tmp

$ type echo
echo is a shell builtin

$ type ls
ls is /bin/ls
```

### I/O Redirection
```bash
# Redirect stdout to file
$ echo "Hello" > output.txt

# Append to file
$ echo "World" >> output.txt

# Redirect stderr
$ invalid_command 2> error.log

# Append stderr
$ another_error 2>> error.log
```

### Quote Handling
```bash
$ echo "Hello World"
Hello World

$ echo 'Single quotes work too'
Single quotes work too

$ echo "Escaped \"quotes\" work"
Escaped "quotes" work
```

### Tab Completion
- Type a partial command and press `Tab` for completion
- Press `Tab` twice to see all available matches
- Works for both built-in commands and external programs in PATH

## 🏗️ Architecture

### Project Structure
```
go_shell/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   └── shell/
│       ├── shell.go         # Main shell logic
│       ├── commands.go      # Built-in commands
│       ├── parser.go        # Command line parsing
│       ├── redirection.go   # I/O redirection handling
│       └── external.go      # External command execution
├── go.mod
└── README.md
```

### Key Components

#### Shell Core (`shell.go`)
- Main shell loop with readline interface
- Command execution orchestration
- Tab completion setup
- History management

#### Command Parser (`parser.go`)
- Tokenizes command line input
- Handles quotes and escape sequences
- Splits commands and arguments

#### I/O Redirection (`redirection.go`)
- Parses redirection operators
- Manages file handles for stdout/stderr redirection
- Supports both overwrite and append modes

#### Built-in Commands (`commands.go`)
- Implements core shell commands
- Extensible command registration system
- Command help and type information

#### External Commands (`external.go`)
- PATH resolution for external programs
- Process creation and management
- Standard I/O piping

## 🔧 Configuration

### Environment Variables
- **`PATH`** - Used for external command lookup
- **`HOME`** - Used for `cd ~` functionality

### Keyboard Shortcuts
- **`Ctrl+C`** - Interrupt current input (continues shell)
- **`Ctrl+D`** - Exit shell (EOF)
- **`Up/Down Arrows`** - Navigate command history
- **`Tab`** - Command completion
- **`Tab Tab`** - Show all matches

### Adding New Built-in Commands

1. Add your command to the `initCommands()` function in `commands.go`
2. Implement the command handler function

Example:
```go
"mycommand": {
    Name:        "mycommand",
    Description: "Does something amazing",
    Handler:     myCommandHandler,
},
```

