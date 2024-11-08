package os

import (
	"fmt"
)

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
	TB = 1024 * GB
)

// ByteSize 字节大小转换辅助函数
type ByteSize int64

func (b ByteSize) ToKB() float64 {
	return float64(b) / KB
}

func (b ByteSize) ToMB() float64 {
	return float64(b) / MB
}

func (b ByteSize) ToGB() float64 {
	return float64(b) / GB
}

func (b ByteSize) String() string {
	switch {
	case b >= TB:
		return fmt.Sprintf("%.2f TB", b.ToGB()/1024)
	case b >= GB:
		return fmt.Sprintf("%.2f GB", b.ToGB())
	case b >= MB:
		return fmt.Sprintf("%.2f MB", b.ToMB())
	case b >= KB:
		return fmt.Sprintf("%.2f KB", b.ToKB())
	default:
		return fmt.Sprintf("%d B", b)
	}
}
