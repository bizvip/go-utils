# Config Package - 通用配置管理器

一个基于Go 1.24泛型的通用配置管理包，支持自动从可执行文件同级目录加载TOML配置文件。

## 特性

- **类型安全**: 使用Go 1.24泛型提供完全的类型安全
- **自动发现**: 支持从可执行文件同级目录自动加载配置文件
- **热重载**: 支持配置文件变化监听和自动重载
- **默认值**: 支持通过接口自定义默认配置
- **线程安全**: 使用读写锁保证并发安全
- **零依赖**: 可以独立使用到任何Go项目

## 快速开始

### 1. 定义配置结构

```go
type AppConfig struct {
    Server   ServerConfig   `toml:"server"`
    Database DatabaseConfig `toml:"database"`
}

type ServerConfig struct {
    Host string `toml:"host"`
    Port int    `toml:"port"`
}

type DatabaseConfig struct {
    URL     string `toml:"url"`
    MaxConn int    `toml:"max_conn"`
}

// 可选：实现默认值接口
func (c *AppConfig) SetDefaults() {
    if c.Server.Host == "" {
        c.Server.Host = "localhost"
    }
    if c.Server.Port == 0 {
        c.Server.Port = 8080
    }
}
```

### 2. 使用方式

#### 基础用法

```go
// 创建配置实例
cfg := &AppConfig{}
manager := config.New(cfg)

// 加载指定路径的配置文件
err := manager.LoadFile("/path/to/conf.toml", false)
if err != nil {
    log.Fatal(err)
}

// 获取配置（类型安全）
config := manager.GetConfig()
fmt.Printf("Server: %s:%d\n", config.Server.Host, config.Server.Port)
```

#### 自动发现配置文件

```go
// 自动从可执行文件同级目录加载 app.toml
cfg := &AppConfig{}
manager, err := config.NewFromExecutable(cfg, "app.toml")
if err != nil {
    log.Fatal(err)
}

// 获取配置
config := manager.GetConfig()
```

#### 启用热重载

```go
// 自动加载并启用文件监听
cfg := &AppConfig{}
manager, err := config.NewFromExecutableWithWatch(cfg, "app.toml")
if err != nil {
    log.Fatal(err)
}

// 配置文件变化时会自动重载
```

#### 原子性更新配置

```go
// 线程安全的配置更新
manager.UpdateConfig(func(cfg *AppConfig) {
    cfg.Server.Port = 9090
    cfg.Database.MaxConn = 20
})
```

### 3. 配置文件示例 (app.toml)

```toml
[server]
host = "0.0.0.0"
port = 8080

[database]
url = "postgresql://localhost/mydb"
max_conn = 10
```

## API 文档

### 类型

#### Manager[T any]
泛型配置管理器，T为配置结构体类型。

#### Config 接口
```go
type Config interface {
    SetDefaults()
}
```
可选接口，用于设置配置的默认值。

### 函数

#### New[T any](config *T) *Manager[T]
创建新的配置管理器实例。

#### NewFromExecutable[T any](config *T, configName string) (*Manager[T], error)
自动从可执行文件同级目录加载指定名称的配置文件。

#### NewFromExecutableWithWatch[T any](config *T, configName string) (*Manager[T], error)
自动加载配置文件并启用文件监听。

### 方法

#### LoadFile(filePath string, watch bool) error
加载指定路径的配置文件，可选择是否启用监听。

#### GetConfig() *T
获取当前配置（类型安全）。

#### UpdateConfig(updateFn func(*T))
原子性更新配置。

## 使用场景

### 1. 命令行工具
```bash
# 可执行文件结构
./myapp
./conf.toml
```

```go
// main.go
func main() {
    cfg := &AppConfig{}
    manager, err := config.NewFromExecutable(cfg, "conf.toml")
    if err != nil {
        log.Fatal(err)
    }
    
    // 使用配置
    config := manager.GetConfig()
    startServer(config.Server.Host, config.Server.Port)
}
```

### 2. 微服务应用
```go
// 支持热重载的微服务配置
cfg := &ServiceConfig{}
manager, err := config.NewFromExecutableWithWatch(cfg, "service.toml")
if err != nil {
    log.Fatal(err)
}

// 配置会自动重载，无需重启服务
```

### 3. 多环境配置
```go
// 根据环境变量加载不同配置
env := os.Getenv("ENV")
if env == "" {
    env = "dev"
}

cfg := &AppConfig{}
configFile := fmt.Sprintf("conf-%s.toml", env)
manager, err := config.NewFromExecutable(cfg, configFile)
```

## 优势

1. **类型安全**: 编译时检查配置字段类型
2. **零配置**: 开箱即用，自动发现配置文件
3. **性能优化**: 使用读写锁，支持高并发读取
4. **内存安全**: 使用指针避免大结构体拷贝
5. **易于集成**: 可以轻松集成到任何Go项目
6. **现代化**: 充分利用Go 1.24+的最新特性

## 依赖

- github.com/spf13/viper (配置解析)
- github.com/fsnotify/fsnotify (文件监听)
- github.com/rs/zerolog (日志记录)

所有依赖都是成熟稳定的库，适合生产环境使用。