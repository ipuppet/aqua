package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"

	"github.com/ipuppet/aqua/server/pkg/infra/config"
	"github.com/ipuppet/aqua/server/pkg/infra/models"
)

// Client wraps *sql.DB and swaps the underlying connection when app config changes.
type Client struct {
	mu       sync.RWMutex
	reloadMu sync.Mutex
	db       *sql.DB
	cfg      models.DatabaseConfig
	logger   *slog.Logger
}

// NewClient creates a reloadable database client and subscribes to app config updates.

func NewClient(watcher *config.Watcher[*models.AppConfig], logger *slog.Logger) (*Client, error) {
	if watcher == nil {
		return nil, fmt.Errorf("database watcher is not configured")
	}
	appCfg := watcher.Value()
	if appCfg == nil {
		return nil, fmt.Errorf("app config is not loaded")
	}

	cfg := appCfg.Database
	db, err := Connect(cfg)
	if err != nil {
		return nil, err
	}

	client := &Client{db: db, cfg: cfg, logger: logger}
	watcher.OnChange(func(nextAppCfg *models.AppConfig) {
		if nextAppCfg == nil {
			return
		}
		client.reloadIfNeeded(nextAppCfg.Database)
	})

	return client, nil
}

// NewClientFromDB wraps an existing *sql.DB (typically from sqlmock in tests).
func NewClientFromDB(db *sql.DB) *Client {
	return &Client{db: db}
}

func (c *Client) current() *sql.DB {
	c.mu.RLock()
	db := c.db
	c.mu.RUnlock()
	return db
}

func (c *Client) reloadIfNeeded(next models.DatabaseConfig) {
	c.reloadMu.Lock()
	defer c.reloadMu.Unlock()

	c.mu.RLock()
	curr := c.cfg
	c.mu.RUnlock()

	if curr == next {
		return
	}

	nextDB, err := Connect(next)
	if err != nil {
		if c.logger != nil {
			c.logger.Error("database reload failed",
				slog.Any("error", err))
		}
		return
	}

	c.mu.Lock()
	oldDB := c.db
	c.db = nextDB
	c.cfg = next
	c.mu.Unlock()

	if oldDB != nil {
		_ = oldDB.Close()
	}

	if c.logger != nil {
		c.logger.Info("database configuration reloaded successfully")
	}
}

// QueryContext proxies QueryContext to the active database handle.
func (c *Client) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return c.current().QueryContext(ctx, query, args...)
}

// QueryRowContext proxies QueryRowContext to the active database handle.
func (c *Client) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return c.current().QueryRowContext(ctx, query, args...)
}

// ExecContext proxies ExecContext to the active database handle.
func (c *Client) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return c.current().ExecContext(ctx, query, args...)
}

// Exec proxies Exec to the active database handle.
func (c *Client) Exec(query string, args ...any) (sql.Result, error) {
	return c.current().Exec(query, args...)
}

// BeginTx starts a transaction from the active database handle.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return c.current().BeginTx(ctx, opts)
}

// Close closes the currently active database handle.
func (c *Client) Close() error {
	db := c.current()
	if db == nil {
		return nil
	}
	return db.Close()
}
