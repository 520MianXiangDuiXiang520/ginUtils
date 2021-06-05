package middleware

import (
	"github.com/520MianXiangDuiXiang520/ginUtils"
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
)

type ExampleUser struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// XExampleAuthFunc 是 AuthFunc 的一个示例函数,它根据 token 检查请求状态
func XExampleAuthFunc(ctx *gin.Context) (interface{}, bool) {
	token := ctx.GetHeader("token")
	if token == "" {
		return nil, false
	}

	// select users by token in the database
	// ...

	return &ExampleUser{
		ID:   1,
		Name: "example",
	}, true
}

// ExampleBaseAuthMiddleware 演示 BaseAuthMiddleware 在原生 gin 中的用法
func ExampleBaseAuthMiddleware() {
	engine := gin.Default()
	defer engine.Run()
	engine.Use(BaseAuthMiddleware(XExampleAuthFunc, nil))
	engine.POST("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, "pong")
	})
}

// ExampleBaseAuthMiddleware_ginUtils 演示 BaseAuthMiddleware 配合
// ginUtils.URLPatterns 的用法.
func ExampleBaseAuthMiddleware_ginUtils() {
	engine := gin.Default()
	defer engine.Run()
	router := func(g *gin.RouterGroup) {
		g.POST("/",
			// auth
			BaseAuthMiddleware(XExampleAuthFunc, nil),
			func(context *gin.Context) {
				context.JSON(200, "pong")
			})
	}
	ginUtils.URLPatterns(engine, "/ping", router)
}

func TestBaseAuthMiddleware(t *testing.T) {
	serverOk := make(chan struct{}, 1)
	// 开启服务端
	go func() {
		engine := gin.Default()
		defer func() {
			serverOk <- struct{}{}
			engine.Run()
		}()
		engine.Use(BaseAuthMiddleware(XExampleAuthFunc, nil))
		engine.POST("/ping", func(context *gin.Context) {
			context.JSON(http.StatusOK, "pong")
		})
	}()

	select {
	case <-serverOk:
		// fail to auth
		resp, err := authPost(false)
		if err != nil {
			t.Error(err)
		}
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("got response code: %d, not 401", resp.StatusCode)
		}

		// success to auth
		resp, err = authPost(true)
		if err != nil {
			t.Error(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("got response code: %d, not 200", resp.StatusCode)
		}
	}
}

var authClient = http.Client{}

func authPost(addToken bool) (*http.Response, error) {
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/ping", nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"+
		" AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.190 Safari/537.36")
	req.Header.Add("X-Forwarded-For", "127.0.0.1")
	req.Header.Add("X-Real-Ip", "127.0.0.1")
	if addToken {
		req.Header.Add("token", "token")
	}
	return authClient.Do(req)
}
