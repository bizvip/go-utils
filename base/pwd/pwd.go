package pwd

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2 参数
const (
	Argon2Time    = 2         // Argon2 参数，迭代次数
	Argon2Memory  = 19 * 1024 // Argon2 参数，内存大小（KB）
	Argon2Threads = 1         // Argon2 参数，线程数
	Argon2KeyLen  = 32        // Argon2 参数，生成的密钥长度
	SaltSize      = 16        // 盐值长度
)

// 定义常见错误消息
var (
	ErrGenerateSaltFailed = errors.New("failed to generate salt")
	ErrDecodeSaltFailed   = errors.New("failed to decode salt")
	ErrSplitHashInvalid   = errors.New("invalid hash format")
	ErrDecodeHashFailed   = errors.New("failed to decode hash")
)

// GenSalt 生成一个随机的盐值，用于密码哈希
func GenSalt() (string, error) {
	salt := make([]byte, SaltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", fmt.Errorf("%w: %v", ErrGenerateSaltFailed, err)
	}
	return base64.RawStdEncoding.EncodeToString(salt), nil
}

// HashConfig 哈希配置参数
type HashConfig struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

// DefaultHashConfig 默认哈希配置
var DefaultHashConfig = HashConfig{
	Time:    Argon2Time,
	Memory:  Argon2Memory,
	Threads: Argon2Threads,
	KeyLen:  Argon2KeyLen,
}

// ToHash 将密码转换为安全的哈希存储格式
// 返回格式: $argon2id$v=19$m=19456,t=2,p=1$<salt>$<hash>
func ToHash(password string) (string, error) {
	return ToHashWithConfig(password, DefaultHashConfig)
}

// ToHashWithConfig 使用自定义配置将密码转换为哈希格式
func ToHashWithConfig(password string, config HashConfig) (string, error) {
	salt, err := GenSalt()
	if err != nil {
		return "", err
	}

	saltBytes, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrDecodeSaltFailed, err)
	}

	hash := argon2.IDKey([]byte(password), saltBytes, config.Time, config.Memory, config.Threads, config.KeyLen)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		config.Memory,
		config.Time,
		config.Threads,
		salt,
		encodedHash,
	), nil
}

// IsCorrect 验证密码是否与存储的哈希匹配
func IsCorrect(password, hashStr string) (bool, error) {
	// 解析哈希字符串
	segments := strings.Split(hashStr, "$")
	if len(segments) != 6 {
		return false, ErrSplitHashInvalid
	}

	// 解析参数 - 修正 parallelism 的类型为 uint8
	var memory, iterations uint32
	var parallelism uint8
	_, err := fmt.Sscanf(segments[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false, fmt.Errorf("failed to parse hash parameters: %w", err)
	}

	salt, err := base64.RawStdEncoding.DecodeString(segments[4])
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrDecodeSaltFailed, err)
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(segments[5])
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrDecodeHashFailed, err)
	}

	// 使用相同参数计算哈希 - parallelism 现在是正确的 uint8 类型
	computedHash := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		parallelism,
		uint32(len(decodedHash)),
	)

	// 常量时间比较防止时序攻击
	return subtle.ConstantTimeCompare(computedHash, decodedHash) == 1, nil
}

// SplitHash 为向后兼容提供的函数 拆分旧格式的哈希字符串
func SplitHash(hash string) []string {
	// 检查是否为新格式
	if strings.HasPrefix(hash, "$argon2id$") {
		segments := strings.Split(hash, "$")
		if len(segments) >= 6 {
			return []string{segments[4], segments[5]}
		}
	}
	// 旧格式处理
	return strings.Split(hash, ":")
}
