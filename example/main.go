package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qf0129/ginz"
	"gorm.io/gorm/logger"
)

type BaseModel struct {
	Id    uint      `gorm:"primaryKey;"`
	Ctime time.Time `gorm:"autoCreateTime;comment:'Created Time'" `
	Utime time.Time `gorm:"autoUpdateTime;comment:'Updated Time'" `
}

type User struct {
	BaseModel
	Name string
	Age  uint
}

func test1(c *gin.Context) (data any, err *ginz.Errors) {
	var users []User
	// new := &User{
	// 	Name: "zs",
	// }
	// ginz.DB.Create(new)
	ginz.DB.Model(&User{}).Where("id=999").Find(&users)
	data = users
	return
}

func main() {
	ginz.Init(&ginz.Option{
		LoadConfigFile:    true,
		ConnectDB:         true,
		AddHealthCheckApi: true,
		DBLogLevel:        logger.Error,
	})
	ginz.MigrateModels(&User{})
	ginz.AddApi("add", "", ginz.CreateOneHandler[User]())
	ginz.AddApi("get", "", test1)
	ginz.Run()
}
