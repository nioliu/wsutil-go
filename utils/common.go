package utils

import (
	"context"
	"go.uber.org/zap"
	"time"
)

func DoWithDeadLine(ctx context.Context, duration time.Duration, funcChan chan int) error {
	ctx, cancle := context.WithTimeout(ctx, duration)
	defer cancle()
	select {
	case <-ctx.Done():
		Logger.Error("do function failed", zap.Error(TimeOutErr))
		return TimeOutErr
	case <-funcChan:
		return nil
	}
}
