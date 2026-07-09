# Jash Functions Reference

Полный справочник всех функций, методов и типов в проекте Jash.

---

## cmd/jash/main.go — точка входа

| Функция | Описание |
|---|---|
| `main()` | Читает `.jash` файл из аргументов CLI, прогоняет через лексер → парсер → evaluator, выводит результат или ошибки. Usage: `jash <file.jash>` |

---

## pkg/token/token.go — токены и ключевые слова

### Функции

| Функция | Описание |
|---|---|
| `LookupIdent(ident string) TokenType` | Проверяет, является ли строка ключевым словом (`def`, `if`, `for`, `true`, `null` и т.д.). Возвращает соответствующий `TokenType` или `IDENT` |

### Типы токенов

| Константа | Описание |
|---|---|
| `EOF` | Конец файла |
| `NEWLINE` | Перенос строки |
| `INDENT` | Увеличение отступа |
| `DEDENT` | Уменьшение отступа |
| `IDENT` | Идентификатор |
| `NUMBER` | Числовой литерал |
| `STRING` | Строковый литерал |
| `DEF` | Ключевое слово `def` |
| `RETURN` | Ключевое слово `return` |
| `IF` | Ключевое слово `if` |
| `ELSE` | Ключевое слово `else` |
| `ELIF` | Ключевое слово `elif` |
| `FOR` | Ключевое слово `for` |
| `IN` | Ключевое слово `in` |
| `WHILE` | Ключевое слово `while` |
| `REPEAT` | Ключевое слово `repeat` |
| `TRUE` | Ключевое слово `true` |
| `FALSE` | Ключевое слово `false` |
| `NULL` | Ключевое слово `null` |
| `AND` | Ключевое слово `and` |
| `OR` | Ключевое слово `or` |
| `NOT` | Ключевое слово `not` |
| `LBRACE` | `{` |
| `RBRACE` | `}` |
| `LBRACKET` | `[` |
| `RBRACKET` | `]` |
| `LPAREN` | `(` |
| `RPAREN` | `)` |
| `COLON` | `:` |
| `COMMA` | `,` |
| `DOT` | `.` |
| `ASSIGN` | `=` |
| `PLUS` | `+` |
| `MINUS` | `-` |
| `STAR` | `*` |
| `SLASH` | `/` |
| `EQ` | `==` |
| `NEQ` | `!=` |
| `LT` | `<` |
| `GT` | `>` |
| `LTE` | `<=` |
| `GTE` | `>=` |

### Структуры

| Тип | Поля |
|---|---|
| `Token` | `Type TokenType`, `Literal string`, `Line int`, `Column int` |

---

## pkg/lexer/lexer.go — лексический анализатор

### Типы

| Тип | Поля |
|---|---|
| `Lexer` | `input string`, `position int`, `readPosition int`, `ch byte`, `line int`, `column int`, `indentStack []int`, `tokens []Token`, `errors []string` |

### Функции

| Функция | Описание |
|---|---|
| `New(input string) *Lexer` | Создаёт новый лексер |
| `(l *Lexer) Tokenize() ([]Token, []string)` | Разбивает исходный код на токены с учётом отступов |
| `(l *Lexer) readChar()` | Читает следующий символ, обновляет строку/колонку |
| `(l *Lexer) peekChar() byte` | Возвращает следующий символ без сдвига |
| `(l *Lexer) countIndent() int` | Считает количество пробелов (табуляция = 4 пробела) |
| `(l *Lexer) handleIndent(indent int)` | Сравнивает с предыдущим отступом, генерирует INDENT/DEDENT |
| `(l *Lexer) emitToken(t TokenType, literal string)` | Добавляет токен в выходной список |
| `(l *Lexer) emitDedent()` | Добавляет токен DEDENT |
| `(l *Lexer) skipComment()` | Пропускает комментарий от `#` до конца строки |
| `(l *Lexer) readIdentifier()` | Читает идентификатор или ключевое слово |
| `(l *Lexer) readNumber()` | Читает целое число или число с плавающей точкой |
| `(l *Lexer) readString()` | Читает строку в двойных кавычках с escape (`\"`, `\\`, `\n`, `\t`, `\r`) |
| `(l *Lexer) readSingleQuotedString()` | Читает строку в одинарных кавычках с escape |
| `isLetter(ch byte) bool` | Проверяет, является ли символ буквой или `_` |
| `isDigit(ch byte) bool` | Проверяет, является ли символ цифрой |

---

## pkg/ast/ast.go — абстрактное синтаксическое дерево

### Интерфейсы

| Интерфейс | Методы | Описание |
|---|---|---|
| `Node` | `TokenLiteral() string`, `String() string` | Базовый узел AST |
| `Statement` | + `statementNode()` | Узел-инструкция |
| `Expression` | + `expressionNode()` | Узел-выражение |

### Инструкции (Statement)

| Тип | Поля | Описание |
|---|---|---|
| `Program` | `Statements []Statement` | Корневой узел программы |
| `BlockStatement` | `Statements []Statement` | Блок кода с отступами |
| `FunctionStatement` | `Name *Identifier`, `Parameters []*Identifier`, `Body *BlockStatement` | Объявление функции |
| `ReturnStatement` | `Value Expression` | Инструкция `return` |
| `AssignStatement` | `Name *Identifier`, `Value Expression` | Присваивание `x = value` |
| `ExpressionStatement` | `Expression Expression` | Выражение, используемое как инструкция |
| `IfStatement` | `Condition Expression`, `Body *BlockStatement`, `ElseBody *BlockStatement` | Условный оператор `if/elif/else` |
| `ForStatement` | `Variable *Identifier`, `Iterable Expression`, `Body *BlockStatement` | Цикл `for x in iterable` |
| `WhileStatement` | `Condition Expression`, `Body *BlockStatement` | Цикл `while condition` |
| `RepeatStatement` | `Count Expression`, `Body *BlockStatement` | Цикл `repeat(n): block` |

### Выражения (Expression)

| Тип | Поля | Описание |
|---|---|---|
| `Identifier` | `Value string` | Идентификатор (имя переменной) |
| `NumberLiteral` | `Value string` | Числовой литерал |
| `StringLiteral` | `Value string` | Строковый литерал |
| `BooleanLiteral` | `Value bool` | Булев литерал `true`/`false` |
| `NullLiteral` | — | Литерал `null` |
| `JSONObject` | `Pairs map[string]Expression` | JSON-объект `{ "key": value }` |
| `JSONArray` | `Elements []Expression` | JSON-массив `[elem, ...]` |
| `CallExpression` | `Function Expression`, `Arguments []Expression` | Вызов функции |
| `MemberAccess` | `Object Expression`, `Member *Identifier` | Доступ к члену `obj.member` |
| `InfixExpression` | `Left Expression`, `Operator string`, `Right Expression` | Бинарная операция |
| `PrefixExpression` | `Operator string`, `Right Expression` | Унарная операция |

---

## pkg/parser/parser.go — Pratt-парсер

### Уровни приоритета (снизу вверх)

| Константа | Значение | Операторы |
|---|---|---|
| `LOWEST` | 1 | — |
| `ASSIGN` | 2 | `=` |
| `OR` | 3 | `or` |
| `AND` | 4 | `and` |
| `EQUALS` | 5 | `==`, `!=` |
| `COMPARISON` | 6 | `<`, `>`, `<=`, `>=` |
| `SUM` | 7 | `+`, `-` |
| `PRODUCT` | 8 | `*`, `/` |
| `PREFIX` | 9 | `-x`, `not x` |
| `CALL` | 10 | `f(x)` |
| `MEMBER` | 11 | `obj.member` |

### Функции

| Функция | Описание |
|---|---|
| `New(tokens []Token) *Parser` | Создаёт парсер, регистрирует prefix/infix функции |
| `(p *Parser) nextToken()` | Переходит к следующему токену |
| `(p *Parser) curTokenIs(t TokenType) bool` | Проверяет текущий токен |
| `(p *Parser) peekTokenIs(t TokenType) bool` | Проверяет следующий токен |
| `(p *Parser) expect(t TokenType) bool` | Ожидает токен; при несовпадении добавляет ошибку |
| `(p *Parser) Errors() []string` | Возвращает список ошибок парсинга |
| `(p *Parser) registerPrefix(t, fn)` | Регистрирует prefix-функцию |
| `(p *Parser) registerInfix(t, fn)` | Регистрирует infix-функцию |
| `(p *Parser) ParseProgram() *Program` | Парсит всю программу |
| `(p *Parser) parseStatement() Statement` | Диспатчер инструкций |
| `(p *Parser) parseFunctionStatement()` | `def name(params): block` |
| `(p *Parser) parseReturnStatement()` | `return expr` |
| `(p *Parser) parseAssignStatement()` | `ident = expr` |
| `(p *Parser) parseIfStatement()` | `if/elif/else` |
| `(p *Parser) parseForStatement()` | `for var in iterable: block` |
| `(p *Parser) parseWhileStatement()` | `while condition: block` |
| `(p *Parser) parseRepeatStatement()` | `repeat(n): block` |
| `(p *Parser) parseExpressionStatement()` | Выражение как инструкция |
| `(p *Parser) parseBlockBody()` | Парсит тело с INDENT → DEDENT |
| `(p *Parser) parseExpression(precedence) Expression` | Pratt-парсинг выражений |
| `(p *Parser) parseIdentifier()` | Идентификатор |
| `(p *Parser) parseNumberLiteral()` | Число |
| `(p *Parser) parseStringLiteral()` | Строка |
| `(p *Parser) parseBooleanLiteral()` | `true`/`false` |
| `(p *Parser) parseNullLiteral()` | `null` |
| `(p *Parser) parseGroupedExpression()` | `(expr)` |
| `(p *Parser) parseJSONObject()` | `{ "key": value }` |
| `(p *Parser) parseJSONArray()` | `[elem, ...]` |
| `(p *Parser) parsePrefixExpression()` | `-x`, `not x` |
| `(p *Parser) parseInfixExpression(left Expression)` | `x + y`, `x == y` и т.д. |
| `(p *Parser) parseCallExpression(function Expression)` | `f(args)` |
| `(p *Parser) parseMemberAccess(obj Expression)` | `obj.member` |
| `(p *Parser) peekPrecedence() int` | Приоритет следующего токена |
| `(p *Parser) curPrecedence() int` | Приоритет текущего токена |
| `precedenceOf(t TokenType) int` | Приоритет для типа токена |

---

## pkg/evaluator/evaluator.go — интерпретатор

### Типы объектов

| Тип | ObjectType | Описание |
|---|---|---|
| `Integer` | `INTEGER` | Целое число (`int64`) |
| `Float` | `FLOAT` | Число с плав. точкой (`float64`) |
| `String` | `STRING` | Строка |
| `Boolean` | `BOOLEAN` | Булево значение |
| `Null` | `NULL` | `null` |
| `JSONObject` | `JSON_OBJECT` | JSON-объект |
| `JSONArray` | `JSON_ARRAY` | JSON-массив |
| `Function` | `FUNCTION` | Пользовательская функция |
| `Builtin` | `BUILTIN` | Встроенная функция |
| `ReturnValue` | `RETURN` | Return-значение |
| `Error` | `ERROR` | Ошибка выполнения |

### Интерфейс Object

```go
type Object interface {
    Type() ObjectType
    Inspect() string
}
```

### Окружение (Environment)

| Функция/Метод | Описание |
|---|---|
| `NewEnvironment() *Environment` | Создаёт глобальное окружение, загружает builtins |
| `NewEnclosedEnvironment(outer *Environment) *Environment` | Создаёт вложенное окружение (для замыканий) |
| `(e *Environment) Get(name string) (Object, bool)` | Получает значение из окружения (рекурсивно вверх) |
| `(e *Environment) Set(name string, val Object) Object` | Устанавливает значение |
| `(e *Environment) loadBuiltins()` | Загружает встроенные функции |

### Основной интерпретатор

| Функция | Описание |
|---|---|
| `Eval(node Node, env *Environment) Object` | Главная функция — вычисляет любой AST-узел |
| `evalProgram(program, env)` | Выполняет программу (возвращает последнее значение) |
| `evalBlockStatement(block, env)` | Выполняет блок кода |
| `evalFunctionStatement(node, env)` | Создаёт объект функции и сохраняет в окружении |
| `evalReturnStatement(node, env)` | Вычисляет return и оборачивает в ReturnValue |
| `evalAssignStatement(node, env)` | Присваивает значение переменной |
| `evalIfStatement(node, env)` | Выполняет `if/elif/else` |
| `evalForStatement(node, env)` | Итерирует по массиву, строке или объекту |
| `evalWhileStatement(node, env)` | Выполняет цикл с условием |
| `evalRepeatStatement(node, env)` | Выполняет цикл `repeat(n)` — повторяет блок n раз |
| `evalIdentifier(node, env)` | Ищет идентификатор в окружении |
| `evalNumberLiteral(node)` | Парсит число в Integer или Float |
| `evalJSONObject(node, env)` | Вычисляет JSON-объект |
| `evalJSONArray(node, env)` | Вычисляет JSON-массив |
| `evalCallExpression(node, env)` | Вызов функции |
| `evalMemberAccess(node, env)` | Доступ к члену: `obj.key`, `arr.length` |
| `evalInfixExpression(node, env)` | Бинарные операции с авто-приведением типов |
| `evalPrefixExpression(node, env)` | Унарные `-` и `not` |
| `evalIntegerInfixExpression(op, left, right)` | Арифметика/сравнение для int |
| `evalFloatInfixExpression(op, left, right)` | Арифметика/сравнение для float |
| `evalStringInfixExpression(op, left, right)` | Конкатенация/сравнение строк |
| `evalMinusPrefixOperatorExpression(right)` | Унарный минус |
| `evalNotPrefixOperatorExpression(right)` | Логическое `not` |
| `evalExpressions(exps, env) []Object` | Вычисляет список выражений (аргументы вызова) |
| `applyFunction(fn, args) Object` | Вызывает пользовательскую или встроенную функцию |
| `nativeBoolToBooleanObject(b bool) *Boolean` | `true`/`false` → объект Boolean |
| `isTruthy(obj Object) bool` | Проверка истинности |
| `isError(obj Object) bool` | Проверка на ошибку |

### Встроенные функции (Builtins)

| Функция | Сигнатура Jash | Описание |
|---|---|---|
| `printFunc` | `print(...)` | Выводит значения в stdout через пробел |
| `lenFunc` | `len(obj)` | Длина строки, массива или объекта |
| `typeFunc` | `type(obj)` | Имя типа в виде строки |
| `aiPredictFunc` | `ai.predict(text)` | Возвращает мок-результат: `prediction`, `confidence`, `model` |
| `ollamaFunc` | `ai.ollama(url)` | Создаёт клиент для Ollama API |
| `serveFunc` | `serve(port, handler)` | Запускает HTTP-сервер |
| `imageASCIIFunc` | `image.ascii(path)` | Конвертирует изображение (файл или URL) в ASCII-арт и выводит в консоль |
| `timeSleepFunc` | `time.sleep(secs)` | Приостанавливает выполнение на указанное количество секунд |

### Вспомогательные функции HTTP-сервера

| Функция | Описание |
|---|---|
| `readBody(r *http.Request) ([]byte, error)` | Читает тело HTTP-запроса |
| `objectToJSON(obj Object) []byte` | Сериализует Object в JSON |
| `writeJSON(buf *bytes.Buffer, obj Object)` | Рекурсивно записывает Object как JSON |

### Ollama-функции

| Функция | Описание |
|---|---|
| `ollamaFunc(args ...Object) Object` | `ai.ollama(url)` — возвращает объект с методами `generate`, `chat`, `list` |
| `makeOllamaGenerate(baseURL string) func` | `client.generate(model, prompt)` — POST `/api/generate` |
| `makeOllamaChat(baseURL string) func` | `client.chat(model, messages)` — POST `/api/chat` |
| `makeOllamaList(baseURL string) func` | `client.list()` — GET `/api/tags` |
| `ollamaRequest(endpoint string, body) Object` | Выполняет HTTP-запрос к Ollama |
| `goToJashObject(v interface{}) Object` | Конвертирует `interface{}` → Jash Object |
| `jashToGoObject(obj Object) interface{}` | Конвертирует Jash Object → `interface{}` |

---

### Модуль `time`

| Функция | Сигнатура Jash | Описание |
|---|---|---|
| `timeSleepFunc` | `time.sleep(seconds)` | Приостанавливает выполнение на `seconds` (int или float). Пример: `time.sleep(1.5)` |

---

### Модуль `image`

| Функция | Сигнатура Jash | Описание |
|---|---|---|
| `imageASCIIFunc` | `image.ascii(source)` | Читает изображение из файла или URL, конвертирует в ASCII-арт (80 символов в ширину) и возвращает строку с ASCII-графикой |

---

## pkg/evaluator/evaluator.go — time.sleep

| Функция | Описание |
|---|---|
| `timeSleepFunc` | `time.sleep(secs)` — приостанавливает выполнение, принимает Integer, Float или String (парсится в число) |

---

## pkg/evaluator/image.go — ASCII art (image)

| Функция | Описание |
|---|---|
| `imageASCIIFunc` | Берёт локальный файл или URL, декодирует PNG/JPEG/GIF, ресайзит под 80 колонок, преобразует яркость пикселей в символы `@%#*+=-:. ` |

### Константы

| Константа | Значение | Описание |
|---|---|---|
| `asciiChars` | `"@%#*+=-:. "` | Набор символов для градаций яркости (от тёмного к светлому) |

---

## Модульная система: `import`

Синтаксис импорта модуля:

```jash
import module_name
```

Доступные модули:

| Модуль      | Функции                                                                 |
|-------------|-------------------------------------------------------------------------|
| `math`      | `sqrt()`, `abs()`, `floor()`, `ceil()`, `sin()`, `cos()` |
| `random`    | `int(min, max)`, `float()`, `choice(array)` |
| `time`      | `sleep(seconds)`, `now()`, `format(layout)` |
| `file`      | `read(path)`, `write(path, content)` |
| `ai`        | `predict(text)` — отправляет текст AI-модели и возвращает ответ |
| `image`     | `ascii(path)` — конвертирует изображение в ASCII-арт |
| `jash_ui`   | `window(title, width?, height?)` — создаёт GUI-окно |

**Встроенные функции** (доступны без импорта): `print()`, `len()`, `type()`, `say()`, `any()`, `serve()`

---

## pkg/evaluator/jit.go — JIT-менеджер

| Тип | Поля | Описание |
|---|---|---|
| `CompiledFunc` | `instructions []Instruction`, `constants []interface{}`, `varNames map[int]string` | Скомпилированная функция |
| `JITManager` | `compiled map[string]*CompiledFunc` | Менеджер JIT-компиляции |

| Функция/Метод | Описание |
|---|---|
| `InitJIT()` | Создаёт глобальный `GlobalJIT` |
| `(jm *JITManager) IsCompiled(name string) bool` | Проверяет, скомпилирована ли функция |
| `(jm *JITManager) Compile(name, params, body) bool` | Компилирует функцию в байткод |
| `(jm *JITManager) Execute(name, args, env) Object` | Выполняет скомпилированную функцию |

---

## pkg/evaluator/jit_opcode.go — JIT opcodes

| Константа | Код | Описание |
|---|---|---|
| `OpConstant` | 0 | Загрузить константу |
| `OpLoad` | 1 | Загрузить переменную |
| `OpStore` | 2 | Сохранить в переменную |
| `OpAdd` / `OpSub` / `OpMul` / `OpDiv` | 3–6 | Арифметика |
| `OpEq` / `OpNeq` / `OpLt` / `OpGt` / `OpLte` / `OpGte` | 7–12 | Сравнение |
| `OpAnd` / `OpOr` / `OpNot` | 13–15 | Логические операции |
| `OpMinus` | 16 | Унарный минус |
| `OpCall` | 17 | Вызов функции |
| `OpReturn` / `OpReturnVal` | 18–19 | Возврат |
| `OpJump` / `OpJumpIfFalse` | 20–21 | Переходы |
| `OpPop` | 22 | Вытолкнуть со стека |
| `OpNewArray` / `OpNewObject` | 23–24 | Создать массив/объект |
| `OpSetMember` / `OpGetMember` | 25–26 | Доступ к элементу |
| `OpNull` / `OpTrue` / `OpFalse` | 27–29 | Литералы |
| `OpDefFunc` | 30 | Объявление функции |

| Тип | Поля | Описание |
|---|---|---|
| `Instruction` | `Op Opcode`, `Arg int`, `Arg2 int`, `Const interface{}` | Инструкция JIT |
| `funcDef` | `Name string`, `Params []*ast.Identifier`, `Body *ast.BlockStatement` | Определение функции для JIT |

---

## pkg/evaluator/jit_compiler.go — JIT-компилятор

| Тип | Поля | Описание |
|---|---|---|
| `compiler` | `instructions []Instruction`, `constants []interface{}`, `varIndex map[string]int`, `varNames map[int]string`, `varCount int` | Компилятор JIT |

| Функция/Метод | Описание |
|---|---|
| `newCompiler() *compiler` | Создаёт новый компилятор |
| `(c *compiler) Compile(node) ([]Instruction, []interface{}, map[int]string, error)` | Компилирует AST в байткод |
| `(c *compiler) compileBlock(block)` | Компилирует блок |
| `(c *compiler) compileStatement(stmt)` | Компилирует инструкцию |
| `(c *compiler) compileIf(stmt)` | Компилирует if/else |
| `(c *compiler) compileFor(stmt)` | Компилирует for |
| `(c *compiler) compileWhile(stmt)` | Компилирует while |
| `(c *compiler) compileNode(node)` | Компилирует выражение |
| `(c *compiler) emit(op, args...)` | Добавляет инструкцию |
| `(c *compiler) addConstant(val) int` | Добавляет константу |
| `(c *compiler) getVarIndex(name) int` | Получает/создаёт индекс переменной |
| `(c *compiler) infixOp(op) Opcode` | Оператор → опкод |

---

## pkg/evaluator/jit_vm.go — JIT-виртуальная машина

| Тип | Поля | Описание |
|---|---|---|
| `vm` | `instructions`, `constants`, `varNames`, `ip`, `stack`, `sp`, `globals`, `env` | Виртуальная машина |

| Функция/Метод | Описание |
|---|---|
| `newVM(instructions, constants, varNames, globals, env) *vm` | Создаёт VM |
| `(vm *vm) Run() Object` | Запускает выполнение байткода |
| `(vm *vm) push(obj Object)` | Кладёт на стек |
| `(vm *vm) pop() Object` | Снимает со стека |
| `(vm *vm) toJashObject(v interface{}) Object` | Конвертирует `interface{}` → Jash Object |
| `evalArith(op, l, r Object) Object` | Арифметика для VM |
| `eq(a, b Object) bool` | Сравнение на равенство |
| `cmp(a, b Object) int` | Сравнение (<0, 0, >0) |

---

## pkg/evaluator/ui.go — GUI (jash_ui)

| Функция | Сигнатура Jash | Описание |
|---|---|---|
| `uiWindowFunc` | `jash_ui.window(title, w, h)` | Создаёт окно, возвращает объект с методами управления |
| `uiMakeAdder(winID, widgetType)` | `win.add_label(text)` / `add_entry(text)` / `add_text(text)` | Добавляет label / entry / text-area |
| `uiMakePhoto(winID)` | `win.add_photo(src, w?, h?)` | Добавляет изображение (URL или файл) |
| `uiMakeButton(winID)` | `win.add_button(text, fn)` | Добавляет кнопку с функцией-колбэком |
| `uiMakeGetValue(winID)` | `win.get_value(widgetID)` | Получает текущее значение поля ввода |
| `uiMakeRun(winID)` | `win.run()` | Открывает окно в браузере, блокирует до закрытия |
| `uiMakeClose(winID)` | `win.close()` | Закрывает окно программно |
| `findAvailablePort() int` | — | Находит свободный TCP-порт |
| `openBrowser(url string)` | — | Открывает URL в браузере (Windows/macOS/Linux) |
| `generateUIHTML(win) string` | — | Генерирует HTML+CSS+JS для окна |
| `isImageURL(s string) bool` | — | Проверяет, является ли строка URL изображения (по расширению) |
| `esc(s string) string` | — | Экранирует HTML-спецсимволы |

### Типы GUI

| Тип | Поля | Описание |
|---|---|---|
| `UIWidget` | `ID string`, `Type string`, `Text string`, `OnClick *Function` | Виджет (label, button, entry, text-area) |
| `UIWindowState` | `ID`, `Title`, `Width`, `Height`, `Widgets`, `Values`, `server`, `closeCh` | Состояние окна |

---

## debug_lexer/main.go — отладка лексера

| Функция | Описание |
|---|---|
| `main()` | Запускает лексер на тестовой строке `def foo()\\n    print("hi")` и выводит все токены в stderr |---

## jashtoexe/main.go — билдер Jash → .exe

| Функция | Описание |
|---|---|
| `main()` | Читает `.jash`, кодирует в base64, генерирует Go-программу со встроенным интерпретатором, компилирует в `.exe` |

Процесс:
1. Читает входной `.jash` файл
2. Кодирует скрипт в base64
3. Создаёт временную директорию с `go.mod` и `main.go`
4. `main.go` содержит полный интерпретатор Jash + закодированный скрипт
5. Запускает `go build -o output.exe -ldflags "-s -w"`
6. Очищает временные файлы
