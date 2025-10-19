package dbs

import (
	"database/sql"
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/configs"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Init() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", configs.DBUrl)
	// db.SetMaxOpenConns(20)
	// db.SetMaxIdleConns(5)
	Init2()
	return db, err
}

var appDB *AppDB

type AppDB struct {
	*sqlx.DB
}

func Init2() error {
	db, err := sqlx.Connect("postgres", configs.DBUrl)
	appDB = &AppDB{DB: db}
	return err
}

func SelectOne[T any](query string, arg any) (T, error) {
	stmt, err := appDB.PrepareNamed(query)
	var zero T
	if err != nil {
		return zero, model.ServerError(err.Error())
	}
	defer stmt.Close()

	var item T
	err = stmt.Get(&item, arg)
	if err != nil {
		if err == sql.ErrNoRows { // todo: debug this
			return zero, nil
		}
		return zero, model.ServerError(err.Error())
	}
	return item, nil
}

func SelectMany[T any](query string, arg any) ([]T, error) {
	stmt, err := appDB.PrepareNamed(query)
	if err != nil {
		return []T{}, model.ServerError(err.Error())
	}
	defer stmt.Close()

	var items []T
	err = stmt.Select(&items, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return []T{}, nil
		}
		return []T{}, model.ServerError(err.Error())
	}
	return items, nil
}

func Exec(query string, arg any) (int64, error) {
	result, err := appDB.NamedExec(query, arg)
	if err != nil {
		return 0, model.ServerError(err.Error())
	}
	res, err := result.RowsAffected()
	if err != nil {
		return 0, model.ServerError(err.Error())
	}
	return res, nil
}

func ExecTx(tx *sqlx.Tx, query string, arg any) (int64, error) {
	result, err := tx.NamedExec(query, arg)
	if err != nil {
		return 0, model.ServerError(err.Error())
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, model.ServerError(err.Error())
	}
	return rows, nil
}

// use it like:
//
//	txFn := func(tx *sqlx.Tx) error {
//		// do db operations using tx
//		return nil
//	}
//
// err := dbs.WithTx(txFn)
func WithTx(fn func(tx *sqlx.Tx) error) error {
	tx, err := appDB.Beginx()
	if err != nil {
		return model.ServerError(err.Error())
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return model.ServerError(err.Error())
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return model.ServerError(err.Error())
	}

	return nil
}
