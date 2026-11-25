package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"aqua/pkg/infra/models"
	"aqua/pkg/infra/registry"
)

// RegisterService registers the gRPC service into etcd with the gRPC listen address.
// 仅在有 gRPC 服务时由 bootstrap.Run 调用；Endpoint 存储 gRPC 地址供客户端发现。
// 注册逻辑也会绑定到 [CleanupQueue] 中，确保进程结束时可以自动从注册中心下线。
func RegisterService(serviceName string, env Env, etcdCfg models.EtcdConfig, cleanup *CleanupQueue, logger *slog.Logger) error {
	reg, err := registry.NewClient(etcdCfg, registry.WithLogger(logger))
	if err != nil {
		return fmt.Errorf("failed to create registry: %v", err)
	}

	grpcAddr := env.GRPCAddress
	host, port, err := net.SplitHostPort(grpcAddr)
	if err != nil {
		return fmt.Errorf("invalid gRPC listen addr %q: %w", grpcAddr, err)
	}
	if host == "" {
		return errors.New("gRPC listen addr is missing host")
	}
	if port == "" {
		return errors.New("gRPC listen addr is missing port")
	}

	service := registry.ServiceInstance{
		ID:       fmt.Sprintf("%s-%s", serviceName, grpcAddr),
		Name:     serviceName,
		Endpoint: grpcAddr,
		Metadata: map[string]string{
			"version": "1.0.0",
			"proto":   "grpc",
		},
	}

	// 限制注册超时时间，防止初始化阶段因 etcd 不可用阻塞主进程启动
	registerCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := reg.Register(registerCtx, service); err != nil {
		if closeErr := reg.Close(); closeErr != nil && logger != nil {
			logger.Error("failed to close registry client", slog.Any("error", closeErr))
		}
		return err
	}

	// 将下线操作推入全局队列中，实现在收到退出信号时自动抛出服务下线请求。
	cleanup.Add(func() {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := reg.Deregister(shutdownCtx, service.Name, service.ID); err != nil && logger != nil {
			logger.Error("failed to deregister service", slog.Any("error", err))
		}

		if err := reg.Close(); err != nil && logger != nil {
			logger.Error("failed to close registry client", slog.Any("error", err))
		}
		if logger != nil {
			logger.Info("service deregistered and registry client closed")
		}
	})

	return nil
}
