/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package lock

import (
	"sync"
	"sync/atomic"
	"time"
)

// Ctrl 全局锁变量掌控，不使用注入方式，只要调用包立刻初始化
var Ctrl *mutexLock

func init() {
	Ctrl = newLocker()
	SetLockerAutoCleanup(30*60, 60*60)
}

type timedMutex struct {
	mutex     sync.Mutex
	lastUsed  int64
	createdAt int64 // 记录锁创建的时间 增加一个多少时间无用的锁，才进行内存清理
}

type mutexLock struct {
	mu sync.Map
}

// 不允许外部使用  只允许通过package初始化调用的方法
func newLocker() *mutexLock {
	return &mutexLock{}
}

// Lock 锁定一个资源，返回该资源的名称作为标识符
func (l *mutexLock) Lock(name string) {
	now := time.Now().Unix()
	tm, _ := l.mu.LoadOrStore(
		name, &timedMutex{
			lastUsed:  now,
			createdAt: now,
		},
	)
	t := tm.(*timedMutex)
	t.mutex.Lock()
	atomic.StoreInt64(&t.lastUsed, now)
}

// Unlock 解锁指定的资源
func (l *mutexLock) Unlock(name string) {
	tm, ok := l.mu.Load(name)
	if !ok {
		return // 如果锁不存在，直接返回
	}
	t := tm.(*timedMutex)
	t.mutex.Unlock()
}

// cleanUp 定期清理超过指定阈值的未使用锁
func (l *mutexLock) cleanUp(threshold int64, minExistTime int64) {
	now := time.Now().Unix()
	l.mu.Range(
		func(key, value interface{}) bool {
			tm := value.(*timedMutex)
			if now-atomic.LoadInt64(&tm.lastUsed) > threshold && now-tm.createdAt > minExistTime {
				l.mu.Delete(key) // 删除条件同时满足最后使用时间和最小存在时间
			}
			return true
		},
	)
}

// SetLockerAutoCleanup 设置自动清理锁的任务
func SetLockerAutoCleanup(threshold int64, minExistTime int64) {
	ticker := time.NewTicker(180 * time.Second)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			Ctrl.cleanUp(threshold, minExistTime)
		}
	}()
}
