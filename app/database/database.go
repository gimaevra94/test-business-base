package database

import (
	"database/sql"

	"github.com/gimaevra94/test-business-base/consts"
	"github.com/gimaevra94/test-business-base/structs"
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

func (db *DB) GetUsers() ([]structs.User, error) {
	rows, _ := db.Query(consts.UsersSelectQuery)
	defer rows.Close()

	var users []structs.User
	for rows.Next() {
		var u structs.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Role); err != nil {
			if err == sql.ErrNoRows {
				return []structs.User{}, errors.WithStack(err)
			}
			return []structs.User{}, errors.WithStack(err)
		}
		users = append(users, u)
	}
	
	return users, nil
}
