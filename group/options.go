package group

import (
	"git.woa.com/nioliu/wsutil-go/utils"
	"git.woa.com/nioliu/wsutil-go/ws"
	"go.uber.org/zap"
	"log"
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

// Option Group configuration
type Option func(group *Group)

func apply(group *Group, opt ...Option) {
	for i := 0; i < len(opt); i++ {
		opt[i](group)
	}
}

// append default configuration, each new group builder need to apply this func.
func appendDefault(opts ...Option) []Option {
	opts = append(opts,
		WithUpgrader(&ws.WrappedGorillaUpgrader{}),
		WithGroupId(""), WithHeartCheck(time.Minute),
		WithMaxConnCnt(100), WithMaxConnDuration(time.Hour*24*30),
		WithGroupMap(nil))

	return opts
}

func WithMaxConnCnt(cnt int) Option {
	if cnt <= 0 {
		log.Fatal("integer can't less than 1")
		return nil
	}
	return func(group *Group) {
		if group.maxConnCnt != 0 {
			return
		}
		group.maxConnCnt = cnt
	}
}

func WithHeartCheck(duration time.Duration) Option {
	return func(group *Group) {
		if group.heartCheck != 0 {
			return
		}
		group.heartCheck = duration
	}
}

func WithHandleBroadcastError(f func(g *Group, conn *ws.SingleConn, err error) error) Option {
	return func(group *Group) {
		if group.handleBroadcastError != nil {
			return
		}
		group.handleBroadcastError = f
	}
}

func WithMaxConnDuration(duration time.Duration) Option {
	return func(group *Group) {
		if group.maxConnDuration != 0 {
			return
		}
		group.maxConnDuration = duration
	}
}

// WithUpgrader for user decide which websocket upgrader to use
func WithUpgrader(upgrader ws.Upgrader) Option {
	t := reflect.ValueOf(upgrader)
	if t.Kind() != reflect.Pointer {
		log.Fatal("need pointer to set upgrader")
	}
	if t.IsNil() {
		log.Fatal("set upgrader failed", zap.Error(utils.InvalidOptionsErr))
		return nil
	}
	return func(group *Group) {
		if group.WsUpgrader != nil {
			return
		}
		group.WsUpgrader = upgrader
	}
}

func WithGroupId(id string) Option {
	if id == "" {
		rand.Seed(time.Now().UnixMilli())
		id = strconv.Itoa(rand.Int())
	}
	return func(group *Group) {
		if group.id != "" {
			return
		}
		group.id = id
	}
}

func WithGroupMap(m Map) Option {
	if m == nil {
		m = Map{}
	}
	return func(group *Group) {
		if group.GetGroupMap() != nil {
			return
		}
		group.groupMap = m
	}
}

func WithBeforeHandleHookFunc(f ws.HandleMsgFunc) Option {
	return func(group *Group) {
		if group.beforeHandleHookFunc != nil {
			return
		}
		group.beforeHandleHookFunc = f
	}
}

func WithAfterHandleHookFunc(f ws.HandleMsgFunc) Option {
	return func(group *Group) {
		if group.afterHandleHookFunc != nil {
			return
		}
		group.afterHandleHookFunc = f
	}
}
