package database

import (
	"database/sql"

	"github.com/gimaevra94/test-business-base/consts"
	"github.com/gimaevra94/test-business-base/structs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type DB struct {
	*sql.DB
}

func DBConn(cfg string) (*DB, error) {
	db, err := sql.Open("mysql", cfg)
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
			return []structs.User{}, errors.WithStack(err)
		}
		users = append(users, u)
	}

	return users, nil
}

func (db *DB) Create(req *structs.Request) (error, bool) {
	tx, err := db.Begin()
	if err != nil {
		return errors.WithStack(err), false
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			return
		}
	}()

	var resultReq structs.Request
	row := tx.QueryRow(consts.RequestsSelectQuery)
	err = row.Scan(&resultReq.ClientName, &resultReq.Phone, &resultReq.Address, &resultReq.ProblemText)
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

func (db *DB) Dashboard(query string, args []any) ([]structs.Request, []structs.User, error) {
	tx, err := db.Begin()
	if err != nil {
		return []structs.Request{}, []structs.User{}, errors.WithStack(err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			return
		}
	}()

	rows, err := tx.Query(query, args...)
	if err != nil {
		tx.Rollback()
		return []structs.Request{}, []structs.User{}, errors.WithStack(err)
	}
	defer rows.Close()

	var reqs []structs.Request
	for rows.Next() {
		var r structs.Request
		if err := rows.Scan(&r.ID, &r.ClientName, &r.Phone, &r.Address, &r.ProblemText, &r.Status, &r.AssignedTo, &r.Version, &r.CreatedAt, &r.UpdatedAt); err != nil {
			tx.Rollback()
			return []structs.Request{}, []structs.User{}, errors.WithStack(err)
		}
		reqs = append(reqs, r)
	}

	rows, err = tx.Query(consts.MastersSelectQuery)
	if err != nil {
		tx.Rollback()
		return []structs.Request{}, []structs.User{}, errors.WithStack(err)
	}
	defer rows.Close()

	var masters []structs.User
	for rows.Next() {
		var master structs.User
		if err := rows.Scan(&master.ID, &master.Name); err != nil {
			return []structs.Request{}, []structs.User{}, errors.WithStack(err)
		}
		masters = append(masters, master)
	}

	if err = tx.Commit(); err != nil {
		return []structs.Request{}, []structs.User{}, errors.WithStack(err)
	}

	return reqs, masters, nil
}

func (db *DB) AssignedStatusUpdate(id int) error {
	_, err := db.Exec(consts.AssignedStatusUpdateQuery, id, id)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (db *DB) CanceledStatusUpdate(rid int) error {
	_, err := db.Exec(consts.CanceledStatusUpdateQuery, rid)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (db *DB) InProgressStatusUpdate(rid int) error {
	res, err := db.Exec(consts.InProgressStatusUpdateQuery, rid)
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

func (db *DB) DoneStatusUpdate(rid int) error {
	_, err := db.Exec(consts.DoneStatusUpdateQuery, rid)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
