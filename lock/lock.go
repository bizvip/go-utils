/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package lock

import (
	"hash/fnv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultShards   = 32
	shardMask       = defaultShards - 1
	concurrentLimit = 1000 // 高并发下自动切换到分片锁的阈值
)

var (
	Adaptive *AdaptiveLock
	initOnce sync.Once
)

func init() {
	initOnce.Do(
		func() {
			Adaptive = newLocker()
			SetLockerAutoCleanup(30*60, 60*60)
		},
	)
}

type timedMutex struct {
	mutex     sync.Mutex
	lastUsed  atomic.Int64
	createdAt int64
}

type lockShard struct {
	sync.RWMutex
	items sync.Map
}

// AdaptiveLock 自适应锁管理器
type AdaptiveLock struct {
	shards    [defaultShards]*lockShard
	simple    sync.Map     // 用于低并发场景
	useShards atomic.Bool  // 是否使用分片模式
	lockCount atomic.Int64 // 活跃锁计数
}

func newLocker() *AdaptiveLock {
	al := &AdaptiveLock{}
	// 初始化分片
	for i := 0; i < defaultShards; i++ {
		al.shards[i] = &lockShard{}
	}
	// 默认使用简单模式
	al.useShards.Store(false)
	return al
}

func (al *AdaptiveLock) getShard(name string) *lockShard {
	h := fnv.New32a()
	_, _ = h.Write([]byte(name))
	return al.shards[h.Sum32()&shardMask]
}

// Lock 锁定一个资源
func (al *AdaptiveLock) Lock(name string) {
	// 更新锁计数并检查是否需要切换模式
	currentCount := al.lockCount.Add(1)
	if currentCount > concurrentLimit && !al.useShards.Load() {
		al.useShards.Store(true)
	}

	now := time.Now().Unix()
	if al.useShards.Load() {
		// 分片模式
		targetShard := al.getShard(name)
		targetShard.RLock()
		tm, loaded := targetShard.items.Load(name)
		targetShard.RUnlock()

		if !loaded {
			targetShard.Lock()
			tm, _ = targetShard.items.LoadOrStore(
				name, &timedMutex{
					createdAt: now,
				},
			)
			targetShard.Unlock()
		}

		t := tm.(*timedMutex)
		t.mutex.Lock()
		t.lastUsed.Store(now)
	} else {
		// 简单模式
		tm, _ := al.simple.LoadOrStore(
			name, &timedMutex{
				createdAt: now,
			},
		)
		t := tm.(*timedMutex)
		t.mutex.Lock()
		t.lastUsed.Store(now)
	}
}

// Unlock 解锁指定的资源
func (al *AdaptiveLock) Unlock(name string) {
	defer al.lockCount.Add(-1)

	if al.useShards.Load() {
		// 分片模式
		targetShard := al.getShard(name)
		targetShard.RLock()
		if tm, ok := targetShard.items.Load(name); ok {
			t := tm.(*timedMutex)
			t.mutex.Unlock()
		}
		targetShard.RUnlock()
	} else {
		// 简单模式
		if tm, ok := al.simple.Load(name); ok {
			t := tm.(*timedMutex)
			t.mutex.Unlock()
		}
	}

	// 如果锁计数降低到阈值以下，切换回简单模式
	if al.lockCount.Load() < concurrentLimit/2 && al.useShards.Load() {
		al.useShards.Store(false)
	}
}

// cleanUp 定期清理超过指定阈值的未使用锁
func (al *AdaptiveLock) cleanUp(threshold, minExistTime int64) {
	now := time.Now().Unix()

	if al.useShards.Load() {
		// 分片模式清理
		var wg sync.WaitGroup
		for i := 0; i < defaultShards; i++ {
			wg.Add(1)
			go func(currentShard *lockShard) {
				defer wg.Done()
				al.cleanupShard(currentShard, now, threshold, minExistTime)
			}(al.shards[i])
		}
		wg.Wait()
	} else {
		// 简单模式清理
		var toDelete []interface{}
		al.simple.Range(
			func(key, value interface{}) bool {
				tm := value.(*timedMutex)
				if tm.mutex.TryLock() {
					if now-tm.lastUsed.Load() > threshold && now-tm.createdAt > minExistTime {
						toDelete = append(toDelete, key)
					}
					tm.mutex.Unlock()
				}
				return true
			},
		)
		for _, key := range toDelete {
			al.simple.Delete(key)
		}
	}
}

func (al *AdaptiveLock) cleanupShard(shard *lockShard, now, threshold, minExistTime int64) {
	var toDelete []interface{}

	shard.RLock()
	shard.items.Range(
		func(key, value interface{}) bool {
			tm := value.(*timedMutex)
			if tm.mutex.TryLock() {
				if now-tm.lastUsed.Load() > threshold && now-tm.createdAt > minExistTime {
					toDelete = append(toDelete, key)
				}
				tm.mutex.Unlock()
			}
			return true
		},
	)
	shard.RUnlock()

	if len(toDelete) > 0 {
		shard.Lock()
		for _, key := range toDelete {
			shard.items.Delete(key)
		}
		shard.Unlock()
	}
}

func SetLockerAutoCleanup(threshold, minExistTime int64) {
	ticker := time.NewTicker(180 * time.Second)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			Adaptive.cleanUp(threshold, minExistTime)
		}
	}()
}

// GetActiveLockCount 获取当前活跃的锁数量
func (al *AdaptiveLock) GetActiveLockCount() int64 {
	return al.lockCount.Load()
}

// IsShardMode 返回当前是否在使用分片模式
func (al *AdaptiveLock) IsShardMode() bool {
	return al.useShards.Load()
}
