package configs

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// connect to sql database
func connectDB() *sql.DB {

	dsn := "root:12345678@tcp(127.0.0.1:3306)/blogoDB"
	DB, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal("ERROR in connection : ", err)
	}

	return DB

}

// This Will Ruturn Connected Database,You Should Close It With defer DB.close()
var DB = connectDB()
