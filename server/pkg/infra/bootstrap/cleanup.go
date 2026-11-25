package bootstrap

import (
	"fmt"
	"log/slog"
	"sync"
)

// CleanupQueue 维护应用退出时需要执行的清理函数列表.
// 采用队列模式可以让各个模块（如数据库、Redis连接等）保持解耦，各自独立地注册清理逻辑。
type CleanupQueue struct {
	mu    sync.Mutex
	funcs []func()
}

// NewCleanupQueue returns a newly initialized [CleanupQueue].
func NewCleanupQueue() *CleanupQueue {
	return &CleanupQueue{
		funcs: make([]func(), 0),
	}
}

// Add appends a cleanup function f to the queue.
func (q *CleanupQueue) Add(f func()) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.funcs = append(q.funcs, f)
}

// Run executes all registered cleanup functions in reverse order.
// logger is used for logging; if nil, slog.Default() is used.
//
// Special cases are:
//
//	If a cleanup function panics, the panic is recovered to ensure subsequent functions still run.
func (q *CleanupQueue) Run(logger *slog.Logger) {
	// 复制函数切片并在不持有锁的情况下执行，避免某些清理函数也需要获取这把锁从而导致死锁。
	q.mu.Lock()
	funcs := make([]func(), len(q.funcs))
	copy(funcs, q.funcs)
	// 清空原队列防止重复执行。
	q.funcs = nil
	q.mu.Unlock()

	var errors []error
	if logger == nil {
		logger = slog.Default()
	}

	// 逆序执行以模拟栈的后进先出 (LIFO) 行为，确保后注册的资源优先释放。
	logger.Info("running cleanup queue", "count", len(funcs))
	for idx := len(funcs) - 1; idx >= 0; idx-- {
		func(index int, f func()) {
			defer func() {
				if r := recover(); r != nil {
					err := fmt.Errorf("cleanup function %d panicked: %v", index, r)
					errors = append(errors, err)
					logger.Error("cleanup function panicked", "index", index, "panic", r)
				}
			}()
			f()
		}(idx, funcs[idx])
	}

	// 如果有错误，记录日志
	if len(errors) > 0 {
		logger.Error("cleanup completed with errors", "error_count", len(errors))
	} else {
		logger.Info("cleanup completed successfully")
	}
}

// RegisterCloserCleanup registers a standardized close hook for shared resources.
func RegisterCloserCleanup(rt *Runtime, resourceName string, resource interface{ Close() error }) {
	if resource == nil {
		return
	}

	rt.Cleanup.Add(func() {
		if err := resource.Close(); err != nil {
			rt.Logger.Error("failed to close "+resourceName, "error", err)
			return
		}
		rt.Logger.Info(resourceName + " closed successfully")
	})
}
