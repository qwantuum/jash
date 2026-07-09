# Jash

<p align="center">
  <img src="assets/logo.png" alt="Jash Logo" width="250">
</p>

## 🕵️‍♂️ Fun Fact & Trivia

Originally, the name **Jash** stood for **JSON + Hash**. However, when the AI was helping to write the first version of this README, it hallucinated and assumed **Jash** meant **JSON + Bash/Python** because of the `sh` at the end. 

The idea sounded so cool that the creator decided to keep the AI's version and actually implemented the system automation features. So yes, this language was partially co-authored by an AI's mistake!

**JSON + Bash/Python** — A clean, indentation-based scripting language with first-class JSON support, built-in HTTP serving, and mock AI predictions.

```jash
def handle_request(req)
    data = { "status": "success", "model": "Jash-AI" }
    return data

serve(8080, handle_request)


```

---

## Features

- **Module system** — `import` statement loads modules: `math`, `random`, `time`, `file`, `ai`, `image`, `jash_ui`
- **Python-like syntax** — Whitespace-sensitive indentation, no semicolons, no curly braces for blocks
- **Native JSON literals** — Write `{ "key": value }` and `[1, 2, 3]` directly in code
- **Built-in HTTP server** — `serve(port, handler)` starts a production-ready HTTP server
- **GUI windows** — `jash_ui.window(title, w, h)` creates browser-based windows with labels, buttons, inputs, and photos
- **ASCII art from images** — `image.ascii(path)` converts any image (file or URL) to console ASCII art
- **Mock AI predictions** — `ai.predict(input)` returns structured prediction results
- **Ollama integration** — `ai.ollama(url)` connects to local or remote LLM instances
- **JIT compilation** — Hot-path functions are JIT-compiled to bytecode for faster execution
- **Type inference** — Integers, floats, strings, booleans, null, JSON objects, and arrays
- **Expressions & operators** — Full arithmetic, comparison, and logical operators
- **Control flow** — `if/else`, `for`, `while`, `repeat(n)` with indentation-based blocks
- **Functions** — Define reusable logic with `def name(params):`
- **Zero external dependencies** — Built entirely on the Go standard library

---

## Quick Start

~~При первом запуске возможна задержка из-за проверки Windows Defender / SmartScreen. Это нормально для неподписанных Go-бинарников~~

### Installation

```bash
# Clone the repository
git clone https://github.com/qwantuum/jash.git
cd jash

# Build the interpreter
go build -o jash ./cmd/jash

# Run a Jash script
./jash script.jash
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

repeat(3)
    print("this runs 3 times")
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

### GUI Windows

Create browser-based windows with labels, buttons, text inputs and photos:

```jash
def on_click(vals)
    print("Button clicked!")

win = jash_ui.window("My App", 500, 400)
win.add_label("Hello, Jash!")
win.add_button("Click me", on_click)
win.add_entry("Enter name")
win.add_photo("logo.png")
win.run()
```

### ASCII Art from Images

Turn any image into console ASCII art:

```jash
# From a local file
art = image.ascii("logo.png")
print(art)

# From a URL
art = image.ascii("https://example.com/photo.jpg")
print(art)
```

### Sleep / Delay

Pause execution with `time.sleep(seconds)`:

```jash
print("waiting...")
time.sleep(1.5)
print("done waiting")
```

### Ollama Integration

Connect to a local or remote [Ollama](https://ollama.ai) instance for LLM inference:

```jash
client = ai.ollama("http://localhost:11434")

# Generate a completion
result = client.generate("llama2", "Hello!")
print(result.response)

# Chat with messages
messages = [
    {"role": "user", "content": "Hi!"}
]
reply = client.chat("llama2", messages)
print(reply.message.content)

# List available models
models = client.list()
print(models)
```

### Built-in Functions

| Function   | Description                              |
|------------|------------------------------------------|
| `print()`  | Prints values to stdout                  |
| `len()`    | Returns length of string, array, or object |
| `type()`   | Returns the type name of a value         |
| `serve()`  | Starts an HTTP server on a given port    |
| `ai.predict()` | Returns a mock AI prediction result  |
| `ai.ollama()` | Creates an Ollama client for LLM inference |
| `image.ascii()` | Converts an image (file path or URL) to ASCII art and returns it |
| `jash_ui.window()` | Creates a GUI window with labels, buttons, entries, text-areas and photos |
| `time.sleep()` | Pauses execution for the given number of seconds (integer or float) |

### Importing Modules

Additional functionality is available through modules loaded with `import`:

```jash
import math
import random
import time
import file
import ai
import image
import jash_ui
```

| Module      | Provides                                                                 |
|-------------|--------------------------------------------------------------------------|
| `math`      | `sqrt()`, `abs()`, `floor()`, `ceil()`, `sin()`, `cos()`                |
| `random`    | `int(min, max)`, `float()`, `choice(array)`                              |
| `time`      | `sleep(seconds)`, `now()`, `format(layout)`                              |
| `file`      | `read(path)`, `write(path, content)`                                      |
| `ai`        | `predict(text)` — sends text to an AI model and returns the response      |
| `image`     | `ascii(path)` — converts an image to ASCII art                            |
| `jash_ui`   | `window(title, width?, height?)` — creates a GUI window                   |

Core builtins (`print()`, `len()`, `type()`, `say()`, `any()`, `serve()`) are always available without any import.

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
  cmd/jash/main.go           # Entry point
  pkg/
    token/token.go            # Token definitions
    lexer/lexer.go            # Lexer with indentation tracking
    ast/ast.go                # Abstract Syntax Tree nodes
    parser/parser.go          # Pratt parser
    evaluator/
      evaluator.go            # Tree-walking interpreter + builtins
      image.go                # image.ascii() — ASCII art from images
      ui.go                   # jash_ui.window() — GUI windows
      jit.go                  # JIT manager
      jit_opcode.go           # JIT opcodes
      jit_compiler.go         # JIT compiler
      jit_vm.go               # JIT virtual machine
```

---

## License

MIT
