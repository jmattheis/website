package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jmattheis/website/content"
	"github.com/jmattheis/website/util"
)

var upgrader = websocket.Upgrader{}

func Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tty := &content.InteractiveText{
			Prompt:     "",
			RemoteAddr: util.GetRemoteAddr(r),
		}
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
