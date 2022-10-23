// Package ws
// Basic websocket configuration
package ws

import (
	"net/http"
)

// Conn websocket basic method
type Conn interface {
	// ReadMessage reading from that reader to a buffer.
	ReadMessage() (messageType int, p []byte, err error)
	// WriteMessage writing the message and closing the writer
	WriteMessage(messageType int, data []byte) error
}

// Upgrader trans http to websocket
type Upgrader interface {
	// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (Conn, error)
}
