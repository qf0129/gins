package ginz

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 连接数据库
func (ginz *Ginz) ConnectDB() {
	gormConf := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(ginz.Option.DBLogLevel),
		// NowFunc: func() time.Time {
		// 	return time.Now().Local()
		// },
	}
	// logrus.Info(fmt.Sprintf("DB log level is %d", app.Option.DBLogLevel))

	var dbConn gorm.Dialector
	if ginz.Config.DbEngine == "mysql" {
		dbUri := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", ginz.Config.DbUser, ginz.Config.DbPsd, ginz.Config.DbHost, ginz.Config.DbPort, ginz.Config.DbDatabase)
		dbConn = mysql.Open(dbUri)
		logrus.Info(fmt.Sprintf("Connected DB on MySQL: %s@%s", ginz.Config.DbUser, ginz.Config.DbHost))
	} else if ginz.Config.DbEngine == "sqlite" {
		dbConn = sqlite.Open(ginz.Config.SqliteFile)
		logrus.Info(fmt.Sprintf("Connected DB on Sqlite: %s", ginz.Config.SqliteFile))
	} else {
		logrus.Panic("InvalidDbType")
	}

	var err error
	ginz.DB, err = gorm.Open(dbConn, gormConf)
	if err != nil {
		panic("Failed to connect to database!")
	}
}

// 迁移数据模型
func (ginz *Ginz) MigrateModels(dst ...any) {
	if err := ginz.DB.AutoMigrate(dst...); err != nil {
		logrus.Panic("AutoMigrateErr:", err)
		return
	}
}
