package models

import (
	"time"

	"gorm.io/gorm"
)

// BookCopyStatus 定义图书副本状态
type BookCopyStatus string

const (
	CopyAvailable BookCopyStatus = "available"
	CopyBorrowed  BookCopyStatus = "borrowed"
	CopyLost      BookCopyStatus = "lost"
	CopyMaintenance BookCopyStatus = "maintenance"
)

// BookCopy 图书副本模型
type BookCopy struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	BookID         uint           `gorm:"not null" json:"book_id"`
	CopyNumber     string         `gorm:"size:20;not null" json:"copy_number"`
	Status         BookCopyStatus `gorm:"size:20;not null;default:available" json:"status"`
	AcquisitionDate time.Time      `json:"acquisition_date"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Book           Book           `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Borrows        []Borrow       `gorm:"foreignKey:BookCopyID" json:"borrows,omitempty"`
}