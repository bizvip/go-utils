package blake3hash

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/zeebo/blake3"
)

// BufSize is the default buffer size (4 MiB). It can be modified (globally)
// before calling, depending on the medium and scenario.
// 1–4 MiB usually provides a good balance between throughput and syscall overhead.
var BufSize = 4 * 1024 * 1024

// SumFile calculates the BLAKE3-256 hash of a specified file in a streaming
// manner and returns a lowercase hexadecimal string.
func SumFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return SumReader(f)
}

// SumReader performs a streaming BLAKE3-256 calculation on any io.Reader
// and returns a lowercase hexadecimal string.
// Note: If the upstream Reader has its own buffering or rate limiting,
// adjust BufSize as needed for better throughput.
func SumReader(r io.Reader) (string, error) {
	// Defensive fallback: prevent BufSize from being accidentally set to a value
	// that is too small or non-positive.
	bufSize := BufSize
	if bufSize <= 0 {
		bufSize = 4 * 1024 * 1024
	}

	h := blake3.New() // Default is 256-bit
	br := bufio.NewReaderSize(r, bufSize)
	buf := make([]byte, bufSize)

	for {
		n, er := br.Read(buf)
		if n > 0 {
			if _, werr := h.Write(buf[:n]); werr != nil {
				return "", werr
			}
		}
		if er != nil {
			if er == io.EOF {
				break
			}
			return "", er
		}
	}

	sum := h.Sum(nil) // 32 bytes
	dst := make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(dst, sum)
	return string(dst), nil
}

// SumBytes calculates the BLAKE3-256 hash of data in memory and returns a
// lowercase hexadecimal string.
// Suitable for small objects or scenarios where data is already in memory;
// SumReader is still recommended for large objects to avoid extra copying.
func SumBytes(b []byte) string {
	h := blake3.New()
	_, _ = h.Write(b)
	sum := h.Sum(nil)
	dst := make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(dst, sum)
	return string(dst)
}

// Must is a helper function: use it when you are certain that no error will
// occur or wish to simplify error handling at the call site.
// For example: digest := blake3hash.Must(blake3hash.SumFile(path))
func Must(hexDigest string, err error) string {
	if err != nil {
		panic(fmt.Errorf("blake3hash: %w", err))
	}
	return hexDigest
}
