# Blood Pressure Bot - Agent Guidelines

## Project Overview
Telegram bot for blood pressure/pulse tracking. Go 1.25.3, PostgreSQL, Echo Charts, Excelize.

## Build & Test Commands

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run a single test file
go test -v ./bot/handlerLog/service_test.go

# Run a single test function
go test -v -run TestLogService_ComputePressureMedian ./bot/handlerLog/

# Run tests matching pattern
go test -v -run "TestLog" ./...

# Build
go build -o bloodpressure .

# Format code
go fmt ./...

# Vet code
go vet ./...
```

## Code Style Guidelines

### Formatting
- Use 4 spaces for indentation (no tabs)
- Group imports: stdlib first, then third-party, then internal packages
- Use `go fmt` before committing

### Naming Conventions
- `PascalCase` for exported types, functions, methods
- `camelCase` for unexported variables and parameters
- Package names: short, lowercase, no underscores
- Interface names: `er` suffix for services (`LogService`, `UserService`)
- Test files: `*_test.go` suffix

### Error Handling
- Return errors as last return value: `func() (result, err error)`
- Use named return parameters for clarity: `func() (result Type, err error) {...}`
- Wrap errors with context: `fmt.Errorf("operation: %w", err)`
- Use `if err != nil` blocks, avoid `else` after error returns
- Handle logged errors inline, don't silently ignore

### Database Patterns
- SQL queries as raw strings with `$1, $2` parameter placeholders
- Close rows with `defer rows.Close()` after `pg.Query()`
- Use `sql.ErrNoRows` to check for missing records
- Transactions for multi-step operations

### Telegram Bot Patterns
- Handlers follow signature: `func(bot *tgbotapi.BotAPI, update tgbotapi.Update)`
- Use alias for telegram-bot-api: `tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"`
- Use `getLogger("HandlerName")` for consistent logging
- Always send user-facing errors as Telegram messages

### Struct Design
- Use pointer receivers for methods that may modify state or need nil checks
- Export structs that cross package boundaries
- JSON tags use camelCase with `omitempty`

### Testing
- Test functions: `TestFunctionName`
- Use table-driven tests with anonymous structs
- Test file lives in same package as code being tested
- Mock database dependencies via interfaces
- Test pure functions first (sorting, calculations, string manipulation)

### Project Structure
```
m/
├── main.go
├── pgsql/           # PostgreSQL connection
├── bot/
│   ├── bot.go       # Main bot setup
│   ├── Config.go
│   ├── core/        # Shared utilities
│   │   ├── RegExpParams.go
│   │   ├── GetUserName.go
│   │   ├── UserService.go
│   │   └── types.go
│   ├── handlerLog/  # Blood pressure logging
│   │   ├── service.go  # Business logic (testable)
│   │   ├── log.go      # Input parsing
│   │   ├── stats.go    # Statistics
│   │   ├── graph.go    # Chart generation
│   │   └── Xlsx.go     # Excel export
│   ├── handlerStart/   # User registration
│   └── callbacks/      # Inline callbacks
```

### Logging
- Use `getLogger("FunctionName")` - returns `func(string)` closure
- Log errors with context before returning
- Log successful operations for debugging

### TODO
- Consider adding `go-sqlmock` for database testing
- Add integration tests with testcontainers
