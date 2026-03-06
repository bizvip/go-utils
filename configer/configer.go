package configer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Decoder interface {
	Decode(data []byte, out any) error
}

type DecoderFunc func(data []byte, out any) error

func (f DecoderFunc) Decode(data []byte, out any) error {
	return f(data, out)
}

type DefaultsSetter interface {
	SetDefaults()
}

type Validatable interface {
	Validate() error
}

type Preprocessor func([]byte) ([]byte, error)

type AfterLoad[T any] func(*T) error

type ResolveOptions struct {
	Flags       []string
	SearchPaths []string
	Fallback    string
}

type loadOptions[T any] struct {
	decoder       Decoder
	searchPaths   []string
	preprocessors []Preprocessor
	afterLoad     []AfterLoad[T]
}

type LoadOption[T any] func(*loadOptions[T])

func WithDecoder[T any](decoder Decoder) LoadOption[T] {
	return func(opts *loadOptions[T]) {
		opts.decoder = decoder
	}
}

func WithSearchPaths[T any](paths ...string) LoadOption[T] {
	return func(opts *loadOptions[T]) {
		opts.searchPaths = append(opts.searchPaths, paths...)
	}
}

func WithPreprocessor[T any](fn Preprocessor) LoadOption[T] {
	return func(opts *loadOptions[T]) {
		opts.preprocessors = append(opts.preprocessors, fn)
	}
}

func WithAfterLoad[T any](fn AfterLoad[T]) LoadOption[T] {
	return func(opts *loadOptions[T]) {
		opts.afterLoad = append(opts.afterLoad, fn)
	}
}

func ResolvePath(args []string, opts ResolveOptions) string {
	flags := opts.Flags
	if len(flags) == 0 {
		flags = []string{"--config", "-c"}
	}

	for i := 0; i < len(args); i++ {
		arg := strings.TrimSpace(args[i])
		for _, flag := range flags {
			if arg == flag && i+1 < len(args) {
				return strings.TrimSpace(args[i+1])
			}
			if strings.HasPrefix(arg, flag+"=") {
				return strings.TrimSpace(strings.TrimPrefix(arg, flag+"="))
			}
		}
	}

	for _, candidate := range opts.SearchPaths {
		candidate = strings.TrimSpace(candidate)
		if candidate == "" {
			continue
		}
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return strings.TrimSpace(opts.Fallback)
}

func Load[T any](path string, options ...LoadOption[T]) (*T, error) {
	opts := loadOptions[T]{}
	for _, option := range options {
		option(&opts)
	}

	resolvedPath := strings.TrimSpace(path)
	if resolvedPath == "" {
		for _, candidate := range opts.searchPaths {
			candidate = strings.TrimSpace(candidate)
			if candidate == "" {
				continue
			}
			if _, err := os.Stat(candidate); err == nil {
				resolvedPath = candidate
				break
			}
		}
	}
	if resolvedPath == "" {
		return nil, fmt.Errorf("config path is required")
	}
	if opts.decoder == nil {
		return nil, fmt.Errorf("config decoder is required")
	}

	raw, err := os.ReadFile(resolvedPath)
	if err != nil {
		return nil, fmt.Errorf("read config file %q failed: %w", resolvedPath, err)
	}
	for _, preprocessor := range opts.preprocessors {
		raw, err = preprocessor(raw)
		if err != nil {
			return nil, fmt.Errorf("preprocess config file %q failed: %w", resolvedPath, err)
		}
	}

	cfg := new(T)
	if defaults, ok := any(cfg).(DefaultsSetter); ok {
		defaults.SetDefaults()
	}
	if err := opts.decoder.Decode(raw, cfg); err != nil {
		return nil, fmt.Errorf("decode config file %q failed: %w", resolvedPath, err)
	}
	for _, hook := range opts.afterLoad {
		if err := hook(cfg); err != nil {
			return nil, fmt.Errorf("post-process config file %q failed: %w", resolvedPath, err)
		}
	}
	if validatable, ok := any(cfg).(Validatable); ok {
		if err := validatable.Validate(); err != nil {
			return nil, fmt.Errorf("validate config file %q failed: %w", resolvedPath, err)
		}
	}
	return cfg, nil
}

func MustLoad[T any](path string, options ...LoadOption[T]) *T {
	cfg, err := Load[T](path, options...)
	if err != nil {
		panic(err)
	}
	return cfg
}

func ResolveExecutableRelative(configName string) (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("resolve executable failed: %w", err)
	}
	resolvedPath := filepath.Join(filepath.Dir(execPath), configName)
	if _, err := os.Stat(resolvedPath); err != nil {
		return "", fmt.Errorf("config file %q not found: %w", resolvedPath, err)
	}
	return resolvedPath, nil
}
