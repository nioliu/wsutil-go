// Package group
// This package is used to build groups and have some common functions.
package group

import (
	"context"
	"git.woa.com/nioliu/wsutil-go/ws"
	"time"
)

// NewDefaultGroupAndUpgrader build websocket group based on gorilla/websocket
func NewDefaultGroupAndUpgrader(opts ...Option) (*Group, error) {
	ctx := context.Background()
	return NewDefaultGroupWithContextAndUpgrader(ctx, opts...)
}

// NewDefaultGroupWithContextAndUpgrader build websocket group based on gorilla/websocket with context
func NewDefaultGroupWithContextAndUpgrader(ctx context.Context, opts ...Option) (*Group, error) {
	g := &Group{}

	// only apply first option for each.
	opts = appendDefault(opts...)

	apply(g, opts...)
	go checkAllInMap(ctx, g)

	return g, nil
}

func NewGroupWithContext(ctx context.Context, upgrader ws.Upgrader, opts ...Option) *Group {
	g := &Group{WsUpgrader: upgrader}
	opts = appendDefault(opts...)

	apply(g, opts...)
	go checkAllInMap(ctx, g)
	return g
}

func checkAllInMap(ctx context.Context, g *Group) {
	ticker := time.NewTicker(g.heartCheck)
	for range ticker.C {
		if err := g.Broadcast(ctx, ws.Msg{}); err != nil {
			return
		}
	}
}
