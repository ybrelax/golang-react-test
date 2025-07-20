package models

import (
	"time"

	"gorm.io/gorm"
)

// Author 作者模型
type Author struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Bio       string         `gorm:"type:text" json:"bio,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Books     []Book         `gorm:"many2many:book_authors" json:"books,omitempty"`
}