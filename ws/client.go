package ws

import (
	"context"
	"git.woa.com/nioliu/wsutil-go/utils"
	"go.uber.org/zap"
	"sync"
)

func NewSingleConn(ctx context.Context, conn Conn, opts ...Option) (*SingleConn, error) {
	if conn == nil {
		utils.Logger.Error("Conn can not be nil", zap.Error(utils.InvalidArgsErr))
		return nil, utils.InvalidArgsErr
	}
	options := appendDefault(opts...)

	s := &SingleConn{ctx: ctx, conn: conn, options: options, serverOnce: sync.Once{}}
	return s, nil
}
