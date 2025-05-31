package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline" // Import the readline library
)

// processExpression encapsulates the calculation and printing logic.
func processExpression(expressionString string) {
	if strings.TrimSpace(expressionString) == "" {
		return // Do nothing for empty input in REPL
	}
	resultStr, err := calculateExpression(expressionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	fmt.Println(resultStr)
}

// startInteractiveMode starts the REPL for toycalc using the readline library.
func startInteractiveMode() {
	fmt.Println("ToyCalc Interactive Mode (v0.2 Stage 1 with Readline)") // Updated version
	fmt.Println("Type 'exit' or 'quit' to leave, or 'help' for assistance.")
	fmt.Println("Use arrow keys for history and line editing.")

	var historyFile string
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not get user home directory for history: %v. Using local history file.\n", err)
		historyFile = ".toycalc_history"
	} else {
		historyFile = filepath.Join(homeDir, ".toycalc_history")
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          ">>> ",
		HistoryFile:     historyFile,
		AutoComplete:    nil, // You can add autocompletion later if needed
		InterruptPrompt: "^C",
		EOFPrompt:       "exit", // What happens when Ctrl+D is pressed

		HistorySearchFold:   true, // Case-insensitive history search
		FuncFilterInputRune: nil,  // No input filtering
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing readline: %v\n", err)
		// Fallback to basic input if readline fails (optional)
		fmt.Println("Falling back to basic input due to readline initialization error.")
		startBasicInteractiveMode() // You might want to keep your old bufio version as a fallback
		return
	}
	defer rl.Close()     // Make sure to close readline instance
	rl.SetPrompt(">>> ") // Redundant if set in NewEx, but can be changed dynamically

	for {
		line, err := rl.Readline()

		if err != nil {
			if err == readline.ErrInterrupt { // Ctrl+C
				// If you want Ctrl+C to exit, uncomment below. Otherwise, it just breaks the current line.
				// fmt.Println("Interrupt received, exiting.")
				// break
				continue // Or print a new prompt
			} else if err == io.EOF { // Ctrl+D
				fmt.Println("Exiting ToyCalc (EOF).")
				break
			}
			// Other read errors
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			break // Exit on other persistent errors
		}

		input := strings.TrimSpace(line)

		if input == "" {
			continue // Skip empty lines
		}

		// Add to history (readline handles this automatically if HistoryFile is set,
		// but you can also manually add if needed for more control using rl.SaveHistory(input))
		// For automatic history saving to file on each command, it's usually handled by rl.Close()
		// or if you need explicit saves, you might do it after each command if not using HistoryFile in config.
		// With HistoryFile in NewEx, readline typically handles loading on start and saving on close.

		lowerInput := strings.ToLower(input)
		if lowerInput == "exit" || lowerInput == "quit" {
			fmt.Println("Exiting ToyCalc.")
			break
		}

		if strings.HasPrefix(lowerInput, "help") {
			parts := strings.Fields(input)
			topic := ""
			if len(parts) > 1 {
				topic = strings.Join(parts[1:], " ")
			}
			displayHelp(topic)
		} else {
			processExpression(input)
		}
	}
	// Save history on exit (readline might do this automatically on Close if HistoryFile is set)
	// err = rl.SaveHistory(historyFile) // This might be redundant if HistoryFile is used in NewEx
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error saving history: %v\n", err)
	// }
}

// startBasicInteractiveMode is your original interactive mode as a fallback
func startBasicInteractiveMode() {
	fmt.Println("ToyCalc Interactive Mode (v0.1 Stage 1 - Basic)")
	fmt.Println("Type 'exit' or 'quit' to leave, or 'help' for assistance.")
	reader := NewStdinReader() // Custom function to create bufio.Reader if you want to keep it
	for {
		fmt.Print(">>> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" { // bufio uses err.Error() == "EOF"
				fmt.Println("Exiting ToyCalc.")
				break
			}
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			continue
		}
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		lowerInput := strings.ToLower(input)
		if lowerInput == "exit" || lowerInput == "quit" {
			fmt.Println("Exiting ToyCalc.")
			break
		}
		if strings.HasPrefix(lowerInput, "help") {
			parts := strings.Fields(input)
			topic := ""
			if len(parts) > 1 {
				topic = strings.Join(parts[1:], " ")
			}
			displayHelp(topic)
		} else {
			processExpression(input)
		}
	}
}

// NewStdinReader is a helper for the basic mode if you keep it.
// This is just to avoid a direct bufio.NewReader in startBasicInteractiveMode
// if you are removing "bufio" import completely when readline works.
// For now, let's assume bufio is still available if needed for the fallback.
//import "bufio" // Keep this import if you use the fallback

func NewStdinReader() *bufio.Reader {
	return bufio.NewReader(os.Stdin)
}

func main() {
	if len(os.Args) < 2 {
		startInteractiveMode()
	} else {
		firstArg := strings.ToLower(os.Args[1])
		if firstArg == "help" {
			topic := ""
			if len(os.Args) > 2 {
				topic = strings.Join(os.Args[2:], " ")
			}
			displayHelp(topic)
		} else {
			expressionString := strings.Join(os.Args[1:], " ")
			resultStr, err := calculateExpression(expressionString)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(resultStr)
		}
	}
}
