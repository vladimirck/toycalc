# ToyCalc - A Complex Number Command-Line and Interactive Calculator

ToyCalc is a versatile calculator written in Go that operates on complex numbers. It provides a command-line interface for quick calculations and an interactive Read-Eval-Print Loop (REPL) for a more engaging experience. All calculations are performed using `complex128` for robust handling of complex arithmetic.

## Current Features

* **Core Arithmetic:** `+`, `-`, `*`, `/`.
* **Power Operator:** `^` (principal value).
* **Modulo Operator (`%`):** Gaussian integer remainder definition.
* **Unary Operators:** `+` (no-op), `-` (negation, handles sign for principal values, e.g., `(-4)^0.5` is `2i`).
* **Implied Multiplication:** Supports common cases like `2(3+4)`, `(1+2)(3+4)`, `3i`, `2log(x)`, `sin(pi)cos(pi)`.
* **Grouping Symbols:** `()`, `[]`, `{}` (interchangeable).
* **Constants:**
    * `i` (imaginary unit).
    * `pi` (mathematical constant $\pi$).
    * `e` (Euler's number).
* **Complex Number Backend:** All calculations use Go's `complex128`.
* **Output Formatting:**
    * Real numbers shown if imaginary part is negligible.
    * Pure imaginary numbers shown as `Ni` (e.g., `2i`, `-3.5i`).
    * Whole numbers formatted without unnecessary decimals.
    * Handles `NaN` and `Infinity`.
* **Two Modes of Operation:**
    1.  **Command-Line (CLI):** Evaluate expressions directly.
    2.  **Interactive (REPL):** With line editing and persistent command history (`~/.toycalc_history`).
* **Core Complex Functions:**
    * `real(x)`, `imag(x)`: Extract real/imaginary parts.
    * `abs(x)`: Magnitude (modulus).
    * `phase(x)`: Argument/angle in radians $(-\pi, \pi]$.
    * `conj(x)`: Complex conjugate.
* **Exponential & Logarithmic Functions (Principal Values):**
    * `exp(x)`: $e^x$.
    * `log(x)`: Natural logarithm.
    * `log10(x)`: Base-10 logarithm.
    * `log2(x)`: Base-2 logarithm.
    * `sqrt(x)`: Principal square root.
* **Trigonometric Functions (Radians, Principal Values for Inverses):**
    * `sin(x)`, `cos(x)`, `tan(x)`
    * `asin(x)`, `acos(x)`, `atan(x)`
* **Hyperbolic Functions (Principal Values for Inverses):**
    * `sinh(x)`, `cosh(x)`, `tanh(x)`
    * `asinh(x)`, `acosh(x)`, `atanh(x)`
* **Angle Conversion Functions (Operate on full complex number):**
    * `degToRad(x)`: Scales complex number by $\pi/180$.
    * `radToDeg(x)`: Scales complex number by $180/\pi$.
* **Component-wise Integer Functions:**
    * `floor(x)`, `ceil(x)`, `round(x)`, `trunc(x)`
* **Integrated Help System:** `help [topic]` available in CLI and REPL.

## Usage

### Command-Line Interface (CLI)

To evaluate an expression directly:

```bash
toycalc <expression>
```

**Examples:**

```bash
toycalc 10 + 2 * 3
toycalc "(1+2i) / (3-i) - (-4)^0.5" # Explicit multiplication for i
toycalc 2(3+1) # Implied multiplication
toycalc 10 % 3.2
toycalc log(-1)
toycalc exp(i*pi)
toycalc sin(pi/2)cos(pi/2) # Implied multiplication between functions
```

* It's **highly recommended to quote expressions** containing spaces or shell special characters (like `*`, `(`, `)`, `^`) to ensure the shell passes the expression to `toycalc` correctly.
    Example: `toycalc "2 * ( (1+i)^2 + log(e) )"`

### Interactive Mode (REPL)

To start the interactive mode, simply run `toycalc` without any arguments:

```bash
./toycalc
```

You will see a prompt:

```
ToyCalc Interactive Mode (v0.2 Stage 2 with Readline & Implied Multiplication)
Type 'exit' or 'quit' to leave, or 'help' for assistance.
Use arrow keys for history and line editing.
>>>
```

Then, type your expressions and press Enter:

```
>>> 10 % -3
1
>>> (1+i)^2
2i
>>> 2(1+i)
2+2i
>>> sin(pi/2)cos(0)
1
>>> help sin
(Help text for sin will be displayed)
>>> exit
Exiting ToyCalc.
```

* Type `exit` or `quit` to leave the interactive mode.
* Type `help` or `help [topic]` for assistance.
* Command history is saved in `~/.toycalc_history`.

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

* **Stage 3: Advanced Operators & Combinatorial Functions:**
    * (Deferred) Factorial operator `x!` and `gamma(x)` for general complex numbers (pending robust complex Gamma solution).
    * Combinatorial functions: `nCr(n,k)` and `nPr(n,k)` (implementation strategy pending Gamma decision; may be restricted to integers initially or deferred).
    * Additional less common mathematical constants (e.g., `phi`).
* **Stage 4: Usability & Parser Enhancements:**
    * Improved error reporting (more context, better positioning).
    * More detailed and categorized help system, potentially with search.
    * Advanced REPL features (e.g., tab completion for functions/constants).
* **Stage 5: Advanced Numeric & Expression Features:**
    * Full complex number input parsing (e.g., "3+2.5i", "1.2e-3 - 4.5j").
    * User-defined variables.
    * (Potentially) User-defined functions.
    * (Potential Revisit) Arbitrary-precision numbers (`big.Float`, `BigComplex`).
* **Stage 6: Comprehensive Multi-Value Exploration Engine:**
    * Mechanisms to explore non-principal values for multi-valued complex functions (e.g., `allRoots(base, n)`, `logBranch(z, k)`).
    * Set-based evaluation for combinatorial results.
    * User controls for exploration depth/criteria.

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
