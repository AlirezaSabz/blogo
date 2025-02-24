package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

type RegisterRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type Articles struct {
	Id      int    `json:"id" form:"id"`
	Title   string `json:"title" form:"title"`
	Content string `json:"content" form:"title"`
}

func main() {
	e := echo.New()

	dsn := "root:12345678@tcp(127.0.0.1:3306)/blogoDB"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("ERROR in connection : ", err)
	}
	defer db.Close()

	CreateUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
	id INT AUTO_INCREMENT PRIMARY KEY,
	username VARCHAR(50) NOT NULL UNIQUE,
	pass_word VARCHAR(50) NOT NULL 
	);`
	_, err = db.Exec(CreateUsersTable)

	if err != nil {
		log.Fatal("ERROR in Creating Users Table: ", err)
	}

	CreatArticlesTable := `
	CREATE TABLE IF NOT EXISTS articles (
    id INT AUTO_INCREMENT PRIMARY KEY ,
	title VARCHAR(100) NOT NULL,
	content TEXT NOT NULL 
	)
	`

	_, err = db.Exec(CreatArticlesTable)
	if err != nil {
		log.Fatal("ERROR in Creating Article Table: ", err)
	}

	e.POST("/login", func(c echo.Context) error {

		var RegisterInfo RegisterRequest

		err := c.Bind(&RegisterInfo)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		}

		var EmailCheck bool
		checkEmailQuery := "SELECT EXISTS( SELECT 1 FROM users WHERE username= ?);"
		err = db.QueryRow(checkEmailQuery, RegisterInfo.Username).Scan(&EmailCheck)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}
		if EmailCheck {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Email already exists"})
		}
		InsertQuery := `
		INSERT INTO users (username , pass_word)
		VALUES (?,?);`
		_, err = db.Exec(InsertQuery, RegisterInfo.Username, RegisterInfo.Password)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error in Inserting Values to Table"})
		}

		fmt.Println("Data successfully stored in database!")

		return c.JSON(http.StatusOK, RegisterInfo)
	})

	e.POST("/api/v1/articles", func(c echo.Context) error {
		var article Articles
		err = c.Bind(&article)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Database error"})
		}
		ArticleInsertQuery := `
			INSERT INTO articles (title, content)
			VALUES (? , ?);`
		_, err = db.Exec(ArticleInsertQuery, article.Title, article.Content)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error ": "Error inserting new articles!"})
		}
		return c.JSON(http.StatusOK, article)

	})

	e.GET("/api/v1/articles", func(c echo.Context) error {
		articles := make([]Articles, 0)

		articleRows, err := db.Query(`SELECT * FROM articles;`)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error ": "Error fetching list of aticles!"})
		}
		var article Articles
		for articleRows.Next() {
			err := articleRows.Scan(&article.Id, &article.Title, &article.Content)
			if err != nil {
				log.Fatal(err)
			}
			articles = append(articles, article)
		}

		return c.JSON(http.StatusOK, articles)
	})

	e.GET("/api/v1/articles/:id", func(c echo.Context) error {

		id := c.Param("id")
		selectRowQuery := `SELECT * FROM articles
		WHERE id = ? ;`
		selectedRow := db.QueryRow(selectRowQuery, id)

		var article Articles
		err = selectedRow.Scan(&article.Id, &article.Title, &article.Content)
		if err != nil {
			log.Fatal(err)
		}

		return c.JSON(http.StatusOK, article)

	})

	e.Start(":8080")
}
