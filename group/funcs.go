package group

import (
	"context"
)

// Operation ws group functions
type Operation interface {
	// Broadcast send msg to everyone in a group
	Broadcast(ctx context.Context, msg []byte) error

	// WorldPing ping everyone in a group
	WorldPing(ctx context.Context) error

	// SendMsgWithIds send msg to conns in a group
	SendMsgWithIds(ctx context.Context, id ...string) error

	// HandleAllReceivedMsg handle msg from client
	HandleAllReceivedMsg(ctx context.Context, f func(ctx context.Context, id string, msgType int, msg []byte)) error
}
