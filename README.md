# 🛠️ Go Utils

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/bizvip/go-utils)

> A comprehensive collection of Go utility libraries and self-developed components with simple and elegant encapsulations.

This package continuously adopts the latest Go versions and cutting-edge technologies to maintain the most modern and advanced features. We prioritize performance optimization while keeping everything simple and reliable.

[中文文档](README_CN.md) | English

## 📋 Table of Contents

- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [Package Overview](#-package-overview)
- [Complete Function Directory](#-complete-function-directory)
- [Features](#-features)
- [Testing](#-testing)
- [Contributing](#-contributing)
- [License](#-license)

## 🚀 Installation

```bash
go get github.com/bizvip/go-utils
```

## 🎯 Quick Start

```go
package main

import (
    "fmt"
    "github.com/bizvip/go-utils/base/num"
    "github.com/bizvip/go-utils/base/str"
)

func main() {
    // Mathematical expression calculation
    result, _ := num.Calc("(2 + 3) * 4")
    fmt.Println("Result:", result) // Result: 20

    // String utilities
    encoded := str.Base62Encode(12345)
    fmt.Println("Encoded:", encoded)
}
```

## 📦 Package Overview

### 🔢 Base Utilities

| Package | Description | Key Features |
|---------|-------------|--------------|
| **base/num** | Numerical operations and calculations | Expression calculator with AST parsing, decimal handling, Sqids hash ID encoding |
| **base/str** | String manipulation utilities | Base26/Base62 encoding, hash calculations (MD5/SHA256/SHA3), Unicode operations |
| **base/str/base26** | Base26 encoding | Numeric/string to base26 conversions, validation |
| **base/str/base62** | Base62 encoding | SHA256 <-> Base62 conversions |
| **base/crypto** | Cryptographic operations | AES-GCM encryption/decryption with PBKDF2 key derivation |
| **base/blake3hash** | BLAKE3 hashing | Stream/file hashing helpers |
| **base/dt** | Date and time utilities | Timezone offset calculations, timestamp manipulation, time comparisons |
| **base/pwd** | Password utilities | Argon2 hashing, security password validation with pattern checks |
| **base/collections** | Generic collections (Go 1.24+) | Filter, Map, Reduce, GroupBy, Chunk with type safety |
| **base/validator** | Validation framework | Generic validators, email/phone/ID card validation |
| **base/id/sqids** | ID generation | Sqids library integration for hash IDs |
| **base/snowflake** | Snowflake IDs | Yitter ID Generator, custom base time support |
| **base/rnd** | Random generation | Cryptographically secure random strings, UUID generation |
| **base/json** | JSON utilities | goccy/go-json integration, pretty formatting |
| **base/reflects** | Reflection helpers | Struct merge and Struct-to-map conversion |
| **base/htm** | HTML utilities | HTML compression and minification |

### 🌐 Network & APIs

| Package | Description | Key Features |
|---------|-------------|--------------|
| **network/google** | Google services integration | Translate API with batch processing via RapidAPI |
| **network/exchange/binance** | Binance exchange API | Market data retrieval with automatic fallback |
| **network/exchange/okx** | OKX exchange API | Fiat/crypto exchange rates |
| **network/hcaptcha** | hCaptcha verification | Server-side token verification |
| **network/httputils** | HTTP utilities | Image download with custom headers, path sanitization |
| **network/ip** | IP utilities | Local/public IP, GeoIP, client info |
| **network/ua** | User-Agent parser | Parses UA strings into device/browser info |

### ☁️ Cloud Services

| Package | Description | Key Features |
|---------|-------------|--------------|
| **cloudservice/wasabi** | S3-compatible storage | Wasabi cloud storage interface |

### 🖼️ Image

| Package | Description | Key Features |
|---------|-------------|--------------|
| **img** | Image utilities | Base64 conversion, resize, image info (optional libvips) |

### 🈯 I18n & Localization

| Package | Description | Key Features |
|---------|-------------|--------------|
| **i18n/goi18n** | i18n helpers | go-i18n integration helpers |
| **i18n/opencc** | Chinese conversion | OpenCC wrapper utilities |

### 📌 Constants

| Package | Description | Key Features |
|---------|-------------|--------------|
| **consts/cryptocurrency** | Crypto constants | Coin/network/token constants, explorers |
| **consts/currencycode** | Currency codes | Fiat currency code constants |

### 🧩 Patterns & Concurrency

| Package | Description | Key Features |
|---------|-------------|--------------|
| **lock** | Adaptive lock | Sharded lock with cleanup |
| **oo/singleton** | Singleton helpers | Lazy initialization, per-key singleton |

### 🧰 Infrastructure & Models

| Package | Description | Key Features |
|---------|-------------|--------------|
| **etcd** | etcd client wrapper | KV ops, leases, locks, watch |
| **ex** | Error model | Structured error type with metadata |

### 🛠️ System & OS

| Package | Description | Key Features |
|---------|-------------|--------------|
| **os** | Byte size helpers | ByteSize conversions and formatting |
| **os/console** | Console utilities | Colored terminal output, progress bars, formatted display |
| **os/em** | Embed helpers | Read embedded files and list directories |
| **os/fs** | File system operations | Cross-platform file handling (Darwin/Linux), hash calculations |
| **os/fsn** | File watch | fsnotify-based directory watcher |
| **os/io/logger** | Logging | Zerolog-based logger with DI and rotation |

### ⚙️ Configuration

| Package | Description | Key Features |
|---------|-------------|--------------|
| **conf** | Generic configuration management | Go 1.24+ generics, Viper integration, file watching |

## 📚 Complete Function Directory

### base/num - Numerical Operations

| Function | Description |
|----------|-------------|
| `ValidateSecPwd(secPwd string) error` | Validates 6-digit security passwords |
| `Int64ToHashId(number int64, minLen uint8) string` | Converts int64 to Sqids hash ID |
| `HashIdToInt64(id string, minLen uint8) (int64, error)` | Converts hash ID back to int64 |
| `MergeToDecimal(number *big.Int, dec int) decimal.Decimal` | Shifts decimal point left |
| `FormatNumStrToDecimalAndShift(number string, decimals uint) decimal.Decimal` | Converts string to decimal with shift |
| `CheckNumStrInRange(s string, min, max float64) (bool, error)` | Checks if number is in range |
| `StrToDecimalTruncate(s string, precision int32) (decimal.Decimal, error)` | Converts with truncation |
| `DecimalFormatBanker(value decimal.Decimal) string` | Banker's rounding format |
| `GetMaxNum(vals ...int) int` | Returns maximum integer |
| `Calc(exp string) (string, error)` | Evaluates mathematical expressions |
| `NewCalculator() *Calculator` | Creates calculator instance |
| `Evaluate(expression string) (float64, error)` | Evaluates using AST parsing |

### base/str - String Utilities

| Function | Description |
|----------|-------------|
| `ToUint32(str string) uint32` | Converts string to uint32 using FNV hash |
| `PadCnSpaceChar(label string, spaces int) string` | Pads with Chinese spaces |
| `UniqueSlice[T comparable](input []T) []T` | Returns unique elements |
| `RegexpMatch(txt, pattern string) (bool, error)` | Regex pattern matching |
| `ParseInt[T ~int...](intStr string) (T, error)` | Generic integer parsing |
| `Md5(input string, useStream bool) (string, error)` | MD5 hash calculation |
| `Sha256(input string, useStream bool, isSha3 bool) (string, error)` | SHA256/SHA3 hash |
| `FilterEmptyChar(str string) string` | Removes empty characters |
| `UnicodeLength(str string) int` | Unicode string length |
| `ToPrettyJson(v interface{}, isProto bool) (string, error)` | JSON formatting |
| `GenSlug(title string) string` | URL-friendly slug generation |

### base/str/base26 - Base26 Encoding

| Function | Description |
|----------|-------------|
| `Uint64ToAlpha(input uint64) (string, error)` | Converts uint64 to base26 |
| `Int64ToAlpha(input int64) (string, error)` | Converts int64 with sign support |
| `StrNumToAlpha(input string) (string, error)` | String number to base26 |
| `ToNum(alphaStr string) (string, error)` | Base26 to decimal string |
| `IsValidBase26(s string) bool` | Validates base26 format |

### base/str/base62 - Base62 Encoding

| Function | Description |
|----------|-------------|
| `SHA256ToBase62(sha256Hash string) (string, error)` | SHA256 to base62 |
| `Base62ToSHA256(base62Str string) (string, error)` | Base62 to SHA256 |

### base/crypto - Cryptographic Functions

| Function | Description |
|----------|-------------|
| `AesEncrypt(text, pass string) (string, error)` | AES-GCM encryption with PBKDF2 |
| `AesDecrypt(cipherText, pass string) (string, error)` | AES-GCM decryption |

### base/dt - Date/Time Utilities

| Function | Description |
|----------|-------------|
| `GetTimezoneOffsetByMillis(millis int64) (string, error)` | Calculates timezone offset |
| `AdjustMilliTimestamp(timestamp uint64, seconds int64) uint64` | Adjusts timestamp |
| `AdjustMilliTimestampByStr(timestamp uint64, shift string) (uint64, error)` | Adjusts by time unit |
| `GetNanoTimestampStr() string` | Current nanosecond timestamp |
| `GetMicroTimestampStr() string` | Current microsecond timestamp |
| `GetMilliTimestampStr() string` | Current millisecond timestamp |
| `ConvertStrMillisToTime(millis string) (time.Time, error)` | String millis to time |
| `SetTimezone(tz ...string)` | Sets timezone (default Shanghai) |
| `CompareTimeStrings(t1, t2, layout string) (int, error)` | Compares time strings |
| `TimeDifference(t1, t2 string) (time.Duration, error)` | Calculates time difference |

### base/pwd - Password Utilities

| Function | Description |
|----------|-------------|
| `GenSalt() (string, error)` | Generates random salt |
| `ToHash(password string) (string, error)` | Argon2 password hashing |
| `ToHashWithConfig(password string, config HashConfig) (string, error)` | Custom config hashing |
| `IsCorrect(password, hashStr string) (bool, error)` | Verifies password |
| `ValidateSixNumberAsPwd(secPwd string, length int) error` | Validates numeric passwords |
| `ValidateSHA256(hash string) error` | Validates SHA256 format |

### base/collections - Generic Collections

| Function | Description |
|----------|-------------|
| `Filter[T any](slice []T, predicate func(T) bool) []T` | Filters by predicate |
| `Map[T, U any](slice []T, mapper func(T) U) []U` | Maps to different type |
| `Reduce[T, U any](slice []T, initialValue U, reducer func(U, T) U) U` | Reduces collection |
| `Find[T any](slice []T, predicate func(T) bool) (T, bool)` | Finds first match |
| `Contains[T comparable](slice []T, target T) bool` | Checks containment |
| `Unique[T comparable](slice []T) []T` | Returns unique elements |
| `SortBy[T any, K cmp.Ordered](slice []T, keyFunc func(T) K)` | Sorts by key |
| `GroupBy[T any, K comparable](slice []T, keyFunc func(T) K) map[K][]T` | Groups by key |
| `Chunk[T any](slice []T, size int) [][]T` | Splits into chunks |
| `Reverse[T any](slice []T)` | Reverses in-place |

### base/validator - Validation Framework

| Function | Description |
|----------|-------------|
| `NewValidator[T any](rules ...ValidationRule[T]) *Validator[T]` | Creates validator |
| `ValidateEmail(email, field string) error` | Validates email format |
| `ValidatePhone(phone, field string) error` | Validates Chinese phone |
| `ValidateIDCard(idCard, field string) error` | Validates Chinese ID card |
| `ValidatePassword(password, field string) error` | Validates password strength |
| `IsValidEmailFormat(email string) bool` | Email format validation |
| `IsDomainResolvable(domain string) bool` | Checks domain resolution |
| `IsEmailAddrValidWithDomain(email string) error` | Email with domain check |
| `IsMd5(input string) error` | Validates MD5 format |
| `IsAlphaNum(str string) bool` | Alphanumeric check |
| `IsLengthBetween(str string, min, max int) bool` | Length range validation |

### base/id/sqids - ID Generation

| Function | Description |
|----------|-------------|
| `ToAlpha(ids []uint64) string` | Converts IDs to alphanumeric |
| `ToInt(ids string) []uint64` | Converts string to ID array |

### base/snowflake - Snowflake IDs

| Function | Description |
|----------|-------------|
| `InitWith(workerId uint16, baseTime *time.Time)` | Initialize with settings |
| `ID() uint64` | Generates new Snowflake ID |
| `QuickID() uint64` | Backward compatible ID generation |

### base/rnd - Random Generation

| Function | Description |
|----------|-------------|
| `RandNumStr(length int) string` | Secure random digits |
| `UUID(isNoDash bool) string` | UUID with/without dashes |
| `GenRandomAlphaNumeric() string` | Random alphanumeric string |
| `GenNumberInRange(min, max int) int` | Random number in range |

### base/json - JSON Utilities

| Function | Description |
|----------|-------------|
| `PrettyFormat(in string) string` | Formats JSON with indentation |
| `ToJsonWithNoErr(payload interface{}, pretty bool) string` | JSON marshaling without errors |

### base/htm - HTML Utilities

| Function | Description |
|----------|-------------|
| `Compress(htmlSrc string, stripScriptStyle bool) (string, error)` | HTML compression |

### network/google - Google Translate

| Function | Description |
|----------|-------------|
| `NewTranslationService(apiKey, apiHost string) *TranslationService` | Creates service |
| `GoogleTranslateToEn(text, source string) (string, error)` | Translates to English |
| `GoogleTranslateToCN(text, source string) (string, error)` | Translates to Chinese |
| `GoogleDetectLang(text string) (string, error)` | Detects language |

### network/httputils - HTTP Utilities

| Function | Description |
|----------|-------------|
| `DownImage(url, name, savePath string) (string, error)` | Downloads images |

### network/exchange/binance - Binance API

| Function | Description |
|----------|-------------|
| `GetApi(query string) interface{}` | Generic GET with fallback |

### cloudservice/wasabi - S3 Storage

| Function | Description |
|----------|-------------|
| `NewWasabiHandler(bucketName, region, endpoint, accessKey, secretKey string) *StorageHandler` | Creates handler |

### os/console - Terminal Output

| Function | Description |
|----------|-------------|
| `Console() *C` | Returns stdout console |
| `ConsoleErr() *C` | Returns stderr console |
| `Black/Red/Green/Yellow/Blue/Magenta/Cyan/White/Gray(txt string)` | Colored output |
| `Success/Error/Warning/Info(txt string)` | Status messages |
| `Progress(current, total, width int)` | Progress bar |
| `Box/Title/Section(txt string)` | Formatted output |

### os/fs - File System Operations

| Function | Description |
|----------|-------------|
| `ComputeFileSHA256(filePath string) (string, error)` | SHA256 of file |
| `GetFileCreationTime(filePath string) (string, time.Time, error)` | File creation time |
| `GetFileNameMd5(filename string) (string, error)` | MD5 of filename |
| `GetSmallFileMd5/GetBigFileMd5(filePath string) (string, error)` | File MD5 |
| `GetCurExeDir() string` | Current executable directory |
| `GetAllFilesByExt(dir, ext string) ([]string, error)` | Files by extension |
| `IsDirAndHasFiles(dirPath string) (bool, bool, error)` | Directory validation |
| `Delete(path string) error` | File/directory deletion |
| `CreateDir/CreateDirIfNotExist(path string) error` | Directory creation |
| `IsFile(path string) (bool, error)` | File type check |
| `DetectFileType(file io.Reader) (string, error)` | MIME type detection |

### conf - Configuration Management

| Function | Description |
|----------|-------------|
| `New[T any](config *T) *Manager[T]` | Creates config manager |
| `NewFromExecutable[T any](config *T, configName string) (*Manager[T], error)` | Auto-loads config |
| `NewFromExecutableWithWatch[T any](config *T, configName string) (*Manager[T], error)` | With file watching |
| `LoadFile(filePath string, watch bool) error` | Loads configuration |
| `GetConfig() *T` | Thread-safe config access |
| `UpdateConfig(updateFn func(*T))` | Atomic config updates |

## ✨ Features

### 🧮 Advanced Calculator
```go
// Supports complex mathematical expressions
result, _ := num.Calc("2 + 3 × (4 ÷ 2)")
fmt.Println(result) // Output: 8
```

### 🔐 Secure Password Validation
```go
// Validates security passwords with custom rules
err := num.ValidateSecPwd("123456")
if err != nil {
    fmt.Println("Invalid password:", err)
}
```

### 🌍 Google Translate Integration
```go
// Batch translation with Google Translate API
translations, _ := google.BatchTranslate(texts, "en", "zh")
```

### 💱 Cryptocurrency Support
```go
// Validate cryptocurrency addresses
isValid := cryptocoin.ValidateBTCAddress("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa")
```

### 🎨 Console Styling
```go
// Colorized console output
console.PrintSuccess("Operation completed successfully!")
console.PrintError("An error occurred!")
```

## 🧪 Testing

We maintain comprehensive test coverage for all packages:

```bash
# Run all tests with formatting
make test

# Run tests with coverage
make test-coverage

# Run benchmarks
make bench

# Full CI pipeline
make ci
```

### Test Structure
```
tests/
├── base/
│   └── num/
│       ├── calculator_test.go
│       └── num_test.go
└── [other packages]/
```

## 🏗️ Project Structure

```
go-utils/
├── base/                    # Core utilities
│   ├── crypto/             # Cryptographic functions
│   ├── dt/                 # Date/time utilities  
│   ├── num/                # Numerical operations
│   ├── str/                # String manipulation
│   └── ...
├── network/                # Network-related packages
│   ├── google/             # Google APIs
│   ├── exchange/           # Crypto exchanges
│   └── ...
├── cloudservice/           # Cloud service integrations
├── tests/                  # Centralized test files
├── Makefile               # Build automation
└── README.md
```

## 📊 Performance

All critical functions are benchmarked:

- **Calculator**: ~500ns per expression evaluation
- **Base62 Encoding**: ~100ns per operation  
- **String Validation**: ~50ns per check

## 🤝 Contributing

We welcome contributions! Please see our [contributing guidelines](CONTRIBUTING.md).

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## ⚠️ Stability Notice

This library is under active development. APIs may change between versions. Please check the changelog before upgrading.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with ❤️ for the Go community
- Inspired by various open-source Go utilities
- Special thanks to all contributors

---

<div align="center">
<p>Made with ❤️ by <a href="https://github.com/bizvip">@bizvip</a></p>
<p>⭐ Star us on GitHub if this project helped you!</p>
</div>
