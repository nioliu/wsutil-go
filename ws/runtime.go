package ws

import (
	"context"
	"git.woa.com/nioliu/wsutil-go/utils"
	"github.com/gorilla/websocket"
	"time"
)

// Serve start listen websocket Msg
func (s *SingleConn) Serve() error {
	s.serverOnce.Do(func() {
		apply(s)
		// control
		s.ctx, s.cancel = context.WithCancel(s.ctx)
		// todo check default
		go s.writePump()
		go s.readPump()

		s.isOn = true // mark
	})
	return nil
}

func (s *SingleConn) writePump() {
	// heart check
	ticker := time.NewTicker(s.heartCheck)
	// control time deadline
	isDone := make(chan int, 1)
	defer func() {
		if !s.closed && s.isOn {
			s.Close()
			s.closed = true
		}
		ticker.Stop()
		close(isDone)
	}()
	for {
		if !s.isOn {
			return
		}
		var msgType int
		var msg []byte
		select {
		case <-ticker.C:
			// avoid competition
			select {
			case sendMsg := <-s.sendChan:
				msgType = sendMsg.MsgType
				msg = sendMsg.Msg
			default:
				msgType = websocket.PingMessage
				msg = nil
			}
		case sendMsg := <-s.sendChan:
			msgType = sendMsg.MsgType
			msg = sendMsg.Msg
		case <-s.ctx.Done():
			return
		}
		var TaskErrs []error
		if s.isOn == false {
			return
		}
		go func() {
			defer func() {
				if !s.closed {
					isDone <- 1
				}
			}()
			if s.beforeHandleSendMsg != nil {
				if err := s.beforeHandleSendMsg(s.ctx, s.id, msgType, msg, TaskErrs); err != nil {
					TaskErrs = append(TaskErrs, err)
				}
			}
			err := s.conn.WriteMessage(msgType, msg)
			if err != nil {
				TaskErrs = append(TaskErrs, err)
			}

			// reset heart check if write successfully
			if msgType != websocket.PingMessage {
				ticker.Reset(s.heartCheck)
			}

			if s.afterHandleSendMsg != nil {
				if err := s.afterHandleSendMsg(s.ctx, s.id, msgType, msg, TaskErrs); err != nil {
					TaskErrs = append(TaskErrs, err)
				}
			}
		}()
		if err := utils.DoWithDeadLine(s.ctx, s.sendTimeOut, isDone); err != nil {
			return
		}
		if TaskErrs != nil {
			if s.handleSendTaskErrors != nil {
				if err := s.handleSendTaskErrors(s.ctx, s.id, TaskErrs); err != nil {
					return
				}
			}
		}
	}
}

// developer need to handle TaskErr, since it will never return to active status.
func (s *SingleConn) readPump() {
	defer func() {
		if !s.closed && s.isOn {
			s.Close()
			s.closed = true
		}
	}()
	for {
		if !s.isOn {
			return
		}
		var TaskErrs []error
		messageType, msg, err := s.conn.ReadMessage()
		if err != nil {
			TaskErrs = append(TaskErrs, err)
			s.isOn = false
			goto handleError
		}

		if s.beforeHandleReceivedMsg != nil {
			if err = s.beforeHandleReceivedMsg(s.ctx, s.id, messageType, msg, TaskErrs); err != nil {
				TaskErrs = append(TaskErrs, err)
			}
		}

		if s.handleReceiveMsg != nil {
			if err = s.handleReceiveMsg(s.ctx, s.id, messageType, msg, TaskErrs); err != nil {
				TaskErrs = append(TaskErrs, err)
			}
		}

		if s.afterHandleReceivedMsg != nil {
			if err = s.afterHandleReceivedMsg(s.ctx, s.id, messageType, msg, TaskErrs); err != nil {
				TaskErrs = append(TaskErrs, err)
			}
		}

	handleError:
		if TaskErrs != nil {
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
		return err
	}
	s.cancel()
	s.isOn = false
	return nil
}

func (s *SingleConn) SendMsg(ctx context.Context, msg Msg) error {
	s.sendChan <- msg
	return nil
}

// GetStatus add some monitor fields.
func (s *SingleConn) GetStatus() bool {
	return s.isOn
}

func (s *SingleConn) GetTags() []string {
	return s.tags
}
