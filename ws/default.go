package ws

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type WrappedGorillaUpgrader struct {
	websocket.Upgrader
}

func (u *WrappedGorillaUpgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (Conn, error) {
	return u.Upgrader.Upgrade(w, r, responseHeader)
}
