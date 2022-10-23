// Package group
// This package is used to build groups and have some common functions.
package group

import (
	"context"
	"go.uber.org/zap"
	"time"
	"wsutil-go/utils"
	"wsutil-go/ws"
)

type Map map[string]interface{}

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

// AddNewConnWithId add new ws Conn in group, id stand for this conn
// the key can be *net.Coon or new *Group
func (g *Group) AddNewConnWithId(id string, key interface{}) (err error) {
	_, ok1 := key.(ws.Conn)
	_, ok2 := key.(*Group)
	if key == nil || (!ok1 && !ok2) {
		utils.Logger.Error("add conn failed", zap.Error(utils.InvalidArgsErr))
		return utils.InvalidArgsErr
	}
	if len(g.groupMap) > g.maxConnCnt {
		utils.Logger.Error("add conn failed", zap.Error(utils.OutOfMaxCntErr))
		return utils.OutOfMaxCntErr
	}
	groupMap := g.GetGroupMap()
	groupMap[id] = key

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

// append default configuration, each new group builder need to apply this func.
func appendDefault(opts ...Option) []Option {
	opts = append(opts,
		WithUpgrader(&ws.WrappedGorillaUpgrader{}),
		WithGroupId(""), WithHeartCheck(time.Minute),
		WithMaxConnCnt(100), WithMaxConnDuration(time.Hour*24*30),
		WithGroupMap(nil))

	return opts
}
