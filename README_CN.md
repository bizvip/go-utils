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
| **os/io/logger** | 日志系统 | Zerolog + 依赖注入 |

### ⚙️ 配置

| 包名 | 描述 | 核心功能 |
|------|------|----------|
| **conf** | 配置管理 | 泛型配置、热更新 |

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
