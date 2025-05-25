# ToyCalc - A Complex Number Command-Line and Interactive Calculator

ToyCalc is a versatile calculator written in Go that operates on complex numbers. It provides a command-line interface for quick calculations and an interactive Read-Eval-Print Loop (REPL) for a more engaging experience. All calculations are performed using `complex128` for robust handling of complex arithmetic.

## Current Features (End of Stage 1)

* **Core Arithmetic:** Supports `+` (addition), `-` (subtraction), `*` (multiplication), `/` (division).
* **Power Operator:** `^` for exponentiation (e.g., `2^3`, `(-4)^0.5`, `(1+i)^2`). Uses principal values.
* **Modulo Operator (`%`):** Implemented for complex numbers using the Gaussian integer remainder definition (remainder `r = a - xb`, where `x` is the complex integer closest to `a/b`). This ensures $|r/b| \le 1/\sqrt{2}$.
* **Unary Operators:**
    * Unary minus (e.g., `-5`, `-(1+2i)`). Correctly handles sign for principal value calculations (e.g., `(-4)^0.5` results in `2i`).
    * Unary plus (e.g., `+5`) is recognized but has no operational effect.
* **Grouping Symbols:** Supports `()`, `[]`, and `{}` interchangeably for grouping expressions.
* **Basic Functions:**
    * `log(x)`: Natural logarithm (principal value, works for complex `x`).
    * `exp(x)`: Exponential function $e^x$ (works for complex `x`).
* **Imaginary Unit `i`:** The constant `i` is recognized as `complex(0, 1)`.
* **Complex Number Backend:** All calculations internally use Go's `complex128` type.
* **Output Formatting:**
    * Results are displayed as real numbers if the imaginary part is negligible (close to zero within a defined epsilon).
    * Otherwise, full complex numbers are displayed (e.g., `3+2i`, `-1-0.5i`).
    * Handles `NaN` and `Infinity` display.
* **Two Modes of Operation:**
    1.  **Command-Line (CLI):** Evaluate expressions directly.
    2.  **Interactive (REPL):** Start `toycalc` without arguments to enter an interactive session.
* **Integrated Help System:** Basic help available via `toycalc help` or `help [topic]` in REPL.

## Usage

### Command-Line Interface (CLI)

To evaluate an expression directly:

```bash
toycalc <expression>
```

**Examples:**

```bash
toycalc 10 + 2 * 3
toycalc "(1+2*i) / (3-i) - ( -4 ) ^ 0.5"
toycalc 10 % 3.2
toycalc log(-1)
toycalc exp(i*pi) # (Note: 'pi' constant will be in Stage 3, use its value for now)
```

* It's **highly recommended to quote expressions** containing spaces or shell special characters (like `*`, `(`, `)`, `^`, `!`) to ensure the shell passes the expression to `toycalc` correctly.
    Example: `toycalc "2 * (3 + 4)!"`

### Interactive Mode (REPL)

To start the interactive mode, simply run `toycalc` without any arguments:

```bash
./toycalc
```

You will see a prompt:

```
ToyCalc Interactive Mode (v0.1 Stage 1)
Type 'exit' or 'quit' to leave, or 'help' for assistance.
>>>
```

Then, type your expressions and press Enter:

```
>>> 10 % -3
1
>>> (1+i)^2
0+2i
>>> log(-1)
0+3.141592653589793i
>>> help log
(Help text for log will be displayed)
>>> exit
Exiting ToyCalc.
```

* Type `exit` or `quit` to leave the interactive mode.
* Type `help` or `help [topic]` for assistance.

## Building from Source

1.  Ensure you have Go installed (version 1.16+ recommended).
2.  Clone the repository (or ensure all source files: `main.go`, `core.go`, `lexer.go`, `parser.go`, `evaluator.go`, `help.go`, `toycalc_test.go` are in the same directory).
3.  Navigate to the project's root directory.
4.  Run:
    ```bash
    go build
    ```
5.  This will create the `toycalc` executable in the current directory.

## Running Tests

To run the test suite:

```bash
go test -v
```

## Planned Future Stages (Roadmap)

* **Stage 2: Standard Mathematical Functions:**
    * Trigonometric functions (sin, cos, tan, asin, acos, atan, atan2).
    * Hyperbolic functions (sinh, cosh, tanh, asinh, acosh, atanh).
    * Additional logarithmic functions (log10, log2).
    * Explicit `sqrt(x)` function.
* **Stage 3: Advanced Operators & Combinatorial Functions:**
    * Factorial operator `x!` (using Gamma function: $\Gamma(x+1)$).
    * Combinatorial functions: `nCr(n,k)` and `nPr(n,k)` (using Gamma functions).
    * Predefined constants: `pi`, `e`.
* **Stage 4: Usability & Parser Enhancements:**
    * Implied multiplication (e.g., `2(3+4)`).
    * Improved error reporting with more precise positions.
    * More detailed and categorized help system.
* **Stage 5: Advanced Numeric Features:**
    * Full complex number input parsing (e.g., "3+2.5i").
    * User-defined variables.
    * (Potentially) User-defined functions.
    * (Potentially) Support for arbitrary-precision numbers (`big.Float`, `big.Complex`).

## Contributing

(Details to be added if the project becomes open for contributions - e.g., coding style, pull request process).

## License

MIT License

Copyright (c) 2025 Vladimir Pérez

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
