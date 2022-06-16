package config

import (
	"fmt"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func GetDb() *gorm.DB {
	if Db != nil {
		return Db
	}
	dbCfg := DbConfig{}
	err := config.GetValueObject("database", &dbCfg)
	if err != nil {
		logger.Warn("读取db配置异常")
		return nil
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/biz_envoy?charset=utf8&parseTime=True&loc=Local", dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port)
	dbTem, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	Db = dbTem
	return Db
}

type DbConfig struct {
	Username string
	Password string
	Host     string
	Port     int
}
