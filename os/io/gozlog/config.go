package gozlog

import "time"

type Level string

const (
	TraceLevel Level = "trace"
	DebugLevel Level = "debug"
	InfoLevel  Level = "info"
	WarnLevel  Level = "warn"
	ErrorLevel Level = "error"
	FatalLevel Level = "fatal"
	PanicLevel Level = "panic"
)

type Config struct {
	Level           Level             `toml:"level" json:"level"`
	Format          string            `toml:"format" json:"format"`
	Output          string            `toml:"output" json:"output"`
	LogDir          string            `toml:"log_dir" json:"log_dir"`
	FileName        string            `toml:"file_name" json:"file_name"`
	TimeFormat      string            `toml:"time_format" json:"time_format"`
	NoColor         bool              `toml:"no_color" json:"no_color"`
	MaxAge          int               `toml:"max_age" json:"max_age"`
	MaxSize         int               `toml:"max_size" json:"max_size"`
	MaxBackups      int               `toml:"max_backups" json:"max_backups"`
	Compress        bool              `toml:"compress" json:"compress"`
	Caller          bool              `toml:"caller" json:"caller"`
	CallerSkip      int               `toml:"caller_skip" json:"caller_skip"`
	Service         string            `toml:"service" json:"service"`
	Fields          map[string]string `toml:"fields" json:"fields"`
	SamplingEnabled bool              `toml:"sampling_enabled" json:"sampling_enabled"`
	SamplingBurst   uint32            `toml:"sampling_burst" json:"sampling_burst"`
	SamplingPeriod  time.Duration     `toml:"sampling_period" json:"sampling_period"`
}

func DefaultConfig() Config {
	return Config{
		Level:      InfoLevel,
		Format:     "console",
		Output:     "stdout",
		FileName:   "app.log",
		TimeFormat: time.RFC3339,
		MaxAge:     30,
		MaxSize:    100,
		MaxBackups: 30,
		Compress:   true,
		Caller:     true,
	}
}

func (c Config) normalized() Config {
	cfg := c
	def := DefaultConfig()
	if cfg.Level == "" {
		cfg.Level = def.Level
	}
	if cfg.Format == "" {
		cfg.Format = def.Format
	}
	if cfg.Output == "" {
		cfg.Output = def.Output
	}
	if cfg.FileName == "" {
		cfg.FileName = def.FileName
	}
	if cfg.TimeFormat == "" {
		cfg.TimeFormat = def.TimeFormat
	}
	if cfg.MaxAge <= 0 {
		cfg.MaxAge = def.MaxAge
	}
	if cfg.MaxSize <= 0 {
		cfg.MaxSize = def.MaxSize
	}
	if cfg.MaxBackups <= 0 {
		cfg.MaxBackups = def.MaxBackups
	}
	if cfg.SamplingEnabled && cfg.SamplingBurst == 0 {
		cfg.SamplingBurst = 100
	}
	if cfg.SamplingEnabled && cfg.SamplingPeriod <= 0 {
		cfg.SamplingPeriod = time.Second
	}
	if cfg.Fields == nil {
		cfg.Fields = map[string]string{}
	}
	return cfg
}
