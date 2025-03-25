/******************************************************************************
 * Copyright (c) 2025. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package snowflake

import (
	"sync"
	"time"

	"github.com/yitter/idgenerator-go/idgen"
)

var (
	// 默认的WorkerId
	defaultWorkerId uint16 = 1

	// 默认基准时间：2025年1月1日
	defaultBaseTime = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	// 确保只初始化一次
	initialized bool = false
	initMutex   sync.Mutex
)

// InitWith 通过指定的WorkerId初始化ID生成器
// 必须在调用QuickID/QuickID2前调用此函数
// 如果baseTime为nil，则使用默认的2025年1月1日作为基准时间
func InitWith(workerId uint16, baseTime *time.Time) {
	initMutex.Lock()
	defer initMutex.Unlock()

	if initialized {
		return // 已经初始化过，忽略后续调用
	}

	// 如果没有传入基准时间，使用默认时间
	actualBaseTime := defaultBaseTime
	if baseTime != nil {
		actualBaseTime = *baseTime
	}

	// 转换为毫秒时间戳
	baseTimeMs := actualBaseTime.UnixMilli()

	// 创建ID生成器配置，使用传入的workerId
	options := idgen.NewIdGeneratorOptions(workerId)
	options.BaseTime = baseTimeMs

	// 设置全局生成器
	idgen.SetIdGenerator(options)
	initialized = true
}

// init 包初始化时使用默认WorkerId和默认基准时间初始化
func init() {
	InitWith(defaultWorkerId, nil)
}

// ID 生成ID
func ID() uint64 {
	return uint64(idgen.NextId())
}

// QuickID 生成ID
func QuickID() uint64 {
	return uint64(idgen.NextId())
}

// QuickID2 为了保持API兼容性，也使用同一个生成器
func QuickID2() uint64 {
	return uint64(idgen.NextId())
}
