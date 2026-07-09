# Learn Jash in 15 Minutes

Jash — Python-подобный язык со встроенным AI, JSON, GUI и JIT-компиляцией.

## 1. Hello, World

```jash
print("Hello, World!")
```

Запуск: `jash script.jash` или `jash` (REPL, пустая строка — выполнить).

## 2. Переменные и типы

```jash
name = "Jash"          # строка
version = 1.0           # число (float)
year = 2026             # целое
is_cool = true          # булево
nothing = null          # null
```

## 3. Арифметика

```jash
print(2 + 2)            # 4
print(10 - 3)           # 7
print(4 * 5)            # 20
print(10 / 3)           # 3 (целое) или 3.333 (float)
print(10 % 3)           # 1 (modulo)
```

Составное присваивание:

```jash
x = 10
x += 5                  # x = 15
x -= 3                  # x = 12
x *= 2                  # x = 24
x /= 4                  # x = 6
```

## 4. Строки

```jash
greeting = "Hello"
name = "Jash"
print(greeting + " " + name)   # Hello Jash
print(len(greeting))           # 5
print(greeting[0])             # H
print(greeting[-1])            # o
```

## 5. Условия

```jash
score = 85

if score >= 90
    print("A")
elif score >= 80
    print("B")
else
    print("C")
```

Логические операторы: `and`, `or`, `not`.

## 6. Циклы

### for — по массиву, строке или объекту:

```jash
for x in [1, 2, 3]
    print(x)

for ch in "hello"
    print(ch)

for key in {"a": 1, "b": 2}
    print(key)
```

### while:

```jash
x = 3
while x > 0
    print(x)
    x = x - 1
```

### repeat:

```jash
repeat(3)
    print("Go!")
```

### break / continue:

```jash
for x in [1,2,3,4,5]
    if x == 3
        break          # выйти из цикла
    if x == 2
        continue       # пропустить итерацию
    print(x)
```

## 7. Функции

```jash
def add(a, b)
    return a + b

result = add(3, 4)
print(result)           # 7
```

## 8. JSON (родной)

```jash
user = {
    "name": "Alice",
    "age": 30,
    "skills": ["Go", "Jash"],
    "meta": {
        "level": "pro"
    }
}

print(user.name)        # Alice
print(user.age)         # 30
print(user.skills[0])   # Go
print(user.meta.level)  # pro
```

## 9. Массивы

```jash
arr = [10, 20, 30]
print(arr[0])           # 10
print(arr[-1])          # 30 (отрицательные индексы)
print(len(arr))         # 3
print(arr.length)       # 3
```

## 10. Модули

```jash
import math
print(math.sqrt(16))    # 4
print(math.abs(-5))     # 5
print(math.floor(3.7))  # 3
print(math.ceil(3.2))   # 4

import random
print(random.int(1, 100))    # случайное целое
print(random.float())        # случайное float 0..1
print(random.choice([1,2,3]))

import time
print(time.now())            # 2026-07-09 15:26:14
time.sleep(1)                # пауза 1 секунда
```

## 11. AI

```jash
import ai

result = ai.predict("What is the meaning of life?")
print(result)
```

`ai.predict()` сам находит самую слабую модель через `ollama list` и шлёт запрос на `localhost:11434`.

Настройка: `OLLAMA_HOST=http://localhost:11434` (по умолчанию).

Для ручного управления:

```jash
client = ai.ollama("http://localhost:11434")
print(client.list())
print(client.generate("llama3.2", "Hello!"))
print(client.chat("llama3.2", [
    {"role": "user", "content": "Hi!"}
]))
```

## 12. Файлы

```jash
import file
file.write("test.txt", "Hello, Jash!")
content = file.read("test.txt")
print(content)
```

## 13. Изображения → ASCII

```jash
import image
art = image.ascii("logo.png")
print(art)
```

## 14. Веб-сервер

```jash
def handler(req)
    return {
        "status": "ok",
        "message": "Hello from Jash!"
    }

serve(3000, handler)
```

`req` содержит: `method`, `path`, `body`, `query`.

## 15. GUI

```jash
import jash_ui

def on_click(vals)
    print("Clicked!")

win = jash_ui.window("My App", 500, 400)
win.add_label("Welcome!")
win.add_button("Click me", on_click)
win.run()
```

## 16. Индексация

```jash
arr = [10, 20, 30]
print(arr[1])            # 20
print(arr[-1])           # 30

str = "hello"
print(str[0])            # h
print(str[-1])           # o
```

## Быстрая шпаргалка

| Конструкция | Синтаксис |
|---|---|
| Переменная | `x = 1` |
| Функция | `def name(a, b)` / `return val` |
| Условие | `if / elif / else` |
| Цикл for | `for x in iterable` |
| Цикл while | `while cond` |
| Цикл repeat | `repeat(n)` |
| Импорт | `import module` |
| JSON | `{"key": val}` |
| Массив | `[1, 2, 3]` |
| Индекс | `arr[0]`, `str[-1]` |
| Составное = | `+= -= *= /=` |
| Модуло | `%` |
| break/continue | `break` / `continue` |
| and/or/not | `and`, `or`, `not` |
| null | `null` |
