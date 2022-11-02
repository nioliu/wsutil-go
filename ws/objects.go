// Package ws
// Basic websocket configuration
package ws

import (
	"context"
	"net/http"
	"sync"
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
type AfterCloseFunc func() error

type Msg struct {
	Msg     []byte
	MsgType int
}

// SingleConnOperations conn basic functions
type SingleConnOperations interface {
	Serve() error
	SendMsg(ctx context.Context, msg []byte, msgType int) error
	writePump()
	readPump()
	Close() error
}

type SingleConn struct {
	// basic conn
	conn       Conn
	id         string
	ctx        context.Context
	cancel     context.CancelFunc
	options    []Option
	closed     bool
	isOn       bool // is running
	serverOnce sync.Once
	tags       []string

	beforeHandleReceivedMsg HandleMsgFunc
	handleReceiveMsg        HandleMsgFunc
	afterHandleReceivedMsg  HandleMsgFunc
	handleReceiveTaskErrors HandleTaskErrsFunc
	isReading               bool

	beforeHandleSendMsg  HandleMsgFunc
	afterHandleSendMsg   HandleMsgFunc
	handleSendTaskErrors HandleTaskErrsFunc
	isWriting            bool

	afterCloseFunc AfterCloseFunc
	sendChan       chan Msg // send Msg to others
	heartCheck     time.Duration
	sendTimeOut    time.Duration
	writeTimeOut   time.Duration
}
