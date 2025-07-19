package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/chzyer/readline"                              // Import the readline library
	toycalc_core "github.com/vladimirck/toycalc/toycalc-core" // Import your toycalc core package
)



// processExpression encapsulates the calculation and printing logic.
func processExpression(expressionString string) {
	if strings.TrimSpace(expressionString) == "" {
		return // Do nothing for empty input in REPL
	}
	resultStr, err := toycalc_core.CalculateExpression(expressionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	fmt.Println(resultStr)
}

// startInteractiveMode starts the REPL for toycalc using the readline library.
func startInteractiveMode() {
	fmt.Println("ToyCalc Interactive Mode (v0.3 Stage 3)") // Updated version
	fmt.Println("Type 'exit', 'quit', or 'help' for assistance.")
	fmt.Println("Use 'set format [auto|fixed N|sci N]' to change output format.")
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
		Prompt:              ">>> ",
		HistoryFile:         historyFile,
		AutoComplete:        nil,
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		HistorySearchFold:   true,
		FuncFilterInputRune: nil,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing advanced readline: %v\n", err)
		fmt.Println("Falling back to basic interactive mode.")
		startBasicInteractiveMode()
		return
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()

		if err != nil { // Handle Readline errors
			if err == readline.ErrInterrupt {
				continue // New prompt
			} else if err == io.EOF {
				fmt.Println("Exiting ToyCalc (EOF).")
				break
			}
			fmt.Fprintf(os.Stderr, "Readline error: %v\n", err)
			break
		}

		input := strings.TrimSpace(line)

		if input == "" {
			continue
		}

		lowerInput := strings.ToLower(input)
		parts := strings.Fields(lowerInput) // Split input into words

		if parts[0] == "exit" || parts[0] == "quit" {
			fmt.Println("Exiting ToyCalc.")
			break
		}

		if parts[0] == "help" {
			topic := ""
			if len(parts) > 1 {
				topic = strings.Join(parts[1:], " ")
			}
			toycalc_core.DisplayHelp(topic)
		} else if parts[0] == "set" {
			if len(parts) == 1 {
				fmt.Println("Usage: set <format|precision> <options>")
				fmt.Println("Example: set format fixed 4")
				fmt.Println("         set precision 6")
				continue
			}
			switch parts[1] {
			case "format":
				if len(parts) < 3 {
					fmt.Println("Usage: set format <auto|fixed N|sci N>")
					fmt.Println("Example: set format fixed 4")
					fmt.Println("         set format sci 6")
					fmt.Println("         set format auto")
					continue
				}
				mode := parts[2]
				precision := toycalc_core.OutputDisplayPrecision // Keep current precision if not specified for auto
				if mode == "fixed" || mode == "sci" {
					if len(parts) < 4 {
						fmt.Printf("Usage: set format %s <N> (where N is number of digits)\n", mode)
						continue
					}
					p, err := strconv.Atoi(parts[3])
					if err != nil || p < 0 || p > 20 { // Set a reasonable max precision
						fmt.Println("Error: Precision N must be a non-negative integer (e.g., 0-20).")
						continue
					}
					precision = p
				} else if mode != "auto" {
					fmt.Printf("Error: Unknown format mode '%s'. Use 'auto', 'fixed N', or 'sci N'.\n", mode)
					continue
				}

				toycalc_core.OutputFormatMode = mode
				toycalc_core.OutputDisplayPrecision = precision
				fmt.Printf("Output format set to: %s", toycalc_core.OutputFormatMode)
				if toycalc_core.OutputFormatMode != "auto" {
					fmt.Printf(", %d digits precision", toycalc_core.OutputDisplayPrecision)
				}
				fmt.Println()

			case "precision":
				if len(parts) < 3 {
					fmt.Println("Usage: set precision <N> (where N is number of digits, e.g., 0-20)")
					continue
				}
				p, err := strconv.Atoi(parts[2])
				if err != nil || p < 0 || p > 20 { // Max precision
					fmt.Println("Error: Precision N must be a non-negative integer (e.g., 0-20).")
					continue
				}
				toycalc_core.OutputDisplayPrecision = p
				fmt.Printf("Display precision set to: %d digits (affects 'fixed', 'sci', and pre-rounding for 'auto' mode)\n", toycalc_core.OutputDisplayPrecision)
			default:
				fmt.Printf("Error: Unknown option for 'set': '%s'. Try 'set format ...' or 'set precision ...'.\n", parts[1])

			}
		} else {
			// Process as mathematical expression
			processExpression(input) // Use the original 'input' not 'lowerInput'
		}
	}
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
			toycalc_core.DisplayHelp(topic)
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
			toycalc_core.DisplayHelp(topic)
		} else {
			expressionString := strings.Join(os.Args[1:], " ")
			resultStr, err := toycalc_core.CalculateExpression(expressionString)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(resultStr)
		}
	}
}
