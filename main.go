package main

import (
	"blogo/configs"
	"blogo/internal/database"
	"blogo/internal/handlers"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func main() {
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	e := echo.New()
	DB := configs.DB
	defer DB.Close()

	database.InitDB(DB)

	e.POST("/register", handlers.Register)
	e.POST("/login", handlers.Login)
	articles := e.Group("/api/v1/articles")
	articles.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: secretKey,
	}))
	articles.POST("", handlers.CreateArticle)
	articles.GET("", handlers.GetArticle)
	articles.GET("/:id", handlers.GetArticleById)
	articles.POST("/:id/comments", handlers.SendComment)
	articles.GET("/:id/comments", handlers.GetAllComments)
	articles.GET("/:id/comments/:commentid", handlers.GetCommentByID)

	e.Start(":8080")
}
