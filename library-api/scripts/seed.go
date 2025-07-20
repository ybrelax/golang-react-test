package main

import (
	"log"
	"time"

	"github.com/example/library-api/config"
	"github.com/example/library-api/database"
	"github.com/example/library-api/models"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 加载配置
	_, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库连接
	database.InitDB()
	db := database.DB

	// 清空现有测试数据
	db.Delete(&models.Borrow{})
	db.Delete(&models.BookCopy{})
	db.Delete(&models.BookAuthor{})
	db.Delete(&models.Book{})
	db.Delete(&models.Author{})
	db.Delete(&models.User{})

	// 创建管理员用户
	adminPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := models.User{
		Username: "admin",
		Email:    "admin@example.com",
		Password: string(adminPassword),
		Role:     models.RoleAdmin,
	}
	db.Create(&admin)

	// 创建普通用户
	userPassword, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	user := models.User{
		Username: "user",
		Email:    "user@example.com",
		Password: string(userPassword),
		Role:     models.RoleUser,
	}
	db.Create(&user)

	// 创建作者
	authors := []models.Author{
		{Name: "J.K. Rowling", Bio: "英国作家，《哈利·波特》系列作者"},
		{Name: "乔治·奥威尔", Bio: "英国小说家、散文家"},
		{Name: "J.R.R.托尔金", Bio: "英国作家、语言学家"},
	}
	for _, a := range authors {
		db.Create(&a)
	}

	// 创建书籍
	books := []models.Book{
		{
			Title:           "哈利·波特与魔法石",
			ISBN:            "9780747532743",
			Description:     "哈利·波特系列第一部",
			Publisher:       "Bloomsbury",
			PublicationDate: time.Date(1997, time.June, 26, 0, 0, 0, 0, time.UTC),
		},
		{
			Title:           "1984",
			ISBN:            "9780451524935",
			Description:     "反乌托邦经典小说",
			Publisher:       "Signet Classic",
			PublicationDate: time.Date(1949, time.June, 8, 0, 0, 0, 0, time.UTC),
		},
		{
			Title:           "魔戒",
			ISBN:            "9780618640157",
			Description:     "奇幻文学经典",
			Publisher:       "Houghton Mifflin",
			PublicationDate: time.Date(1954, time.July, 29, 0, 0, 0, 0, time.UTC),
		},
	}
	for i, b := range books {
		db.Create(&b)
		// 创建书籍-作者关联
		db.Create(&models.BookAuthor{BookID: b.ID, AuthorID: uint(i + 1)})
	}

	// 创建图书副本
	for bookID := 1; bookID <= 3; bookID++ {
		for i := 1; i <= 3; i++ {
			copy := models.BookCopy{
				BookID:          uint(bookID),
				CopyNumber:      string(i),
				Status:          "available", // 假设书籍副本状态为字符串 "available"，请根据实际模型定义修改
				AcquisitionDate: time.Now().AddDate(-i, 0, 0),
			}
			db.Create(&copy)
		}
	}

	// 创建借阅记录
	borrows := []models.Borrow{
		{
			UserID:     2,
			BookCopyID: 1,
			BorrowDate: time.Now().AddDate(0, 0, -14),
			DueDate:    time.Now().AddDate(0, 0, 7),
			Status:     "active", // 假设状态字符串为 "active"，需根据实际情况替换为正确的状态值
		},
		{
			UserID:     2,
			BookCopyID: 4,
			BorrowDate: time.Now().AddDate(0, 0, -7),
			DueDate:    time.Now().AddDate(0, 0, 14),
			Status:     "active",
		},
	}
	for _, br := range borrows {
		db.Create(&br)
	}

	log.Println("假数据插入成功!")
}
