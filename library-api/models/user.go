package models

import (
	"time"

	"gorm.io/gorm"
)

// UserRole 定义用户角色
type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	GoogleID  string         `gorm:"size:100;uniqueIndex" json:"google_id,omitempty"`
	Username  string         `gorm:"size:50;not null;unique" json:"username"`
	Email     string         `gorm:"size:100;not null;unique" json:"email"`
	Password  string         `gorm:"size:100;not null" json:"-"` // 密码不返回给前端
	Role      UserRole       `gorm:"size:20;not null;default:user" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Borrows   []Borrow       `gorm:"foreignKey:UserID" json:"borrows,omitempty"`
}
