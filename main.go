package main

import (
	_ "Projectmugen/docs"
	"Projectmugen/internal/controllers"
	"Projectmugen/internal/services"

	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Документация для API
// @version         1.0
// @description     Документация моего API

func main() {
	services.InitDB()
	router := gin.Default()

	router.GET("/swagger/*any", gin.WrapF(httpSwagger.WrapHandler))

	router.POST("/login", controllers.Login)
	router.POST("/register", controllers.Register)
	router.POST("/refresh", controllers.Refresh)

	protected := router.Group("/")
	protected.Use(controllers.AuthMiddleware())
	{
		protected.GET("/books", controllers.GetBooks)

		protected.GET("/books/:id", controllers.GetBookByID)

		protected.GET("/books/year-range", controllers.GetBooksByYearRange)

		protected.GET("/books/count-by-author", controllers.CountBooksByAuthor)

		protected.POST("/books/publisher", controllers.UpdateBooksPublisher)

		protected.POST("/books", controllers.RoleMiddleware("admin"), controllers.CreateBook)

		protected.PUT("/books/:id", controllers.RoleMiddleware("admin"), controllers.UpdateBook)

		protected.DELETE("/books/:id", controllers.RoleMiddleware("admin"), controllers.DeleteBook)

	}
	router.Run(":8080")
}
