// help.go
package toycalc_core

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
		"Supported features include:\n" + // Changed heading slightly
		"- Basic arithmetic: +, -, *, /\n" +
		"- Power: ^\n" +
		"- Modulo: % (Gaussian integer remainder)\n" +
		"- Unary plus (+) and minus (-)\n" +
		"- Grouping: (), [], {}\n" +
		"- Constants: i, pi, e (see 'help constants')\n" +
		"- A wide range of mathematical functions including logarithmic, exponential, trigonometric,\n" +
		"  hyperbolic, complex component manipulation, angle conversion, and rounding.\n" +
		"  (Type 'help functions' for a full list).\n\n" +
		"Type 'help <topic>' for more information on a specific feature (e.g., 'help +', 'help log').",

	"set format": "Command: set format <mode> [N]\n" +
		"  Sets the output display format for numbers.\n" +
		"  Modes:\n" +
		"    auto         : Default. Concise output, integers shown without decimals (e.g., 5, 3.14, 1.2e-5).\n" +
		"                   Uses the current display precision for pre-rounding.\n" +
		"    fixed <N>    : Fixed-point decimal notation with N digits after the decimal point.\n" +
		"                   Example: set format fixed 4  (Output for pi: 3.1416)\n" +
		"    sci <N>      : Scientific notation with N digits after the decimal point for the significand.\n" +
		"                   Example: set format sci 6  (Output for pi: 3.141590e+00)\n" +
		"  N is an integer, typically 0-20. This N also updates the general display precision.",

	"set precision": "Command: set precision <N>\n" +
		"  Sets the number of decimal places (N) to which numbers are rounded for display purposes\n" +
		"  before being formatted according to the current format mode ('auto', 'fixed', or 'sci').\n" +
		"  This also sets the precision used by 'fixed N' and 'sci N' formats directly.\n" +
		"  N is an integer, typically 0-20.\n" +
		"    Example: set precision 9 (default for 'auto' pre-rounding)\n" +
		"    Example: set format fixed 2 (equivalent to 'set format fixed' then 'set precision 2' for fixed mode)",
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

	"functions": "Supported functions (all operate on complex numbers):\n" + // Emphasize complex operation
		"  Core: real(x), imag(x), abs(x), phase(x), conj(x)\n" +
		"  Log/Exp: exp(x), log(x) (natural), log10(x), log2(x)\n" +
		"  Power/Root: sqrt(x) (Note: '^' is the power operator)\n" +
		"  Trigonometric: sin(x), cos(x), tan(x)\n" +
		"  Inverse Trig: asin(x), acos(x), atan(x)\n" +
		"  Hyperbolic: sinh(x), cosh(x), tanh(x)\n" +
		"  Inverse Hyperbolic: asinh(x), acosh(x), atanh(x)\n" +
		"  Angle Conversion: degToRad(x), radToDeg(x)\n" +
		"  Rounding/Truncation: floor(x), ceil(x), round(x), trunc(x)\n\n" +
		"Type 'help <function_name>' for more details (e.g., 'help sin').",

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
		"  Results are displayed as complex numbers. Formatting can be controlled.\n" +
		"  See 'help set format' and 'help set precision' for details.\n\n" +
		"  Default ('auto' mode) behavior:\n" +
		"  - If imaginary part is negligible, only the real part is displayed.\n" +
		"  - If real part is negligible, output is like '2i' or '-i'.\n" +
		"  - Whole numbers formatted without unnecessary decimals (e.g., '5').\n" +
		"  - Handles 'NaN' and complex 'Inf' representations.",

	"constants": "Supported constants:\n" +
		"  i  : The imaginary unit, complex(0, 1).\n" +
		"  pi : The mathematical constant π (Pi), approx. 3.1415926535...\n" +
		"  e  : Euler's number (base of natural logarithm), approx. 2.7182818284...\n\n" +
		"Type 'help <constant_name>' for more details (e.g., 'help pi').",

	"pi": "Constant: pi\n" +
		"  Represents the mathematical constant π (Pi), the ratio of a circle's circumference to its diameter.\n" +
		"  Value: " + fmt.Sprintf("%.10f...", math.Pi) + "\n" + // Show some precision
		"    Example: sin(pi/2)    (Result: 1)\n" +
		"    Example: 2*pi         (Result: " + fmt.Sprintf("%g", 2*math.Pi) + ")",

	"e": "Constant: e\n" +
		"  Represents Euler's number, the base of the natural logarithm.\n" +
		"  Value: " + fmt.Sprintf("%.10f...", math.E) + "\n" +
		"    Example: log(e)         (Result: 1)\n" +
		"    Example: e^2            (Result: " + fmt.Sprintf("%g", math.E*math.E) + ")",
	"sin": "Function: sin(x)\n" +
		"  Calculates the trigonometric sine of the complex number x.\n" +
		"  x is assumed to be in radians.\n" +
		"    Example: sin(0)         (Result: 0)\n" +
		"    Example: sin(" + fmt.Sprintf("%g", math.Pi/2) + ") (Result: 1)\n" +
		"    Example: sin(i)         (Result: " + fmt.Sprintf("%gi", math.Sinh(1)) + ") (since sin(ix) = i*sinh(x))",

	"cos": "Function: cos(x)\n" +
		"  Calculates the trigonometric cosine of the complex number x.\n" +
		"  x is assumed to be in radians.\n" +
		"    Example: cos(0)         (Result: 1)\n" +
		"    Example: cos(" + fmt.Sprintf("%g", math.Pi) + ")   (Result: -1)\n" +
		"    Example: cos(i)         (Result: " + fmt.Sprintf("%g", math.Cosh(1)) + ") (since cos(ix) = cosh(x))",

	"tan": "Function: tan(x)\n" +
		"  Calculates the trigonometric tangent of the complex number x (sin(x)/cos(x)).\n" +
		"  x is assumed to be in radians.\n" +
		"  Result may be Inf or NaN if cos(x) is zero (e.g., at pi/2, 3pi/2).\n" +
		"    Example: tan(0)         (Result: 0)\n" +
		"    Example: tan(" + fmt.Sprintf("%g", math.Pi/4) + ") (Result: 1)",

	"asin": "Function: asin(x)\n" +
		"  Calculates the principal value of the inverse trigonometric sine (arcsine) of x.\n" +
		"    Example: asin(0)        (Result: 0)\n" +
		"    Example: asin(1)        (Result: " + fmt.Sprintf("%g", math.Pi/2) + ")",

	"acos": "Function: acos(x)\n" +
		"  Calculates the principal value of the inverse trigonometric cosine (arccosine) of x.\n" +
		"    Example: acos(1)        (Result: 0)\n" +
		"    Example: acos(0)        (Result: " + fmt.Sprintf("%g", math.Pi/2) + ")",

	"atan": "Function: atan(x)\n" +
		"  Calculates the principal value of the inverse trigonometric tangent (arctangent) of x.\n" +
		"    Example: atan(0)        (Result: 0)\n" +
		"    Example: atan(1)        (Result: " + fmt.Sprintf("%g", math.Pi/4) + ")",

	"sinh": "Function: sinh(x)\n" +
		"  Calculates the hyperbolic sine of the complex number x.\n" +
		"    Example: sinh(0)        (Result: 0)\n" +
		"    Example: sin(i*(" + fmt.Sprintf("%g", math.Pi/2) + ")) (Result: i)", // sin(i*x) = i*sinh(x) so sinh(x) = -i*sin(ix)

	"cosh": "Function: cosh(x)\n" +
		"  Calculates the hyperbolic cosine of the complex number x.\n" +
		"    Example: cosh(0)        (Result: 1)\n" +
		"    Example: cos(i)         (Result: " + fmt.Sprintf("%g", math.Cosh(1)) + ")", // cos(ix) = cosh(x)

	"tanh": "Function: tanh(x)\n" +
		"  Calculates the hyperbolic tangent of the complex number x (sinh(x)/cosh(x)).\n" +
		"    Example: tanh(0)        (Result: 0)",

	"asinh": "Function: asinh(x)\n  Calculates the principal value of the inverse hyperbolic sine of x.\n    Example: asinh(0) (Result: 0)",
	"acosh": "Function: acosh(x)\n  Calculates the principal value of the inverse hyperbolic cosine of x.\n    Example: acosh(1) (Result: 0)",
	"atanh": "Function: atanh(x)\n  Calculates the principal value of the inverse hyperbolic tangent of x.\n    Example: atanh(0) (Result: 0)",

	"log10": "Function: log10(x)\n" +
		"  Calculates the base-10 logarithm of the complex number x.\n" +
		"  Returns the principal value.\n" +
		"    Example: log10(100)     (Result: 2)\n" +
		"    Example: log10(1)       (Result: 0)",

	"log2": "Function: log2(x)\n" +
		"  Calculates the base-2 logarithm of the complex number x.\n" +
		"  Returns the principal value.\n" +
		"    Example: log2(8)        (Result: 3)\n" +
		"    Example: log2(1)        (Result: 0)",

	"sqrt": "Function: sqrt(x)\n" +
		"  Calculates the principal value of the square root of the complex number x.\n" +
		"  Equivalent to x^0.5.\n" +
		"    Example: sqrt(4)        (Result: 2)\n" +
		"    Example: sqrt(-1)       (Result: i)\n" + // Output format will show 'i'
		"    Example: sqrt(2i)       (Result: 1+1i)", // sqrt(2i) = 1+i
	"real": "Function: real(x)\n" +
		"  Returns the real part of the complex number x, as a complex number with a zero imaginary part.\n" +
		"    Example: real(3+4*i)    (Result: 3)\n" +
		"    Example: real(5)        (Result: 5)\n" +
		"    Example: real(2*i)      (Result: 0)",

	"imag": "Function: imag(x)\n" +
		"  Returns the imaginary part of the complex number x, as a complex number with a zero imaginary part.\n" +
		"  Note: This returns the coefficient of 'i'. For the complex number 'i' itself, use the constant 'i'.\n" +
		"    Example: imag(3+4*i)    (Result: 4)\n" +
		"    Example: imag(5)        (Result: 0)\n" +
		"    Example: imag(2*i)      (Result: 2)",

	"abs": "Function: abs(x)\n" +
		"  Calculates the absolute value (or modulus/magnitude) of the complex number x.\n" +
		"  This is a non-negative real number, returned as complex(abs_value, 0).\n" +
		"    Example: abs(3+4*i)    (Result: 5)\n" +
		"    Example: abs(-5)       (Result: 5)\n" +
		"    Example: abs(i)        (Result: 1)",

	"phase": "Function: phase(x)\n" +
		"  Calculates the argument (or phase/angle) of the complex number x.\n" +
		"  The result is in radians, in the interval (-π, π].\n" +
		"  Returned as complex(angle_value, 0).\n" +
		"    Example: phase(1+i)    (Result: " + fmt.Sprintf("%g", math.Pi/4) + ")\n" +
		"    Example: phase(-1)     (Result: " + fmt.Sprintf("%g", math.Pi) + ")\n" +
		"    Example: phase(i)      (Result: " + fmt.Sprintf("%g", math.Pi/2) + ")\n" +
		"    Example: phase(0)      (Result: 0)",

	"conj": "Function: conj(x)\n" +
		"  Calculates the complex conjugate of x.\n" +
		"  If x = a+bi, conj(x) = a-bi.\n" +
		"    Example: conj(3+4*i)    (Result: 3-4i)\n" +
		"    Example: conj(5)        (Result: 5)\n" +
		"    Example: conj(2*i)      (Result: -2i)",

	"degtorad": "Function: degToRad(x)\n" +
		"  Converts the complex number x from degrees to radians.\n" +
		"  The entire complex number (both real and imaginary parts) is scaled by π/180.\n" +
		"    Example: degToRad(180)          (Result: " + fmt.Sprintf("%g", math.Pi) + ")\n" +
		"    Example: degToRad(90+180*i)  (Result: " + fmt.Sprintf("%g+%gi", math.Pi/2, math.Pi) + ")",

	"radtodeg": "Function: radToDeg(x)\n" +
		"  Converts the complex number x from radians to degrees.\n" +
		"  The entire complex number (both real and imaginary parts) is scaled by 180/π.\n" +
		"    Example: radToDeg(pi)           (Result: 180)\n" +
		"    Example: radToDeg(pi/2 + i)     (Result: " + fmt.Sprintf("%g+%gi", 90.0, 180.0/math.Pi) + ")",

	"floor": "Function: floor(x)\n" +
		"  Computes the floor of the complex number x component-wise.\n" +
		"  Result: complex(math.Floor(real(x)), math.Floor(imag(x)))\n" +
		"    Example: floor(3.7+2.3*i)   (Result: 3+2i)\n" +
		"    Example: floor(-3.7-2.3*i)  (Result: -4-3i)",

	"ceil": "Function: ceil(x)\n" +
		"  Computes the ceiling of the complex number x component-wise.\n" +
		"  Result: complex(math.Ceil(real(x)), math.Ceil(imag(x)))\n" +
		"    Example: ceil(3.2+2.8*i)    (Result: 4+3i)\n" +
		"    Example: ceil(-3.2-2.8*i)   (Result: -3-2i)",

	"round": "Function: round(x)\n" +
		"  Rounds the complex number x to the nearest integer component-wise.\n" +
		"  Uses Go's math.Round (rounds half to even).\n" +
		"  Result: complex(math.Round(real(x)), math.Round(imag(x)))\n" +
		"    Example: round(3.5+2.5*i)   (Result: 4+2i)\n" +
		"    Example: round(3.7+2.3*i)   (Result: 4+2i)",

	"trunc": "Function: trunc(x)\n" +
		"  Truncates the complex number x towards zero component-wise.\n" +
		"  Result: complex(math.Trunc(real(x)), math.Trunc(imag(x)))\n" +
		"    Example: trunc(3.7+2.3*i)   (Result: 3+2i)\n" +
		"    Example: trunc(-3.7-2.3*i)  (Result: -3-2i)",
}

// displayHelp shows help information.
// If topic is empty, it shows general help or a list of topics.
// If topic is specified, it shows help for that topic.
func DisplayHelp(topic string) {
	topic = strings.ToLower(strings.TrimSpace(topic))
	availableTopics := []string{
		"usage", "general", "operators", "unary", "+", "-", "*", "/", "%", "^", "grouping",
		"functions", "constants", "output", "i", "pi", "e",
		"log", "exp", "sin", "cos", "tan", "asin", "acos", "atan",
		"sinh", "cosh", "tanh", "asinh", "acosh", "atanh",
		"log10", "log2", "sqrt",
		"real", "imag", "abs", "phase", "conj",
		"degtorad", "radtodeg",
		"floor", "ceil", "round", "trunc",
	} // Ensure all helpTopics keys are listable here if desired for discoverability

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
