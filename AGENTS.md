# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

Always use context7 when I need code generation, setup or configuration steps, or
library/API documentation. This means you should automatically use the Context7 MCP
tools to resolve library id and get library docs without me having to explicitly ask.

## Project Purpose

**go-utils 是一个快捷工具封装库，专门为其他 Go 项目提供开箱即用的工具函数集合。**

本项目的定位：
- 作为其他项目的依赖库存在，而非独立应用
- 提供高度封装的常用功能，减少重复代码
- 所有函数都经过优化和测试，可直接在生产环境使用
- 遵循 Go 最佳实践，保持简洁、高效、可靠

使用方式：
```go
import "github.com/bizvip/go-utils/base/num"
import "github.com/bizvip/go-utils/base/str"
// 直接调用封装好的工具函数
```

## Build and Test Commands

### Essential Commands
```bash
# Run tests with automatic formatting (default)
make test

# Run single test file
go test -v ./tests/base/num/calculator_test.go

# Run specific test function
go test -v -run TestCalculator_Evaluate ./tests/base/num/

# Full CI pipeline (tidy, vet, test)
make ci

# Run tests with coverage
make test-coverage

# Run benchmarks
make bench

# Clean up dependencies
make tidy

# Run go vet with formatting
make vet
```

### Important: All test commands automatically run `go fmt` first. Tests are centralized in the `tests/` directory, mirroring the source structure.

## Architecture Overview

### Core Design Principles
- **Utility Library Focus**: 作为工具库为其他项目服务，不是独立应用
- **Latest Go Version**: Uses Go 1.24+ with cutting-edge features
- **Performance First**: All critical functions are benchmarked
- **Simple & Reliable**: Minimal dependencies, clear interfaces
- **Test Centralization**: All tests in `tests/` directory, not alongside source files
- **Zero Config**: 大部分功能无需配置即可使用
- **Dependency Injection**: 需要状态管理的模块支持依赖注入

### Package Structure

The codebase follows a modular architecture with clear separation of concerns:

```
base/           # Core utilities - stateless, zero external dependencies
├── num/        # Numerical operations (Calculator, decimal handling, ID encoding)
├── str/        # String manipulation (Base26/62 encoding, validation)
├── crypto/     # Cryptographic operations (AES)
├── dt/         # Date/time utilities
├── pwd/        # Password generation and validation
├── collections/# Generic collections with Go 1.24+ features
├── validator/  # Generic validation framework
├── id/         # ID generation utilities
├── snowflake/  # Distributed ID generation
├── rnd/        # Random generation utilities
├── json/       # JSON utilities
└── htm/        # HTML utilities

network/        # External API integrations
├── google/     # Google Translate API with batch processing
├── exchange/   # Cryptocurrency exchange APIs (Binance, OKX)
└── httputils/  # HTTP helpers and download utilities

cloudservice/   # Cloud provider integrations
└── wasabi/     # S3-compatible storage interface

os/             # System-level utilities
├── console/    # Terminal output with colors
├── fs/         # Cross-platform file operations (Darwin/Linux specific)
├── fsn/        # File system notifications
└── io/
    └── logger/ # Structured logging with dependency injection

conf/           # Configuration management
└── config.go   # Generic configuration with file watching

tests/          # Centralized test directory
└── [mirrors source structure]
```

### Key Architectural Patterns

1. **Calculator Implementation** (`base/num/calculator.go`)
   - Uses Go AST for expression parsing instead of regex
   - Global instance pattern with both struct methods and package-level functions
   - Backward compatibility maintained through wrapper functions

2. **Cross-Platform File Operations** (`os/fs/`)
   - Separate implementations for Darwin and Linux
   - Uses build tags for conditional compilation
   - Consistent function naming across platforms (e.g., `GetBigFileMd5`)

3. **Modern Generic Collections** (`base/collections/`)
   - Type-safe generic operations using Go 1.24+ features
   - Functional programming patterns (Filter, Map, Reduce)
   - Performance optimized with slices standard library

4. **Validation System** (`base/validator/`)
   - Generic validation rules with `ValidationRule[T any]` interface
   - Composable validators with type safety
   - Built-in validators for common use cases

5. **ID Generation Strategy**
   - Sqids for hash IDs (`base/id/sqids/`)
   - Snowflake variants (`base/snowflake/`)
   - Base26/62 encoding for custom formats

6. **Error Handling Pattern**
   - Pre-defined error variables in packages (e.g., `ErrInvalidSecPwdLength`)
   - Modern error wrapping with `fmt.Errorf("%w: %w", ...)`
   - Structured error types with field information

7. **Logger System** (`os/io/logger/`)
   - Dependency injection design with `Manager` pattern
   - Daily log rotation with configurable retention (default 30 days)
   - Default log directory: `runtime/logs` (configurable)
   - Supports structured logging with zerolog
   - No global state or init functions

8. **Configuration Management** (`conf/`)
   - Generic configuration with type safety
   - Viper integration for multiple format support
   - File watching for hot-reload capability
   - Thread-safe configuration updates

## Recent Updates

### Logger Refactoring (os/io/logger)
- **Removed init function**: 改为依赖注入模式，通过 `NewManager` 创建管理器
- **Daily log rotation**: 使用 lumberjack 实现按天轮转，默认保留30天
- **Configurable log directory**: 默认 `runtime/logs`，可通过配置自定义
- **Usage example**:
  ```go
  config := logger.DefaultConfig()
  config.Output = "file"
  config.LogDir = "/custom/logs"
  
  logManager, err := logger.NewManager(config)
  log := logManager.GetLogger()
  ```

## Testing Strategy

### Test Organization
- Tests use `package_test` naming convention for black-box testing
- Table-driven tests with subtests for comprehensive coverage
- Benchmark tests included for performance-critical functions

### Running Tests
```bash
# Test a specific package
go test -v ./tests/base/num/...

# Run with race detector
go test -race -v ./tests/...

# Generate coverage for specific package
go test -coverprofile=cover.out ./tests/base/num/
go tool cover -html=cover.out
```

## Code Style Requirements

1. **No Comments**: Do not add comments unless explicitly requested
2. **Formatting**: All code must pass `go fmt` (automatically run before tests)
3. **Chinese Error Messages**: Some functions use Chinese error messages (e.g., "无效的表达式")
4. **Import Organization**: Standard library first, then external, then internal packages
5. **Go 1.24+ Features**: Use modern Go features - generics, error wrapping, slices package
6. **Cross-Platform Compatibility**: Ensure function names are consistent across Darwin/Linux builds
7. **Package Independence**: base/ packages should have minimal external dependencies

## Common Gotchas

1. **Test Package Names**: Tests must use `package_test` suffix and explicit imports
2. **Makefile Indentation**: Must use tabs, not spaces
3. **Expression Calculator**: Supports `×` and `÷` Unicode operators, automatically converted to `*` and `/`
4. **Security Password Validation**: 6 digits, no three consecutive incremental numbers
5. **Logger Dependency**: Requires `gopkg.in/natefinch/lumberjack.v2` for log rotation

## GitHub Actions

Two workflows are configured:
- **ci.yml**: Runs on PRs and pushes, executes tests and quality checks
- **release.yml**: Auto-creates releases on master push with semantic versioning

Uses `actions/checkout@v4` and requires `GH_TOKEN` secret for releases.

## Integration Guidelines

当其他项目使用 go-utils 时：
1. 通过 `go get github.com/bizvip/go-utils` 引入
2. 直接导入需要的包，无需初始化
3. 大部分函数都是无状态的，可直接调用
4. 需要状态管理的模块（如 logger）使用依赖注入模式
5. 遵循每个包的错误处理约定