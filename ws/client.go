package ws

import (
	"context"
	"go.uber.org/zap"
	"wsutil-go/utils"
)

func NewSingleConn(ctx context.Context, conn Conn, opts ...Option) (*SingleConn, error) {
	if conn == nil {
		utils.Logger.Error("Conn can not be nil", zap.Error(utils.InvalidArgsErr))
		return nil, utils.InvalidArgsErr
	}
	options := appendDefault(opts...)

	s := &SingleConn{conn: conn, options: options}
	return s, nil
}
