package middleware

import (
	"github.com/520MianXiangDuiXiang520/ginUtils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// accessList: 允许访问的白名单
func CorsHandler(accessList []string) gin.HandlerFunc {
	return func(context *gin.Context) {
		origin := context.GetHeader("Origin")
		method := context.Request.Method
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		for _, allow := range accessList {
			if allow == origin {
				context.Header("Access-Control-Allow-Origin", origin)
				break
			}
		}
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		context.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,"+
			"session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT,"+
			" X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type,"+
			" Pragma,token,openid,opentoken, cookies, Cookies, cookie, Cookies")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, "+
			"Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
		context.Header("Access-Control-Max-Age", "172800")
		context.Header("Access-Control-Allow-Credentials", "true")
		context.Set("content-type", "application/json")
		// 设置返回格式是json
		if method == "OPTIONS" {
			context.Abort()
			context.JSON(http.StatusOK, ginUtils.SuccessRespHeader)
		}
		context.Next()
	}
}
