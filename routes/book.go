package routes

import (
	"go-api/controllers"
	"go-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupBookRoutes(router *gin.Engine) {
	books := router.Group("/api/books")
	books.Use(middleware.AuthMiddleware())
	{
		books.POST("", controllers.CreateBook)
		books.GET("", controllers.GetBooks)
		books.GET("/:id", controllers.GetBook)
		books.PUT("/:id", controllers.UpdateBook)
		books.DELETE("/:id", controllers.DeleteBook)
	}
}
