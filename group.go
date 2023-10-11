package ginz

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/ginz/pkg/strs"
)

// type ApiHandler func(c *gin.Context) (any, *errs.Err)
// type ApiHandler func(c *Context) (any, *errs.Err)
type ApiHandler func(c *Context)

type Api struct {
	Method   string
	Path     string
	Handler  ApiHandler
	UseCache bool
}

type ApiGroup struct {
	BasePath    string
	RouterGroup *gin.RouterGroup
	Apis        []*Api
}

// 使用中间件
func (group *ApiGroup) Use(middleware Middleware) {
	group.RouterGroup.Use(func(ctx *gin.Context) {
		middleware(&Context{C: ctx, ReqId: strs.UUID()})
	})
}

func (group *ApiGroup) Handle(httpMethod, relativePath string, handler ApiHandler) gin.IRoutes {
	group.Apis = append(group.Apis, &Api{
		Method:  httpMethod,
		Path:    relativePath,
		Handler: handler,
	})
	return group.RouterGroup.Handle(httpMethod, relativePath, func(ctx *gin.Context) {
		handler(&Context{C: ctx, ReqId: strs.UUID()})
	})
}

func (group *ApiGroup) GET(relativePath string, handler ApiHandler) gin.IRoutes {
	return group.Handle(http.MethodGet, relativePath, handler)
}

func (group *ApiGroup) POST(relativePath string, handler ApiHandler) gin.IRoutes {
	return group.Handle(http.MethodPost, relativePath, handler)
}

func (group *ApiGroup) DELETE(relativePath string, handlers ApiHandler) gin.IRoutes {
	return group.Handle(http.MethodDelete, relativePath, handlers)
}

func (group *ApiGroup) PATCH(relativePath string, handlers ApiHandler) gin.IRoutes {
	return group.Handle(http.MethodPatch, relativePath, handlers)
}

func (group *ApiGroup) PUT(relativePath string, handlers ApiHandler) gin.IRoutes {
	return group.Handle(http.MethodPut, relativePath, handlers)
}

func (group *ApiGroup) OPTIONS(relativePath string, handlers ApiHandler) gin.IRoutes {
	return group.Handle(http.MethodOptions, relativePath, handlers)
}

func (group *ApiGroup) HEAD(relativePath string, handlers ApiHandler) gin.IRoutes {
	return group.Handle(http.MethodHead, relativePath, handlers)
}
