// Package group
// This package is used to build groups and have some common functions.
package group

import (
	"context"
	"git.woa.com/nioliu/wsutil-go/utils"
	"git.woa.com/nioliu/wsutil-go/ws"
	"go.uber.org/zap"
)

var groups map[string]*Group

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
	return g, nil
}

func NewGroupWithContext(ctx context.Context, upgrader ws.Upgrader, opts ...Option) (*Group, error) {
	g := &Group{WsUpgrader: upgrader}
	opts = appendDefault(opts...)

	return g, nil
}

// RegisterGroup register group to map
// developer are responsible for this maintenance
func RegisterGroup(ctx context.Context, group *Group) error {
	_, exist := groups[group.id]
	if exist {
		utils.Logger.Error("register gruop failed", zap.Error(utils.DuplicatedIdErr))
		return utils.DuplicatedIdErr
	}
	groups[group.id] = group
	return nil
}
