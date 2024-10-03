package storage

import (
	"database/sql"
	"fmt"
	"log"
)

var db *sql.DB

// Initialize MySQL connection
func InitDB() {
	var err error
	dataSourceName := "user:password@tcp(127.0.0.1:3306)/blogdb" // Replace with your MySQL user, password, and db name
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	// Ping to check if database is accessible
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database:", err)
	}

	fmt.Println("Connected to MySQL database!")
}
