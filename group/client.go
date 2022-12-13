// Package group
// This package is used to build groups and have some common functions.
package group

import (
	"context"
	"git.woa.com/nioliu/wsutil-go/ws"
	"time"
)

// New build websocket group based on gorilla/websocket
func New(opts ...Option) (*Group, error) {
	ctx := context.Background()
	return NewWithContext(ctx, opts...)
}

// NewWithContext build websocket group based on gorilla/websocket with context
func NewWithContext(ctx context.Context, opts ...Option) (*Group, error) {
	g := &Group{}
	// only apply first option for each.
	opts = appendDefault(opts...)

	apply(g, opts...)
	checkAllInMap(ctx, g) // check status

	return g, nil
}

// heart beat
func checkAllInMap(ctx context.Context, g *Group) {
	ticker := time.NewTicker(g.heartCheck)
	for range ticker.C {
		if err := g.Broadcast(ctx, ws.Msg{}); err != nil {
			return
		}
	}
}
