package user

import (
	"dev-go-apis/internal/models"
	"slices"

	"github.com/google/uuid"
)

type IUserRepository interface {
	GetUserByID(id uuid.UUID) (*models.UserWithAccounts, error)
	GetUserPermissionsByRoleID(roleID int) ([]models.Permission, error)
}

type IRoleRepository interface {
	GetRoleById(roleWithPermissions *models.RoleWithPermissions) (*models.RoleWithPermissions, error)
}

type UserService struct {
	UserRepo IUserRepository
	RoleRepo IRoleRepository
}

func NewUserService(userRepo IUserRepository, roleRepo IRoleRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
		RoleRepo: roleRepo,
	}
}

func (s *UserService) CheckUserRole(id uuid.UUID, requiredPermissions []string) (bool, error) {
	user, err := s.GetUserByID(id)
	if err != nil {
		return false, err
	}

	roleWithPermissions := &models.RoleWithPermissions{
		Role: models.Role{
			ID: user.RoleID,
		},
	}

	_, err = s.RoleRepo.GetRoleById(roleWithPermissions)
	if err != nil {
		return false, err
	}

	hasPermission := slices.ContainsFunc(requiredPermissions, func(r string) bool {
		return slices.Contains(roleWithPermissions.Permissions, r)
	})

	return hasPermission, nil
}

func (s *UserService) GetUserByID(id uuid.UUID) (*models.UserWithAccounts, error) {
	return s.UserRepo.GetUserByID(id)
}

func (s *UserService) GetUserPermissionsByRoleID(roleID int) ([]models.Permission, error) {
	return s.UserRepo.GetUserPermissionsByRoleID(roleID)
}
