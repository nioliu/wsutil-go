package group

import (
	"context"
	"git.woa.com/nioliu/wsutil-go/utils"
	"git.woa.com/nioliu/wsutil-go/ws"
	"github.com/gorilla/websocket"
	"sort"
	"sync"
	"time"
)

// Operation ws group functions
type Operation interface {
	// Broadcast send msg to everyone in a group
	Broadcast(ctx context.Context, msg ws.Msg) error

	// WorldPing ping everyone in a group
	WorldPing(ctx context.Context) error

	// SendMsgWithIds send msg to conns in a group
	SendMsgWithIds(ctx context.Context, msg ws.Msg, to ...string) error

	// SendMsgWithTags send msg to specified tags
	SendMsgWithTags(ctx context.Context, msg ws.Msg, strict bool, tags ...string) error

	// AddNewSingleConn add single connector to group
	AddNewSingleConn(conn *ws.SingleConn) error

	// DeleteConnById delete object in map
	DeleteConnById(ctx context.Context, id string) error
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
	// WsUpgrader
	WsUpgrader ws.Upgrader
	// safe concurrent operation
	mu sync.RWMutex

	// beforeHandleHookFunc is applied before handle received msg
	beforeHandleHookFunc ws.HandleMsgFunc
	// afterHandleHookFunc is applied after handle received msg
	afterHandleHookFunc ws.HandleMsgFunc

	// wordCheckInterval check if all connectors in groupMap is active
	wordCheckInterval time.Time

	// handleBroadcastError handle broadcast error after send msg failed, if this return error,
	// broadcast will be stopped and return
	handleBroadcastError func(g *Group, conn *ws.SingleConn, err error) error
}

func (g *Group) SendMsgWithTags(ctx context.Context, msg ws.Msg, strict bool, tags ...string) error {
	sort.Strings(tags)

	g.mu.RLock()
	defer g.mu.RUnlock()

	if !strict {
		for tag := range tags {
			for _, v := range g.GetGroupMap() {
				if subG, is := v.(*Group); is {
					if err := subG.SendMsgWithTags(ctx, msg, strict, tags...); err != nil {
						return err
					}
				} else {
					singleConn := v.(*ws.SingleConn)
					for t := range singleConn.GetTags() {
						if t == tag {
							if err := singleConn.SendMsg(ctx, msg); err != nil {
								return err
							}
						} else if t > tag {
							break
						}
					}
				}
			}
		}
	} else {
		for _, v := range g.GetGroupMap() {
			if subG, is := v.(*Group); is {
				if err := subG.SendMsgWithTags(ctx, msg, strict, tags...); err != nil {
					return err
				}
			} else {
				singleConn := v.(*ws.SingleConn)
				if len(tags) > len(singleConn.GetTags()) {
					break
				}
				var match int
				currentTags := singleConn.GetTags()

				// check if all tags in current singleConn
				for i, j := 0, 0; i < len(tags) && j < len(currentTags); {
					if tags[i] == currentTags[j] {
						match++
						i++
					} else if tags[i] < currentTags[j] {
						break
					}
					j++
				}
				if match == len(tags) {
					if err := singleConn.SendMsg(ctx, msg); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (g *Group) Broadcast(ctx context.Context, msg ws.Msg) error {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for _, v := range g.GetGroupMap() {
		if subG, is := v.(*Group); is {
			if err := subG.Broadcast(ctx, msg); err != nil {
				return err
			}
		} else {
			singleConn := v.(*ws.SingleConn)
			if !singleConn.GetStatus() {
				if err := g.DeleteConnById(ctx, singleConn.GetId()); err != nil {
					return err
				}
				continue
			}
			if msg.MsgType == 0 {
				// just for check active
				continue
			}
			if err := singleConn.SendMsg(ctx, msg); err != nil {
				if g.handleBroadcastError != nil {
					return g.handleBroadcastError(g, singleConn, err)
				}
			}
		}
	}
	return nil
}

func (g *Group) WorldPing(ctx context.Context) error {
	return g.Broadcast(ctx, ws.Msg{
		Msg:     nil,
		MsgType: websocket.PingMessage,
	})
}

func (g *Group) SendMsgWithIds(ctx context.Context, msg ws.Msg, to ...string) error {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for i := 0; i < len(to); i++ {
		c, err := g.GetConnById(to[i])
		if err != nil {
			return err
		}
		if subG, is := c.(*Group); is {
			if err = subG.Broadcast(ctx, msg); err != nil {
				return err
			}
		}
		singleConn := c.(*ws.SingleConn)

		// check status
		if !singleConn.GetStatus() {
			if err = g.DeleteConnById(ctx, singleConn.GetId()); err != nil {
				return err
			}
			continue
		}

		if err = singleConn.SendMsg(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}

// AddNewSingleConn add new ws Conn in group, id stand for this conn
// the key can be *net.Coon or new *Group
func (g *Group) AddNewSingleConn(singleConn *ws.SingleConn) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if singleConn == nil {
		return utils.InvalidArgsErr
	}
	// check limit
	if len(g.groupMap)+1 > g.maxConnCnt {
		// check active connections and delete inactive connections
		if err := g.Broadcast(context.Background(), ws.Msg{}); err != nil {
			return err
		}
		if len(g.groupMap)+1 > g.maxConnCnt {
			return utils.OutOfMaxCntErr
		}
	}

	if !singleConn.GetStatus() {
		if err := singleConn.Serve(); err != nil {
			return err
		}
	}

	groupMap := g.GetGroupMap()
	groupMap[singleConn.GetId()] = singleConn

	return nil
}

func (g *Group) AddSubGroup(ctx context.Context, id string, group *Group) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if group == nil {
		return utils.InvalidArgsErr
	}
	if len(g.groupMap)+len(group.groupMap) > g.maxConnCnt {
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
	g.mu.Lock()
	defer g.mu.Unlock()

	groupMap := g.GetGroupMap()
	singleConn, exist := groupMap[id]
	if !exist {
		return utils.IdNotFoundErr
	}
	if subG, is := singleConn.(*Group); is {
		if err := subG.deleteAllInMap(ctx); err != nil {
			return err
		}
	}
	sc, is := singleConn.(*ws.SingleConn)
	if !is {
		return utils.InvalidArgsErr
	}

	if sc.GetStatus() {
		if err := sc.Close(); err != nil {
			return err
		}
	}
	delete(groupMap, id)
	g.groupMap = groupMap

	return nil
}

func (g *Group) deleteAllInMap(ctx context.Context) error {
	for k, _ := range g.groupMap {
		if err := g.DeleteConnById(ctx, k); err != nil {
			return err
		}
	}
	return nil
}

// GetConnById This method is used for get conn in Group map,
// the return may be a subgroup or net.Conn, developer need to
// charge with this.
func (g *Group) GetConnById(id string) (interface{}, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	groupMap := g.GetGroupMap()
	i, exist := groupMap[id]
	if !exist {
		return nil, utils.IdNotFoundErr
	}
	return i, nil
}

func (g *Group) GetGroupMap() Map {
	return g.groupMap
}
