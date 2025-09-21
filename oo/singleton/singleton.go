// Package singleton 提供“只初始化一次”的通用工具。
// 兼容 Go 1.21+（推荐使用 1.25 或更新版本）。
package singleton

import (
	"sync"
)

// -----------------------------
// 1) 全局惰性单例（零参数工厂）
// -----------------------------

// Lazy 返回一个闭包 get()，第一次调用时执行 initFn 并缓存其返回值，之后重复返回同一实例。
// 线程安全；若 initFn panic，后续调用也会以相同的 panic 失败（与 sync.OnceValue 语义一致）。
func Lazy[T any](initFn func() T) func() T {
	// Go 1.21+ 提供的 OnceValue：将“只初始化一次并返回值”的模式标准化。
	return sync.OnceValue(initFn)
}

// LazyPtr 是 Lazy 的指针版（常见于需要 *Client/*Service 的场景）。
func LazyPtr[T any](initFn func() *T) func() *T {
	return sync.OnceValue(initFn)
}

// -----------------------------
// 2) 按 Key 的单例（多配置/多租户）
//    - 并发首次访问同一 Key 时只初始化一次
//    - 无第三方依赖，内置 singleflight
// -----------------------------

// PerKey 管理“按键单例”，例如：按 DSN 建立并复用连接客户端。
type PerKey[K comparable, V any] struct {
	mu   sync.Mutex
	objs map[K]V

	inflightMu sync.Mutex
	inflight   map[K]*call[V]
}

type call[V any] struct {
	wg  sync.WaitGroup
	val V
	err error
}

// NewPerKey 创建一个按键单例容器。
func NewPerKey[K comparable, V any]() *PerKey[K, V] {
	return &PerKey[K, V]{
		objs:     make(map[K]V),
		inflight: make(map[K]*call[V]),
	}
}

// Get 返回 key 对应的实例；若不存在，则用 initFn(key) 初始化一次并缓存。
// 并发同一时刻对同一 key 的首次请求只会触发一次 initFn。
func (p *PerKey[K, V]) Get(key K, initFn func(K) (V, error)) (V, error) {
	// 1) 快速路径：已存在
	p.mu.Lock()
	if v, ok := p.objs[key]; ok {
		p.mu.Unlock()
		return v, nil
	}
	p.mu.Unlock()

	// 2) 并发抑制：同一 key 的初始化只做一次
	p.inflightMu.Lock()
	if c, ok := p.inflight[key]; ok {
		p.inflightMu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := &call[V]{}
	c.wg.Add(1)
	p.inflight[key] = c
	p.inflightMu.Unlock()

	// 3) 真正初始化
	v, err := initFn(key)
	if err == nil {
		p.mu.Lock()
		p.objs[key] = v
		p.mu.Unlock()
	}

	// 4) 广播结果并清理 inflight
	c.val, c.err = v, err
	c.wg.Done()

	p.inflightMu.Lock()
	delete(p.inflight, key)
	p.inflightMu.Unlock()

	return v, err
}

// Has 报告是否已存在 key 对应的实例（不触发初始化）。
func (p *PerKey[K, V]) Has(key K) bool {
	p.mu.Lock()
	_, ok := p.objs[key]
	p.mu.Unlock()
	return ok
}

// Delete 删除已缓存的 key 实例（不做资源关闭；若需要请在外部先 Close 再删）。
func (p *PerKey[K, V]) Delete(key K) {
	p.mu.Lock()
	delete(p.objs, key)
	p.mu.Unlock()
}

// Range 遍历已缓存的实例（只读回调）。
func (p *PerKey[K, V]) Range(fn func(key K, val V)) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for k, v := range p.objs {
		fn(k, v)
	}
}
