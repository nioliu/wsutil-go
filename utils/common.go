package utils

import (
	"context"
	"time"
)

func DoWithDeadLine(ctx context.Context, duration time.Duration, funcChan chan int) error {
	ctx, cancle := context.WithTimeout(ctx, duration)
	defer cancle()
	select {
	case <-ctx.Done():
		return TimeOutErr
	case <-funcChan:
		return nil
	}
}
