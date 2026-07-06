package database

import (
	"database/sql"

	"github.com/bytedance/gopkg/util/logger"
	_ "github.com/lib/pq"
)

func NewPostgres(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)

	if err != nil {
		db.Close()
		return nil, err
	}

	pingErr := db.Ping()

	if pingErr != nil {
		db.Close()
		return nil, pingErr
	}

	logger.Info("Database connected !!")

	return db, nil

}
