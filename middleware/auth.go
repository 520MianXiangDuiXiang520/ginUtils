package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthFunc 用于身份验证，如果验证通过，应该将用户对象返回，反之第二个参数应该
// 返回 false，该类型只用于作为参数传递给认证中间件，具体用法请参考
// BaseAuthMiddleware 的示例。
type AuthFunc func(ctx *gin.Context) (user interface{}, authed bool)

// BaseAuthMiddleware 拦截请求，并使用 f 验证请求状态，如果验证不通过（f 返回 nil, false）
// 会直接响应 401(Unauthorized), 反之，验证通过该中间件会将请求状态（由 f 返回）保存在上下文的
// user 属性中，在应用逻辑中，你可以使用 context.Get("user") 获取到。
// 如果传入的 f 为空，该中间件不起任何作用；
// errResp 用于自定义认证失败时的响应数据，默认为 {"header": {"code": 401, "msg": "Unauthorized"}}.
func BaseAuthMiddleware(f AuthFunc, errResp interface{}) gin.HandlerFunc {
	return func(context *gin.Context) {
		if f == nil {
			context.Next()
			return
		}
		if errResp == nil {
			errResp = map[string]interface{}{
				"header": map[string]interface{}{
					"code": http.StatusUnauthorized,
					"msg":  "Unauthorized",
				},
			}
		}
		user, ok := f(context)
		if !ok {
			context.Abort()
			context.JSON(http.StatusUnauthorized, errResp)
			return
		}
		if user != nil {
			context.Set("user", user)
		}
		context.Next()
	}
}

// Deprecated: 为了兼容旧版本而存在，请使用 BaseAuthMiddleware
func Auth(af AuthFunc) gin.HandlerFunc {
	return BaseAuthMiddleware(af, nil)
}
