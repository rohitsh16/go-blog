package db

import (
	"database/sql"
	"errors"
	"sync"
)

type Database interface {
	Connect() (*sql.DB, error)
}

var (
	once     sync.Once
	instance *sql.DB
	initErr  error
)

func GetInstance(dbType string, cfg any) (*sql.DB, error) {
	once.Do(func() {
		var strategy Database
		strategy, initErr = NewDatabase(dbType, cfg)
		if initErr != nil {
			return
		}
		instance, initErr = strategy.Connect()
	})

	return instance, initErr
}

func NewDatabase(dbType string, cfg any) (Database, error) {
	switch dbType {
	case "mysql":
		dbConfig := cfg.(*DatabaseConfig)
		return NewMySQLDatabase(&dbConfig.MySql), nil
	default:
		return nil, errors.New("unsupported DB type")
	}
}
