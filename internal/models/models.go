package models

type User struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}
type Article struct {
	Id       int    `json:"id" form:"id"`
	Title    string `json:"title" form:"title"`
	Content  string `json:"content" form:"content"`
	AuthorID int    `json:"author_id" form:"author_id"`
}

type Comment struct {
	ID        int    `json:"id" form:"id"`
	Content   string `json:"content" form:"content"`
	UserID    int    `json:"user_id" form:"user_id"`
	ArticleID int    `json:"article_id" form:"article_id"`
}
