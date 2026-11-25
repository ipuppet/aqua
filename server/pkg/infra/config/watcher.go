package config

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Watcher tracks a physical configuration file and manages atomic hot-reloading notifications to downstream subscribers.
// 其底层借用 [atomic.Value] 维护快照，这有助于在读多写少的配置读取场景下彻底避开互斥锁的资源竞争导致阻塞。
type Watcher[T any] struct {
	v           *viper.Viper
	value       atomic.Value
	path        string
	subscribers []func(T)
	mu          sync.RWMutex
	logger      *slog.Logger
}

// NewWatcher creates a [Watcher] for the config file at path (relative to loader.BasePath).
func NewWatcher[T any](loader *Loader, path string, logger *slog.Logger) (*Watcher[T], error) {
	v := viper.New()
	absPath := loader.resolvePath(path)
	v.SetConfigFile(absPath)
	if ext := strings.TrimPrefix(filepath.Ext(absPath), "."); ext != "" {
		v.SetConfigType(ext)
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config %s: %w", absPath, err)
	}

	w := &Watcher[T]{v: v, path: path, logger: logger}
	if err := w.reload(); err != nil {
		return nil, err
	}

	v.OnConfigChange(func(fsnotify.Event) {
		if err := w.reload(); err != nil {
			if l := w.getLogger(); l != nil {
				l.Error("failed to reload configuration file",
					slog.String("config.path", path),
					slog.Any("error", err))
			}
			return
		}
		w.dispatch()
	})
	v.WatchConfig()

	return w, nil
}

func (w *Watcher[T]) reload() error {
	if err := w.v.ReadInConfig(); err != nil {
		return fmt.Errorf("read config %s: %w", w.path, err)
	}
	var cfg T
	if err := w.v.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("unmarshal %s: %w", w.path, err)
	}
	w.value.Store(cfg)
	if l := w.getLogger(); l != nil {
		l.Info("configuration file reloaded successfully",
			slog.String("config.path", w.path))
	}
	return nil
}

// SetLogger dynamically updates the logger instance.
func (w *Watcher[T]) SetLogger(logger *slog.Logger) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.logger = logger
}

func (w *Watcher[T]) getLogger() *slog.Logger {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.logger
}

// Value safely extracts the latest generic structure bound to this Watcher.
func (w *Watcher[T]) Value() T {
	val, _ := w.value.Load().(T)
	return val
}

// OnChange appends a subscriber hook that fires via dedicated goroutines upon successful reloads.
// 注册配置变动的异步回调。
func (w *Watcher[T]) OnChange(fn func(T)) {
	if fn == nil {
		return
	}
	w.mu.Lock()
	w.subscribers = append(w.subscribers, fn)
	w.mu.Unlock()
}

func (w *Watcher[T]) dispatch() {
	cfg := w.Value()
	w.mu.RLock()
	subs := append([]func(T){}, w.subscribers...)
	w.mu.RUnlock()
	for _, sub := range subs {
		go sub(cfg)
	}
}
