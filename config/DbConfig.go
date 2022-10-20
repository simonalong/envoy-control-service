package config

import (
	"github.com/isyscore/isc-gobase/logger"
	"github.com/isyscore/isc-gobase/orm"
	"gorm.io/gorm"
)

var Db *gorm.DB

func init() {
	logger.Info("连接数据库")
	dbTem, err := orm.GetGormDb()
	if err != nil {
		logger.Error("连接数据库失败, %v", err.Error())
	}
	logger.Info("连接成功")
	Db = dbTem
}
