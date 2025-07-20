package models

// BookAuthor 书籍与作者的多对多关联模型
type BookAuthor struct {
	BookID   uint `gorm:"primaryKey" json:"book_id"`
	AuthorID uint `gorm:"primaryKey" json:"author_id"`
	Book     Book `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Author   Author `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
}