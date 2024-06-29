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

// GlobalLocker 全局锁变量掌控，不使用注入方式，只要调用包立刻初始化
var GlobalLocker *Locker

func init() {
	GlobalLocker = newGlobalLocker()
}

type timedMutex struct {
	mutex     sync.Mutex
	lastUsed  int64
	createdAt int64 // 记录锁创建的时间 增加一个多少时间无用的锁，才进行内存清理
}

type Locker struct {
	mu sync.Map
}

// 不允许外部使用  只允许通过package调用的方法
func newGlobalLocker() *Locker {
	return &Locker{}
}

func (l *Locker) Lock(name string) *sync.Mutex {
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
	return &t.mutex
}

func (l *Locker) Unlock(lock *sync.Mutex) {
	lock.Unlock()
}

func (l *Locker) cleanUp(threshold int64, minExistTime int64) {
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

func SetLockerAutoCleanup(threshold int64, minExistTime int64) {
	ticker := time.NewTicker(180 * time.Second)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			GlobalLocker.cleanUp(threshold, minExistTime)
		}
	}()
}
