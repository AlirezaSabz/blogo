package models

type Users struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}
type Articles struct {
	Id      int    `json:"id" form:"id"`
	Title   string `json:"title" form:"title"`
	Content string `json:"content" form:"title"`
}
