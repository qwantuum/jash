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

## 🔢 2. Advanced Mathematical Engine (PyPy3-style)

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

## 🤖 4. Built-in AI Core Module (`ai`)

Jash is an **AI-native language**, meaning AI interaction capability is baked right into the core global runtime.

### `ai.set_key(key?)`
Initializes security credentials for the model. 
* If a string is passed: `ai.set_key("sk-...")` uses that explicit key.
* If left empty: `ai.set_key()` programmatically inspects Go's `os.Getenv("OPENAI_API_KEY")` to securely fetch it from your Windows system environment. If empty, it drops a clean `[Jash Error]: API key not found` warning.

### `ai.predict(prompt)`
Sends a raw text prompt directly to the built-in AI interface and returns a clean execution-ready string payload.
```jash
ai.set_key()
slogan = ai.predict("Give me a slogan for Jash language")
print(slogan)
```

---

## 🌐 5. Web Backend Engine (`http`)

Jash exposes underlying high-performance native Go network routines (`net/http`) through ultra-simplified bindings.

### `http.json(status, object)`
A high-level utility function that serializes a Jash/Python native object into a minified, valid JSON string and applies the standard `application/json` Content-Type header header automatically.

### `http.listen(port, handler_function)`
Spins up a multi-threaded, persistent HTTP production-ready backend web server running locally on the specified port.

```jash
def api_handler()
    payload = { "framework": "Jash-Core", "status": "online" }
    return http.json(200, payload)

# Starts server on http://localhost:8080
http.listen(8080, api_handler)
```
