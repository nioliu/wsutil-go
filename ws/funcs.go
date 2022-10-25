package ws

import (
	"context"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"time"
	"wsutil-go/utils"
)

// Serve start listen websocket Msg
func (s *SingleConn) Serve() error {
	apply(s)
	// todo check default
	go s.writePump()
	go s.readPump()

	return nil
}

func (s *SingleConn) writePump() {
	// heart check
	ticker := time.NewTicker(s.heartCheck)
	defer func() {
		s.conn.Close()
		ticker.Stop()
	}()
	// control time deadline
	isDone := make(chan int, 1)
	for {
		var msgType int
		var msg []byte
		select {
		case <-ticker.C:
			msgType = websocket.PingMessage
			msg = nil
		case sendMsg := <-s.sendChan:
			msgType = sendMsg.MsgType
			msg = sendMsg.Msg
		case <-s.closeWriteChan:
			return
		}
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
			err := s.conn.WriteMessage(msgType, nil)
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
		}
		if err := s.handleSendTaskErrors(s.ctx, s.id, TaskErrs); err != nil {
			return
		}
	}
}

func (s *SingleConn) readPump() {
	defer s.Close()
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
		if err = s.handleReceiveMsg(s.ctx, s.id, messageType, msg, TaskErrs); err != nil {
			TaskErrs = append(TaskErrs, err)
			utils.Logger.Error("execute handleMsg failed", zap.Error(err))
		}
		if err = s.afterHandleReceivedMsg(s.ctx, s.id, messageType, msg, TaskErrs); err != nil {
			TaskErrs = append(TaskErrs, err)
			utils.Logger.Error("execute AfterHandleMsg hook failed", zap.Error(err))
		}
		if err = s.handleReceiveTaskErrors(s.ctx, s.id, TaskErrs); err != nil {
			return
		}
	}
}

func (s *SingleConn) GetId() string {
	return s.id
}

func (s *SingleConn) Close() error {
	s.closeWriteChan <- 1
	close(s.sendChan)
	if err := s.conn.Close(); err != nil {
		utils.Logger.Error("close basic conn failed", zap.Error(err))
		return err
	}
	return nil
}

func (s *SingleConn) SendMsg(ctx context.Context, msg Msg) error {
	s.sendChan <- msg
	return nil
}
