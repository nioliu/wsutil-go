package group

import (
	"context"
	"go.uber.org/zap"
	"time"
	"wsutil-go/utils"
	"wsutil-go/ws"
)

// Operation ws group functions
type Operation interface {
	// Broadcast send msg to everyone in a group
	Broadcast(ctx context.Context, msg []byte) error

	// WorldPing ping everyone in a group
	WorldPing(ctx context.Context) error

	// SendMsgWithIds send msg to conns in a group
	SendMsgWithIds(ctx context.Context, from string, to ...string) error

	// AddNewSingleConnWithId add single connector to group
	AddNewSingleConnWithId(id string, conn *ws.SingleConn) error
}

type Map map[string]interface{}

// Group basic group struct
type Group struct {
	// current group id
	id string
	// current group connections
	groupMap Map
	// max connector number
	maxConnCnt int
	// heart check duration
	heartCheck time.Duration
	// max conn duration
	maxConnDuration time.Duration
	// upgrader
	WsUpgrader ws.Upgrader

	// beforeHandleHookFunc is applied before handle received msg
	beforeHandleHookFunc ws.HandleMsgFunc
	// afterHandleHookFunc is applied after handle received msg
	afterHandleHookFunc ws.HandleMsgFunc
}

func (g *Group) Broadcast(ctx context.Context, msg []byte) error {
	//TODO implement me
	panic("implement me")
}

func (g *Group) WorldPing(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (g *Group) SendMsgWithIds(ctx context.Context, from string, to ...string) error {
	//TODO implement me
	panic("implement me")
}

// AddNewSingleConnWithId add new ws Conn in group, id stand for this conn
// the key can be *net.Coon or new *Group
func (g *Group) AddNewSingleConnWithId(id string, singleConn *ws.SingleConn) error {
	if singleConn == nil {
		utils.Logger.Error("add singleConn failed", zap.Error(utils.InvalidArgsErr))
		return utils.InvalidArgsErr
	}
	// check limit
	if len(g.groupMap)+1 > g.maxConnCnt {
		utils.Logger.Error("add singleConn failed", zap.Error(utils.OutOfMaxCntErr))
		return utils.OutOfMaxCntErr
	}
	groupMap := g.GetGroupMap()
	groupMap[id] = singleConn

	return singleConn.Serve()
}

func (g *Group) AddSubGroup(ctx context.Context, id string, group *Group) error {
	if group == nil {
		utils.Logger.Error("add group failed", zap.Error(utils.InvalidArgsErr))
		return utils.InvalidArgsErr
	}
	if len(g.groupMap)+len(group.groupMap) > g.maxConnCnt {
		utils.Logger.Error("add group failed", zap.Error(utils.OutOfMaxCntErr))
		return utils.OutOfMaxCntErr
	}
	// check status
	if err := group.WorldPing(ctx); err != nil {
		return err
	}
	groupMap := g.GetGroupMap()
	groupMap[id] = group

	return nil
}

func (g *Group) DeleteConnById(ctx context.Context, id string) error {
	groupMap := g.GetGroupMap()
	singleConn, exist := groupMap[id]
	if !exist {
		utils.Logger.Error("delete failed", zap.Error(utils.IdNotFoundErr))
		return utils.IdNotFoundErr
	}
	if subG, is := singleConn.(*Group); is {
		if err := subG.DeleteAllInMap(ctx); err != nil {
			return err
		}
	}
	sc, is := singleConn.(*ws.SingleConn)
	if !is {
		utils.Logger.Error("invalid type for current id", zap.Error(utils.InvalidArgsErr))
		return utils.InvalidArgsErr
	}
	if err := sc.Close(); err != nil {
		utils.Logger.Error("close single conn failed", zap.Error(err))
		return err
	}
	delete(groupMap, id)

	return nil
}

func (g *Group) DeleteAllInMap(ctx context.Context) error {
	for k, _ := range g.groupMap {
		if err := g.DeleteConnById(ctx, k); err != nil {
			utils.Logger.Error("close failed", zap.Error(err), zap.String("id", k))
		}
	}
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
