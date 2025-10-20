package repo

import (
	"github.com/google/uuid"
)

type SuggestedEmailRepo interface {
	Create(userID uuid.UUID, email string)
	Delete(id uuid.UUID) int64
	DeleteManyByEmail(email string) int64
	GetMany(userID uuid.UUID) []SuggestedEmail
	GetOne(id uuid.UUID) *SuggestedEmail
}

type suggestedEmailRepo struct {
	db *AppDB
}

func NewSuggestedEmailRepo(db *AppDB) SuggestedEmailRepo {
	return &suggestedEmailRepo{db: db}
}

func (r *suggestedEmailRepo) Create(userID uuid.UUID, email string) {
	query := `
		INSERT INTO suggested_emails (id, user_id, email, created_at)
		VALUES (:id, :user_id, :email, current_timestamp)
		ON CONFLICT (user_id, email) DO NOTHING`

	id := uuid.New()
	params := Param{"id": id, "user_id": userID, "email": email}
	execute(r.db, query, params)
}

func (r *suggestedEmailRepo) Delete(id uuid.UUID) int64 {
	query := `DELETE FROM suggested_emails WHERE id = :id`
	param := Param{"id": id}
	return execute(r.db, query, param)
}

func (r *suggestedEmailRepo) DeleteManyByEmail(email string) int64 {
	query := `DELETE FROM suggested_emails WHERE email = :email`
	param := Param{"email": email}
	return execute(r.db, query, param)
}

func (r *suggestedEmailRepo) GetMany(userID uuid.UUID) []SuggestedEmail {
	query := `
		SELECT id, user_id, email, created_at
		FROM suggested_emails
		WHERE user_id = :user_id
		ORDER BY created_at DESC`
	param := Param{"user_id": userID}
	return selectMany[SuggestedEmail](r.db, query, param)
}

func (r *suggestedEmailRepo) GetOne(id uuid.UUID) *SuggestedEmail {
	query := `
		SELECT id, user_id, email, created_at
		FROM suggested_emails
		WHERE id = :id`
	param := Param{"id": id}
	email := selectOne[SuggestedEmail](r.db, query, param)
	return &email
}
