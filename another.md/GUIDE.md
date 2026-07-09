# 📘 Jash Language: Full Documentation & Feature Review

Welcome to the official feature guide for **Jash** — a high-performance, Python-inspired programming language built on top of **Go**. It is designed specifically for rapid development of AI integrations, web backends, and lightweight data structures.

---

## 🔑 1. Core Syntax & Variables

Jash uses **whitespace indentation** for code blocks (no curly braces `{}`) and strictly avoids semicolons `;`. 

### Variables & Data Types
Variables are dynamically typed and can hold Integers, Strings, Booleans, or Native JSON Objects.
```jash
username = "qwantuum"
version = 1.0
is_active = true
```

### Functions (`def`)
Defined using the `def` keyword followed by the function name, arguments, and a colon `:`. Blocks must be indented.
```jash
def greet(name)
    print("Welcome to Jash, ", name)
    return true
```

---

## 🔢 2. Loops: `while`, `for`, `repeat`

**`while`** — repeats while condition is truthy:
```jash
x = 3
while x > 0
    print(x)
    x = x - 1
```

**`for ... in`** — iterates over an array, string, or object:
```jash
for i in [1, 2, 3]
    print(i)
```

**`repeat(n)`** — repeats a block exactly `n` times:
```jash
repeat(5)
    print("Hello!")
```

---

## 🔢 3. Advanced Mathematical Engine (PyPy3-style)

Jash supports basic infix operations (`+`, `-`, `*`, `/`) with standard mathematical operator precedence (multiplication and division happen before addition and subtraction).

### ⚡ Compiler Optimization: Constant Folding
Jash implements a smart compiler optimization inspired by PyPy3. If an expression consists entirely of fixed numbers, Jash calculates the result **at compile-time** rather than dragging it to runtime.

```jash
# Slow in standard interpreters, but Jash compiles this instantly into 'result = 25'
result = 5 + 10 * 2 
```

---

## 📦 3. Native JSON Support ("First-Class Citizen")

Unlike Go or C++, where JSON parsing requires external packages and complex map structures, **Jash treats JSON objects natively**, just like standard Python dictionaries.

```jash
# Direct JSON instantiation without extra quotes or strings
my_data = {
    "status": "success",
    "code": 200,
    "meta": {
        "author": "qwantuum",
        "age": 11
    }
}
```

---

## 📥 4. Module Import System

Jash has a modular architecture. Core builtins (`print()`, `len()`, `type()`, `say()`, `any()`, `serve()`) are always available, but additional functionality is loaded via the `import` statement:

```jash
import math
import random
import time
import file
import ai
import image
import jash_ui
```

Simply write `import module_name` at the top of your script to enable that module's features.

### Available Modules

| Module      | Contents                                                                 |
|-------------|--------------------------------------------------------------------------|
| `math`      | `sqrt()`, `abs()`, `floor()`, `ceil()`, `sin()`, `cos()`                |
| `random`    | `int(min, max)`, `float()`, `choice(array)`                              |
| `time`      | `sleep(seconds)`, `now()`, `format(layout)`                              |
| `file`      | `read(path)`, `write(path, content)`                                      |
| `ai`        | `predict(text)` — sends text to an AI model and returns the response      |
| `image`     | `ascii(path)` — converts an image to ASCII art                            |
| `jash_ui`   | `window(title, width?, height?)` — creates a GUI window                   |

---

## 🤖 5. Built-in AI Core Module (`ai`)

Jash is an **AI-native language**, meaning AI interaction capability is baked right into the core global runtime.

### `ai.predict(text)`
Sends a text to the built-in mock AI interface and returns a structured result with `prediction`, `confidence`, and `model` fields.
```jash
result = ai.predict("Great product!")
print(result.prediction)  # "positive"
print(result.confidence)  # 0.9532
```

### `ai.ollama(url)`
Creates an Ollama client connected to a local or remote Ollama instance:
```jash
client = ai.ollama("http://localhost:11434")
result = client.generate("llama2", "Hello!")
print(result.response)
```

---

## 🌐 6. Web Backend (`serve`)

Jash exposes Go's `net/http` through a simple `serve(port, handler)` function:

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

The handler receives an object with `method`, `path`, `body`, and `query` fields. The return value is serialized to JSON automatically.

---

## 🖼️ 7. ASCII Art from Images (`image`)

Convert any image (local file or URL) to ASCII art in the console:

```jash
art = image.ascii("logo.png")
print(art)

art = image.ascii("https://example.com/photo.jpg")
print(art)
```

Uses PNG/JPEG/GIF decoding, resizes to 80 columns, and maps brightness to `@%#*+=-:. ` characters.

---

## 🪟 8. GUI Windows (`jash_ui`)

Create browser-based GUI windows with the `jash_ui` module:

```jash
def on_click(vals)
    print("Button clicked!")

win = jash_ui.window("My App", 500, 400)
win.add_label("Welcome to Jash!")
win.add_button("Click me", on_click)
win.add_entry("Enter your name")
win.add_photo("logo.png")
win.run()
```

Methods available on the window object:
- `add_label(text)` — static text
- `add_button(text, callback)` — clickable button
- `add_entry(text)` — single-line input (returns widget ID)
- `add_text(text)` — multi-line text area (returns widget ID)
- `add_photo(src, width?, height?)` — image from file or URL
- `get_value(widgetID)` — get current value of an entry or text-area
- `run()` — open the window in a browser and block until closed
- `close()` — close the window programmatically
