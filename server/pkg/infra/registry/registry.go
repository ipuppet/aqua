// Package registry provides service registration and discovery functionality using etcd.
package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/ipuppet/aqua/server/pkg/infra/models"
)

// ServiceInstance represents a registered service instance with its connection details.
type ServiceInstance struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Endpoint string            `json:"endpoint"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// Client provides service registration and discovery backed by etcd.
type Client struct {
	logger          *slog.Logger
	client          *clientv3.Client
	leaseID         clientv3.LeaseID
	keepAliveCancel context.CancelFunc
}

type Option func(*Client)

func WithLogger(logger *slog.Logger) Option {
	return func(r *Client) {
		r.logger = logger
	}
}

const (
	etcdConnectMaxRetries = 3
	etcdConnectTimeout    = 5 * time.Second
	etcdConnectBaseDelay  = time.Second
)

// NewClient creates a Client using the provided app config.
// etcd 连通性验证会重试最多 3 次，每次带 5s 超时，间隔指数退避 (1s/2s/4s)。
func NewClient(etcdCfg models.EtcdConfig, opts ...Option) (*Client, error) {
	endpoints := append([]string{}, etcdCfg.Endpoints...)
	if len(endpoints) == 0 {
		return nil, fmt.Errorf("no etcd endpoints configured")
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %v", err)
	}

	// 主动验证 etcd 可达性，避免 clientv3.New 成功但实际无法连通的情况。
	// 重试最多 3 次，指数退避，以应对启动期间 etcd 短暂不可用。
	var syncErr error
	delay := etcdConnectBaseDelay
	for attempt := 1; attempt <= etcdConnectMaxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), etcdConnectTimeout)
		syncErr = cli.Sync(ctx)
		cancel()
		if syncErr == nil {
			break
		}
		if attempt < etcdConnectMaxRetries {
			time.Sleep(delay)
			delay *= 2
		}
	}
	if syncErr != nil {
		_ = cli.Close()
		return nil, fmt.Errorf("etcd unreachable after %d attempts: %w", etcdConnectMaxRetries, syncErr)
	}

	r := &Client{client: cli}

	for _, opt := range opts {
		opt(r)
	}

	return r, nil
}

// Register registers a service instance in etcd with a lease and starts a background keep-alive goroutine.
func (r *Client) Register(ctx context.Context, service ServiceInstance) error {
	// 创建租约
	leaseResp, err := r.client.Grant(ctx, 30) // 30秒租约
	if err != nil {
		return fmt.Errorf("failed to grant lease: %v", err)
	}

	r.leaseID = leaseResp.ID

	// 序列化服务信息
	serviceData, err := json.Marshal(service)
	if err != nil {
		return fmt.Errorf("failed to marshal service data: %v", err)
	}

	// 注册服务
	key := fmt.Sprintf("/services/%s/%s", service.Name, service.ID)
	_, err = r.client.Put(ctx, key, string(serviceData), clientv3.WithLease(r.leaseID))
	if err != nil {
		return fmt.Errorf("failed to register service: %v", err)
	}

	r.logger.Info("service registered", "service", service.Name, "key", key)

	// 保持租约活跃
	if r.keepAliveCancel != nil {
		r.keepAliveCancel()
	}

	kaCtx, kaCancel := context.WithCancel(context.Background())
	r.keepAliveCancel = kaCancel

	go r.keepAlive(kaCtx)

	return nil
}

// keepAlive continuously maintains the etcd lease until the background context is canceled.
func (r *Client) keepAlive(ctx context.Context) {
	keepAliveCh, err := r.client.KeepAlive(ctx, r.leaseID)
	if err != nil {
		r.logger.Error("failed to keep lease alive", "error", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			r.logger.Info("keepAlive context canceled", "leaseID", r.leaseID)
			return
		case ka, ok := <-keepAliveCh:
			if !ok || ka == nil {
				r.logger.Info("keepAlive channel closed", "leaseID", r.leaseID)
				return
			}
		}
	}
}

// Deregister removes a service instance from etcd by revoking its lease.
func (r *Client) Deregister(ctx context.Context, serviceName, serviceID string) error {
	if r.keepAliveCancel != nil {
		r.keepAliveCancel()
		r.keepAliveCancel = nil
	}

	// 撤销租约，etcd 会自动删除所有关联的 key
	_, err := r.client.Revoke(ctx, r.leaseID)
	if err != nil {
		r.logger.Error("failed to revoke lease", "leaseID", r.leaseID, "error", err)
		return fmt.Errorf("failed to revoke lease: %v", err)
	}

	r.logger.Info("service deregistered by revoking lease", "service", serviceName, "serviceID", serviceID, "leaseID", r.leaseID)
	r.leaseID = 0
	return nil
}

// Close stops active keep-alive operations and closes the etcd client connection.
func (r *Client) Close() error {
	if r.keepAliveCancel != nil {
		r.keepAliveCancel()
		r.keepAliveCancel = nil
	}
	return r.client.Close()
}

// GetService retrieves a list of active ServiceInstance matching the given service name.
func (r *Client) GetService(ctx context.Context, serviceName string) ([]ServiceInstance, error) {
	prefix := fmt.Sprintf("/services/%s/", serviceName)

	resp, err := r.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %v", err)
	}

	var services []ServiceInstance
	for _, kv := range resp.Kvs {
		var service ServiceInstance
		if err := json.Unmarshal(kv.Value, &service); err != nil {
			r.logger.Warn("failed to unmarshal service data", "error", err)
			continue
		}
		services = append(services, service)
	}

	return services, nil
}
