package database

import (
	"database/sql"
	"log"
)

// This Func Will Create Users and Articles Table in Database IF NOT EXISTS
func InitDB(DB *sql.DB) {
	CreateUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
	id INT AUTO_INCREMENT PRIMARY KEY,
	username VARCHAR(50) NOT NULL UNIQUE,
	pass_word VARCHAR(50) NOT NULL 
	);`
	_, err := DB.Exec(CreateUsersTable)

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

	_, err = DB.Exec(CreatArticlesTable)
	if err != nil {
		log.Fatal("ERROR in Creating Article Table: ", err)
	}
}
