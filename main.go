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

func main() {
	e := echo.New()

	dsn := "root:12345678@tcp(127.0.0.1:3306)/blogoDB"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("ERROR in connection : ", err)
	}
	defer db.Close()

	CreateTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
	id INT AUTO_INCREMENT PRIMARY KEY,
	username VARCHAR(50) NOT NULL UNIQUE,
	pass_word VARCHAR(50) NOT NULL 
	);`
	_, err = db.Exec(CreateTableQuery)

	if err != nil {
		log.Fatal("ERROR in Creating Users Table: ", err)
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

	e.Start(":8080")
}
