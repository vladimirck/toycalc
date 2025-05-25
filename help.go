// help.go
package main

import (
	"fmt"
	"math" // Required for Pi in one example, will use its value directly
	"math/cmplx"
	"strings"
)

// helpTopics stores the detailed help information for each topic.
// The key is the topic name (lowercase).
var helpTopics = map[string]string{
	"usage": "Usage:\n" +
		"  toycalc <expression>\n" +
		"  toycalc \"<expression with spaces or special characters>\"\n" +
		"  toycalc help [topic]\n\n" +
		"If no arguments are provided, toycalc starts in interactive mode (REPL).\n" +
		"In REPL, type an expression and press Enter, or type 'help [topic]', 'exit', or 'quit'.",

	"general": "ToyCalc is a command-line and interactive calculator that works with complex numbers.\n" +
		"All calculations use complex128 arithmetic. Results with a negligible imaginary part\n" +
		"(near zero) are displayed as real numbers. Pure imaginary numbers are shown like '2i' or '-3.5i'.\n\n" +
		"Supported features in this version (Stage 1):\n" +
		"- Basic arithmetic: +, -, *, /\n" +
		"- Power: ^\n" +
		"- Modulo: %\n" +
		"- Unary plus (+) and minus (-)\n" +
		"- Grouping: (), [], {}\n" +
		"- Functions: log(x) (natural), exp(x)\n" +
		"- Constant: i (imaginary unit, use as 'i', e.g., '5*i')\n\n" +
		"Type 'help <topic>' for more information on a specific feature (e.g., 'help +', 'help log').",

	"operators": "Supported operators:\n" +
		"  +  : Addition (binary)\n" +
		"  -  : Subtraction (binary) / Unary Minus (prefix)\n" +
		"  * : Multiplication (binary)\n" +
		"  /  : Division (binary)\n" +
		"  %  : Modulo (binary)\n" +
		"  ^  : Power (binary)\n\n" +
		"See 'help <operator_symbol>' or 'help unary' or 'help modulo' for details.",

	"unary": "Unary Plus and Minus:\n" +
		"  -x : Negation. Example: -5, -(1+2*i)\n" +
		"       For real numbers like -4, this is treated as complex(-4, +0.0)\n" +
		"       for consistent principal value results in functions like power (e.g. (-4)^0.5 results in 2i).\n" +
		"  +x : Unary Plus. Example: +5. This operator is recognized but has no\n" +
		"       effect on the value (e.g., +5 evaluates to 5).",

	"+": "Operator: + (Addition / Unary Plus)\n" +
		"  Binary Addition: Adds two complex numbers.\n" +
		"    Example: 3+2*i + (1-1*i)  (Result: 4+1i)\n" + // Corrected example
		"    Example: 5 + 2          (Result: 7)\n" +
		"  Unary Plus: Indicates a positive number. It has no mathematical effect.\n" +
		"    Example: +5             (Result: 5)\n" +
		"    Example: 10 * +2        (Result: 20)",

	"-": "Operator: - (Subtraction / Unary Minus)\n" +
		"  Binary Subtraction: Subtracts the second complex number from the first.\n" +
		"    Example: (3+2*i) - (1-1*i)  (Result: 2+3i)\n" + // Corrected example
		"    Example: 5 - 2              (Result: 3)\n" +
		"  Unary Minus: Negates a complex number.\n" +
		"    Example: -5                 (Result: -5)\n" +
		"    Example: -(1+2*i)           (Result: -1-2i)\n" + // Corrected example
		"    Example: 10 * -2            (Result: -20)",

	"*": "Operator: * (Multiplication)\n" +
		"  Multiplies two complex numbers.\n" +
		"    Example: (1+i) * (1-i)    (Result: 2)\n" +
		"    Example: 3 * 4            (Result: 12)\n" +
		"    Example: 5*i              (Result: 5i)", // Corrected example

	"/": "Operator: / (Division)\n" +
		"  Divides the first complex number by the second.\n" +
		"    Example: (1+i) / i        (Result: 1-1i)\n" +
		"    Example: 10 / 4           (Result: 2.5)\n" +
		"  Division by zero (0+0i) results in NaN or Inf components, as per complex arithmetic rules.",

	"%": "Operator: % (Modulo)\n" +
		"  Calculates the modulo of two complex numbers a % b = r, such that r = a - x*b,\n" +
		"  where x is the complex integer (Gaussian integer) closest to a/b.\n" +
		"  This definition ensures the remainder r has a small magnitude relative to b.\n" +
		"    Example: 10 % 3           (Result: 1)\n" +
		"    Example: 10.5 % 3.2       (Result: 0.9)\n" +
		"    Example: (5+3*i) % (2+i)  (Result: -1)\n" + // Corrected example
		"  Modulo by zero results in an error.",

	"^": "Operator: ^ (Power)\n" +
		"  Raises a complex base to a complex exponent (base^exponent).\n" +
		"  Returns the principal value.\n" +
		"    Example: 2^3              (Result: 8)\n" +
		"    Example: 16^0.5           (Result: 4)\n" +
		"    Example: (-4)^0.5          (Result: 2i)\n" + // Corrected example
		"    Example: i^2              (Result: -1)",

	"grouping": "Grouping Symbols: (), [], {}\n" +
		"  Parentheses `()`, square brackets `[]`, and curly braces `{}` can all be used\n" +
		"  interchangeably to group sub-expressions and control the order of operations.\n" +
		"  They must be correctly matched.\n" +
		"    Example: (1 + 2) * 3\n" +
		"    Example: {[ (10 - 2) / 4 ] + 1}^2",

	"functions": "Supported functions (Stage 1):\n" +
		"  log(x)   : Natural logarithm (base e), principal value.\n" +
		"  exp(x)   : Exponential function (e^x).\n\n" +
		"Type 'help <function_name>' for more details (e.g., 'help log').",

	"log": "Function: log(x)\n" +
		"  Calculates the natural logarithm (base e) of the complex number x.\n" +
		"  Returns the principal value. The imaginary part of the result is in (-π, π].\n" +
		"    Example: log(exp(2))      (Result: 2)\n" +
		"    Example: log(-1)           (Result: " + fmt.Sprintf("%gi", math.Pi) + ")\n" + // Corrected output
		"    Example: log(i)            (Result: " + fmt.Sprintf("%gi", math.Pi/2) + ")\n" + // Corrected output
		"  log(0) results in " + fmt.Sprintf("%v", cmplx.Log(0)) + ".", // Show actual Inf/NaN output

	"exp": "Function: exp(x)\n" +
		"  Calculates the exponential function e^x, where e is Euler's number, for the complex number x.\n" +
		"    Example: exp(0)             (Result: 1)\n" +
		"    Example: exp(1)             (Result: " + fmt.Sprintf("%g", math.E) + ")\n" +
		"    Example: exp(log(5))        (Result: 5)\n" +
		"    Example: exp(i * " + fmt.Sprintf("%g", math.Pi) + ") (Result: -1, Euler's Identity)",

	"i": "Constant: i\n" +
		"  The imaginary unit, evaluated as complex(0, 1).\n" +
		"  Must be used with multiplication operator if scaling, e.g., '5*i'.\n" +
		"    Example: i*i              (Result: -1)\n" +
		"    Example: 2+3*i            (Result: 2+3i)\n" +
		"    Example: exp(i*" + fmt.Sprintf("%g", math.Pi/2) + ")    (Result: i)",

	"output": "Output Formatting:\n" +
		"  Results are displayed as complex numbers (e.g., a+bi or a-bi).\n" +
		"  If the imaginary part of a result is very close to zero (negligible based on internal epsilon),\n" +
		"  only the real part is displayed.\n" +
		"  If the real part is negligible and the imaginary part is not, the output is like '2i' or '-3.5i'.\n" +
		"  Whole numbers (in real or imaginary parts) are displayed without unnecessary decimal points (e.g., '5' instead of '5.0').\n" +
		"  Handles 'NaN' (Not a Number) and standard complex 'Inf' (Infinity) representations for results where applicable.",
}

// displayHelp shows help information.
// If topic is empty, it shows general help or a list of topics.
// If topic is specified, it shows help for that topic.
func displayHelp(topic string) {
	topic = strings.ToLower(strings.TrimSpace(topic))
	availableTopics := []string{"usage", "general", "operators", "unary", "+", "-", "*", "/", "%", "^", "grouping", "functions", "log", "exp", "i", "output"}

	if topic == "" {
		fmt.Println(helpTopics["general"])
		fmt.Println("\nAvailable topics (type 'help <topic>'):")
		// Simple way to list topics, can be formatted better if many topics
		topicsStr := strings.Join(availableTopics, ", ")
		fmt.Println("  " + topicsStr)
		return
	}

	content, found := helpTopics[topic]
	if found {
		fmt.Println(content)
	} else {
		fmt.Printf("Sorry, no help found for topic: '%s'\n", topic)
		fmt.Println("Type 'help' for a list of available topics.")
	}
}
