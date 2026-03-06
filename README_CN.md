# 🛠️ Go Utils

[![Go 版本](https://img.shields.io/badge/Go-%3E%3D%201.21-blue.svg)](https://golang.org/)
[![许可证](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![构建状态](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/bizvip/go-utils)

> 一个全面的 Go 工具库集合，包含自主开发的组件和优雅简洁的封装。

持续不断使用最新版本的golang版本和技术，来保持最现代化、最前卫的特性更新，并持续优先关注性能优化，保持简单可靠。

[English](README.md) | 中文文档

## 📋 目录

- [安装](#-安装)
- [快速开始](#-快速开始)
- [包概览](#-包概览)
- [功能函数目录](#-功能函数目录)
- [功能特性](#-功能特性)
- [测试](#-测试)
- [贡献](#-贡献)
- [许可证](#-许可证)

## 🚀 安装

```bash
go get github.com/bizvip/go-utils
```

## 🎯 快速开始

```go
package main

import (
    "fmt"
    "github.com/bizvip/go-utils/base/num"
    "github.com/bizvip/go-utils/base/str"
)

func main() {
    // 数学表达式计算
    result, _ := num.Calc("(2 + 3) * 4")
    fmt.Println("结果:", result) // 结果: 20

    // 字符串工具
    encoded := str.Base62Encode(12345)
    fmt.Println("编码结果:", encoded)
}
```

## 📦 包概览

### 🔢 基础工具

| 包名 | 描述 | 核心功能 |
|------|------|----------|
| **base/num** | 数值操作和计算 | 表达式计算器、小数处理、ID编码 |
| **base/str** | 字符串操作工具 | Base26/Base62编码、字符串验证 |
| **base/str/base26** | Base26 编码 | 数字/字符串与 Base26 互转 |
| **base/str/base62** | Base62 编码 | SHA256 与 Base62 互转 |
| **base/crypto** | 加密操作 | AES加密/解密 |
| **base/blake3hash** | BLAKE3 哈希 | 文件/流哈希计算 |
| **base/dt** | 日期时间工具 | 日期格式化、解析、计算 |
| **base/pwd** | 密码工具 | 生成、验证、安全检查 |
| **base/collections** | 泛型集合 | Filter/Map/Reduce/GroupBy 等 |
| **base/validator** | 验证框架 | 邮箱/手机号/身份证等验证 |
| **base/id/sqids** | ID 生成 | Sqids 编码/解码 |
| **base/snowflake** | 雪花ID | 分布式 ID 生成 |
| **base/rnd** | 随机工具 | 安全随机数/UUID |
| **base/json** | JSON 工具 | JSON 格式化 |
| **base/reflects** | 反射工具 | 结构体合并/转 map |
| **base/htm** | HTML工具 | 压缩、清理、优化 |

### 🌐 网络与API

| 包名 | 描述 | 核心功能 |
|------|------|----------|
| **network/google** | Google服务集成 | 翻译API与批量处理 |
| **network/exchange/binance** | Binance API | 行情获取 |
| **network/exchange/okx** | OKX API | 汇率与报价 |
| **network/hcaptcha** | hCaptcha 验证 | 服务端验证 |
| **network/httputils** | HTTP工具 | 下载助手、请求构建器 |
| **network/ip** | IP地址工具 | IP 信息与 GeoIP |
| **network/ua** | User-Agent 解析 | 浏览器/设备信息 |

### ☁️ 云服务

| 包名 | 描述 | 核心功能 |
|------|------|----------|
| **cloudservice/wasabi** | Wasabi云存储 | 文件上传、下载、管理 |

### 🖼️ 图像

| 包名 | 描述 | 核心功能 |
|------|------|----------|
| **img** | 图像处理工具包 | 调整大小、格式转换、优化 |

### 🈯 国际化

| 包名 | 描述 | 核心功能 |
|------|------|----------|
| **i18n/goi18n** | i18n 工具 | go-i18n 封装 |
| **i18n/opencc** | 中文转换 | OpenCC 封装 |

### 📌 常量

| 包名 | 描述 | 核心功能 |
|------|------|----------|
| **consts/cryptocurrency** | 加密货币常量 | 币种/链/协议/浏览器 |
| **consts/currencycode** | 法币代码 | 常用货币代码 |

### 🧩 模式与并发

| 包名 | 描述 | 核心功能 |
|------|------|----------|
| **lock** | 自适应锁 | 分片锁与清理 |
| **oo/singleton** | 单例工具 | 延迟初始化/按 key 单例 |

### 🧰 基础设施与模型

| 包名 | 描述 | 核心功能 |
|------|------|----------|
| **etcd** | ETCD 客户端 | KV/租约/锁/监听 |
| **ex** | 错误模型 | 结构化错误与元数据 |

### 🛠️ 系统与操作系统

| 包名 | 描述 | 核心功能 |
|------|------|----------|
| **os** | ByteSize 工具 | 字节大小换算 |
| **os/console** | 控制台工具 | 彩色输出、格式化 |
| **os/em** | Embed 工具 | 读取嵌入文件 |
| **os/fs** | 文件系统操作 | 跨平台文件处理 |
| **os/fsn** | 文件监听 | 基于 fsnotify |
| **os/io/gozlog** | 日志系统 | Zerolog + 依赖注入 |

### ⚙️ 配置

| 包名 | 描述 | 核心功能 |
|------|------|----------|
| **conf** | 配置管理 | 泛型配置、热更新 |

## 📚 功能函数目录

以下为所有对外导出的函数/方法及其源码路径。

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

### os/io/gozlog
- `DefaultConfig` — `os/io/gozlog/logger.go`
- `Logger.Debug` — `os/io/gozlog/logger.go`
- `Logger.Debugf` — `os/io/gozlog/logger.go`
- `Logger.Error` — `os/io/gozlog/logger.go`
- `Logger.Errorf` — `os/io/gozlog/logger.go`
- `Logger.Fatal` — `os/io/gozlog/logger.go`
- `Logger.Fatalf` — `os/io/gozlog/logger.go`
- `Logger.Info` — `os/io/gozlog/logger.go`
- `Logger.Infof` — `os/io/gozlog/logger.go`
- `Logger.Panic` — `os/io/gozlog/logger.go`
- `Logger.Panicf` — `os/io/gozlog/logger.go`
- `Logger.Warn` — `os/io/gozlog/logger.go`
- `Logger.Warnf` — `os/io/gozlog/logger.go`
- `Logger.With` — `os/io/gozlog/logger.go`
- `Logger.WithFields` — `os/io/gozlog/logger.go`
- `LoggerMiddleware` — `os/io/gozlog/example.go`
- `Manager.GetLogger` — `os/io/gozlog/logger.go`
- `Manager.GetModuleLogger` — `os/io/gozlog/logger.go`
- `Manager.GetServiceLogger` — `os/io/gozlog/logger.go`
- `NewManager` — `os/io/gozlog/logger.go`
- `NewUserService` — `os/io/gozlog/example.go`
- `UserService.CreateUser` — `os/io/gozlog/example.go`

## ✨ 功能特性

### 🧮 高级计算器
```go
// 支持复杂数学表达式
result, _ := num.Calc("2 + 3 × (4 ÷ 2)")
fmt.Println(result) // 输出: 8
```

### 🔐 安全密码验证
```go
// 使用自定义规则验证安全密码
err := num.ValidateSecPwd("123456")
if err != nil {
    fmt.Println("密码无效:", err)
}
```

### 🌍 谷歌翻译集成
```go
// 使用谷歌翻译API进行批量翻译
translations, _ := google.BatchTranslate(texts, "en", "zh")
```

### 💱 加密货币支持
```go
// 验证加密货币地址
isValid := cryptocoin.ValidateBTCAddress("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa")
```

### 🎨 控制台样式
```go
// 彩色控制台输出
console.PrintSuccess("操作成功完成！")
console.PrintError("发生错误！")
```

## 🧪 测试

我们为所有包维护全面的测试覆盖：

```bash
# 运行所有测试（包含格式化）
make test

# 运行覆盖率测试
make test-coverage

# 运行基准测试
make bench

# 完整CI流水线
make ci
```

### 测试结构
```
tests/
├── base/
│   └── num/
│       ├── calculator_test.go
│       └── num_test.go
└── [其他包]/
```

## 🏗️ 项目结构

```
go-utils/
├── base/                    # 核心工具
│   ├── crypto/             # 加密功能
│   ├── dt/                 # 日期/时间工具  
│   ├── num/                # 数值操作
│   ├── str/                # 字符串操作
│   └── ...
├── network/                # 网络相关包
│   ├── google/             # Google APIs
│   ├── exchange/           # 加密货币交易所
│   └── ...
├── cloudservice/           # 云服务集成
├── tests/                  # 集中测试文件
├── Makefile               # 构建自动化
└── README.md
```

## 📊 性能

所有关键功能都经过基准测试：

- **计算器**: 每次表达式评估约500纳秒
- **Base62编码**: 每次操作约100纳秒  
- **字符串验证**: 每次检查约50纳秒

## 🤝 贡献

欢迎贡献！请查看我们的[贡献指南](CONTRIBUTING.md)。

1. Fork 这个仓库
2. 创建您的功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启一个 Pull Request

## ⚠️ 稳定性说明

此库正在积极开发中。版本之间API可能会发生变化。升级前请检查更新日志。

## 📄 许可证

此项目在MIT许可证下授权 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- 用 ❤️ 为Go社区构建
- 受各种开源Go工具启发
- 特别感谢所有贡献者

## 📚 详细功能说明

### 数值计算包 (base/num)
- **表达式计算器**: 基于Go AST的强大数学表达式解析器
- **交易密码验证**: 6位数字密码，防止连续数字
- **哈希ID转换**: Sqids编码的安全ID转换
- **小数处理**: 精确的小数运算和格式化

### 字符串处理包 (base/str)  
- **Base26编码**: 适用于序号编码的Base26算法
- **Base62编码**: 高效的数字到字符串编码
- **字符串验证**: 多种格式验证工具

### 网络工具包 (network/*)
- **谷歌翻译**: 支持批量翻译的Google Translate API封装
- **交易所API**: Binance和OKX加密货币市场数据
- **HTTP工具**: 文件下载、请求构建等实用工具

### 系统工具包 (os/*)
- **控制台工具**: 彩色输出、进度条、格式化显示
- **文件系统**: 跨平台文件操作，支持macOS和Linux
- **嵌入文件**: Go embed文件管理工具

---

<div align="center">
<p>由 <a href="https://github.com/bizvip">@bizvip</a> 用 ❤️ 制作</p>
<p>⭐ 如果这个项目对您有帮助，请在GitHub上给我们一个星标！</p>
</div>
