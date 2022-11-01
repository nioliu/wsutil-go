package ws

import (
	"context"
	"git.woa.com/nioliu/wsutil-go/utils"
	"sync"
)

func NewSingleConn(ctx context.Context, conn Conn, opts ...Option) (*SingleConn, error) {
	if conn == nil {
		return nil, utils.InvalidArgsErr
	}
	options := appendDefault(opts...)

	s := &SingleConn{ctx: ctx, conn: conn, options: options, serverOnce: sync.Once{}}
	return s, nil
}
