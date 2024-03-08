package model

import (
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func NewGorm(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		logx.Error("mysql数据库连接失败：", err)
		return nil
	}
	logx.Info("mysql数据库连接成功！")

	AutoMigrate(db)

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Minute * 5)

	return db
}

func AutoMigrate(db *gorm.DB) {
	_ = db.AutoMigrate(
		&TdFirm{},
	)
}

var ModelList = map[string]interface{}{
	"TdFirm": TdFirm{},
}
