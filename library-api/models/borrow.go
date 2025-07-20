package models

import (
	"time"

	"gorm.io/gorm"
)

// BorrowStatus 定义借阅状态
type BorrowStatus string

const (
	BorrowActive   BorrowStatus = "active"
	BorrowReturned BorrowStatus = "returned"
	BorrowOverdue  BorrowStatus = "overdue"
	BorrowLost     BorrowStatus = "lost"
)

// Borrow 借阅模型
type Borrow struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"not null" json:"user_id"`
	BookCopyID    uint           `gorm:"not null" json:"book_copy_id"`
	BorrowDate    time.Time      `json:"borrow_date"`
	DueDate       time.Time      `json:"due_date"`
	ReturnDate    *time.Time     `json:"return_date,omitempty"`
	Status        BorrowStatus   `gorm:"size:20;not null;default:active" json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	User          User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	BookCopy      BookCopy       `gorm:"foreignKey:BookCopyID" json:"book_copy,omitempty"`
}