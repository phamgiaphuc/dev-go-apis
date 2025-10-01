package role

import (
	"dev-go-apis/internal/models"
	"fmt"
	"log"
	"time"
)

type IRoleRepository interface {
	GetRoleList() (*models.RoleList, error)
	CreateRole(roleWithPermissions *models.RoleWithPermissions) (*models.RoleWithPermissions, error)
	UpdateRoleById(roleWithPermissions *models.RoleWithPermissions) (*models.RoleWithPermissions, error)
	DeleteRole(roleIds *models.RoleIDs) error
	GetRoleById(roleWithPermissions *models.RoleWithPermissions) (*models.RoleWithPermissions, error)
}

type ICacheRepository interface {
	SetValue(string, interface{}, time.Duration) error
}

type RoleService struct {
	RoleRepository  IRoleRepository
	CacheRepository ICacheRepository
}

func NewRoleService(roleRepo IRoleRepository, cacheRepo ICacheRepository) *RoleService {
	return &RoleService{
		RoleRepository:  roleRepo,
		CacheRepository: cacheRepo,
	}
}

func (s *RoleService) DeleteRole(req *models.DeleteRolesRequest) error {
	return s.RoleRepository.DeleteRole(&req.RoleIDs)
}

func (s *RoleService) GetRoleList() (*models.RoleList, error) {
	list, err := s.RoleRepository.GetRoleList()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *RoleService) CreateRole(req *models.CreateRoleRequest) (*models.RoleWithPermissions, error) {
	roleWithPermissions := &models.RoleWithPermissions{
		Role: models.Role{
			Name:        req.Name,
			Description: req.Description,
		},
	}

	role, err := s.RoleRepository.CreateRole(roleWithPermissions)
	if err != nil {
		return nil, fmt.Errorf("failed to create a role: %s", err.Error())
	}

	return role, nil
}

func (s *RoleService) GetRoleById(req *models.GetRoleByIdRequest) (*models.RoleWithPermissions, error) {
	roleWithPermissions := &models.RoleWithPermissions{
		Role: models.Role{
			ID: req.ID,
		},
	}
	_, err := s.RoleRepository.GetRoleById(roleWithPermissions)
	if err != nil {
		log.Printf("%v\n", err.Error())
		return nil, err
	}
	return roleWithPermissions, nil
}

func (s *RoleService) UpdateRoleById(req *models.UpdateRoleRequest) (*models.RoleWithPermissions, error) {
	roleWithPermissions := &models.RoleWithPermissions{
		Role:          req.Role,
		PermissionIDs: req.PermissionIDs,
	}
	return s.RoleRepository.UpdateRoleById(roleWithPermissions)
}
