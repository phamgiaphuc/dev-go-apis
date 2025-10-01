package session

import (
	"context"
	"time"

	"dev-go-apis/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SessionRepository struct {
	DBClient *sqlx.DB
}

func NewSessionRepository(dbClient *sqlx.DB) *SessionRepository {
	return &SessionRepository{
		DBClient: dbClient,
	}
}

func (repo *SessionRepository) CreateSession(session *models.Session) (*models.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := repo.DBClient.GetContext(ctx, session, `
		INSERT INTO "sessions" (user_id, ip_address, user_agent, expired_at)
		VALUES ($1, $2, $3, $4)
		RETURNING *
	`, session.UserID, session.IpAddress, session.UserAgent, session.ExpiredAt)

	return session, err
}

func (repo *SessionRepository) GetSessionById(id uuid.UUID) (*models.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	session := &models.Session{}
	if err := repo.DBClient.GetContext(ctx, session, "SELECT * FROM sessions WHERE id = $1 LIMIT 1", id.String()); err != nil {
		return nil, err
	}

	return session, nil
}
