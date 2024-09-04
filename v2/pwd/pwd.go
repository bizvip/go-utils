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
	Argon2Time    = 4         // Argon2 参数，迭代次数
	Argon2Memory  = 64 * 1024 // Argon2 参数，内存大小（KB）
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

// GenerateSalt 生成一个随机盐值
func GenerateSalt() (string, error) {
	salt := make([]byte, SaltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", fmt.Errorf("%w: %v", ErrGenerateSaltFailed, err)
	}
	return base64.RawStdEncoding.EncodeToString(salt), nil
}

// HashPassword 使用 Argon2id 对密码进行加盐哈希
func HashPassword(password, salt string) (string, error) {
	saltBytes, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrDecodeSaltFailed, err)
	}

	hash := argon2.IDKey([]byte(password), saltBytes, Argon2Time, Argon2Memory, Argon2Threads, Argon2KeyLen)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)
	return fmt.Sprintf("%s$%s", salt, encodedHash), nil
}

// IsPasswordCorrect 用于比较输入的密码和存储的散列值是否匹配
func IsPasswordCorrect(password, hashStr string) bool {
	parts := SplitHash(hashStr)
	if len(parts) != 2 {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		fmt.Printf("%v: %v\n", ErrDecodeSaltFailed, err)
		return false
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		fmt.Printf("%v: %v\n", ErrDecodeHashFailed, err)
		return false
	}

	hash := argon2.IDKey([]byte(password), salt, Argon2Time, Argon2Memory, Argon2Threads, Argon2KeyLen)

	return subtle.ConstantTimeCompare(hash, expectedHash) == 1
}

// SplitHash 分割存储的散列值和盐值
func SplitHash(hash string) []string {
	return strings.Split(hash, "$")
}
