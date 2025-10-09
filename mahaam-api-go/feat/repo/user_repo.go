package repo

import (
	"mahaam-api/infra/dbs"

	"github.com/google/uuid"
)

type UserRepo interface {
	Create(tx Tx) UUID
	UpdateName(id UUID, name string) int64
	UpdateEmail(tx Tx, id UUID, email string) int64
	GetOneByEmail(email string) *User
	GetOne(id UUID) *User
	Delete(tx Tx, id UUID) int64
}

type userRepo struct {
}

func NewUserRepo() UserRepo {
	return &userRepo{}
}

func (r *userRepo) Create(tx Tx) UUID {
	id := uuid.New()
	query := `INSERT INTO users (id, created_at) VALUES (:id, current_timestamp)`
	params := Param{"id": id}
	rows := dbs.ExecTx(tx, query, params)
	if rows != 1 {
		panic("Error creating user")
	}
	return id
}

func (r *userRepo) UpdateName(id UUID, name string) int64 {
	query := `UPDATE users SET name = :name, updated_at = current_timestamp WHERE id = :id`
	params := Param{"id": id, "name": name}
	return dbs.Exec(query, params)
}

func (r *userRepo) UpdateEmail(tx Tx, id UUID, email string) int64 {
	query := `UPDATE users SET email = :email, updated_at = current_timestamp WHERE id = :id`
	params := Param{"id": id, "email": email}
	return dbs.ExecTx(tx, query, params)
}

func (r *userRepo) GetOneByEmail(email string) *User {
	query := `SELECT id, name, email FROM users WHERE email = :email`
	params := Param{"email": email}
	user := dbs.SelectOne[User](query, params)
	if user.ID == uuid.Nil {
		return nil
	}
	return &user
}

func (r *userRepo) GetOne(id UUID) *User {
	query := `SELECT id, name, email FROM users WHERE id = :id`
	params := Param{"id": id}
	user := dbs.SelectOne[User](query, params)
	if user.ID == uuid.Nil {
		return nil
	}
	return &user
}

func (r *userRepo) Delete(tx Tx, id UUID) int64 {
	query := ` DELETE FROM users WHERE id = :id`
	params := Param{"id": id}
	return dbs.ExecTx(tx, query, params)
}
