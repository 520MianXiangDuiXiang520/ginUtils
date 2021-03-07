package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type WSHandlerFunc func(ws *websocket.Conn)
type WSConnWrapper func(ws *websocket.Conn)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// TransformToWS 将 Gin 的 HTTP 连接转换为 WebSocket 连接
// 他接受一个 WSHandlerFunc 用来处理具体的 WebSocket 事务。
func TransformToWS(wsFunc WSHandlerFunc, responseHeader http.Header, wrappers []WSConnWrapper) gin.HandlerFunc {
	return func(context *gin.Context) {
		ws, err := upGrader.Upgrade(context.Writer, context.Request, responseHeader)
		for _, wrapper := range wrappers {
			wrapper(ws)
		}
		if err != nil {
			log.Println(err)
			return
		}
		defer ws.Close()
		wsFunc(ws)
	}
}
