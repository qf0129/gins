package main

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/ginz"
	"github.com/qf0129/ginz/crud"
	"github.com/qf0129/ginz/simple"
	"gorm.io/gorm/logger"
)

var App *ginz.Ginz

func main() {
	App = ginz.Init(&ginz.Option{
		LoadConfigFile:    true,
		ConnectDB:         true,
		AddHealthCheckApi: true,
		DBLogLevel:        logger.Error,
		// Middlewares:       []ginz.Middleware{simple.RequireTokenFromCookie()},
	})

	App.MigrateModels(&simple.User{})

	group1 := App.Group("/api")
	group1.AddApi("login", simple.UserLoginHandler(App.Config.SecretKey, App.Config.TokenExpiredTime))
	group1.AddApi("register", simple.UserRegisterHandler())

	group2 := App.Group("/api")
	group2.Use(simple.RequireTokenFromCookie(App.Config.SecretKey, App.Config.TokenExpiredTime))
	group2.AddApi("test", func(c *gin.Context) (data any, err *ginz.Err) {
		data, er := crud.QueryAll[simple.User](nil, nil)
		if er != nil {
			err = ginz.ErrDBQueryFailed.Add(er.Error())
		}
		return
	})
	App.Run()
}
