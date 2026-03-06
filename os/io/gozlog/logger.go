package gozlog

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var disabledZerolog = zerolog.New(io.Discard)

type contextKey struct{}

type Logger struct {
	provider func() *zerolog.Logger
}

func NewLazy(provider func() *zerolog.Logger) *Logger {
	return &Logger{provider: provider}
}

func Wrap(logger zerolog.Logger) *Logger {
	log := logger
	return WrapPtr(&log)
}

func WrapPtr(logger *zerolog.Logger) *Logger {
	return NewLazy(func() *zerolog.Logger {
		if logger == nil {
			return &disabledZerolog
		}
		return logger
	})
}

func (l *Logger) Zerolog() *zerolog.Logger {
	if l == nil || l.provider == nil {
		return &disabledZerolog
	}
	logger := l.provider()
	if logger == nil {
		return &disabledZerolog
	}
	return logger
}

func (l *Logger) With() zerolog.Context { return l.Zerolog().With() }
func (l *Logger) Trace() *zerolog.Event { return l.Zerolog().Trace() }
func (l *Logger) Debug() *zerolog.Event { return l.Zerolog().Debug() }
func (l *Logger) Info() *zerolog.Event  { return l.Zerolog().Info() }
func (l *Logger) Warn() *zerolog.Event  { return l.Zerolog().Warn() }
func (l *Logger) Error() *zerolog.Event { return l.Zerolog().Error() }
func (l *Logger) Fatal() *zerolog.Event { return l.Zerolog().Fatal() }
func (l *Logger) Panic() *zerolog.Event { return l.Zerolog().Panic() }
func (l *Logger) Printf(format string, args ...interface{}) {
	l.Info().Msgf(format, args...)
}
func (l *Logger) Debugf(format string, args ...interface{}) { l.Debug().Msgf(format, args...) }
func (l *Logger) Infof(format string, args ...interface{})  { l.Info().Msgf(format, args...) }
func (l *Logger) Warnf(format string, args ...interface{})  { l.Warn().Msgf(format, args...) }
func (l *Logger) Errorf(format string, args ...interface{}) { l.Error().Msgf(format, args...) }

func IntoContext(ctx context.Context, logger *Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if logger == nil {
		return ctx
	}
	return context.WithValue(ctx, contextKey{}, logger.Zerolog())
}

func FromContext(ctx context.Context) *Logger {
	if ctx == nil {
		return nil
	}
	value := ctx.Value(contextKey{})
	if logger, ok := value.(*zerolog.Logger); ok && logger != nil {
		return WrapPtr(logger)
	}
	return nil
}

type Manager struct {
	cfg   Config
	root  zerolog.Logger
	mu    sync.RWMutex
	cache map[string]*zerolog.Logger
}

func NewManager(cfg Config) (*Manager, error) {
	cfg = cfg.normalized()
	writer, err := buildWriter(cfg)
	if err != nil {
		return nil, err
	}
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}
	base := zerolog.New(writer).Level(level).With().Timestamp()
	if cfg.Caller {
		base = base.CallerWithSkipFrameCount(cfg.CallerSkip)
	}
	if cfg.Service != "" {
		base = base.Str("service", cfg.Service)
	}
	for key, value := range cfg.Fields {
		if key == "" || value == "" {
			continue
		}
		base = base.Str(key, value)
	}
	root := base.Logger()
	if cfg.SamplingEnabled {
		root = root.Sample(&zerolog.BurstSampler{Burst: cfg.SamplingBurst, Period: cfg.SamplingPeriod, NextSampler: &zerolog.BasicSampler{N: 1}})
	}
	return &Manager{cfg: cfg, root: root, cache: make(map[string]*zerolog.Logger)}, nil
}

func (m *Manager) Config() Config              { return m.cfg }
func (m *Manager) Logger() *Logger             { return WrapPtr(&m.root) }
func (m *Manager) Root() *zerolog.Logger       { return &m.root }
func (m *Manager) Named(name string) *Logger   { return m.scoped("logger", name) }
func (m *Manager) Service(name string) *Logger { return m.scoped("service", name) }
func (m *Manager) Module(name string) *Logger  { return m.scoped("module", name) }
func (m *Manager) IntoContext(ctx context.Context, logger *Logger) context.Context {
	return IntoContext(ctx, logger)
}
func (m *Manager) FromContext(ctx context.Context) *Logger {
	if logger := FromContext(ctx); logger != nil {
		return logger
	}
	return m.Logger()
}

func (m *Manager) scoped(field string, value string) *Logger {
	return NewLazy(func() *zerolog.Logger {
		return m.resolveScoped(field, value)
	})
}

func (m *Manager) resolveScoped(field string, value string) *zerolog.Logger {
	if value == "" {
		return &m.root
	}
	key := field + ":" + value
	m.mu.RLock()
	if logger, ok := m.cache[key]; ok {
		m.mu.RUnlock()
		return logger
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()
	if logger, ok := m.cache[key]; ok {
		return logger
	}
	child := m.root.With().Str(field, value).Logger()
	m.cache[key] = &child
	return &child
}

func buildWriter(cfg Config) (io.Writer, error) {
	var writer io.Writer
	switch cfg.Output {
	case "stderr":
		writer = os.Stderr
	case "file":
		if cfg.LogDir == "" {
			return nil, fmt.Errorf("log_dir is required when output=file")
		}
		if err := os.MkdirAll(cfg.LogDir, 0o755); err != nil {
			return nil, fmt.Errorf("create log dir: %w", err)
		}
		writer = &lumberjack.Logger{
			Filename:   filepath.Join(cfg.LogDir, cfg.FileName),
			MaxSize:    cfg.MaxSize,
			MaxAge:     cfg.MaxAge,
			MaxBackups: cfg.MaxBackups,
			Compress:   cfg.Compress,
			LocalTime:  true,
		}
	default:
		writer = os.Stdout
	}
	if cfg.Format == "console" {
		return zerolog.ConsoleWriter{Out: writer, TimeFormat: cfg.TimeFormat, NoColor: cfg.NoColor}, nil
	}
	return writer, nil
}

func parseLevel(level Level) (zerolog.Level, error) {
	switch level {
	case TraceLevel:
		return zerolog.TraceLevel, nil
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
		return zerolog.InfoLevel, fmt.Errorf("unsupported log level %q", level)
	}
}
