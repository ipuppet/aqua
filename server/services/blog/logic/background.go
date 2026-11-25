package logic

import (
	"context"
	"log/slog"
	"time"
)

func enqueueBackground(
	ctx context.Context,
	timeout time.Duration,
	job string,
	logger func(context.Context) *slog.Logger,
	run func(context.Context) error,
) {
	go func() {
		jobCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), timeout)
		defer cancel()
		if err := run(jobCtx); err != nil {
			logger(jobCtx).Error("background job failed",
				slog.String("job", job),
				slog.Any("error", err))
		}
	}()
}
