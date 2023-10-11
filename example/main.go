package main

import (
	"github.com/qf0129/ginz"
	"github.com/qf0129/ginz/dao"
	"github.com/qf0129/ginz/pkg/errs"
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
		Models:            []any{&simple.User{}},
		// Middlewares:       []ginz.Middleware{simple.RequireTokenFromCookie()},
	})

	group1 := App.Group("/api")
	group1.POST("login", simple.UserLoginHandler())
	group1.POST("register", simple.UserRegisterHandler())
	simple.CreateRestApis[simple.User](group1)

	group2 := App.Group("/api")
	group2.Use(simple.RequireTokenFromCookie())
	group2.POST("test1", func(c *ginz.Context) {
		data, er := dao.QueryAll[simple.User]()
		if er != nil {
			c.ReturnErr(errs.RetrieveDataFailed.Add(er.Error()))
			return
		}
		c.ReturnOk(data)
	})
	App.Run()
}
