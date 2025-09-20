package etcd

import (
	"context"
	"errors"
	"fmt"
	"grpc-demo-server/util"
	"sync/atomic"

	"github.com/archine/gin-plus/v4/component/config"
	"github.com/archine/gin-plus/v4/component/gplog"
	"github.com/archine/gin-plus/v4/exception"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// Manager etcd客户端管理器
type Manager struct {
	cfg         *conf
	cli         *clientv3.Client
	key         string           // etcd key
	leaseId     clientv3.LeaseID // 租约ID
	unRegSignal chan struct{}    // 注销信号
	unRegFlag   atomic.Bool      // 注销标志
	regFlag     atomic.Bool      // 注册标志
}

func NewManager(cp config.Provider) (*Manager, error) {
	var cfg conf
	if err := cp.Unmarshal("etcd", &cfg); err != nil {
		return nil, err
	}
	cfg.verify()

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: cfg.DialTimeout,
	})
	if err != nil {
		return nil, err
	}

	m := &Manager{
		cli:         cli,
		cfg:         &cfg,
		unRegSignal: make(chan struct{}),
	}
	m.unRegFlag.Store(false)
	m.regFlag.Store(false)

	return m, nil
}

// Register 注册etcd服务
func (m *Manager) Register(ctx context.Context, svc, addr string) error {
	if !m.regFlag.CompareAndSwap(false, true) {
		return errors.New("etcd service already registered")
	}

	m.key = svc + "/" + util.GenerateUUID()

	lease, err := m.cli.Grant(ctx, m.cfg.TTL)
	if err != nil {
		return exception.Wrap(err, "etcd grant failed")
	}
	m.leaseId = lease.ID

	_, err = m.cli.Put(ctx, m.key, fmt.Sprintf(`{"Addr": "%s"}`, addr), clientv3.WithLease(lease.ID))
	if err != nil {
		m.remove()
		return exception.Wrap(err, "etcd put failed")
	}

	// 设置心跳续约
	kaCtx, kaCancel := context.WithCancel(ctx)

	kaCh, err := m.cli.KeepAlive(kaCtx, lease.ID)
	if err != nil {
		kaCancel()
		m.remove()
		return exception.Wrap(err, "etcd keepalive failed")
	}
	// 监听续约情况
	go func() {
		defer func() {
			m.remove()
			gplog.Info("etcd service unregistered")
		}()

		for {
			select {
			case <-m.unRegSignal:
				// 被外部调用了取消函数
				kaCancel()
				return
			case _, ok := <-kaCh:
				if !ok {
					return
				}
			}
		}

	}()

	return nil
}

// Unregister 发送取消注册etcd服务的信号
func (m *Manager) Unregister() {
	if m.unRegFlag.CompareAndSwap(false, true) {
		close(m.unRegSignal) // 触发注销
	}
}

// remove 注销服务
func (m *Manager) remove() {
	_, _ = m.cli.Revoke(context.Background(), m.leaseId) // 心跳失败，撤销租约
	_, _ = m.cli.Delete(context.Background(), m.key)     // 删除注册的服务
	_ = m.cli.Close()
}
