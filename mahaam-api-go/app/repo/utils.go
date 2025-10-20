package repo

import (
	"database/sql"
	"fmt"
	"mahaam-api/app/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type AppDB struct {
	*sqlx.DB
}

func NewAppDB(dbUrl string) (*AppDB, error) {
	db, err := sqlx.Connect("postgres", dbUrl)
	if err != nil {
		return nil, err
	}
	return &AppDB{DB: db}, nil
}

func selectOne[T any](db *AppDB, query string, arg any) T {
	stmt, err := db.PrepareNamed(query)
	if err != nil {
		panic(models.ServerError(err.Error()))
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
		panic(models.ServerError(err.Error()))
	}
	return item
}

func selectMany[T any](db *AppDB, query string, arg any) []T {
	stmt, err := db.PrepareNamed(query)
	if err != nil {
		panic(models.ServerError(err.Error()))
	}
	defer stmt.Close()

	var items []T
	err = stmt.Select(&items, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return []T{}
		}
		panic(models.ServerError(err.Error()))
	}
	return items
}

func execute(db *AppDB, query string, arg any) int64 {
	result, err := db.NamedExec(query, arg)
	if err != nil {
		fmt.Println("err", err)
		panic(models.ServerError(err.Error()))
	}
	res, err := result.RowsAffected()
	if err != nil {
		panic(models.ServerError(err.Error()))
	}
	return res
}

func executeTransaction(tx *sqlx.Tx, query string, arg any) int64 {
	result, err := tx.NamedExec(query, arg)
	if err != nil {
		panic(models.ServerError(err.Error()))
	}
	rows, err := result.RowsAffected()
	if err != nil {
		panic(models.ServerError(err.Error()))
	}
	return rows
}

func WithTransaction(db *AppDB, fn func(tx *sqlx.Tx) error) error {
	tx, err := db.Beginx()
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
