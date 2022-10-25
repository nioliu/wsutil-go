package ws

import (
	"context"
	"fmt"
	"git.woa.com/nioliu/wsutil-go/utils"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"time"
)

// Serve start listen websocket Msg
func (s *SingleConn) Serve() error {
	apply(s)
	// control
	s.ctx, s.cancel = context.WithCancel(s.ctx)
	// todo check default
	go s.writePump()
	go s.readPump()

	return nil
}

func (s *SingleConn) writePump() {
	// heart check
	ticker := time.NewTicker(s.heartCheck)
	// control time deadline
	isDone := make(chan int, 1)
	defer func() {
		if !s.closed {
			s.Close()
			s.closed = true
		}
		ticker.Stop()
		close(isDone)
	}()
	for {
		var msgType int
		var msg []byte
		select {
		case <-ticker.C:
			msgType = websocket.BinaryMessage
			msg = nil
		case sendMsg := <-s.sendChan:
			msgType = sendMsg.MsgType
			msg = sendMsg.Msg
		case <-s.ctx.Done():
			return
		}
		fmt.Printf("msgType:%d, msg:%s", msgType, string(msg))
		var TaskErrs []error
		go func() {
			defer func() {
				isDone <- 1
			}()
			if s.beforeHandleSendMsg != nil {
				if err := s.beforeHandleSendMsg(s.ctx, s.id, msgType, msg, TaskErrs); err != nil {
					TaskErrs = append(TaskErrs, err)
					utils.Logger.Error("execute before hook failed", zap.Error(err))
				}
			}
			err := s.conn.WriteMessage(msgType, msg)
			if err != nil {
				TaskErrs = append(TaskErrs, err)
				// todo add handle error func
				utils.Logger.Error("send Msg failed", zap.Error(err))
			}

			if s.afterHandleSendMsg != nil {
				if err = s.afterHandleSendMsg(s.ctx, s.id, msgType, msg, TaskErrs); err != nil {
					TaskErrs = append(TaskErrs, err)
					utils.Logger.Error("execute afterHook failed", zap.Error(err))
				}
			}
		}()
		if err := utils.DoWithDeadLine(s.ctx, s.sendTimeOut, isDone); err != nil {
			utils.Logger.Error("send Msg failed", zap.Error(err))
			return
		}
		if s.handleSendTaskErrors != nil {
			if err := s.handleSendTaskErrors(s.ctx, s.id, TaskErrs); err != nil {
				return
			}
		}
	}
}

func (s *SingleConn) readPump() {
	defer func() {
		if !s.closed {
			s.Close()
			s.closed = true
		}
	}()
	for {
		var TaskErrs []error
		messageType, msg, err := s.conn.ReadMessage()
		if err != nil {
			utils.Logger.Error("read Msg failed", zap.Error(err))
			return
		}

		if s.beforeHandleReceivedMsg != nil {
			if err = s.beforeHandleReceivedMsg(s.ctx, s.id, messageType, msg, TaskErrs); err != nil {
				TaskErrs = append(TaskErrs, err)
				utils.Logger.Error("execute before handle hook failed", zap.Error(err))
			}
		}

		if s.handleReceiveMsg != nil {
			if err = s.handleReceiveMsg(s.ctx, s.id, messageType, msg, TaskErrs); err != nil {
				TaskErrs = append(TaskErrs, err)
				utils.Logger.Error("execute handleMsg failed", zap.Error(err))
			}
		}

		if s.afterHandleReceivedMsg != nil {
			if err = s.afterHandleReceivedMsg(s.ctx, s.id, messageType, msg, TaskErrs); err != nil {
				TaskErrs = append(TaskErrs, err)
				utils.Logger.Error("execute AfterHandleMsg hook failed", zap.Error(err))
			}
		}

		if s.handleReceiveTaskErrors != nil {
			if err = s.handleReceiveTaskErrors(s.ctx, s.id, TaskErrs); err != nil {
				return
			}
		}
	}
}

func (s *SingleConn) GetId() string {
	return s.id
}

func (s *SingleConn) Close() error {
	close(s.sendChan)
	if err := s.conn.Close(); err != nil {
		utils.Logger.Error("close basic conn failed", zap.Error(err))
		return err
	}
	s.cancel()
	return nil
}

func (s *SingleConn) SendMsg(ctx context.Context, msg Msg) error {
	s.sendChan <- msg
	return nil
}
