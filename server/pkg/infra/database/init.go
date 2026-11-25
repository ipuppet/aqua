package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"aqua/pkg/infra/models"
)

const (
	dbConnectMaxRetries = 3
	dbConnectTimeout    = 5 * time.Second
	dbConnectBaseDelay  = time.Second
)

// Connect 根据配置参数构建连接字符串并初始化数据库连接。
// 连接验证会重试最多 3 次，每次带 5s 超时，间隔指数退避 (1s/2s/4s)。
func Connect(cfg models.DatabaseConfig) (*sql.DB, error) {
	parts := []string{
		"user=" + cfg.Username,
		"password=" + cfg.Password,
		"dbname=" + cfg.Database,
		"host=" + cfg.Host,
		"port=" + cfg.Port,
	}

	if cfg.SSLMode != "" {
		parts = append(parts, "sslmode="+cfg.SSLMode)

		if cfg.SSLMode == "verify-ca" || cfg.SSLMode == "verify-full" {
			if cfg.SSLRootCert == "" {
				return nil, fmt.Errorf("sslmode is %s but sslrootcert is not provided", cfg.SSLMode)
			}
			parts = append(parts, "sslrootcert="+cfg.SSLRootCert)
		}

		// 客户端证书和私钥必须成对出现，否则无法完成双向 TLS 认证。
		if cfg.SSLCert != "" && cfg.SSLKey != "" {
			parts = append(parts, "sslcert="+cfg.SSLCert, "sslkey="+cfg.SSLKey)
		} else if cfg.SSLCert != "" || cfg.SSLKey != "" {
			return nil, fmt.Errorf("both sslcert and sslkey must be provided for client authentication, or neither")
		}
	} else {
		parts = append(parts, "sslmode=disable")
	}

	connStr := strings.Join(parts, " ")

	db, err := sql.Open(cfg.Driver, connStr)
	if err != nil {
		return nil, fmt.Errorf("connect to %s failed: %w", cfg.Database, err)
	}

	// 主动调用 Ping 以验证数据库连接配置和网络状态是否正常，而不是等到第一次查询出现报错。
	// 重试最多 3 次，每次带独立超时，间隔指数退避，以应对启动阶段数据库短暂不可用的情况。
	var lastErr error
	delay := dbConnectBaseDelay
	for attempt := 1; attempt <= dbConnectMaxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), dbConnectTimeout)
		lastErr = db.PingContext(ctx)
		cancel()
		if lastErr == nil {
			return db, nil
		}
		if attempt < dbConnectMaxRetries {
			time.Sleep(delay)
			delay *= 2
		}
	}

	_ = db.Close()
	return nil, fmt.Errorf("unable to ping database after %d attempts: %w", dbConnectMaxRetries, lastErr)
}
