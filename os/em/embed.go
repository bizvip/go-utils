/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package em

import (
	"embed"
	"log"
)

// GetFileByPath 读取文件内容
func GetFileByPath(eb embed.FS, path string) string {
	data, err := eb.ReadFile(path)
	if err != nil {
		log.Fatalf("读取文件失败: %v", err)
	}
	return string(data)
}

// GetFileList 递归读取目录并返回嵌套的 map 结构 !!! 极其复杂嵌套的文件目录影响性能 !!!
func GetFileList(eb embed.FS, dir string) map[string]interface{} {
	// 读取目录
	files, err := eb.ReadDir(dir)
	if err != nil {
		log.Fatalf("读取目录失败: %v", err)
	}

	result := make(map[string]interface{})

	// 遍历目录中的所有文件和子目录
	for _, file := range files {
		if file.IsDir() {
			// 递归处理子目录，嵌套子目录到 map 中
			result[file.Name()] = GetFileList(eb, dir+"/"+file.Name())
		} else {
			// 如果是文件，直接存储文件名
			result[file.Name()] = nil
		}
	}

	return result
}
