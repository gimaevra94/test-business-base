package database

import (
	"database/sql"

	"github.com/gimaevra94/test-business-base/consts"
	"github.com/gimaevra94/test-business-base/structs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
	rows, err := db.Query(consts.UsersSelectQuery)
	if err != nil {
		return []structs.User{}, errors.WithStack(err)
	}
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

func (db *DB) Create(req *structs.Request) (error, bool) {

	tx, err := db.Begin()
	if err != nil {
		logrus.Panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logrus.Panic(r)
		}
	}()

	var resultReq structs.Request
	row := tx.QueryRow(consts.RequestsSelectQuery)
	row.Scan(&resultReq.ClientName, &resultReq.Phone, &resultReq.Address, &resultReq.ProblemText)
	if err != nil {
		if err == sql.ErrNoRows {
			if _, err := tx.Exec(consts.RequestInsertQuery, req.ClientName, req.Phone, req.Address, req.ProblemText); err != nil {
				tx.Rollback()
				return errors.WithStack(err), false
			}

			if err = tx.Commit(); err != nil {
				return errors.WithStack(err), false

			}
			
		} else {
			tx.Rollback()
			return errors.WithStack(err), false
		}
	}
	return nil, true
}
