package handlers

import (
	"blogo/configs"
	"blogo/internal/models"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var DB = configs.DB

func Register(c echo.Context) error {

	var RegisterInfo models.Users

	err := c.Bind(&RegisterInfo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	var UserCheck bool
	checkUserQuery := "SELECT EXISTS( SELECT 1 FROM users WHERE username= ?);"
	err = DB.QueryRow(checkUserQuery, RegisterInfo.Username).Scan(&UserCheck)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if UserCheck {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Username already exists"})
	}
	InsertQuery := `
	INSERT INTO users (username , pass_word)
	VALUES (?,?);`
	_, err = DB.Exec(InsertQuery, RegisterInfo.Username, RegisterInfo.Password)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error in Inserting Values to Table"})
	}

	fmt.Println("Data successfully stored in database!")

	return c.JSON(http.StatusOK, RegisterInfo)
}

func Login(c echo.Context) error {
	var user models.Users
	var passInDB string
	c.Bind(&user)
	var query = `SELECT pass_word FROM users WHERE username=? `
	err := DB.QueryRow(query, user.Username).Scan(&passInDB)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "username or password is incorrect "})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	if passInDB == user.Password {
		secretKey := os.Getenv("JWT_SECRET")
		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": user.Username,
			"exp":      time.Now().Add(time.Hour).Unix(),
		})
		signedjwtToken, err := jwtToken.SignedString([]byte(secretKey))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"Error ": err.Error()})

		} else {
			return c.JSON(http.StatusAccepted, map[string]string{"Your JWT ": signedjwtToken})

		}

	} else {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "password is incorrect "})
	}
}

func CreateArticle(c echo.Context) error {
	var article models.Articles
	err := c.Bind(&article)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Database error"})
	}
	ArticleInsertQuery := `
		INSERT INTO articles (title, content)
		VALUES (? , ?);`
	_, err = DB.Exec(ArticleInsertQuery, article.Title, article.Content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error ": "Error inserting new articles!"})
	}
	selectRowQuery := `SELECT id FROM  articles ORDER BY id DESC LIMIT 1`
	id := DB.QueryRow(selectRowQuery)
	err = id.Scan(&article.Id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error ": "Error inserting new articles!"})
	}

	return c.JSON(http.StatusOK, article)

}

func GetArticle(c echo.Context) error {
	articles := make([]models.Articles, 0)

	articleRows, err := DB.Query(`SELECT * FROM articles;`)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error ": "Error fetching list of aticles!"})
	}
	var article models.Articles
	for articleRows.Next() {
		err := articleRows.Scan(&article.Id, &article.Title, &article.Content)
		if err != nil {
			log.Fatal(err)
		}
		articles = append(articles, article)
	}

	return c.JSON(http.StatusOK, articles)
}

func GetArticleById(c echo.Context) error {

	id := c.Param("id")
	selectRowQuery := `SELECT * FROM articles
	WHERE id = ? ;`
	selectedRow := DB.QueryRow(selectRowQuery, id)

	var article models.Articles
	err := selectedRow.Scan(&article.Id, &article.Title, &article.Content)
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, article)

}
