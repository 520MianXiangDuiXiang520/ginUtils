package ginUtils

import (
	"github.com/gin-gonic/gin"
	"testing"
	"time"
)

type resp struct {
	msg string
}

func TestURLPatterns(t *testing.T) {
	// /ping
	// /v1
	//    /ping
	//    /pong
	// v2/ping/
	//    /index
	//    /home
	go func() {
		engine := gin.Default()
		defer engine.Run()
		URLPatterns(engine, "/ping", func(g *gin.RouterGroup) {
			g.POST("/", func(context *gin.Context) {
				context.JSON(200, resp{msg: "pong"})
			})
		})

		URLPatterns(engine, "v1/", func(g *gin.RouterGroup) {
			URLPatterns(g, "ping/", func(g *gin.RouterGroup) {
				g.POST("/", func(context *gin.Context) {
					context.JSON(200, resp{msg: "v1:pong"})
				})
			})
			URLPatterns(g, "pong/", func(g *gin.RouterGroup) {
				g.POST("/", func(context *gin.Context) {
					context.JSON(200, resp{msg: "v1:ping"})
				})
			})
		})

		URLPatterns(engine, "v2/ping/", func(g *gin.RouterGroup) {
			g.POST("index/", func(context *gin.Context) {
				context.JSON(200, resp{msg: "v2:ping index"})
			})
			g.POST("home/", func(context *gin.Context) {
				context.JSON(200, resp{msg: "v2:ping home"})
			})
		})
	}()
	time.Sleep(time.Second * 3)

}
