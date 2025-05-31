package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024, // TODO: move to config
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		hub.RegisterClient(conn)
		defer hub.UnRegisterClient(conn)
		defer conn.Close()

		// Ждём закрытия соединения
		for {
			if _, _, err := conn.NextReader(); err != nil {
				break
			}
		}
	}
}
