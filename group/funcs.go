package group

import (
	"context"
	"time"
	"wsutil-go/ws"
)

// Operation ws group functions
type Operation interface {
	// Run launch all the conns in a group
	Run(ctx context.Context) error

	// Broadcast send msg to everyone in a group
	Broadcast(ctx context.Context, msg []byte) error

	// WorldPing ping everyone in a group
	WorldPing(ctx context.Context) error

	// SendMsgWithIds send msg to conns in a group
	SendMsgWithIds(ctx context.Context, from string, to ...string) error
}

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
	beforeHandleHookFunc ws.HandleReceiveMsg
	// afterHandleHookFunc is applied after handle received msg
	afterHandleHookFunc ws.HandleReceiveMsg
}

func (g *Group) Broadcast(ctx context.Context, msg []byte) error {
	//TODO implement me
	panic("implement me")
}

func (g *Group) WorldPing(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (g *Group) SendMsgWithIds(ctx context.Context, id ...string) error {
	//TODO implement me
	panic("implement me")
}

func (g *Group) HandleAllReceivedMsg(ctx context.Context, f func(ctx context.Context, id string, msgType int, msg []byte)) error {

}

func (g *Group) Run(ctx context.Context) {
	for k, v := range g.groupMap {

	}
}

func waitAndRead(ctx context.Context, conn ws.Conn) error {

}

func waitAndWrite(ctx context.Context, conn ws.Conn) error {

}
