# Jash

<p align="center">
  <img src="assets/logo.png" alt="Jash Logo" width="250">
</p>


Язык программирования **Jash** — чистая, инденто-зависимая (как Python) скриптовая платформа со встроенной поддержкой JSON, HTTP-сервера, GUI-окон и AI-интеграций.

## Возможности

- **Модульная система** — `import module_name` для загрузки модулей: `math`, `random`, `time`, `file`, `ai`, `image`, `jash_ui`
- **Синтаксис как в Python** — отступы вместо скобок, нет точек с запятой
- **Нативные JSON-литералы** — `{ "key": value }` и `[1, 2, 3]` прямо в коде
- **Встроенный HTTP-сервер** — `serve(port, handler)`
- **GUI-окна** — `jash_ui.window(title, width?, height?)` — окна в браузере
- **ASCII-арт из изображений** — `image.ascii(path)`
- **AI-модуль** — `ai.predict(text)` отправляет текст AI-модели
- **JIT-компиляция** — горячие функции компилируются в байткод
- **Управляющие конструкции** — `if/else`, `for`, `while`, `repeat(n)`
- **Функции** — `def name(params):`

## Импорт модулей

Дополнительные возможности подключаются через `import`:

```jash
import math
import random
import time
import file
import ai
import image
import jash_ui
```

| Модуль      | Функции                                                                 |
|-------------|-------------------------------------------------------------------------|
| `math`      | `sqrt()`, `abs()`, `floor()`, `ceil()`, `sin()`, `cos()`                |
| `random`    | `int(min, max)`, `float()`, `choice(array)`                              |
| `time`      | `sleep(seconds)`, `now()`, `format(layout)`                              |
| `file`      | `read(path)`, `write(path, content)`                                      |
| `ai`        | `predict(text)` — отправляет текст AI-модели и возвращает ответ           |
| `image`     | `ascii(path)` — конвертирует изображение в ASCII-арт                      |
| `jash_ui`   | `window(title, width?, height?)` — создаёт GUI-окно                       |

**Встроенные функции** (доступны без импорта): `print()`, `len()`, `type()`, `say()`, `any()`, `serve()`

## Пример

```jash
import math
import ai

print(math.sqrt(16))       # 4
result = ai.predict("Привет!")
print(result)
```

## Быстрый старт

```bash
git clone https://github.com/qwantuum/jash.git
cd jash
go build -o jash ./cmd/jash
./jash examples/hello.jash
```
