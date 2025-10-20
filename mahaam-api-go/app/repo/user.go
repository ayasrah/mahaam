package repo

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	Create(tx *sqlx.Tx) uuid.UUID
	UpdateName(id uuid.UUID, name string) int64
	UpdateEmail(tx *sqlx.Tx, id uuid.UUID, email string) int64
	GetOneByEmail(email string) *User
	GetOne(id uuid.UUID) *User
	Delete(tx *sqlx.Tx, id uuid.UUID) int64
}

type userRepo struct {
	db *AppDB
}

func NewUserRepo(db *AppDB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(tx *sqlx.Tx) uuid.UUID {
	id := uuid.New()
	query := `INSERT INTO users (id, created_at) VALUES (:id, current_timestamp)`
	params := Param{"id": id}
	rows := executeTransaction(tx, query, params)
	if rows != 1 {
		panic("Error creating user")
	}
	return id
}

func (r *userRepo) UpdateName(id uuid.UUID, name string) int64 {
	query := `UPDATE users SET name = :name, updated_at = current_timestamp WHERE id = :id`
	params := Param{"id": id, "name": name}
	return execute(r.db, query, params)
}

func (r *userRepo) UpdateEmail(tx *sqlx.Tx, id uuid.UUID, email string) int64 {
	query := `UPDATE users SET email = :email, updated_at = current_timestamp WHERE id = :id`
	params := Param{"id": id, "email": email}
	return executeTransaction(tx, query, params)
}

func (r *userRepo) GetOneByEmail(email string) *User {
	query := `SELECT id, name, email FROM users WHERE email = :email`
	params := Param{"email": email}
	user := selectOne[User](r.db, query, params)
	if user.ID == uuid.Nil {
		return nil
	}
	return &user
}

func (r *userRepo) GetOne(id uuid.UUID) *User {
	query := `SELECT id, name, email FROM users WHERE id = :id`
	params := Param{"id": id}
	user := selectOne[User](r.db, query, params)
	if user.ID == uuid.Nil {
		return nil
	}
	return &user
}

func (r *userRepo) Delete(tx *sqlx.Tx, id uuid.UUID) int64 {
	query := ` DELETE FROM users WHERE id = :id`
	params := Param{"id": id}
	return executeTransaction(tx, query, params)
}
