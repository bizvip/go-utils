# configer

`configer` 是一个显式加载、显式校验、格式无关的 Go 配置引擎。

它不依赖 `init()`、不依赖全局单例、也不强绑定 `viper`。核心思路是：

- 调用方自己决定配置文件路径
- 调用方自己提供解码器
- 引擎只负责路径解析、文件读取、预处理、默认值、校验、热重载管理

## 设计原则

- 显式优于隐式
- 小接口优于大而全框架
- 配置 schema 由业务项目自己定义
- 引擎只做通用流程，不做业务字段假设

## 核心能力

- `ResolvePath` 解析 `--config/-c` 参数和候选路径
- `Load[T]` / `MustLoad[T]` 泛型加载
- `DefaultsSetter` 默认值注入
- `Validatable` 加载后统一校验
- `Preprocessor` 支持解密/解压/环境替换等前置处理
- `Manager[T]` 支持内存持有和文件监听热重载

## 示例

```go
type AppConfig struct {
    Name string `toml:"name"`
}

func (c *AppConfig) Validate() error {
    if c.Name == "" {
        return fmt.Errorf("name is required")
    }
    return nil
}

cfg, err := configer.Load[AppConfig]("env.toml",
    configer.WithDecoder[AppConfig](configer.DecoderFunc(func(data []byte, out any) error {
        _, err := toml.Decode(string(data), out)
        return err
    })),
)
```
