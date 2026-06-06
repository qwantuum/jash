<p align="center">
  <img src="assets/logo.png" alt="Jash Logo" width="250">
</p>


# Jash Language

**JSON + Bash/Python** — A clean, indentation-based scripting language with first-class JSON support, built-in HTTP serving, and mock AI predictions.

```jash
def handle_request(req)
    data = { "status": "success", "model": "Jash-AI" }
    return data

serve(8080, handle_request)
```

---

## Features

- **Python-like syntax** — Whitespace-sensitive indentation, no semicolons, no curly braces for blocks
- **Native JSON literals** — Write `{ "key": value }` and `[1, 2, 3]` directly in code
- **Built-in HTTP server** — `serve(port, handler)` starts a production-ready HTTP server
- **Mock AI predictions** — `ai.predict(input)` returns structured prediction results
- **Type inference** — Integers, floats, strings, booleans, null, JSON objects, and arrays
- **Expressions & operators** — Full arithmetic, comparison, and logical operators
- **Control flow** — `if/else`, `for`, `while` with indentation-based blocks
- **Functions** — Define reusable logic with `def name(params):`
- **Zero external dependencies** — Built entirely on the Go standard library

---

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/qwantuum/jash.git
cd jash

# Build the interpreter
go build -o jash ./cmd/jash

# Run a Jash script
./jash examples/hello.jash
```

### Hello, World

```jash
print("Hello, Jash!")
```

Save as `hello.jash` and run:

```bash
jash hello.jash
```

---

## Language Guide

### Variables

```jash
name = "Jash"
version = 1.0
active = true
count = 42
```

### JSON Literals

```jash
user = {
    "name": "Alice",
    "age": 30,
    "skills": ["Go", "Python", "Jash"],
    "active": true,
    "address": null
}
```

### Functions

```jash
def greet(name)
    return "Hello, " + name + "!"

print(greet("Jash"))
```

### Control Flow

```jash
x = 10

if x > 5
    print("x is greater than 5")
else
    print("x is 5 or less")

for i in [1, 2, 3]
    print(i)

while x > 0
    print(x)
    x = x - 1
```

### HTTP Server

```jash
def handler(req)
    return {
        "status": "ok",
        "method": req.method,
        "path": req.path,
        "message": "Hello from Jash!"
    }

serve(3000, handler)
```

### AI Predictions

```jash
def analyze(text)
    result = ai.predict(text)
    return {
        "text": text,
        "prediction": result.prediction,
        "confidence": result.confidence,
        "model": result.model
    }

data = analyze("Great product!")
print(data)
```

### Built-in Functions

| Function   | Description                              |
|------------|------------------------------------------|
| `print()`  | Prints values to stdout                  |
| `len()`    | Returns length of string, array, or object |
| `type()`   | Returns the type name of a value         |
| `serve()`  | Starts an HTTP server on a given port    |
| `ai.predict()` | Returns a mock AI prediction result  |

---

## Examples

### JSON API Server

```jash
def get_user(id)
    return {
        "id": id,
        "name": "User " + id,
        "email": "user" + id + "@example.com"
    }

def handler(req)
    if req.path == "/users"
        return {
            "users": [
                get_user("1"),
                get_user("2"),
                get_user("3")
            ]
        }

    return { "error": "not found" }

serve(8080, handler)
```

### Data Analysis Mock

```jash
def analyze(item)
    pred = ai.predict(item)
    return {
        "input": item,
        "result": pred.prediction,
        "score": pred.confidence
    }

print(analyze("Great product!"))
```

---

## Building from Source

```bash
git clone https://github.com/qwantuum/jash.git
cd jash
go build -o jash ./cmd/jash
```

The binary requires no dependencies beyond the Go standard library.

---

## Project Structure

```
jash/
  go.mod
  README.md
  cmd/jash/main.go          # Entry point
  pkg/
    token/token.go           # Token definitions
    lexer/lexer.go           # Lexer with indentation tracking
    ast/ast.go               # Abstract Syntax Tree nodes
    parser/parser.go         # Pratt parser
    evaluator/evaluator.go   # Tree-walking interpreter
```

---

## License

MIT
