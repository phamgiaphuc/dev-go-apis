package session

import (
	"dev-go-apis/internal/models"

	"github.com/google/uuid"
)

type ISessionRepository interface {
	CreateSession(session *models.Session) (*models.Session, error)
	GetSessionById(id uuid.UUID) (*models.Session, error)
}

type SessionService struct {
	SessionRepo ISessionRepository
}

func NewSessionService(sessionRepo ISessionRepository) *SessionService {
	return &SessionService{
		SessionRepo: sessionRepo,
	}
}

func (s *SessionService) CreateSession(session *models.Session) (*models.Session, error) {
	return s.SessionRepo.CreateSession(session)
}

func (s *SessionService) GetSessionById(id uuid.UUID) (*models.Session, error) {
	return s.SessionRepo.GetSessionById(id)
}
