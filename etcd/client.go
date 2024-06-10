package etcd

import (
	"context"
	"fmt"
	"time"

	cliv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type Client struct {
	cli     *cliv3.Client
	timeout time.Duration
}

// Connect 建立连接
func (c *Client) Connect(endpoints []string, dialTimeout time.Duration) error {
	var err error
	c.cli, err = cliv3.New(cliv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to etcd: %w", err)
	}
	return nil
}

// Close 关闭连接
func (c *Client) Close() error {
	if c.cli != nil {
		return c.cli.Close()
	}
	return nil
}

func (c *Client) withTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// Put 放置/更新
func (c *Client) Put(key, value string, timeout time.Duration) error {
	ctx, cancel := c.withTimeout(timeout)
	defer cancel()

	_, err := c.cli.Put(ctx, key, value)
	if err != nil {
		return fmt.Errorf("failed to put key-value: %w", err)
	}
	return nil
}

// Get 读取
func (c *Client) Get(key string, timeout time.Duration) (string, error) {
	ctx, cancel := c.withTimeout(timeout)
	defer cancel()

	resp, err := c.cli.Get(ctx, key)
	if err != nil {
		return "", fmt.Errorf("failed to get key: %w", err)
	}

	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("key not found")
	}
	return string(resp.Kvs[0].Value), nil
}

// Txn 事务示例代码
func (c *Client) Txn() error {
	//ctx, cancel := c.withTimeout(5 * time.Second)
	//defer cancel()
	//
	//txn := c.cli.Txn(ctx)
	//txn.If(
	//	cliv3.Compare(cliv3.Value("/example/key"), "=", "value1"),
	//).Then(
	//	cliv3.OpPut("/example/key", "value2"),
	//).Else(
	//	cliv3.OpPut("/example/key", "value3"),
	//)
	//_, err := txn.Commit()
	//return err
	return nil
}

// CreateLease 租约管理
func (c *Client) CreateLease(ttl int64) (cliv3.LeaseID, error) {
	ctx, cancel := c.withTimeout(5 * time.Second)
	defer cancel()

	leaseResp, err := c.cli.Grant(ctx, ttl)
	if err != nil {
		return 0, err
	}
	return leaseResp.ID, nil
}

// KeepAliveLease 保持租约激活
func (c *Client) KeepAliveLease(id cliv3.LeaseID) error {
	ch, err := c.cli.KeepAlive(context.Background(), id)
	if err != nil {
		return err
	}
	go func() {
		for {
			<-ch
		}
	}()
	return nil
}

// AcquireLock 分布式锁
func (c *Client) AcquireLock(lockName string) (*concurrency.Mutex, *concurrency.Session, error) {
	s, err := concurrency.NewSession(c.cli)
	if err != nil {
		return nil, nil, err
	}
	m := concurrency.NewMutex(s, lockName)
	err = m.Lock(context.Background())
	if err != nil {
		return nil, nil, err
	}
	return m, s, nil
}

// ReleaseLock 释放锁
func (c *Client) ReleaseLock(m *concurrency.Mutex, s *concurrency.Session) error {
	err := m.Unlock(context.Background())
	if err != nil {
		return err
	}
	return s.Close()
}

// ListMembers 成员管理
func (c *Client) ListMembers() error {
	ctx, cancel := c.withTimeout(5 * time.Second)
	defer cancel()

	resp, err := c.cli.MemberList(ctx)
	if err != nil {
		return err
	}
	for _, member := range resp.Members {
		fmt.Printf("Member: %v\n", member)
	}
	return nil
}

// RegisterService 注册服务
func (c *Client) RegisterService(serviceName, serviceAddr string, ttl int64) (cliv3.LeaseID, error) {
	// 创建租约
	leaseID, err := c.CreateLease(ttl)
	if err != nil {
		return 0, fmt.Errorf("failed to create lease: %w", err)
	}

	// 存储带有租约的键值对
	ctx, cancel := c.withTimeout(5 * time.Second)
	defer cancel()

	key := fmt.Sprintf("/services/%s", serviceName)
	_, err = c.cli.Put(ctx, key, serviceAddr, cliv3.WithLease(leaseID))
	if err != nil {
		return 0, fmt.Errorf("failed to register service: %w", err)
	}

	// 保持租约
	err = c.KeepAliveLease(leaseID)
	if err != nil {
		return 0, fmt.Errorf("failed to keep alive lease: %w", err)
	}

	return leaseID, nil
}

// Watch 监听 withPrefix用来监听服务:"/services/xxx"
func (c *Client) Watch(key string, withPrefix bool) {
	var rch cliv3.WatchChan
	if withPrefix {
		rch = c.cli.Watch(context.Background(), key, cliv3.WithPrefix())
	} else {
		rch = c.cli.Watch(context.Background(), key)
	}

	for watchResponse := range rch {
		for _, ev := range watchResponse.Events {
			fmt.Printf("Type: %s Key: %q Value: %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}

// GetService 获取服务
func (c *Client) GetService(serviceName string) (map[string]string, error) {
	ctx, cancel := c.withTimeout(5 * time.Second)
	defer cancel()

	resp, err := c.cli.Get(ctx, fmt.Sprintf("/services/%s", serviceName), cliv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	services := make(map[string]string)
	for _, kv := range resp.Kvs {
		services[string(kv.Key)] = string(kv.Value)
	}

	return services, nil
}
