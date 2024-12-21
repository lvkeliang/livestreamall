package dao

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"livestreamall/model"
	"log"
)

var DB *gorm.DB

func InitDB() {
	// 配置数据库连接字符串
	// 开发时用本地的数据库
	// dsn := "root:123456@tcp(127.0.0.1:3306)/livestream?charset=utf8mb4&parseTime=True&loc=Local"
	// 部署时使用
	dsn := "LvKeliang:lkl12190711@tcp(localhost:3306)/livestream?charset=utf8mb4&parseTime=True&loc=Local"
	// 开发时用云端的数据库
	// dsn := "LvKeliang:lkl12190711@tcp(120.79.27.213:3306)/livestream?charset=utf8mb4&parseTime=True&loc=Local"

	var err error

	// 使用 GORM 打开 MySQL 连接
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info), // 设置日志级别
		DisableForeignKeyConstraintWhenMigrating: false,
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 检查连接是否成功
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("failed to get database instance: %v", err)
	}
	if err = sqlDB.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	fmt.Println("数据库连接成功")

	MigrateModels()
}

// MigrateModels 自动迁移模型到数据库
func MigrateModels() {
	err := DB.AutoMigrate(&model.User{}, &model.LiveRoom{}, &model.Message{})
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}
}
