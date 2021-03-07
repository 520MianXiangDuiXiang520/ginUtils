package ginUtils

import (
	"github.com/gin-gonic/gin"
)

type DoChildRouteFunc func(g *gin.RouterGroup) // 子路由

// 用于适配 Engine 和 RouteGroup
type RouteAdapter interface {
	Use(middleware ...gin.HandlerFunc) gin.IRoutes
	Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup
}

// URL 调度器
func URLPatterns(route RouteAdapter, path string, childRoute DoChildRouteFunc, middles ...gin.HandlerFunc) {
	group := route.Group(path)
	for _, mid := range middles {
		group.Use(mid)
	}
	childRoute(group)
}
