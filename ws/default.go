package ws

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type WrappedGorillaUpgrader struct {
	websocket.Upgrader
}

// Upgrade transform http to websocket
func (u *WrappedGorillaUpgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (Conn, error) {
	return u.Upgrader.Upgrade(w, r, responseHeader)
}
