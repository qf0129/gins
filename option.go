package ginz

import (
	"gorm.io/gorm/logger"
)

type Option struct {
	LoadConfigFile     bool
	ConnectDB          bool
	AddHealthCheckApi  bool
	DBLogLevel         logger.LogLevel
	DefaultGroupPrefix string
	PrimaryKey         string
	DefaultPageSize    int
}

func (option *Option) Load() {
	if option.DefaultGroupPrefix == "" {
		option.DefaultGroupPrefix = "/api"
	}
	if option.DBLogLevel == 0 {
		option.DBLogLevel = logger.Info
	}
}
