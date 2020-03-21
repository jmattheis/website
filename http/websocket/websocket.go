package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/jmattheis/website/content"
	"net/http"
)

var upgrader = websocket.Upgrader{}

func Handle(port string) http.HandlerFunc {
	tty := &content.InteractiveText{
		Prompt:   "",
		Protocol: "websocket",
		Port:     port,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		go func() {
			defer conn.Close()
			data, _ := tty.Exec("start")
			err := conn.WriteMessage(websocket.TextMessage, []byte(data))

			if err != nil {
				return
			}

			for {

				mt, message, err := conn.ReadMessage()

				if err != nil || mt != websocket.TextMessage {
					return
				}

				data, exit := tty.Exec(string(message))

				if err := conn.WriteMessage(websocket.TextMessage, []byte(data)); err != nil {
					return
				}

				if exit {
					break
				}
			}
		}()
	}
}
