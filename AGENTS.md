# AGENTS.md

本文件是给自动化代理和协作者使用的项目工作指南。修改代码前先按这里的约束理解项目定位、命令、包结构和常见坑点。

当任务涉及代码生成、项目搭建、配置步骤、库/API 文档时，必须使用 Context7：先 resolve library id，再查询对应文档。不要等用户额外要求。

## 项目定位

`go-utils` 是给其他 Go 项目直接依赖的工具封装库，不是独立应用。

- 模块路径：`github.com/bizvip/go-utils`
- Go 版本：`1.26.4+`
- 目标：提供开箱即用的常用工具函数和轻量组件，减少业务项目重复代码
- 设计取向：简洁、可靠、性能优先、尽量少配置
- 使用方式：按需导入具体包，大多数工具无需初始化

示例：

```go
import "github.com/bizvip/go-utils/base/num"
import "github.com/bizvip/go-utils/base/str"
```

## 常用命令

```bash
# 默认测试，会先 go fmt ./...
make test

# 完整 CI：tidy、vet、test
make ci

# vet，会先 go fmt ./...
make vet

# 覆盖率，会先 go fmt ./...
make test-coverage

# benchmark，会先 go fmt ./...
make bench

# 整理依赖
make tidy

# 单个测试包
go test -v ./tests/base/num

# 单个测试函数
go test -v ./tests/base/num -run TestCalculator_Evaluate

# img 纯 Go 默认路径
go test -v ./tests/img

# img 可选 libvips 后端
go test -tags libvips -v ./tests/img
```

注意：`make test`、`make vet`、`make test-coverage`、`make bench` 都会先执行 `go fmt ./...`，会修改源码格式。只想验证单点行为时优先用 `go test` 包级命令。

## 包结构

当前主要包：

```text
base/                    核心基础工具
  blake3hash/            BLAKE3 文件/流哈希
  collections/           泛型集合工具
  crypto/                AES 等加密工具
  dt/                    日期时间工具
  htm/                   HTML 压缩
  id/bizid/              业务 ID 工具
  id/sqids/              Sqids 编码
  json/                  JSON 工具
  num/                   数值、表达式计算、小数处理
  pwd/                   密码哈希与安全密码校验
  reflects/              反射辅助
  rnd/                   随机数、UUID、随机中文名
  snowflake/             短版雪花 ID
  str/                   字符串、哈希、slug、Unicode 工具
  str/base26/            Base26 编码
  str/base62/            Base62 编码
  validator/             泛型验证框架

cloudservice/wasabi/     Wasabi/S3 兼容存储
configer/                显式配置加载、校验、热重载
consts/                  货币、加密货币等常量
cryptocoin/              加密货币地址识别与校验
etcd/                    etcd KV、租约、锁、watch 封装
ex/                      结构化错误模型
i18n/                    go-i18n 与 OpenCC 封装
img/                     图片 Base64、缩放、元信息读取
lock/                    自适应锁
network/                 外部 API、HTTP、IP、UA 工具
oo/singleton/            Lazy 与按 key 单例
os/                      ByteSize 与系统工具
os/console/              控制台输出
os/em/                   embed.FS 读取工具
os/fs/                   跨平台文件系统工具
os/fsn/                  fsnotify 文件监听
os/io/gozlog/            zerolog 日志管理器
tests/                   集中测试目录
```

## 架构约定

- `base/` 下的包应尽量保持轻依赖、无状态、可直接调用。
- 需要状态或生命周期的模块使用显式构造和依赖注入，例如 `configer`、`os/io/gozlog`、`etcd`。
- 不要引入全局隐式初始化，除非现有包已有明确模式。
- 新增公共函数时优先放在语义最窄的包内，避免把跨领域功能塞进 `base/str` 或 `base/num`。
- 保持向后兼容。已有导出函数、类型、错误变量、包路径不能随意改名或删除。
- 错误处理优先使用明确错误变量、结构化错误或 `fmt.Errorf("%w", err)` 包装。
- `img` 默认是纯 Go 实现；`-tags libvips` 才启用 govips/libvips 路径，并要求 `CGO_ENABLED=1`、系统安装 libvips 和 pkg-config。
- `configer` 是实际配置包名；如果其他文档出现 `conf`，以源码目录 `configer` 为准。

## 关键实现模式

- `base/num` 表达式计算器使用 Go AST 解析，支持包级函数和 `Calculator` 实例方法，兼容 `×`、`÷` 到 `*`、`/` 的转换。
- `os/fs` 使用 Darwin/Linux 分文件实现，跨平台导出函数名必须保持一致。
- `base/collections` 和 `base/validator` 使用 Go 泛型，保持类型安全和小接口。
- `base/id/sqids`、`base/snowflake`、`base/str/base26`、`base/str/base62` 分别承担不同 ID/编码场景，不要混用语义。
- `os/io/gozlog` 使用 `NewManager` 显式创建，支持 service/module scoped logger、stdout/stderr/file 输出和 lumberjack 轮转。
- `configer` 由调用方提供 decoder，支持默认值、校验、预处理和 watch，不绑定 viper。

## 测试策略

- 测试主要集中在 `tests/`，按源码结构镜像组织；新增测试优先放到对应 `tests/...` 路径。
- 测试包名通常使用 `_test` 后缀，做黑盒测试并显式导入目标包。
- 已存在包内示例或特殊测试时可以沿用，例如 `configer/example_test.go`。
- 数值、加密、验证、编码、图片处理、并发锁等风险较高的逻辑需要表驱动测试覆盖边界条件。
- 修改 `img` 时至少覆盖默认纯 Go 路径；涉及 libvips 专属行为时补充 `-tags libvips` 测试说明或测试。
- 修改跨平台文件逻辑时检查 Darwin/Linux build tag 和导出函数一致性。

## 代码风格

1. 不要添加注释，除非用户明确要求，或复杂逻辑没有注释会显著降低可维护性。
2. 所有 Go 代码必须通过 `gofmt`。
3. import 分组按标准库、外部依赖、内部包顺序。
4. 保留现有中文错误信息风格；不要无理由改成英文。
5. 使用 Go 1.26.4 可用的现代特性，但不要为了新语法牺牲可读性。
6. `base/` 包新增依赖要特别克制，避免把基础工具变重。
7. Makefile 命令必须使用 tab 缩进。

## 常见坑点

- 当前 `go.mod` 和 README 要求 Go `1.26.4+`；`.github/workflows/ci.yml` 里仍配置 `actions/setup-go` 的 `1.21`，这是已知不一致。不要在无关任务里顺手改 CI。
- `make` 测试类目标会自动格式化仓库，可能产生与任务无关的格式化 diff。
- `README_CN.md` 的配置概览可能写成 `conf`，实际包是 `configer`。
- 安全密码校验要求 6 位数字，并禁止连续递增等弱模式。
- `network/` 包可能依赖外部服务，测试时避免引入不稳定网络调用，除非已有 mock 或明确集成测试意图。
- `etcd`、`cloudservice/wasabi`、`network/exchange/*` 属于外部系统封装，改动时注意超时、上下文、凭据和错误包装。
- `os/io/gozlog` 依赖 `gopkg.in/natefinch/lumberjack.v2` 做日志轮转。

## GitHub Actions

- `ci.yml`：PR 和 push 到 `master`/`main` 时运行格式检查、vet、测试和覆盖率；Markdown 等文档变更被 push 触发忽略。
- `release.yml`：push 到 `master` 时自动递增 `vMAJOR.MINOR.PATCH` patch 版本并创建 GitHub Release。

## 集成使用原则

其他项目使用本库时：

1. 通过 `go get github.com/bizvip/go-utils` 引入。
2. 直接导入需要的包，不要导入根模块。
3. 无状态工具函数直接调用。
4. 有状态模块显式创建 manager/client/verifier。
5. 遵循各包已有错误返回和生命周期约定。
