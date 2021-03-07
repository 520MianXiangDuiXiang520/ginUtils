package middleware

import (
	"github.com/520MianXiangDuiXiang520/ginUtils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PermiterFunc func(context *gin.Context) bool

// 鉴权中间件， 注册该中间件后， 如果 license 返回 true 则表示
// 有权限访问该组（个）接口， 否则响应 403
func Permiter(license PermiterFunc) gin.HandlerFunc {
	return func(context *gin.Context) {
		if !license(context) {
			context.Abort()
			context.JSON(http.StatusForbidden, ginUtils.ForbiddenErrorRespHeader)
		}
	}
}
