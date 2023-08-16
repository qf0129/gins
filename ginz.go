package ginz

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/ginz/pkg/dao"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Middleware func(*Context)

type Option struct {
	LoadConfigFile     bool
	ConnectDB          bool
	AddHealthCheckApi  bool
	DBLogLevel         logger.LogLevel
	DefaultGroupPrefix string
	PrimaryKey         string
	DefaultPageSize    int
	Middlewares        []Middleware
	Models             []any
}

func (option *Option) InitValue() {
	if option.DBLogLevel == 0 {
		option.DBLogLevel = logger.Info
	}
}

// 初始化
func Init(option *Option) (ginz *Ginz) {
	ginz = &Ginz{Engine: gin.New(), Option: option}
	ginz.ApiGroup = ginz.Group(option.DefaultGroupPrefix)

	option.InitValue()
	LoadLogger("debug")
	if option.LoadConfigFile {
		LoadConfig()
		LoadLogger(Config.LogLevel)
	}

	gin.SetMode(Config.AppMode)

	if option.ConnectDB {
		ginz.ConnectDB()
		dao.Init(&dao.DaoOption{
			DB:              ginz.DB,
			DefaultPageSize: Config.DefaultPageSize,
			QueryPrimaryKey: Config.QueryPrimaryKey,
		})
		ginz.MigrateModels(option.Models...)
	}

	if len(option.Middlewares) > 0 {
		for _, mid := range option.Middlewares {
			ginz.Use(mid)
		}
	}
	ginz.Engine.Use(gin.Logger(), gin.Recovery())

	if option.AddHealthCheckApi {
		ginz.ApiGroup.GET("/health", func(c *Context) { c.ReturnOk("ok") })
	}

	return ginz
}

type Ginz struct {
	Engine *gin.Engine
	DB     *gorm.DB
	Option *Option
	// Config    *Configuration
	ApiGroup  *ApiGroup
	ApiGroups []*ApiGroup
}

// 运行服务
func (ginz *Ginz) Run() {
	listenAddr := fmt.Sprintf("%s:%d", Config.AppHost, Config.AppPort)
	svr := &http.Server{
		Handler:      ginz.Engine,
		Addr:         listenAddr,
		ReadTimeout:  time.Duration(Config.AppTimeout) * time.Second,
		WriteTimeout: time.Duration(Config.AppTimeout) * time.Second,
	}
	logrus.Info("Run with " + Config.AppMode + " mode ")
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
