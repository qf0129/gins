package ginz

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

	if option.ConnectDB {
		App.ConnectDB()
	}

	gin.SetMode(Config.AppMode)

	App.Engine.Use(gin.Logger(), gin.Recovery())
	DefaultGroup = App.AddGroup(option.DefaultGroupPrefix)
	if option.AddHealthCheckApi {
		addHealthCheckApi()
	}
}

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

func AddApi(name string, info string, handler ApiHandler) {
	DefaultGroup.Add(&Api{
		Name:    name,
		Info:    info,
		Method:  http.MethodPost,
		Handler: handler,
	})
}

func addHealthCheckApi() {
	DefaultGroup.Add(&Api{
		Name:   "health",
		Info:   "HealthCheck",
		Method: http.MethodGet,
		Handler: func(c *gin.Context) (data any, err *Errors) {
			data = "ok"
			return
		},
	})
}
