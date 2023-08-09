package main

import (
	"github.com/gin-gonic/gin"
	"github.com/qf0129/ginz"
	"github.com/qf0129/ginz/pkg/dao"
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
	group1.AddApi("login", simple.UserLoginHandler(App))
	group1.AddApi("register", simple.UserRegisterHandler(App))
	group1.AddApi("test2", func(c *gin.Context) (data any, err *ginz.Err) {

		return
	})

	group2 := App.Group("/api")
	group2.Use(simple.RequireTokenFromCookie(App))
	group2.AddApi("test1", func(c *gin.Context) (data any, err *ginz.Err) {
		data, er := dao.QueryAll[simple.User](nil)
		if er != nil {
			err = ginz.ErrDBQueryFailed.Add(er.Error())
		}
		return
	})
	App.Run()
}
