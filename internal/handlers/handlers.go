package handlers

import (
	"blogo/configs"
	"blogo/internal/models"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var DB = configs.DB

func Register(c echo.Context) error {

	var RegisterInfo models.User

	err := c.Bind(&RegisterInfo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}
	if RegisterInfo.Password == "" || RegisterInfo.Username == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Usernsme and Password are required"})
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
	var user models.User
	var passInDB string
	var userId int
	c.Bind(&user)
	var findPasswordQuery = `SELECT pass_word FROM users WHERE username=? `
	err := DB.QueryRow(findPasswordQuery, user.Username).Scan(&passInDB)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "username or password is incorrect "})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	if passInDB == user.Password {
		var findIdQuery = `SELECT id FROM users WHERE username=?`
		DB.QueryRow(findIdQuery, user.Username).Scan(&userId)
		secretKey := os.Getenv("JWT_SECRET")
		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": user.Username,
			"user_id":  userId,
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
	var article models.Article
	err := c.Bind(&article)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Database error"})
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(float64)

	article.AuthorID = int(userId)

	ArticleInsertQuery := `
		INSERT INTO articles (title, content,author_id )
		VALUES (? , ? , ?);`
	_, err = DB.Exec(ArticleInsertQuery, article.Title, article.Content, userId)
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

	articles := make([]models.Article, 0)

	articleRows, err := DB.Query(`SELECT id , title ,content ,author_id FROM articles;`)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error ": "Error fetching list of aticles!"})
	}
	var article models.Article
	for articleRows.Next() {
		err := articleRows.Scan(&article.Id, &article.Title, &article.Content, &article.AuthorID)
		if err != nil {
			log.Fatal(err)
		}
		articles = append(articles, article)
	}

	return c.JSON(http.StatusOK, articles)
}

func GetArticleById(c echo.Context) error {

	id := c.Param("id")
	selectRowQuery := `SELECT id, title, content ,author_id FROM articles WHERE id = ?`

	selectedRow := DB.QueryRow(selectRowQuery, id)

	var article models.Article
	err := selectedRow.Scan(&article.Id, &article.Title, &article.Content, &article.AuthorID)
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, article)

}

func SendComment(c echo.Context) error {
	var comment models.Comment

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	comment.UserID = int(claims["user_id"].(float64))

	comment.ArticleID, _ = strconv.Atoi(c.Param("id"))

	err := c.Bind(&comment)
	if err != nil {
		return c.JSON(http.StatusBadRequest, comment)
	}
	commentInsertQuery := `
	INSERT INTO comments ( content ,user_id,article_id )
	VALUES ( ? , ?, ?);`
	_, err = DB.Exec(commentInsertQuery, comment.Content, comment.UserID, comment.ArticleID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error ": "Error inserting new articles!"})
	}

	selectRowQuery := `SELECT comment_id FROM  comments ORDER BY comment_id DESC LIMIT 1`
	id := DB.QueryRow(selectRowQuery)
	err = id.Scan(&comment.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error ": "Error Sending Comment !"})
	}

	return c.JSON(http.StatusAccepted, comment)
}

func GetAllComments(c echo.Context) error {
	comments := make([]models.Comment, 0)
	var comment models.Comment

	getCommentsQuery := `SELECT comment_id , content, article_id ,user_id FROM comments`
	commentRows, err := DB.Query(getCommentsQuery)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error ": "Error Fetching Comments 1 !"})
	}

	for commentRows.Next() {
		err := commentRows.Scan(&comment.ID, &comment.Content, &comment.ArticleID, &comment.UserID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error ": "Error Fetching Comments 2!"})
		}
		comments = append(comments, comment)
	}

	return c.JSON(http.StatusOK, comments)

}

func GetCommentByID(c echo.Context) error {
	var comment models.Comment
	commentID := c.Param("commentid")
	GetCommentByIDQuery := `SELECT comment_id , content, article_id ,user_id FROM comments WHERE comment_id= ?`
	commentRow := DB.QueryRow(GetCommentByIDQuery, commentID)
	err := commentRow.Scan(&comment.ID, &comment.Content, &comment.ArticleID, &comment.UserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error ": "Error Fetching Comment!"})
	}
	return c.JSON(http.StatusOK, comment)

}
