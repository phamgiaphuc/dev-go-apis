package auth

import (
	"context"
	"dev-go-apis/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	DBClient *sqlx.DB
}

func NewAuthRepository(dbClient *sqlx.DB) *AuthRepository {
	return &AuthRepository{
		DBClient: dbClient,
	}
}

func (repo *AuthRepository) CreateVerification(verification *models.Verification) (*models.Verification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO verifications (user_id, identifier, value, expired_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, identifier)
		DO UPDATE
		SET 
			value = EXCLUDED.value,
			expired_at = EXCLUDED.expired_at,
			updated_at = NOW()
		RETURNING *;

	`
	if err := repo.DBClient.GetContext(ctx, verification, query, verification.UserID, verification.Identifier, verification.Value, verification.ExpiredAt); err != nil {
		return nil, err
	}

	return verification, nil
}
