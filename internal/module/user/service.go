package user

import (
	"dev-go-apis/internal/models"

	"github.com/google/uuid"
)

type IUserRepository interface {
	GetUserByID(id uuid.UUID) (*models.UserWithAccounts, error)
}

type UserService struct {
	UserRepo IUserRepository
}

func NewUserService(userRepo IUserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (s *UserService) GetUserByID(id uuid.UUID) (*models.UserWithAccounts, error) {
	return s.UserRepo.GetUserByID(id)
}
