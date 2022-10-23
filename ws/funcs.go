package ws

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"time"
	"wsutil-go/utils"
)

// Serve start listen websocket conn
func (s *SingleConn) Serve() {
	go s.writePump()
	go s.readPump()
}

func (s *SingleConn) writePump() {
	// heart check
	ticker := time.NewTicker(s.HeartCheck)
	defer func() {
		s.Conn.Close()
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
		case sendMsg := <-s.SendChan:
			msgType = sendMsg.msgType
			msg = sendMsg.msg
		}
		var TaskErrs []error
		go func() {
			defer func() {
				isDone <- 1
			}()
			if err := s.BeforeHandleSendMsg(s.ctx, s.id, msgType, msg, TaskErrs); err != nil {
				TaskErrs = append(TaskErrs, err)
				utils.Logger.Error("execute before hook failed", zap.Error(err))
			}
			err := s.Conn.WriteMessage(msgType, nil)
			if err != nil {
				TaskErrs = append(TaskErrs, err)
				// todo add handle error func
				utils.Logger.Error("send msg failed", zap.Error(err))
			}
			if s.AfterHandleSendMsg != nil {
				if err = s.AfterHandleSendMsg(s.ctx, s.id, msgType, msg, TaskErrs); err != nil {
					TaskErrs = append(TaskErrs, err)
					utils.Logger.Error("execute afterHook failed", zap.Error(err))
				}
			}
		}()
		if err := utils.DoWithDeadLine(s.ctx, s.SendTimeOut, isDone); err != nil {
			utils.Logger.Error("send msg failed", zap.Error(err))
		}
		if err := s.HandleSendTaskErrors(s.ctx, s.id, TaskErrs); err != nil {
			return
		}
	}
}

func (s *SingleConn) readPump() {
	defer s.Conn.Close()
	for {
		var TaskErrs []error
		messageType, msg, err := s.Conn.ReadMessage()
		if err != nil {
			utils.Logger.Error("read msg failed", zap.Error(err))
			return
		}
		if s.BeforeHandleReceivedMsg != nil {
			if err = s.BeforeHandleReceivedMsg(s.ctx, s.id, messageType, msg, TaskErrs); err != nil {
				TaskErrs = append(TaskErrs, err)
				utils.Logger.Error("execute before handle hook failed", zap.Error(err))
			}
		}
		if err = s.HandleReceiveMsg(s.ctx, s.id, messageType, msg, TaskErrs); err != nil {
			TaskErrs = append(TaskErrs, err)
			utils.Logger.Error("execute handleMsg failed", zap.Error(err))
		}
		if err = s.AfterHandleReceivedMsg(s.ctx, s.id, messageType, msg, TaskErrs); err != nil {
			TaskErrs = append(TaskErrs, err)
			utils.Logger.Error("execute AfterHandleMsg hook failed", zap.Error(err))
		}
		if err = s.HandleSendTaskErrors(s.ctx, s.id, TaskErrs); err != nil {
			return
		}
	}
}
