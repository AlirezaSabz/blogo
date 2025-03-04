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
	content TEXT NOT NULL ,
	author_id INT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	status ENUM('draft', 'published', 'archived') DEFAULT 'draft',
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
	)
	`

	_, err = DB.Exec(CreatArticlesTable)
	if err != nil {
		log.Fatal("ERROR in Creating Article Table: ", err)
	}
}
