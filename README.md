# ascii-art

A command line tool written in Go that turns plain text into large, stylized ASCII banner art rendered entirely in your terminal.

```
    /\      /\      /\   
   /  \    /  \    /  \  
  / /\ \  / /\ \  / /\ \ 
 / ____ \/ ____ \/ ____ \
/_/    \_/_/    \_/_/    \_\
```

---

## What it does

You give it a word (or a few), it gives you back that text printed big character by character using a prebuilt font stored in a plain text file (`standard.txt`). It also handles multi-line input using `\n` as a line break.

---

## How to run it

Make sure you have Go installed, then:

```bash
go run . "Hello"
```

Print multiple lines:

```bash
go run . "Hello\nWorld"
```

That `\n` in the argument is treated as a line separator so each word gets its own 8 row block.

---

## How it works

The font lives in `standard.txt`. Every printable ASCII character (from space `32` to tilde `126`) has a slot in that file exactly 9 lines per character. When you pass in a string, the program:

1. Reads and splits `standard.txt` into lines
2. Splits your input on `\n` to get words
3. For each word, loops through 8 rows and prints each character's glyph side by side
4. Skips any character outside the printable ASCII range

The position of each character's glyph in the file is calculated as:

```
start = (charCode - 32) * 9
```

Then it reads `lines[start]` through `lines[start + 7]` for the 8 visual rows.

---

## Project structure

```
.
├── main.go          # All the logic lives here
├── standard.txt     # The font  one glyph per ASCII character, 9 lines each
└── ascii_test.go    # Tests covering output structure, widths, and edge cases
```

---

## Running the tests

```bash
go test -v
```

The test suite checks:

- `standard.txt` exists and has at least 855 lines (95 chars × 9 lines)
- Single input → exactly 8 output rows
- `\n` splits produce the right number of row blocks
- Empty words between `\n\n` print a blank line
- Row widths are consistent across all 8 rows per word
- Character glyph widths match what's stored in `standard.txt`
- The same input always produces the same output

---

## Edge cases handled

- **Empty word** between newlines  prints a blank line and moves on
- **Out of range characters** (below 32 or above 126)  silently skipped
- **Windows line endings** (`\r\n`)  normalized to `\n` before processing

---

## Usage error

If you forget to pass an argument (or pass too many), the program tells you:

```bash
go run . 
# Usage: go run . "text"
```

---

## Requirements

- Go 1.18 or later
- `standard.txt` in the same directory as `main.go`

---

## Author

Built by Solomon Adanu a self taught software engineer who enjoys the satisfying overlap between code and visual design.
