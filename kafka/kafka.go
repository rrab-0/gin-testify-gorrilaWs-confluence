package kafka

import (
	"example/unit-test-hello-world/websocket"
	"net/http"

	"github.com/gin-gonic/gin"
)

type KafkaMsg struct {
	Message string `form:"msg"`
}

func ReadMessageThenSendToSlashWS(c *gin.Context) {
	var msg KafkaMsg

	if c.ShouldBind(&msg) == nil {
		websocket.SendMessage(msg.Message)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "msg from kafka to ws sent.",
	})
}
