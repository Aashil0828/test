package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// ConnectDatabase initializes the database connection
func ConnectDatabase() {
	var err error
	dsn := fmt.Sprintf("%v:%v@tcp(%v:3306)/%v?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_URL"), os.Getenv("DB_NAME"))

	// Open a connection
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	// Ping the database
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	log.Println("Database connection established")
}
