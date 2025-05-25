// help.go
package main

import "fmt"

// displayHelp (stub) - Shows help information
func displayHelp(topic string) {
	if topic == "" {
		fmt.Println("General help for toycalc (Stage 0)")
		fmt.Println("---------------------------------")
		fmt.Println("Usage: toycalc \"<expression>\"")
		fmt.Println("       toycalc help [topic_name]")
		fmt.Println("\nCalculation functionality is very limited in this stage.")
		fmt.Println("Available help topics: (none yet)")
		// In later stages, list main topics here
	} else {
		fmt.Printf("Help for topic: '%s' (Stage 0)\n", topic)
		fmt.Println("---------------------------------")
		fmt.Println("Detailed documentation for specific topics is not yet implemented.")
		// In later stages, find and display info for 'topic'
	}
}
