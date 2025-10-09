package repo

import (
	"mahaam-api/infra/dbs"

	"github.com/google/uuid"
)

type SuggestedEmailRepo interface {
	Create(userID UUID, email string) UUID
	Delete(id UUID) int64
	DeleteManyByEmail(email string) int64
	GetMany(userID UUID) []SuggestedEmail
	GetOne(id UUID) *SuggestedEmail
}

type suggestedEmailRepo struct {
}

func NewSuggestedEmailRepo() SuggestedEmailRepo {
	return &suggestedEmailRepo{}
}

func (r *suggestedEmailRepo) Create(userID UUID, email string) UUID {
	query := `
		INSERT INTO suggested_emails (id, user_id, email, created_at)
		VALUES (:id, :user_id, :email, current_timestamp)
		ON CONFLICT (user_id, email) DO NOTHING`

	id := uuid.New()
	params := Param{"id": id, "user_id": userID, "email": email}
	updated := dbs.Exec(query, params)
	if updated > 0 {
		return id
	}
	return UUID{}
}

func (r *suggestedEmailRepo) Delete(id UUID) int64 {
	query := `DELETE FROM suggested_emails WHERE id = :id`
	param := Param{"id": id}
	return dbs.Exec(query, param)
}

func (r *suggestedEmailRepo) DeleteManyByEmail(email string) int64 {
	query := `DELETE FROM suggested_emails WHERE email = :email`
	param := Param{"email": email}
	return dbs.Exec(query, param)
}

func (r *suggestedEmailRepo) GetMany(userID UUID) []SuggestedEmail {
	query := `
		SELECT id, user_id, email, created_at
		FROM suggested_emails
		WHERE user_id = :user_id
		ORDER BY created_at DESC`
	param := Param{"user_id": userID}
	return dbs.SelectMany[SuggestedEmail](query, param)
}

func (r *suggestedEmailRepo) GetOne(id UUID) *SuggestedEmail {
	query := `
		SELECT id, user_id, email, created_at
		FROM suggested_emails
		WHERE id = :id`
	param := Param{"id": id}
	email := dbs.SelectOne[SuggestedEmail](query, param)
	return &email
}
