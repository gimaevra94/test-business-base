package database

import (
	"database/sql"

	"github.com/pkg/errors"
)

type DB struct {
	*sql.DB
}

func DBConn(driver, cfg string) (*DB, error) {
	db, err := sql.Open(driver, cfg)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := db.Ping(); err != nil {
		return nil, errors.WithStack(err)
	}
	return &DB{db}, nil
}
