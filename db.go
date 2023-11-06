package ginz

import (
	"fmt"
	"strings"

	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	GormConf *gorm.Config
	DB       *gorm.DB
)

// 连接数据库
func (ginz *Ginz) ConnectDB() {
	GormConf = &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		// Logger: logger.Default.LogMode(Config.DbLogLevel),
		// NowFunc: func() time.Time {
		// 	return time.Now().Local()
		// },
	}
	if strings.Contains(Config.DbLogLevel, "info") {
		GormConf.Logger = logger.Default.LogMode(logger.Info)
	} else if strings.Contains(Config.DbLogLevel, "error") {
		GormConf.Logger = logger.Default.LogMode(logger.Error)
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

	var err error
	DB, err = gorm.Open(dbConn, GormConf)
	if err != nil {
		panic("Failed to connect to database!")
	}
}

// 迁移数据模型
func (ginz *Ginz) MigrateModels(dst ...any) {
	if err := DB.AutoMigrate(dst...); err != nil {
		logrus.Panic("AutoMigrateErr:", err)
		return
	}
}
