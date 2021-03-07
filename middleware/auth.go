package middleware

import (
	"github.com/520MianXiangDuiXiang520/ginUtils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserBase interface {
	GetID() int
}

// 检查授权，如果授权通过，返回授权用户，否则第二个参数返回 false
type AuthFunc func(context *gin.Context) (UserBase, bool)

// 授权中间件，注册使用中间件后，如果授权未通过（af() return nil, false）
// 请求会被在此拦截并响应 301，反之， 如果授权通过，会在请求上下文对象 context
// 中添加一个 user 字段，保存授权的用户信息，授权后可以使用 ctx.Get("user")
// 经过类型转换后获取到该 User 对象
func Auth(af AuthFunc) gin.HandlerFunc {
	return func(context *gin.Context) {
		user, ok := af(context)
		if !ok {
			context.Abort()
			context.JSON(http.StatusUnauthorized,
				ginUtils.UnauthorizedRespHeader)
		}
		context.Set("user", user)
	}
}
