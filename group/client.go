// Package group
// This package is used to build groups and have some common functions.
package group

import (
	"context"
	"go.uber.org/zap"
	"wsutil-go/utils"
	"wsutil-go/ws"
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

// AddNewSingleConnWithId add new ws Conn in group, id stand for this conn
// the key can be *net.Coon or new *Group
func (g *Group) AddNewSingleConnWithId(id string, conn *ws.SingleConn) error {
	if conn == nil {
		utils.Logger.Error("add conn failed", zap.Error(utils.InvalidArgsErr))
		return utils.InvalidArgsErr
	}
	if len(g.groupMap) > g.maxConnCnt {
		utils.Logger.Error("add conn failed", zap.Error(utils.OutOfMaxCntErr))
		return utils.OutOfMaxCntErr
	}
	groupMap := g.GetGroupMap()
	groupMap[id] = conn

	return nil
}

func (g *Group) AddSubGroup(id string, group *Group) error {
	if group == nil {
		utils.Logger.Error("add group failed", zap.Error(utils.InvalidArgsErr))
		return utils.InvalidArgsErr
	}
	if len(g.groupMap) > g.maxConnCnt {
		utils.Logger.Error("add group failed", zap.Error(utils.OutOfMaxCntErr))
		return utils.OutOfMaxCntErr
	}
	groupMap := g.GetGroupMap()
	groupMap[id] = group

	return nil
}

func (g *Group) DeleteConnById(id string) error {
	groupMap := g.GetGroupMap()
	_, exist := groupMap[id]
	if !exist {
		utils.Logger.Error("delete failed", zap.Error(utils.IdNotFoundErr))
		return utils.IdNotFoundErr
	}
	delete(groupMap, id)

	return nil
}

// GetConnById This method is used for get conn in Group map,
// the return may be a subgroup or net.Conn, developer need to
// charge with this.
func (g *Group) GetConnById(id string) (interface{}, error) {
	groupMap := g.GetGroupMap()
	i, exist := groupMap[id]
	if !exist {
		utils.Logger.Error("get id failed", zap.Error(utils.IdNotFoundErr))
		return nil, utils.IdNotFoundErr
	}
	return i, nil
}

func (g *Group) GetGroupMap() Map {
	return g.groupMap
}
