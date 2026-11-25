package bootstrap

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

// Env holds startup parameters parsed from flags.
// 这些参数需要在配置文件加载之前就确定，否则会形成循环依赖.
type Env struct {
	ConfigPath  string // 配置文件目录
	StorePath   string // 存储目录
	HTTPAddress string // HTTP 监听地址
	GRPCAddress string // 内部 gRPC 监听地址
}

// ParseFlags parses command-line flags and returns Env.
func ParseFlags() (Env, error) {
	var env Env
	configDefault := envFirst("AQUA_CONFIG_DIR", "./config")
	storeDefault := envFirst("AQUA_STORE_DIR", "./storage")
	httpAddrDefault := envFirst("AQUA_HTTP_ADDR", "0.0.0.0:8080")
	grpcAddrDefault := envFirst("AQUA_GRPC_ADDR", "0.0.0.0:18080")

	flag.StringVar(&env.ConfigPath, "config-dir", configDefault, "Config directory path.")
	flag.StringVar(&env.StorePath, "store-dir", storeDefault, "Store directory path.")
	flag.StringVar(&env.HTTPAddress, "http-addr", httpAddrDefault, "HTTP listen address.")
	flag.StringVar(&env.GRPCAddress, "grpc-addr", grpcAddrDefault, "Internal gRPC listen address.")

	flag.Parse()

	if err := validateEnv(env); err != nil {
		return Env{}, err
	}
	return env, nil
}

func validateEnv(env Env) error {
	if strings.TrimSpace(env.ConfigPath) == "" {
		return fmt.Errorf("config path cannot be empty")
	}
	if strings.TrimSpace(env.StorePath) == "" {
		return fmt.Errorf("store path cannot be empty")
	}
	if _, _, err := net.SplitHostPort(env.HTTPAddress); err != nil {
		return fmt.Errorf("invalid http-addr %q: %w", env.HTTPAddress, err)
	}
	if _, _, err := net.SplitHostPort(env.GRPCAddress); err != nil {
		return fmt.Errorf("invalid grpc-addr %q: %w", env.GRPCAddress, err)
	}
	return nil
}

func envFirst(keys ...string) string {
	for i := 0; i < len(keys)-1; i++ {
		if value, ok := os.LookupEnv(keys[i]); ok && strings.TrimSpace(value) != "" {
			return value
		}
	}
	return keys[len(keys)-1]
}
