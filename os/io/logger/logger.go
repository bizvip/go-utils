package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogLevel 定义日志级别
type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	FatalLevel LogLevel = "fatal"
	PanicLevel LogLevel = "panic"
)

// Config 日志配置
type Config struct {
	Level      LogLevel `toml:"level" json:"level"`             // 日志级别
	Format     string   `toml:"format" json:"format"`           // 输出格式: json, console
	Output     string   `toml:"output" json:"output"`           // 输出目标: stdout, stderr, file
	LogDir     string   `toml:"log_dir" json:"log_dir"`         // 日志目录 (当output为file时)
	TimeFormat string   `toml:"time_format" json:"time_format"` // 时间格式
	NoColor    bool     `toml:"no_color" json:"no_color"`       // 禁用颜色输出
	MaxAge     int      `toml:"max_age" json:"max_age"`         // 日志保留天数
	MaxSize    int      `toml:"max_size" json:"max_size"`       // 单个日志文件最大大小(MB)
	MaxBackups int      `toml:"max_backups" json:"max_backups"` // 最多保留的日志文件数量
	Compress   bool     `toml:"compress" json:"compress"`       // 是否压缩旧日志
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		Level:      InfoLevel,
		Format:     "console",
		Output:     "stdout",
		LogDir:     "runtime/logs",
		TimeFormat: time.RFC3339,
		NoColor:    false,
		MaxAge:     30,   // 默认保留30天
		MaxSize:    100,  // 默认单个文件100MB
		MaxBackups: 30,   // 默认最多保留30个文件
		Compress:   true, // 默认压缩旧日志
	}
}

// Logger 统一日志器
type Logger struct {
	*zerolog.Logger
	config Config
	writer io.Writer
}

// Manager 日志管理器，用于依赖注入
type Manager struct {
	defaultLogger *Logger
	loggers       map[string]*Logger
}

// NewManager 创建日志管理器
func NewManager(config Config) (*Manager, error) {
	defaultLogger, err := newLogger(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create default logger: %w", err)
	}

	return &Manager{
		defaultLogger: defaultLogger,
		loggers:       make(map[string]*Logger),
	}, nil
}

// GetLogger 获取默认日志器
func (m *Manager) GetLogger() *Logger {
	return m.defaultLogger
}

// GetServiceLogger 获取服务专用日志器
func (m *Manager) GetServiceLogger(serviceName string) *Logger {
	key := "service:" + serviceName
	if logger, exists := m.loggers[key]; exists {
		return logger
	}

	// 创建带服务名的日志器
	zLogger := m.defaultLogger.With().Str("service", serviceName).Logger()
	logger := &Logger{
		Logger: &zLogger,
		config: m.defaultLogger.config,
		writer: m.defaultLogger.writer,
	}
	m.loggers[key] = logger
	return logger
}

// GetModuleLogger 获取模块专用日志器
func (m *Manager) GetModuleLogger(moduleName string) *Logger {
	key := "module:" + moduleName
	if logger, exists := m.loggers[key]; exists {
		return logger
	}

	// 创建带模块名的日志器
	zLogger := m.defaultLogger.With().Str("module", moduleName).Logger()
	logger := &Logger{
		Logger: &zLogger,
		config: m.defaultLogger.config,
		writer: m.defaultLogger.writer,
	}
	m.loggers[key] = logger
	return logger
}

// newLogger 创建新的日志器实例
func newLogger(config Config) (*Logger, error) {
	// 设置全局日志级别
	level, err := parseLogLevel(config.Level)
	if err != nil {
		return nil, err
	}
	zerolog.SetGlobalLevel(level)

	// 配置输出目标
	var writer io.Writer
	switch config.Output {
	case "stderr":
		writer = os.Stderr
	case "file":
		writer, err = createRotatingFileWriter(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create rotating file writer: %w", err)
		}
	default: // stdout
		writer = os.Stdout
	}

	// 配置输出格式
	if config.Format == "console" {
		writer = zerolog.ConsoleWriter{
			Out:        writer,
			TimeFormat: config.TimeFormat,
			NoColor:    config.NoColor,
		}
	}

	// 创建logger
	logger := zerolog.New(writer).With().
		Timestamp().
		Caller().
		Logger()

	return &Logger{
		Logger: &logger,
		config: config,
		writer: writer,
	}, nil
}

// createRotatingFileWriter 创建按天轮转的日志文件写入器
func createRotatingFileWriter(config Config) (io.Writer, error) {
	// 确保日志目录存在
	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// 生成日志文件名（包含日期）
	filename := fmt.Sprintf("app-%s.log", time.Now().Format("2006-01-02"))
	logPath := filepath.Join(config.LogDir, filename)

	// 使用lumberjack进行日志轮转
	return &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    config.MaxSize,    // 单个文件最大大小（MB）
		MaxAge:     config.MaxAge,     // 保留天数
		MaxBackups: config.MaxBackups, // 最多保留文件数
		LocalTime:  true,              // 使用本地时间
		Compress:   config.Compress,   // 压缩旧日志
	}, nil
}

// parseLogLevel 解析日志级别
func parseLogLevel(level LogLevel) (zerolog.Level, error) {
	switch level {
	case DebugLevel:
		return zerolog.DebugLevel, nil
	case InfoLevel:
		return zerolog.InfoLevel, nil
	case WarnLevel:
		return zerolog.WarnLevel, nil
	case ErrorLevel:
		return zerolog.ErrorLevel, nil
	case FatalLevel:
		return zerolog.FatalLevel, nil
	case PanicLevel:
		return zerolog.PanicLevel, nil
	default:
		return zerolog.InfoLevel, nil
	}
}

// Logger方法

// Debug 记录调试日志
func (l *Logger) Debug() *zerolog.Event {
	return l.Logger.Debug()
}

// Info 记录信息日志
func (l *Logger) Info() *zerolog.Event {
	return l.Logger.Info()
}

// Warn 记录警告日志
func (l *Logger) Warn() *zerolog.Event {
	return l.Logger.Warn()
}

// Error 记录错误日志
func (l *Logger) Error() *zerolog.Event {
	return l.Logger.Error()
}

// Fatal 记录致命错误日志并退出程序
func (l *Logger) Fatal() *zerolog.Event {
	return l.Logger.Fatal()
}

// Panic 记录panic日志并触发panic
func (l *Logger) Panic() *zerolog.Event {
	return l.Logger.Panic()
}

// 格式化日志方法

// Debugf 格式化调试日志
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Logger.Debug().Msgf(format, v...)
}

// Infof 格式化信息日志
func (l *Logger) Infof(format string, v ...interface{}) {
	l.Logger.Info().Msgf(format, v...)
}

// Warnf 格式化警告日志
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Logger.Warn().Msgf(format, v...)
}

// Errorf 格式化错误日志
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Logger.Error().Msgf(format, v...)
}

// Fatalf 格式化致命错误日志并退出程序
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Logger.Fatal().Msgf(format, v...)
}

// Panicf 格式化panic日志并触发panic
func (l *Logger) Panicf(format string, v ...interface{}) {
	l.Logger.Panic().Msgf(format, v...)
}

// WithFields 创建带字段的日志器
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	ctx := l.Logger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	zLogger := ctx.Logger()
	return &Logger{
		Logger: &zLogger,
		config: l.config,
		writer: l.writer,
	}
}

// With 创建子日志器
func (l *Logger) With() zerolog.Context {
	return l.Logger.With()
}
