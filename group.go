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
	Name      string
	Info      string
	Method    string
	Handler   ApiHandler
	WithCache bool
}

type ApiGroup struct {
	BasePath    string
	RouterGroup *gin.RouterGroup
	Apis        []*Api
}

// 使用中间件
func (group *ApiGroup) Use(middleware Middleware) {
	group.RouterGroup.Use(func(ctx *gin.Context) {
		middleware(&Context{C: ctx})
	})
}

func (group *ApiGroup) Handle(httpMethod, relativePath string, handler ApiHandler) gin.IRoutes {
	return group.RouterGroup.Handle(httpMethod, relativePath, func(ctx *gin.Context) {
		handler(&Context{C: ctx, ReqId: strs.UUID()})
		return
		// if err != nil {
		// 	c.ReturnErr(err)
		// } else {
		// 	c.ReturnOk(data)
		// }
	})
}

func (group *ApiGroup) GET(relativePath string, handler ApiHandler) gin.IRoutes {
	return group.Handle(http.MethodGet, relativePath, handler)
}

func (group *ApiGroup) POST(relativePath string, handler ApiHandler) gin.IRoutes {
	return group.Handle(http.MethodPost, relativePath, handler)
}

// 添加api对象
func (group *ApiGroup) Add(api *Api) {
	group.Apis = append(group.Apis, api)
	group.Handle(api.Method, api.Name, api.Handler)
}

// 添加api
func (group *ApiGroup) AddApi(name string, handler ApiHandler) {
	group.Add(&Api{
		Name:    name,
		Method:  http.MethodPost,
		Handler: handler,
	})
}
