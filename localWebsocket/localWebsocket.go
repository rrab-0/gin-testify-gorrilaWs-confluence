package localWebsocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	WebsocketConn *websocket.Conn
)

func Listener(c *gin.Context) {
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
	}
	defer wsConn.Close()

	WebsocketConn = wsConn // store ws conn globally

	for {
		_, _, err := wsConn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
	}
}

func Writer(message string) {
	if WebsocketConn != nil {
		err := WebsocketConn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("Error sending message to websocket:", err)
		}
	}
}
