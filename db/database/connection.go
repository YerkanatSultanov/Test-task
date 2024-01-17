package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"

	"test-task/config"
)

type DataBase struct {
	db *sql.DB
}

func NewDataBase() (*DataBase, error) {
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("error to load configs: %s", err)
	}

	db, err := sql.Open(os.Getenv("DB_DRIVER"), fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL"),
	))

	if err != nil {
		return nil, err
	}

	return &DataBase{db: db}, nil
}

func (d *DataBase) Close() {
	err := d.db.Close()
	if err != nil {
		return
	}
}

func (d *DataBase) GetDB() *sql.DB {
	return d.db
}
