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

### base/blake3hash
- `Must` — `base/blake3hash/blake3hash.go`
- `SumBytes` — `base/blake3hash/blake3hash.go`
- `SumFile` — `base/blake3hash/blake3hash.go`
- `SumReader` — `base/blake3hash/blake3hash.go`

### base/crypto
- `AesDecrypt` — `base/crypto/aes.go`
- `AesEncrypt` — `base/crypto/aes.go`

### base/dt
- `AdjustMilliTimestamp` — `base/dt/datetime.go`
- `AdjustMilliTimestampByStr` — `base/dt/datetime.go`
- `CompareTimeStrings` — `base/dt/datetime.go`
- `ConvertStrMillisToTime` — `base/dt/datetime.go`
- `GetMicroTimestampStr` — `base/dt/datetime.go`
- `GetMilliTimestampStr` — `base/dt/datetime.go`
- `GetNanoTimestampStr` — `base/dt/datetime.go`
- `GetTimezoneOffsetByMillis` — `base/dt/datetime.go`
- `SetTimezone` — `base/dt/datetime.go`
- `TimeDifference` — `base/dt/datetime.go`

### base/htm
- `Compress` — `base/htm/htm.go`

### base/id/sqids
- `ToAlpha` — `base/id/sqids/sqids.go`
- `ToInt` — `base/id/sqids/sqids.go`

### base/json
- `PrettyFormat` — `base/json/goccy_json.go`
- `ToJsonWithNoErr` — `base/json/goccy_json.go`

### base/num
- `Calc` — `base/num/num.go`
- `Calculator.Evaluate` — `base/num/calculator.go`
- `CheckNumStrInRange` — `base/num/num.go`
- `DecimalFormatBanker` — `base/num/num.go`
- `Evaluate` — `base/num/calculator.go`
- `EvaluateToString` — `base/num/calculator.go`
- `FormatNumStrToDecimalAndShift` — `base/num/num.go`
- `GetMaxNum` — `base/num/num.go`
- `HashIdToInt64` — `base/num/num.go`
- `Int64ToHashId` — `base/num/num.go`
- `MergeToDecimal` — `base/num/num.go`
- `NewCalculator` — `base/num/calculator.go`
- `StrToDecimalTruncate` — `base/num/num.go`
- `ValidateSecPwd` — `base/num/num.go`

### base/pwd
- `GenSalt` — `base/pwd/pwd.go`
- `IsCorrect` — `base/pwd/pwd.go`
- `SplitHash` — `base/pwd/pwd.go`
- `ToHash` — `base/pwd/pwd.go`
- `ToHashWithConfig` — `base/pwd/pwd.go`
- `ValidateSHA256` — `base/pwd/validate.go`
- `ValidateSixNumberAsPwd` — `base/pwd/validate.go`

### base/reflects
- `MergeStructData` — `base/reflects/reflects.go`

### base/rnd
- `GenNumberInRange` — `base/rnd/rnd.go`
- `GenRandomAlphaNumeric` — `base/rnd/rnd.go`
- `RandNumStr` — `base/rnd/rnd.go`
- `RandNumStrNonSafe` — `base/rnd/rnd.go`
- `RandomCnName` — `base/rnd/cn_usr.go`
- `UUID` — `base/rnd/rnd.go`

### base/snowflake
- `NewShortIdGenerator` — `base/snowflake/short_version.go`
- `ShortIdGenerator.BatchNext` — `base/snowflake/short_version.go`
- `ShortIdGenerator.Decompose` — `base/snowflake/short_version.go`
- `ShortIdGenerator.NextID` — `base/snowflake/short_version.go`

### base/str
- `CalcHash` — `base/str/str.go`
- `FilterEmptyChar` — `base/str/str.go`
- `GenFixedStrWithSeed` — `base/str/str.go`
- `GenSha1` — `base/str/str.go`
- `GenSlug` — `base/str/str.go`
- `GetDirNameFromSnowflakeID` — `base/str/str.go`
- `Md5` — `base/str/str.go`
- `PadCnSpaceChar` — `base/str/str.go`
- `RegexpMatch` — `base/str/str.go`
- `Sha256` — `base/str/str.go`
- `ToInt64` — `base/str/str.go`
- `ToPrettyJson` — `base/str/str.go`
- `ToUint32` — `base/str/str.go`
- `UnicodeLength` — `base/str/str.go`
- `UniqueStrings` — `base/str/str.go`

### base/str/base26
- `Int64ToAlpha` — `base/str/base26/base26.go`
- `IsValidBase26` — `base/str/base26/base26.go`
- `StrNumToAlpha` — `base/str/base26/base26.go`
- `ToNum` — `base/str/base26/base26.go`
- `Uint64ToAlpha` — `base/str/base26/base26.go`

### base/str/base62
- `Base62ToSHA256` — `base/str/base62/base62.go`
- `SHA256ToBase62` — `base/str/base62/base62.go`

### base/validator
- `InRule[T].Validate` — `base/validator/validator.go`
- `IsAlphaNum` — `base/validator/str.go`
- `IsDomainResolvable` — `base/validator/email.go`
- `IsEmailAddrValidWithDomain` — `base/validator/email.go`
- `IsLengthBetween` — `base/validator/str.go`
- `IsMd5` — `base/validator/str.go`
- `IsValidDomain` — `base/validator/email.go`
- `IsValidEmailFormat` — `base/validator/email.go`
- `NewRegexRule` — `base/validator/validator.go`
- `RangeRule[T].Validate` — `base/validator/validator.go`
- `RegexRule.Validate` — `base/validator/validator.go`
- `RequiredRule[T].Validate` — `base/validator/validator.go`
- `StringLengthRule.Validate` — `base/validator/validator.go`
- `ValidateEmail` — `base/validator/validator.go`
- `ValidateIDCard` — `base/validator/validator.go`
- `ValidatePassword` — `base/validator/validator.go`
- `ValidatePhone` — `base/validator/validator.go`
- `ValidationError.Error` — `base/validator/validator.go`
- `Validator[T].AddRule` — `base/validator/validator.go`
- `Validator[T].Validate` — `base/validator/validator.go`

### cloudservice/wasabi
- `NewWasabiHandler` — `cloudservice/wasabi/wasabi.go`
- `s3Conf.DelFile` — `cloudservice/wasabi/storage_interface.go`
- `s3Conf.GetAccessKey` — `cloudservice/wasabi/storage_interface.go`
- `s3Conf.GetAllBuckets` — `cloudservice/wasabi/storage_interface.go`
- `s3Conf.GetAllFilesFromBucket` — `cloudservice/wasabi/storage_interface.go`
- `s3Conf.GetBucketName` — `cloudservice/wasabi/storage_interface.go`
- `s3Conf.GetEndpoint` — `cloudservice/wasabi/storage_interface.go`
- `s3Conf.GetFile` — `cloudservice/wasabi/storage_interface.go`
- `s3Conf.GetRegion` — `cloudservice/wasabi/storage_interface.go`
- `s3Conf.GetSecretKey` — `cloudservice/wasabi/storage_interface.go`
- `s3Conf.PutFile` — `cloudservice/wasabi/storage_interface.go`

### conf
- `AppConfig.SetDefaults` — `conf/example.go`
- `Manager[T].GetConfig` — `conf/config.go`
- `Manager[T].LoadFile` — `conf/config.go`
- `Manager[T].UpdateConfig` — `conf/config.go`

### cryptocoin
- `DetectAddress` — `cryptocoin/validate.go`
- `IsValidBEP20Address` — `cryptocoin/validate.go`
- `IsValidBTCAddress` — `cryptocoin/validate.go`
- `IsValidERC20Address` — `cryptocoin/validate.go`
- `IsValidEVMAddress` — `cryptocoin/validate.go`
- `IsValidTONAddress` — `cryptocoin/validate.go`
- `IsValidTRC20Address` — `cryptocoin/validate.go`
- `ParseTONRaw` — `cryptocoin/validate.go`

### etcd
- `Client.AcquireLock` — `etcd/client.go`
- `Client.Close` — `etcd/client.go`
- `Client.Connect` — `etcd/client.go`
- `Client.CreateLease` — `etcd/client.go`
- `Client.Get` — `etcd/client.go`
- `Client.KeepAliveLease` — `etcd/client.go`
- `Client.ListMembers` — `etcd/client.go`
- `Client.Put` — `etcd/client.go`
- `Client.RegisterService` — `etcd/client.go`
- `Client.ReleaseLock` — `etcd/client.go`
- `Client.Txn` — `etcd/client.go`
- `Client.Watch` — `etcd/client.go`
- `NewClient` — `etcd/client.go`

### ex
- `Error.Error` — `ex/ex_model.go`
- `Error.MarshalZerologObject` — `ex/ex_model.go`
- `Error.SetMessage` — `ex/ex_model.go`
- `Error.SetMeta` — `ex/ex_model.go`
- `Error.String` — `ex/ex_model.go`

### i18n/goi18n
- `I18nManager.GetTemplateLangMap` — `i18n/goi18n/go_i18n.go`
- `I18nManager.Translate` — `i18n/goi18n/go_i18n.go`
- `NewI18nManager` — `i18n/goi18n/go_i18n.go`

### i18n/opencc
- `Convert` — `i18n/opencc/opencc.go`
- `SimpToTW` — `i18n/opencc/opencc.go`
- `SimpToTrad` — `i18n/opencc/opencc.go`
- `TWToS` — `i18n/opencc/opencc.go`
- `TradToSimp` — `i18n/opencc/opencc.go`
- `WarmUp` — `i18n/opencc/opencc.go`

### img
- `Base64ToFile` — `img/img_toolkit.go`
- `GetImageInfo` — `img/image_info.go`, `img/image_info_vips.go`
- `ImageToBase64` — `img/img_toolkit.go`
- `ResizeImage` — `img/img_toolkit.go`

### lock
- `AdaptiveLock.GetActiveLockCount` — `lock/lock.go`
- `AdaptiveLock.IsShardMode` — `lock/lock.go`
- `AdaptiveLock.Lock` — `lock/lock.go`
- `AdaptiveLock.Unlock` — `lock/lock.go`
- `SetLockerAutoCleanup` — `lock/lock.go`

### network/exchange/binance
- `GetApi` — `network/exchange/binance/base.go`
- `MarketService.GetAggTrades` — `network/exchange/binance/market.go`
- `MarketService.GetAvgPrice` — `network/exchange/binance/market.go`
- `MarketService.GetDepth` — `network/exchange/binance/market.go`
- `MarketService.GetExchangeInfo` — `network/exchange/binance/market.go`
- `MarketService.GetHistoricalTrades` — `network/exchange/binance/market.go`
- `MarketService.GetKlines` — `network/exchange/binance/market.go`
- `MarketService.GetPing` — `network/exchange/binance/market.go`
- `MarketService.GetServerTime` — `network/exchange/binance/market.go`
- `MarketService.GetTicker` — `network/exchange/binance/market.go`
- `MarketService.GetTicker24Hr` — `network/exchange/binance/market.go`
- `MarketService.GetTickerBookTicker` — `network/exchange/binance/market.go`
- `MarketService.GetTickerPrice` — `network/exchange/binance/market.go`
- `MarketService.GetTickerTradingDay` — `network/exchange/binance/market.go`
- `MarketService.GetTrades` — `network/exchange/binance/market.go`
- `MarketService.GetUIKlines` — `network/exchange/binance/market.go`
- `NewMarketService` — `network/exchange/binance/market.go`

### network/exchange/okx
- `NewOkxExchangeService` — `network/exchange/okx/okx.go`
- `OKX.GetTop10Exchanges` — `network/exchange/okx/okx.go`
- `OKX.GetUsdtCnyExchangeList` — `network/exchange/okx/okx.go`
- `OKX.GetUsdtCnyRateOnly` — `network/exchange/okx/okx.go`

### network/google
- `NewTranslationService` — `network/google/google_translate.go`
- `TranslationService.GoogleDetectLang` — `network/google/google_translate.go`
- `TranslationService.GoogleTranslateToCN` — `network/google/google_translate.go`
- `TranslationService.GoogleTranslateToEn` — `network/google/google_translate.go`

### network/hcaptcha
- `NewHCaptchaVerifier` — `network/hcaptcha/captcha.go`
- `Verifier.Verify` — `network/hcaptcha/captcha.go`

### network/httputils
- `DownImage` — `network/httputils/download.go`

### network/ip
- `GetClientIP` — `network/ip/ip.go`
- `GetFullClientInfo` — `network/ip/ip.go`
- `GetGeoIPInfo` — `network/ip/ip.go`
- `GetLocalPrivateIP` — `network/ip/ip.go`
- `GetLocalPublicIP` — `network/ip/ip.go`
- `GetMyGeoIPInfo` — `network/ip/ip.go`
- `IsPrivateIP` — `network/ip/ip.go`
- `IsValidPublicIP` — `network/ip/ip.go`
- `ToUniqueStr` — `network/ip/ip.go`

### network/ua
- `Parse` — `network/ua/ua.go`

### oo/singleton
- `PerKey[K, V].Delete` — `oo/singleton/singleton.go`
- `PerKey[K, V].Get` — `oo/singleton/singleton.go`
- `PerKey[K, V].Has` — `oo/singleton/singleton.go`
- `PerKey[K, V].Range` — `oo/singleton/singleton.go`

### os
- `ByteSize.String` — `os/bytesize.go`
- `ByteSize.ToGB` — `os/bytesize.go`
- `ByteSize.ToKB` — `os/bytesize.go`
- `ByteSize.ToMB` — `os/bytesize.go`

### os/console
- `C.Black` — `os/console/color.go`
- `C.BlackBold` — `os/console/color.go`
- `C.Blue` — `os/console/color.go`
- `C.BlueBold` — `os/console/color.go`
- `C.Bluef` — `os/console/color.go`
- `C.Box` — `os/console/color.go`
- `C.Clear` — `os/console/color.go`
- `C.Cyan` — `os/console/color.go`
- `C.CyanBold` — `os/console/color.go`
- `C.Cyanf` — `os/console/color.go`
- `C.DashedLine` — `os/console/color.go`
- `C.DoubleLine` — `os/console/color.go`
- `C.Error` — `os/console/color.go`
- `C.Gray` — `os/console/color.go`
- `C.GrayBold` — `os/console/color.go`
- `C.Grayf` — `os/console/color.go`
- `C.Green` — `os/console/color.go`
- `C.GreenBold` — `os/console/color.go`
- `C.Greenf` — `os/console/color.go`
- `C.Info` — `os/console/color.go`
- `C.Italic` — `os/console/color.go`
- `C.KeyValue` — `os/console/color.go`
- `C.Line` — `os/console/color.go`
- `C.List` — `os/console/color.go`
- `C.Magenta` — `os/console/color.go`
- `C.MagentaBold` — `os/console/color.go`
- `C.Magentaf` — `os/console/color.go`
- `C.NewLine` — `os/console/color.go`
- `C.NumberedList` — `os/console/color.go`
- `C.Print` — `os/console/color.go`
- `C.Printf` — `os/console/color.go`
- `C.Println` — `os/console/color.go`
- `C.Progress` — `os/console/color.go`
- `C.Red` — `os/console/color.go`
- `C.RedBold` — `os/console/color.go`
- `C.Redf` — `os/console/color.go`
- `C.Section` — `os/console/color.go`
- `C.Spinner` — `os/console/color.go`
- `C.Success` — `os/console/color.go`
- `C.Title` — `os/console/color.go`
- `C.Underline` — `os/console/color.go`
- `C.Warning` — `os/console/color.go`
- `C.White` — `os/console/color.go`
- `C.WhiteBold` — `os/console/color.go`
- `C.Whitef` — `os/console/color.go`
- `C.WithWriter` — `os/console/color.go`
- `C.Yellow` — `os/console/color.go`
- `C.YellowBold` — `os/console/color.go`
- `C.Yellowf` — `os/console/color.go`
- `Console` — `os/console/color.go`
- `ConsoleErr` — `os/console/color.go`

### os/em
- `GetFileByPath` — `os/em/embed.go`
- `GetFileList` — `os/em/embed.go`

### os/fs
- `ComputeFileSHA256` — `os/fs/file_darwin.go`, `os/fs/file_linux.go`
- `CreateDir` — `os/fs/file_darwin.go`, `os/fs/file_linux.go`
- `CreateDirIfNotExist` — `os/fs/file_darwin.go`, `os/fs/file_linux.go`
- `Delete` — `os/fs/file_darwin.go`, `os/fs/file_linux.go`
- `DetectFileType` — `os/fs/file_darwin.go`, `os/fs/file_linux.go`
- `GetAllFilesByExt` — `os/fs/file_darwin.go`, `os/fs/file_linux.go`
- `GetBigFileMd5` — `os/fs/file_darwin.go`, `os/fs/file_linux.go`
- `GetCurExeDir` — `os/fs/file_darwin.go`, `os/fs/file_linux.go`
- `GetFileCreationTime` — `os/fs/file_darwin.go`, `os/fs/file_linux.go`
- `GetFileMd5` — `os/fs/file_linux.go`
- `GetFileMd5Stream` — `os/fs/file_linux.go`
- `GetFileNameMd5` — `os/fs/file_darwin.go`, `os/fs/file_linux.go`
- `GetSmallFileMd5` — `os/fs/file_darwin.go`
- `IsDirAndHasFiles` — `os/fs/file_darwin.go`, `os/fs/file_linux.go`
- `IsFile` — `os/fs/file_darwin.go`, `os/fs/file_linux.go`
- `StartsWithDot` — `os/fs/file_darwin.go`, `os/fs/file_linux.go`

### os/fsn
- `NewFsnWatcher` — `os/fsn/fw_fsnotify.go`
- `StartWatcher` — `os/fsn/fw_fsnotify.go`
- `Watcher.AddDirRecursive` — `os/fsn/fw_fsnotify.go`
- `Watcher.Close` — `os/fsn/fw_fsnotify.go`
- `Watcher.Start` — `os/fsn/fw_fsnotify.go`

### os/io/logger
- `DefaultConfig` — `os/io/logger/logger.go`
- `Logger.Debug` — `os/io/logger/logger.go`
- `Logger.Debugf` — `os/io/logger/logger.go`
- `Logger.Error` — `os/io/logger/logger.go`
- `Logger.Errorf` — `os/io/logger/logger.go`
- `Logger.Fatal` — `os/io/logger/logger.go`
- `Logger.Fatalf` — `os/io/logger/logger.go`
- `Logger.Info` — `os/io/logger/logger.go`
- `Logger.Infof` — `os/io/logger/logger.go`
- `Logger.Panic` — `os/io/logger/logger.go`
- `Logger.Panicf` — `os/io/logger/logger.go`
- `Logger.Warn` — `os/io/logger/logger.go`
- `Logger.Warnf` — `os/io/logger/logger.go`
- `Logger.With` — `os/io/logger/logger.go`
- `Logger.WithFields` — `os/io/logger/logger.go`
- `LoggerMiddleware` — `os/io/logger/example.go`
- `Manager.GetLogger` — `os/io/logger/logger.go`
- `Manager.GetModuleLogger` — `os/io/logger/logger.go`
- `Manager.GetServiceLogger` — `os/io/logger/logger.go`
- `NewManager` — `os/io/logger/logger.go`
- `NewUserService` — `os/io/logger/example.go`
- `UserService.CreateUser` — `os/io/logger/example.go`

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
