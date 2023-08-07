package ginz

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/ginz/crud"
	"github.com/qf0129/ginz/pkg/strs"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// 初始化
func Init(option *Option) (ginz *Ginz) {
	ginz = &Ginz{Engine: gin.New(), Option: option}
	ginz.ApiGroup = ginz.Group(option.DefaultGroupPrefix)

	option.InitValue()
	LoadLogger("debug")
	if option.LoadConfigFile {
		ginz.LoadConfig()
		LoadLogger(ginz.Config.LogLevel)
	}

	gin.SetMode(ginz.Config.AppMode)
	if option.ConnectDB {
		ginz.ConnectDB()
		crud.Init(ginz.DB)
	}

	if len(option.Middlewares) > 0 {
		for _, mid := range option.Middlewares {
			ginz.Engine.Use(func(ctx *gin.Context) { mid(ctx) })
		}
	}
	ginz.Engine.Use(func(ctx *gin.Context) {
		ctx.Set(REQUEST_KEY_ID, strs.UUID())
		// ctx.Next()
	})
	ginz.Engine.Use(gin.Logger(), gin.Recovery())

	if option.AddHealthCheckApi {
		ginz.Engine.GET("/health", func(ctx *gin.Context) { RespOk(ctx, "ok") })
	}
	return ginz
}

type Ginz struct {
	Engine    *gin.Engine
	DB        *gorm.DB
	Option    *Option
	Config    *Configuration
	ApiGroup  *ApiGroup
	ApiGroups []*ApiGroup
}

// 运行服务
func (ginz *Ginz) Run() {
	listenAddr := fmt.Sprintf("%s:%d", ginz.Config.AppHost, ginz.Config.AppPort)
	svr := &http.Server{
		Handler:      ginz.Engine,
		Addr:         listenAddr,
		ReadTimeout:  time.Duration(ginz.Config.AppTimeout) * time.Second,
		WriteTimeout: time.Duration(ginz.Config.AppTimeout) * time.Second,
	}
	logrus.Info("Run with " + ginz.Config.AppMode + " mode ")
	logrus.Info("Server is listening " + listenAddr)
	svr.ListenAndServe()
}

// 添加接口组
func (ginz *Ginz) Group(basePath string) *ApiGroup {
	group := &ApiGroup{
		BasePath:    basePath,
		RouterGroup: ginz.Engine.Group(basePath),
	}
	ginz.ApiGroups = append(ginz.ApiGroups, group)
	return group
}

// 默认接口组-添加中间件
func (ginz *Ginz) Use(middleware Middleware) {
	ginz.ApiGroup.Use(middleware)
}

// 默认接口组-添加接口
func (ginz *Ginz) AddApi(name string, handler ApiHandler) {
	ginz.ApiGroup.AddApi(name, handler)
}
