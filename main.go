// main.go
package main

import (
	"bufio" // For reading interactive input
	"fmt"
	"os"
	"strings"
)

// processExpression encapsulates the calculation and printing logic.
// This can be used by both CLI mode and interactive mode.
func processExpression(expressionString string) {
	if strings.TrimSpace(expressionString) == "" {
		return // Do nothing for empty input in REPL
	}

	resultStr, err := calculateExpression(expressionString)
	if err != nil {
		// For REPL, print error to stdout and continue. For CLI, it exits in main.
		// We'll handle CLI exit in main based on this function returning an error.
		// Let's modify calculateExpression to return error, and this func to also return error.
		// Or, for simplicity now, print error here and main decides to exit or not.
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	fmt.Println(resultStr)
}

// startInteractiveMode starts the REPL for toycalc.
func startInteractiveMode() {
	fmt.Println("ToyCalc Interactive Mode (v0.1 Stage 1)") // Adjust version/stage as you like
	fmt.Println("Type 'exit' or 'quit' to leave, or 'help' for assistance.")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">>> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			// Handle EOF (Ctrl+D) gracefully
			if err.Error() == "EOF" {
				fmt.Println("\nExiting ToyCalc.")
				break
			}
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			continue // Or break, depending on desired robustness
		}

		input = strings.TrimSpace(input) // Remove leading/trailing whitespace, including the newline

		if input == "" {
			continue // Skip empty lines
		}

		lowerInput := strings.ToLower(input)
		if lowerInput == "exit" || lowerInput == "quit" {
			fmt.Println("Exiting ToyCalc.")
			break
		}

		if strings.HasPrefix(lowerInput, "help") {
			parts := strings.Fields(input) // Split "help topic"
			topic := ""
			if len(parts) > 1 {
				topic = strings.Join(parts[1:], " ")
			}
			displayHelp(topic) // Call your existing help display function
		} else {
			// Process the expression
			processExpression(input)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		// No arguments provided, start interactive mode
		startInteractiveMode()
	} else {
		// Arguments provided, process as CLI command
		firstArg := strings.ToLower(os.Args[1])

		if firstArg == "help" {
			topic := ""
			if len(os.Args) > 2 {
				// For 'help some topic with spaces', join args after 'help'
				topic = strings.Join(os.Args[2:], " ")
			}
			displayHelp(topic)
		} else {
			// If the first argument is not 'help', assume all arguments from index 1 onwards
			// form the expression.
			expressionString := strings.Join(os.Args[1:], " ")

			// Call calculateExpression and handle error for CLI mode (exit on error)
			resultStr, err := calculateExpression(expressionString)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1) // Exit for CLI mode on error
			}
			fmt.Println(resultStr)
		}
	}
}
