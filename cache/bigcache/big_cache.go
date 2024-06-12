/******************************************************************************
 * Copyright (c) Archer++ 2024.                                               *
 ******************************************************************************/

package bigcache

import (
	"context"
	"log"
	"time"

	"github.com/allegro/bigcache/v3"
)

var Cache *bigcache.BigCache

type bigCache struct{}

func InitConfig(config *bigcache.Config) {
	if config == nil {
		config = &bigcache.Config{
			Shards:     1024,
			LifeWindow: 60 * 24 * 365 * time.Minute, //100年
			// 含义：清理间隔时间
			// 作用：设置多长时间清理一次过期条目。如果设置为 <= 0，则不会进行自动清理
			// 设置：例如 5 * time.Minute 表示每 5 分钟清理一次过期条目
			CleanWindow: 0 * time.Minute,
			// 含义：在生命周期窗口内的最大条目数
			// 作用：用于初始内存分配的估算
			MaxEntriesInWindow: 1000 * 10 * 60,
			// 含义：单个条目的最大字节大小
			// 作用：用于初始内存分配的估算
			// 设置：例如 500 表示每个条目的最大大小为 500 字节
			MaxEntrySize: 128 * 1024,
			// 是否打印调试信息
			Verbose: true,
			// 含义：缓存的硬最大内存大小（单位：MB）
			// 作用：设置缓存使用的最大内存限制。一旦达到这个限制，新的条目会覆盖最旧的条目
			// 设置：例如 8192 表示最大内存大小为 8192 MB（8 GB）0代表无限制
			HardMaxCacheSize: 16384,
			// 含义：当条目因过期或无空间而被移除时的回调函数
			// 作用：当条目被删除时触发回调，便于执行额外操作
			// 设置：例如 nil 表示没有回调
			OnRemove: nil,
			// OnRemoveWithReason 是一个回调函数，当最旧的条目因到期时间、没有空间存放新条目或调用了删除操作而被移除时触发
			// 将传递一个表示原因的常量
			// 默认值为 nil，表示没有回调，并且阻止解包最旧的条目
			// 如果指定了 OnRemove，此设置将被忽略
			OnRemoveWithReason: nil,
		}
	}
	var err error
	Cache, err = bigcache.New(context.Background(), *config)
	if err != nil {
		log.Fatal(err)
	}
}
