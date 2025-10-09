package dbs

import (
	"database/sql"
	"fmt"
	"mahaam-api/feat/models"
	"mahaam-api/infra/configs"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBErr = models.DBErr

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

func SelectOne[T any](query string, arg any) T {
	stmt, err := appDB.PrepareNamed(query)
	if err != nil {
		panic(&DBErr{Err: err})
	}
	defer stmt.Close()

	var item T
	err = stmt.Get(&item, arg)
	if err != nil {
		if err == sql.ErrNoRows { // todo: debug this
			var zero T
			return zero
			// return *new(T)
		}
		panic(&DBErr{Err: err})
	}
	return item
}

func SelectMany[T any](query string, arg any) []T {
	stmt, err := appDB.PrepareNamed(query)
	if err != nil {
		panic(&DBErr{Err: err})
	}
	defer stmt.Close()

	var items []T
	err = stmt.Select(&items, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return []T{}
		}
		panic(&DBErr{Err: err})
	}
	return items
}

func Exec(query string, arg any) int64 {
	result, err := appDB.NamedExec(query, arg)
	if err != nil {
		fmt.Println("err", err)
		panic(DBErr{Err: err})
	}
	res, err := result.RowsAffected()
	if err != nil {
		panic(DBErr{Err: err})
	}
	return res
}

func ExecTx(tx *sqlx.Tx, query string, arg any) int64 {
	result, err := tx.NamedExec(query, arg)
	if err != nil {
		panic(DBErr{Err: err})
	}
	rows, err := result.RowsAffected()
	if err != nil {
		panic(DBErr{Err: err})
	}
	return rows
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
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	return fn(tx)
}
