package ws

import (
	"context"
	"time"
)

type Option func(conn *SingleConn)

func apply(conn *SingleConn, opt ...Option) {
	for i := 0; i < len(opt); i++ {
		opt[i](conn)
	}
}

func WithContext(ctx context.Context) Option {
	return func(conn *SingleConn) {
		if conn.ctx != nil {
			return
		}
		conn.ctx = ctx
	}
}

func WithSendTimeOut(duration time.Duration) Option {
	return func(conn *SingleConn) {
		if conn.SendTimeOut != 0 {
			return
		}
		conn.SendTimeOut = duration
	}
}

func WithSendChan(sendChan chan Msg) Option {
	return func(conn *SingleConn) {
		if conn.SendChan != nil {
			return
		}
		conn.SendChan = sendChan
	}
}
