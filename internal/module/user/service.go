package user

import (
	"dev-go-apis/internal/models"

	"github.com/google/uuid"
)

type IUserRepository interface {
	GetUserByID(id uuid.UUID) (*models.UserWithAccounts, error)
	GetUserPermissionsByRoleID(roleID int) ([]models.Permission, error)
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

func (s *UserService) GetUserPermissionsByRoleID(roleID int) ([]models.Permission, error) {
	return s.UserRepo.GetUserPermissionsByRoleID(roleID)
}
