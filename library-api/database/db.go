package database

import (
	"log"
	"github.com/example/library-api/models"
	"github.com/example/library-api/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("无法加载配置: %v", err)
	}

	// 连接SQLite数据库
	DB, err = gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("无法连接数据库: %v", err)
	}

	// 自动迁移数据表
	err = DB.AutoMigrate(
		&models.User{},
		&models.Book{},
		&models.Author{},
		&models.BookAuthor{},
		&models.BookCopy{},
		&models.Borrow{},
	)
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	log.Println("数据库连接和迁移成功")
}