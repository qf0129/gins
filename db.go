package ginz

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

type DbOption struct {
	ShowLog bool
}

func (app *Ginz) ConnectDB() {
	var database *gorm.DB

	gormConf := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(app.Option.DBLogLevel),
		// NowFunc: func() time.Time {
		// 	return time.Now().Local()
		// },
	}
	// logrus.Info(fmt.Sprintf("DB log level is %d", app.Option.DBLogLevel))

	var dbConn gorm.Dialector
	if Config.DbEngine == "mysql" {
		dbUri := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", Config.DbUser, Config.DbPsd, Config.DbHost, Config.DbPort, Config.DbDatabase)
		dbConn = mysql.Open(dbUri)
		logrus.Info(fmt.Sprintf("Connected DB on MySQL: %s@%s", Config.DbUser, Config.DbHost))
	} else if Config.DbEngine == "sqlite" {
		dbConn = sqlite.Open(Config.SqliteFile)
		logrus.Info(fmt.Sprintf("Connected DB on Sqlite: %s", Config.SqliteFile))
	} else {
		logrus.Panic("InvalidDbType")
	}

	database, err := gorm.Open(dbConn, gormConf)
	if err != nil {
		panic("Failed to connect to database!")
	}

	DB = database
}

func MigrateModels(dst ...any) {
	if err := DB.AutoMigrate(dst...); err != nil {
		logrus.Panic("AutoMigrateErr:", err)
		return
	}
}
