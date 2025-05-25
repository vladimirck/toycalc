// main.go
package main

import (
	"fmt"
	"os"
	"strings" // Useful for joining expression arguments
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: toycalc <expression> | help [topic]")
		fmt.Println("Example: toycalc 1 + 2 * 3")
		fmt.Println("         toycalc \"2 * (3 + 4)\"  (quotes recommended for complex expressions or those with shell special characters)")
		fmt.Println("         toycalc help log")
		os.Exit(1)
	}

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

		// It's still a good idea to warn users about shell interpretation if no quotes are used,
		// especially if common shell metacharacters are detected (more advanced).
		// For now, the usage message provides this guidance.

		resultStr, err := calculateExpression(expressionString)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(resultStr)
	}
}
