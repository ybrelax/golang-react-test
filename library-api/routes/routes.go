package routes

import (
	"github.com/example/library-api/controllers"
	"github.com/example/library-api/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有API路由
func RegisterRoutes(router *gin.Engine) {
	// 公开路由
	public := router.Group("/")
	{
		// 用户认证路由
		auth := public.Group("auth")
		{
			public.GET("google/login", controllers.GoogleLogin)
			public.GET("google/callback", controllers.GoogleLoginCallback)
			auth.POST("register", controllers.Register)
			auth.POST("login", controllers.Login)
		}
	}

	// 需要认证的路由
	api := router.Group("/api")
	api.Use(middleware.JWTMiddleware())
	{
		// 用户路由
		user := api.Group("user")
		{
			user.GET("borrows", controllers.GetMyBorrows)
		}

		// 图书路由
		books := api.Group("books")
		{
			// 所有用户可访问的图书路由
			books.GET("", controllers.GetBooks)
			books.GET("/:id", controllers.GetBook)
			books.POST("borrow", controllers.BorrowBook)
			books.POST("return", controllers.ReturnBook)

			// 管理员路由
			admin := books.Group("")
			admin.Use(middleware.AdminRequired())
			{
				admin.POST("", controllers.CreateBook)
				admin.PUT("/:id", controllers.UpdateBook)
				admin.DELETE("/:id", controllers.DeleteBook)
				admin.POST("/:id/copies", controllers.AddBookCopies)
			}
		}
	}
}
