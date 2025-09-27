package role

import (
	"dev-go-apis/internal/models"
	"fmt"
)

type IRoleRepository interface {
	GetRoleList() ([]models.RoleList, error)
	CreateRole(role *models.Role) (*models.Role, error)
}

type RoleService struct {
	RoleRepository IRoleRepository
}

func NewRoleService(roleRepo IRoleRepository) *RoleService {
	return &RoleService{
		RoleRepository: roleRepo,
	}
}

func (s *RoleService) GetRoleList() ([]models.RoleList, error) {
	return s.RoleRepository.GetRoleList()
}

func (s *RoleService) CreateRole(req *models.CreateRoleRequest) (*models.Role, error) {
	role := &models.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	role, err := s.RoleRepository.CreateRole(role)
	if err != nil {
		return nil, fmt.Errorf("failed to create a role: %s", err.Error())
	}

	return role, nil
}
