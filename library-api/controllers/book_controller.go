package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/example/library-api/database"
	"github.com/example/library-api/models"
)

// BookRequest 图书创建/更新请求结构
type BookRequest struct {
	Title           string    `json:"title" binding:"required"`
	ISBN            string    `json:"isbn" binding:"required"`
	Description     string    `json:"description"`
	Publisher       string    `json:"publisher"`
	PublicationDate string    `json:"publication_date"`
	AuthorIDs       []uint    `json:"author_ids" binding:"required"`
}

// BookCopyRequest 图书副本创建请求结构
type BookCopyRequest struct {
	BookID     uint `json:"book_id" binding:"required"`
	CopiesCount int `json:"copies_count" binding:"required,min=1"`
}

// BorrowBookRequest 借阅图书请求结构
type BorrowBookRequest struct {
	BookID uint `json:"book_id" binding:"required"`
	Days   int  `json:"days" binding:"required,min=1,max=30"`
}

// ReturnBookRequest 归还图书请求结构
type ReturnBookRequest struct {
	BorrowID uint `json:"borrow_id" binding:"required"`
}

// CreateBook 添加新图书
func CreateBook(c *gin.Context) {
	var req BookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 解析出版日期
	publicationDate, err := time.Parse("2006-01-02", req.PublicationDate)
	if err != nil && req.PublicationDate != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid publication date format, should be YYYY-MM-DD"})
		return
	}

	// 检查ISBN是否已存在
	var existingBook models.Book
	if result := database.DB.Where("isbn = ?", req.ISBN).First(&existingBook); result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "book with this ISBN already exists"})
		return
	}

	// 创建图书
	book := models.Book{
		Title:           req.Title,
		ISBN:            req.ISBN,
		Description:     req.Description,
		Publisher:       req.Publisher,
		PublicationDate: publicationDate,
	}

	// 开始事务
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建图书
	if err := tx.Create(&book).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create book"})
		return
	}

	// 添加作者关联
	var authors []models.Author
	if err := tx.Find(&authors, req.AuthorIDs).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find authors"})
		return
	}

	if len(authors) != len(req.AuthorIDs) {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "one or more author IDs are invalid"})
		return
	}

	if err := tx.Model(&book).Association("Authors").Append(authors); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add authors to book"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	// 预加载作者信息并返回
	tx.Preload("Authors").First(&book)
	c.JSON(http.StatusCreated, book)
}

// GetBooks 获取图书列表
func GetBooks(c *gin.Context) {
	var books []models.Book
	result := database.DB.Preload("Authors").Find(&books)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch books"})
		return
	}

	c.JSON(http.StatusOK, books)
}

// GetBook 获取单本图书详情
func GetBook(c *gin.Context) {
	id := c.Param("id")

	var book models.Book
	result := database.DB.Preload("Authors").Preload("Copies").First(&book, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// UpdateBook 更新图书信息
func UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var req BookRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找图书
	var book models.Book
	if result := database.DB.First(&book, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	// 解析出版日期
	publicationDate, err := time.Parse("2006-01-02", req.PublicationDate)
	if err != nil && req.PublicationDate != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid publication date format, should be YYYY-MM-DD"})
		return
	}

	// 检查ISBN是否已被其他图书使用
	if book.ISBN != req.ISBN {
		var existingBook models.Book
		if result := database.DB.Where("isbn = ? AND id != ?", req.ISBN, id).First(&existingBook); result.Error == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "book with this ISBN already exists"})
			return
		}
	}

	// 开始事务
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新图书信息
	book.Title = req.Title
	book.ISBN = req.ISBN
	book.Description = req.Description
	book.Publisher = req.Publisher
	book.PublicationDate = publicationDate

	if err := tx.Save(&book).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update book"})
		return
	}

	// 更新作者关联
	var authors []models.Author
	if err := tx.Find(&authors, req.AuthorIDs).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find authors"})
		return
	}

	if len(authors) != len(req.AuthorIDs) {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "one or more author IDs are invalid"})
		return
	}

	// 替换作者关联
	if err := tx.Model(&book).Association("Authors").Replace(authors); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update authors"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	// 预加载作者信息并返回
	database.DB.Preload("Authors").First(&book)
	c.JSON(http.StatusOK, book)
}

// DeleteBook 删除图书
func DeleteBook(c *gin.Context) {
	id := c.Param("id")

	// 查找图书
	var book models.Book
	if result := database.DB.First(&book, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	// 检查是否有关联的借阅记录
	var borrowCount int64
	if err := database.DB.Model(&models.Borrow{}).
		Joins("JOIN book_copies ON borrows.book_copy_id = book_copies.id").
		Where("book_copies.book_id = ? AND borrows.status = ?", id, models.BorrowActive).
		Count(&borrowCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check active borrows"})
		return
	}

	if borrowCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "cannot delete book with active borrows"})
		return
	}

	// 删除图书
	if err := database.DB.Delete(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book deleted successfully"})
}

// AddBookCopies 添加图书副本
func AddBookCopies(c *gin.Context) {
	var req BookCopyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查图书是否存在
	var book models.Book
	if result := database.DB.First(&book, req.BookID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	// 获取当前副本数量
	var copyCount int64
	database.DB.Model(&models.BookCopy{}).Where("book_id = ?", req.BookID).Count(&copyCount)

	// 创建新副本
	copies := make([]models.BookCopy, req.CopiesCount)
	for i := 0; i < req.CopiesCount; i++ {
		copies[i] = models.BookCopy{
			BookID:         req.BookID,
			CopyNumber:     book.ISBN + "." + string(rune(i+1+int(copyCount))),
			Status:         models.CopyAvailable,
			AcquisitionDate: time.Now(),
		}
	}

	if err := database.DB.Create(&copies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create book copies"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "book copies added successfully",
		"book_id":       req.BookID,
		"copies_added":  req.CopiesCount,
		"total_copies":  copyCount + int64(req.CopiesCount),
	})
}

// BorrowBook 借阅图书
func BorrowBook(c *gin.Context) {
	var req BorrowBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 开始事务
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查找可用的图书副本
	var bookCopy models.BookCopy
	result := tx.Where("book_id = ? AND status = ?", req.BookID, models.CopyAvailable).First(&bookCopy)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "no available copies of this book"})
		return
	}

	// 检查用户是否已借阅此书
	var activeBorrowCount int64
	tx.Model(&models.Borrow{}).
		Where("user_id = ? AND book_copy_id IN (SELECT id FROM book_copies WHERE book_id = ?) AND status = ?",
		userID, req.BookID, models.BorrowActive).
		Count(&activeBorrowCount)

	if activeBorrowCount > 0 {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "you already have an active borrow for this book"})
		return
	}

	// 更新副本状态
	bookCopy.Status = models.CopyBorrowed
	if err := tx.Save(&bookCopy).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update book copy status"})
		return
	}

	// 创建借阅记录
	borrow := models.Borrow{
		UserID:     userID.(uint),
		BookCopyID: bookCopy.ID,
		BorrowDate: time.Now(),
		DueDate:    time.Now().AddDate(0, 0, req.Days),
		Status:     models.BorrowActive,
	}

	if err := tx.Create(&borrow).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create borrow record"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	// 预加载相关信息
	database.DB.Preload("BookCopy").Preload("BookCopy.Book").Preload("BookCopy.Book.Authors").First(&borrow)
	c.JSON(http.StatusOK, borrow)
}

// ReturnBook 归还图书
func ReturnBook(c *gin.Context) {
	var req ReturnBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 开始事务
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查找借阅记录
	var borrow models.Borrow
	result := tx.Where("id = ? AND user_id = ? AND status = ?", req.BorrowID, userID, models.BorrowActive).First(&borrow)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "active borrow record not found for this user"})
		return
	}

	// 更新借阅记录
	now := time.Now()
	borrow.ReturnDate = &now
	borrow.Status = models.BorrowReturned
	if err := tx.Save(&borrow).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update borrow record"})
		return
	}

	// 更新图书副本状态
	var bookCopy models.BookCopy
	if err := tx.First(&bookCopy, borrow.BookCopyID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find book copy"})
		return
	}

	bookCopy.Status = models.CopyAvailable
	if err := tx.Save(&bookCopy).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update book copy status"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book returned successfully"})
}