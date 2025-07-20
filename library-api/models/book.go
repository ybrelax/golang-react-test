package models

import (
	"time"

	"gorm.io/gorm"
)

// BookStatus 定义图书状态
type BookStatus string

const (
	StatusAvailable BookStatus = "available"
	StatusBorrowed  BookStatus = "borrowed"
	StatusLost      BookStatus = "lost"
)

// Book 图书模型
type Book struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Title           string         `gorm:"size:200;not null" json:"title"`
	ISBN            string         `gorm:"size:20;not null;unique" json:"isbn"`
	Description     string         `gorm:"type:text" json:"description,omitempty"`
	Publisher       string         `gorm:"size:100" json:"publisher,omitempty"`
	PublicationDate time.Time      `json:"publication_date,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Authors         []Author       `gorm:"many2many:book_authors" json:"authors,omitempty"`
	Copies          []BookCopy     `gorm:"foreignKey:BookID" json:"copies,omitempty"`
}