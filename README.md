# 🛠️ Go Utils

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/bizvip/go-utils)

> A comprehensive collection of Go utility libraries and self-developed components with simple and elegant encapsulations.

[中文文档](README_CN.md) | English

## 📋 Table of Contents

- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [Package Overview](#-package-overview)
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
| **base/num** | Numerical operations and calculations | Expression calculator, decimal handling, ID encoding |
| **base/str** | String manipulation utilities | Base26/Base62 encoding, string validation |
| **base/crypto** | Cryptographic operations | AES encryption/decryption |
| **base/dt** | Date and time utilities | Date formatting, parsing, calculations |
| **base/pwd** | Password utilities | Generation, validation, security checks |

### 🌐 Network & APIs

| Package | Description | Key Features |
|---------|-------------|--------------|
| **network/google** | Google services integration | Translate API with batch processing |
| **network/exchange** | Cryptocurrency exchange APIs | Binance, OKX market data |
| **network/httputils** | HTTP utilities | Download helpers, request builders |
| **network/ip** | IP address utilities | Geolocation, validation |

### ☁️ Cloud Services

| Package | Description | Key Features |
|---------|-------------|--------------|
| **cloudservice/wasabi** | Wasabi cloud storage | File upload, download, management |

### 🖼️ Media & Processing

| Package | Description | Key Features |
|---------|-------------|--------------|
| **img** | Image processing toolkit | Resize, format conversion, optimization |
| **i18n** | Internationalization | Multi-language support, OpenCC integration |

### 🛠️ System & OS

| Package | Description | Key Features |
|---------|-------------|--------------|
| **os/console** | Console utilities | Colored output, formatting |
| **os/fs** | File system operations | Cross-platform file handling |
| **lock** | Concurrency utilities | Atomic locks, synchronization |

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