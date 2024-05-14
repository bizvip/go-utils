//go:build darwin

package fs

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

func GetFileCreationTime(filePath string) (string, time.Time, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", time.Time{}, err
	}

	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return "", time.Time{}, fmt.Errorf("failed to get raw syscall.Stat_t data")
	}

	unixTime := time.Unix(int64(stat.Ctimespec.Sec), int64(stat.Ctimespec.Nsec))
	formattedTime := unixTime.Format("2006-01-02 15:04:05")

	return formattedTime, unixTime, nil
}

// StartsWithDot 检查文件名是否以 '.' 开头
func StartsWithDot(fileName string) bool {
	return len(fileName) > 0 && fileName[0] == '.'
}

func GetFileNameMd5(filename string) (string, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return "", err
	}

	// 计算文件名的 MD5 值
	hash := md5.Sum([]byte(fileInfo.Name()))
	md5Value := hex.EncodeToString(hash[:])

	return md5Value, nil
}

func GetFileMd5(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:]), nil
}

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

func GetCurExeDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}

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

// IsDirAndHasFiles 是否dir 是否有内容 错误
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

func Delete(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}
