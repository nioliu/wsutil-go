// Package ws
// Basic websocket configuration
package ws

import (
	"context"
	"net/http"
	"time"
	"wsutil-go/utils"
)

// Conn websocket basic method
type Conn interface {
	// ReadMessage reading from that reader to a buffer.
	ReadMessage() (messageType int, p []byte, err error)
	// WriteMessage writing the message and closing the writer
	WriteMessage(messageType int, data []byte) error
	// Close close a conn
	Close() error
}

// Upgrader trans http to websocket
type Upgrader interface {
	// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (Conn, error)
}

type HandleReceiveMsg func(ctx context.Context, id string, msgType int, msg []byte) error

type Msg struct {
	msg     []byte
	msgType int
}
type SingleConn struct {
	// basic conn
	Conn
	ctx                     context.Context
	BeforeHandleReceivedMsg HandleReceiveMsg
	AfterHandleReceivedMsg  HandleReceiveMsg
	Send                    chan Msg // send msg to others
	HeartCheck              time.Duration
	SendTimeOut             time.Duration
	WriteTimeOut            time.Duration
}

// SingleConnOperations conn basic functions
type SingleConnOperations interface {
}

// Serve start listen websocket conn
func (s *SingleConn) Serve() {

}

func (s *SingleConn) writePump() {
	// heart check
	ticker := time.NewTicker(s.HeartCheck)
	defer func() {
		s.Conn.Close()
		ticker.Stop()
	}()
	select {
	case <-ticker.C:
		isDone := make(chan int, 1)
		go func() {
			s.Conn.
		}()
		if err := utils.DoWithDeadLine(s.ctx, s.SendTimeOut, isDone); err != nil {

		}
	}
}

func (s *SingleConn) readPump() {

}
