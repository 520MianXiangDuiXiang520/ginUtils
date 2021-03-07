package websocket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/robfig/cron/v3"
	"log"
)

func ExampleTransformToWS() {
	go func() {
		c := cron.New(cron.WithSeconds())
		_, _ = c.AddFunc("* * * * * *", task)
		c.Start()
	}()
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var router = gin.Default()
	router.GET("/ws", TransformToWS(ping, nil, []WSConnWrapper{SaveConn}))

	if err := router.Run(":8888"); err != nil {
		log.Println(err)
	}
}

func SaveConn(ws *websocket.Conn) {
	webSocketConnList = append(webSocketConnList, ws)
}

var webSocketConnList = make([]*websocket.Conn, 0)

var canvas string
var canvasVersion int

func ping(ws *websocket.Conn) {
	for {
		_, message, err := ws.ReadMessage()
		fmt.Println(string(message))
		if err != nil {
			break
		}
		canvas += string(message)
		canvasVersion++
	}
}

var version = 0

func task() {
	if version == canvasVersion {
		return
	}
	for _, conn := range webSocketConnList {
		fmt.Println(canvas)
		err := conn.WriteMessage(websocket.TextMessage, []byte(canvas))
		if err != nil {
			break
		}
		version = canvasVersion
	}
}
