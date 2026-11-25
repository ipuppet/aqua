package discovery

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"github.com/ipuppet/aqua/server/pkg/infra/config"
	"github.com/ipuppet/aqua/server/pkg/infra/models"
	"github.com/ipuppet/aqua/server/pkg/infra/registry"
)

const (
	defaultServiceTimeout = 30
	serviceCacheTTL       = 10 * time.Second
)

var (
	randomSource = rand.New(rand.NewSource(time.Now().UnixNano()))
	randomLock   sync.Mutex
)

// ServiceInfo describes the discovered metadata for a given service.
type ServiceInfo struct {
	BaseURL string
	Timeout int
}

// ServiceRegistry describes the discovery contract needed by dependent layers.
type ServiceRegistry interface {
	Close() error
	GetServiceInfo(serviceName string) (ServiceInfo, error)
}

type etcdServiceRegistry struct {
	discovery *registry.Client
	cacheLock sync.RWMutex
	cache     map[string]cachedServiceInfo
}

var _ ServiceRegistry = (*etcdServiceRegistry)(nil)

type cachedServiceInfo struct {
	info      ServiceInfo
	expiresAt time.Time
}

// NewEtcdServiceRegistry initializes a registry backed by etcd discovery.
func NewEtcdServiceRegistry(configProvider config.EtcdConfigProvider, logger *slog.Logger) (ServiceRegistry, error) {
	if configProvider == nil {
		configProvider = func() models.EtcdConfig { return models.EtcdConfig{} }
	}
	etcdCfg := configProvider()
	discoveryClient, err := registry.NewClient(etcdCfg, registry.WithLogger(logger))
	if err != nil {
		return nil, fmt.Errorf("failed to create service discovery client: %w", err)
	}

	return &etcdServiceRegistry{
		discovery: discoveryClient,
		cache:     make(map[string]cachedServiceInfo),
	}, nil
}

func (r *etcdServiceRegistry) Close() error {
	if r.discovery == nil {
		return nil
	}
	return r.discovery.Close()
}

func (r *etcdServiceRegistry) GetServiceInfo(serviceName string) (ServiceInfo, error) {
	if info, ok := r.getCached(serviceName); ok {
		return info, nil
	}

	info, err := r.resolveService(serviceName)
	if err != nil {
		return ServiceInfo{}, err
	}
	r.setCache(serviceName, info)
	return info, nil
}

func (r *etcdServiceRegistry) resolveService(serviceName string) (ServiceInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	instances, err := r.discovery.GetService(ctx, serviceName)
	if err != nil {
		return ServiceInfo{}, fmt.Errorf("failed to query service %s: %w", serviceName, err)
	}
	if len(instances) == 0 {
		return ServiceInfo{}, fmt.Errorf("no available instances for service %s", serviceName)
	}

	instance := pickInstance(instances)
	return ServiceInfo{BaseURL: instance.Endpoint, Timeout: defaultServiceTimeout}, nil
}

func (r *etcdServiceRegistry) getCached(serviceName string) (ServiceInfo, bool) {
	r.cacheLock.RLock()
	defer r.cacheLock.RUnlock()
	entry, ok := r.cache[serviceName]
	if !ok {
		return ServiceInfo{}, false
	}
	if time.Now().After(entry.expiresAt) {
		return ServiceInfo{}, false
	}
	return entry.info, true
}

func (r *etcdServiceRegistry) setCache(serviceName string, info ServiceInfo) {
	r.cacheLock.Lock()
	defer r.cacheLock.Unlock()
	r.cache[serviceName] = cachedServiceInfo{
		info:      info,
		expiresAt: time.Now().Add(serviceCacheTTL),
	}
}

func pickInstance(instances []registry.ServiceInstance) registry.ServiceInstance {
	if len(instances) == 1 {
		return instances[0]
	}
	randomLock.Lock()
	defer randomLock.Unlock()
	idx := randomSource.Intn(len(instances))
	return instances[idx]
}
