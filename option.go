package ginz

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
)

type Middleware func(*gin.Context)

type Option struct {
	LoadConfigFile     bool
	ConnectDB          bool
	AddHealthCheckApi  bool
	DBLogLevel         logger.LogLevel
	DefaultGroupPrefix string
	PrimaryKey         string
	DefaultPageSize    int
	Middlewares        []Middleware
}

func (option *Option) Load() {
	if option.DefaultGroupPrefix == "" {
		option.DefaultGroupPrefix = "/api"
	}
	if option.DBLogLevel == 0 {
		option.DBLogLevel = logger.Info
	}
}
