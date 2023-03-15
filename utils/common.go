package utils

import (
	"context"
	"time"
)

func DoWithDeadLine(ctx context.Context, duration time.Duration, funcChan chan int) error {
	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()
	select {
	case <-ctx.Done():
		return TimeOutErr
	case <-funcChan:
		return nil
	}
}
