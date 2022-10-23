// Package ws
// Basic websocket configuration
package ws

import (
	"context"
	"net/http"
	"time"
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

type HandleMsgFunc func(ctx context.Context, id string, msgType int, msg []byte, err []error) error
type HandleTaskErrsFunc func(ctx context.Context, id string, err []error) error

type Msg struct {
	msg     []byte
	msgType int
}

type SingleConn struct {
	// basic conn
	Conn
	id  string
	ctx context.Context

	BeforeHandleReceivedMsg HandleMsgFunc
	HandleReceiveMsg        HandleMsgFunc
	AfterHandleReceivedMsg  HandleMsgFunc
	HandleReceiveTaskErrors HandleTaskErrsFunc

	BeforeHandleSendMsg  HandleMsgFunc
	AfterHandleSendMsg   HandleMsgFunc
	HandleSendTaskErrors HandleTaskErrsFunc

	SendChan     chan Msg // send msg to others
	HeartCheck   time.Duration
	SendTimeOut  time.Duration
	WriteTimeOut time.Duration
}

// SingleConnOperations conn basic functions
type SingleConnOperations interface {
}
