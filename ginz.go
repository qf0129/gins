package ginz

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/ginz/crud"
	"github.com/sirupsen/logrus"
)

type Ginz struct {
	Engine *gin.Engine
	Option *Option
}

func (ginz *Ginz) AddGroup(basePath string) *ApiGroup {
	return &ApiGroup{
		BasePath:    basePath,
		RouterGroup: ginz.Engine.Group(basePath),
	}
}

var App *Ginz
var DefaultGroup *ApiGroup

// 初始化
func Init(option *Option) {
	option.Load()
	LoadLogger()
	if option.LoadConfigFile {
		LoadConfigFile()
		LoadLogger()
	}

	App = &Ginz{
		Engine: gin.New(),
		Option: option,
	}

	gin.SetMode(Config.AppMode)
	if option.ConnectDB {
		App.ConnectDB()
	}
	crud.Init(DB)

	if len(option.Middlewares) > 0 {
		for _, mid := range option.Middlewares {
			App.Engine.Use(func(ctx *gin.Context) { mid(ctx) })
		}
	}
	// App.Engine.Use(gin.Logger(), gin.Recovery())
	DefaultGroup = App.AddGroup(option.DefaultGroupPrefix)

	if option.AddHealthCheckApi {
		addHealthCheckApi()
	}
}

// 使用中间件
func Use(middleware ...gin.HandlerFunc) {
	App.Engine.Use(middleware...)
}

// 运行服务
func Run() {
	listenAddr := fmt.Sprintf("%s:%d", Config.AppHost, Config.AppPort)
	svr := &http.Server{
		Handler:      App.Engine,
		Addr:         listenAddr,
		ReadTimeout:  time.Duration(Config.AppTimeout) * time.Second,
		WriteTimeout: time.Duration(Config.AppTimeout) * time.Second,
	}
	logrus.Info("Run with " + Config.AppMode + " mode ")
	logrus.Info("Server is listening " + listenAddr)
	svr.ListenAndServe()
}

// 默认路由组添加接口
func AddApi(name string, handler ApiHandler) {
	DefaultGroup.Add(&Api{
		Name:    name,
		Method:  http.MethodPost,
		Handler: handler,
	})
}

// 添加健康检测接口
func addHealthCheckApi() {
	DefaultGroup.Add(&Api{
		Name:   "health",
		Info:   "HealthCheck",
		Method: http.MethodGet,
		Handler: func(c *gin.Context) (data any, err *Err) {
			data = "ok"
			return
		},
	})
}
