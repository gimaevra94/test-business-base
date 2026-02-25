package database

import (
	"database/sql"

	"github.com/gimaevra94/test-business-base/consts"
	"github.com/gimaevra94/test-business-base/structs"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	sql  *sql.DB
	gorm *gorm.DB
}

func (db *DB) GORM() *gorm.DB {
	return db.gorm
}

func DBConn(cfg string) (*DB, error) {
	sql, sqlErr := sql.Open("mysql", cfg)
	gorm, gormErr := gorm.Open(mysql.Open(cfg), &gorm.Config{})

	if sqlErr != nil {
		return nil, errors.WithStack(sqlErr)
	}

	if gormErr != nil {
		return nil, errors.WithStack(gormErr)
	}

	if err := sql.Ping(); err != nil {
		return nil, errors.WithStack(err)
	}
	return &DB{sql, gorm}, nil
}

func (db *DB) GetUsers() ([]structs.User, error) {
	rows, err := db.sql.Query(consts.UsersSelectQuery)
	if err != nil {
		return []structs.User{}, errors.WithStack(err)
	}
	defer rows.Close()

	var users []structs.User
	for rows.Next() {
		var u structs.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Role); err != nil {
			return []structs.User{}, errors.WithStack(err)
		}
		users = append(users, u)
	}

	return users, nil
}

func (db *DB) Create(req *structs.Request) error {
	_, err := db.sql.Exec(consts.RequestInsertQuery, req.ClientName, req.Phone, req.Address, req.ProblemText)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (db *DB) AssignedStatusUpdate(mid, rid int) error {
	_, err := db.sql.Exec(consts.AssignedStatusUpdateQuery, mid, rid)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (db *DB) CanceledStatusUpdate(rid int) error {
	_, err := db.sql.Exec(consts.CanceledStatusUpdateQuery, rid)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (db *DB) InProgressStatusUpdate(rid, mid int) error {
	res, err := db.sql.Exec(consts.InProgressStatusUpdateQuery, rid, mid)
	if err != nil {
		return errors.WithStack(err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return errors.WithStack(err)
	}

	if rows == 0 {
		err = sql.ErrNoRows
		return errors.WithStack(err)
	}

	return nil
}

func (db *DB) DoneStatusUpdate(rid, mid int) error {
	_, err := db.sql.Exec(consts.DoneStatusUpdateQuery, rid, mid)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
