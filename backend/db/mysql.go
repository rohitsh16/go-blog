package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLDatabase struct {
	cfg *MySQLConfig
}

func NewMySQLDatabase(cfg *MySQLConfig) *MySQLDatabase {
	return &MySQLDatabase{cfg: cfg}
}

func (m *MySQLDatabase) Connect() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		m.cfg.User,
		m.cfg.Password,
		m.cfg.Host,
		m.cfg.Port,
		m.cfg.Database,
	)
	dbInstance, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf(" Failed to open MySQL: %v", err)
	}

	dbInstance.SetMaxOpenConns(25)
	dbInstance.SetMaxIdleConns(25)
	dbInstance.SetConnMaxLifetime(5 * time.Minute)

	if err := dbInstance.Ping(); err != nil {
		log.Fatalf("Could not connect to MySQL: %v", err)
	}

	log.Println("Connected to MySQL")

	return dbInstance, nil
}
