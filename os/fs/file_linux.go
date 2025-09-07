//go:build linux

package fs

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// ComputeFileSHA256 计算文件的SHA-256哈希
func ComputeFileSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) { _ = file.Close() }(file)

	hash := sha256.New()
	buffer := make([]byte, 1024*1024) // 1MB 缓冲区
	if _, err := io.CopyBuffer(hash, file, buffer); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	return hashString, nil
}

// GetFileCreationTime 获取文件创建时间
func GetFileCreationTime(filePath string) (string, time.Time, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", time.Time{}, err
	}

	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return "", time.Time{}, fmt.Errorf("failed to get raw syscall.Stat_t data")
	}

	unixTime := time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))
	formattedTime := unixTime.Format("2006-01-02 15:04:05")

	return formattedTime, unixTime, nil
}

// StartsWithDot 检查文件名是否以 '.' 开头
func StartsWithDot(fileName string) bool {
	return len(fileName) > 0 && fileName[0] == '.'
}

// GetFileNameMd5 计算文件名的 MD5 值
func GetFileNameMd5(filename string) (string, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		log.Error().Err(err).Msg("Error getting file info")
		return "", err
	}

	hash := md5.Sum([]byte(fileInfo.Name()))
	md5Value := hex.EncodeToString(hash[:])

	return md5Value, nil
}

// GetSmallFileMd5 计算文件的 MD5 值
func GetFileMd5(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:]), nil
}

// GetBigFileMd5 通过流的方式计算文件的 MD5 值
func GetFileMd5Stream(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) { _ = file.Close() }(file)

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashValue := hash.Sum(nil)
	hashString := hex.EncodeToString(hashValue)
	return hashString, nil
}

// GetCurExeDir 获取当前执行文件所在目录
func GetCurExeDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}

// GetAllFilesByExt 根据文件扩展名获取所有文件
func GetAllFilesByExt(dir string, ext string) ([]string, error) {
	var files []string
	err := filepath.Walk(
		dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == ext {
				files = append(files, path)
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error walking the path %v: %v", dir, err)
	}
	return files, nil
}

// IsDirAndHasFiles 检查目录是否存在并且包含文件
func IsDirAndHasFiles(dirPath string) (bool, bool, error) {
	info, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, false, nil
		}
		return false, false, err
	}

	if !info.IsDir() {
		return false, false, nil
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return true, false, err
	}

	return true, len(files) > 0, nil
}

// Delete 删除指定路径的文件或目录
func Delete(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

// CreateDirIfNotExist 在当前执行文件所在目录下创建指定的目录，如果目录不存在
func CreateDirIfNotExist(relativePath string) error {
	// 获取当前执行文件所在目录
	curDir := GetCurExeDir()

	// 拼接完整的目录路径
	fullPath := filepath.Join(curDir, relativePath)

	// 检查目录是否已经存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		err := os.MkdirAll(fullPath, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}

	return nil
}

// IsFile 路径是否是个文件
func IsFile(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return !info.IsDir(), nil
}

// CreateDir 在指定路径下创建目录，如果目录不存在
func CreateDir(dirPath string) error {
	var err error
	// 检查目录是否已经存在
	if _, err = os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}
	return nil
}

// DetectFileType 读取文件的前 512 个字节并检测文件的 MIME 类型
func DetectFileType(file io.Reader) (string, error) {
	// 创建缓冲区以存储前 512 个字节
	buffer := make([]byte, 512)
	// 读取前 512 个字节
	if _, err := file.Read(buffer); err != nil {
		return "", err
	}
	// 重置文件读取位置
	if seeker, ok := file.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	}
	// 检测文件类型
	mimeType := http.DetectContentType(buffer)
	return mimeType, nil
}
