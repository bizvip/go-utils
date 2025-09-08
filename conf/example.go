package conf

// AppConfig 示例应用配置
type AppConfig struct {
	Server   ServerConfig   `toml:"server"`
	Database DatabaseConfig `toml:"database"`
	Logger   LoggerConfig   `toml:"logger"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	URL         string `toml:"url"`
	MaxConn     int    `toml:"max_conn"`
	MaxIdleConn int    `toml:"max_idle_conn"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level  string `toml:"level"`
	Format string `toml:"format"`
}

// SetDefaults 实现Config接口，设置默认值
func (c *AppConfig) SetDefaults() {
	if c.Server.Host == "" {
		c.Server.Host = "localhost"
	}
	if c.Server.Port == 0 {
		c.Server.Port = 8080
	}
	if c.Database.MaxConn == 0 {
		c.Database.MaxConn = 10
	}
	if c.Database.MaxIdleConn == 0 {
		c.Database.MaxIdleConn = 5
	}
	if c.Logger.Level == "" {
		c.Logger.Level = "info"
	}
	if c.Logger.Format == "" {
		c.Logger.Format = "json"
	}
}
