package ws

import (
	"context"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

type Option func(conn *SingleConn)

func apply(conn *SingleConn) {
	for i := 0; i < len(conn.options); i++ {
		conn.options[i](conn)
	}
}

func appendDefault(opts ...Option) []Option {
	opts = append(opts, WithContext(context.Background()), WithId(""), WithHeartCheck(0),
		WithSendChan(make(chan Msg, 1)), WithSendTimeOut(0), WithWriteTimeOut(0),
		WithReceiveTaskErrors(nil), WithSendTaskErrors(nil))

	return opts
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
		if conn.sendTimeOut != 0 {
			return
		} else if duration == 0 {
			duration = time.Second * 5
		}
		conn.sendTimeOut = duration
	}
}

func WithWriteTimeOut(duration time.Duration) Option {
	return func(conn *SingleConn) {
		if conn.writeTimeOut != 0 {
			return
		} else if duration == 0 {
			duration = time.Second * 5
		}
		conn.writeTimeOut = duration
	}
}

func WithSendChan(sendChan chan Msg) Option {
	return func(conn *SingleConn) {
		if conn.sendChan != nil {
			return
		}
		conn.sendChan = sendChan
	}
}

func WithAfterHandleReceivedMsg(f HandleMsgFunc) Option {
	return func(conn *SingleConn) {
		if conn.afterHandleReceivedMsg != nil {
			return
		}
		conn.afterHandleReceivedMsg = f
	}
}

func WithBeforeHandleReceivedMsg(f HandleMsgFunc) Option {
	return func(conn *SingleConn) {
		if conn.beforeHandleReceivedMsg != nil {
			return
		}
		conn.beforeHandleReceivedMsg = f
	}
}

func WithAfterHandleSendMsg(f HandleMsgFunc) Option {
	return func(conn *SingleConn) {
		if conn.afterHandleSendMsg != nil {
			return
		}
		conn.afterHandleSendMsg = f
	}
}

func WithBeforeHandleSendMsg(f HandleMsgFunc) Option {
	return func(conn *SingleConn) {
		if conn.beforeHandleSendMsg != nil {
			return
		}
		conn.beforeHandleSendMsg = f
	}
}

func WithId(id string) Option {
	return func(conn *SingleConn) {
		if conn.id != "" {
			return
		} else if id == "" {
			rand.Seed(time.Now().Unix())
			id = strconv.Itoa(rand.Int())
		}
		conn.id = id
	}
}

func WithHeartCheck(d time.Duration) Option {
	return func(conn *SingleConn) {
		if conn.heartCheck != 0 {
			return
		} else if d == 0 {
			d = time.Minute
		}
		conn.heartCheck = d
	}
}

func WithHandleReceiveMsg(f HandleMsgFunc) Option {
	return func(conn *SingleConn) {
		if conn.handleReceiveMsg != nil {
			return
		}
		conn.handleReceiveMsg = f
	}
}

func WithSendTaskErrors(f HandleTaskErrsFunc) Option {
	return func(conn *SingleConn) {
		if conn.handleSendTaskErrors != nil {
			return
		}
		if f == nil {
			f = func(ctx context.Context, id string, err []error) error {
				return nil
			}
		}
		conn.handleSendTaskErrors = f
	}
}

func WithReceiveTaskErrors(f HandleTaskErrsFunc) Option {
	return func(conn *SingleConn) {
		if conn.handleReceiveTaskErrors != nil {
			return
		}
		if f == nil {
			f = func(ctx context.Context, id string, err []error) error {
				return err[len(err)-1]
			}
		}
		conn.handleReceiveTaskErrors = f
	}
}

func WithAfterCloseFunc(f AfterCloseFunc) Option {
	return func(conn *SingleConn) {
		if conn.afterCloseFunc != nil {
			return
		}
		conn.afterCloseFunc = f
	}
}

func WithTags(tags ...string) Option {
	return func(conn *SingleConn) {
		conn.tags = append(conn.tags, tags...)
		sort.Strings(conn.tags)
	}
}
