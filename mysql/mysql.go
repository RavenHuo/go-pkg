package mysql

import (
	"context"
	"github.com/RavenHuo/go-pkg/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var mysqlDB *gorm.DB

func Init(config *Config) {
	connectDB(config)
}

func connectDB(config *Config) {
	dsn := config.GetDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:         "",
			SingularTable:       true,
			NoLowerCase:         false,
			IdentifierMaxLength: 0,
		},
	})
	ctx := context.Background()
	if err != nil {
		log.Errorf(ctx, "connect failed dsn:%s, err:%s", dsn, err)
		panic("connect failed, error=" + err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf(ctx, "getDB failed dsn:%s, err:%s", dsn, err)
		panic("getDB, error=" + err.Error())
	}
	//设置数据库连接池参数
	maxOpenConns := 100
	if config.MaxOpenConns != 0 {
		maxOpenConns = config.MaxOpenConns
	}
	maxIdleConns := 20
	if config.MaxIdleConns != 0 {
		maxIdleConns = config.MaxIdleConns
	}
	sqlDB.SetMaxOpenConns(maxOpenConns) //设置数据库连接池最大连接数
	sqlDB.SetMaxIdleConns(maxIdleConns) //连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于最大空闲数，超过的连接会被连接池关闭。

	mysqlDB = db
}

func GetDB() *gorm.DB {
	return mysqlDB
}
