package main

import (
	"blogo/configs"
	"blogo/internal/database"
	"blogo/internal/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	DB := configs.DB
	defer DB.Close()

	database.InitDB(DB)

	e.POST("/register", handlers.Register)
	articles := e.Group("/api/v1/articles")
	articles.POST("", handlers.CreateArticle)
	articles.GET("", handlers.GetArticle)
	articles.GET("/:id", handlers.GetArticleById)

	e.Start(":8080")
}
