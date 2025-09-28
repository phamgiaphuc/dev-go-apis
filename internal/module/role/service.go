package role

import (
	"dev-go-apis/internal/models"
	"fmt"
	"time"
)

type IRoleRepository interface {
	GetRolePermissionsList() (*models.RolePermissionsList, error)
	CreateRole(role *models.Role) (*models.Role, error)
	GetRolePermissionsByID(role *models.Role) (*models.RolePermissions, error)
	UpdateRolePermissionsByID(rolePermissions *models.RolePermissions) (*models.RolePermissions, error)
	DeleteRolePermissions(roleIds *models.RoleIDs) error
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

func (s *RoleService) DeleteRolePermissions(req *models.DeleteRolePermissionsRequest) error {
	return s.RoleRepository.DeleteRolePermissions(&req.RoleIDs)
}

func (s *RoleService) GetRolePermissionsList() (*models.RolePermissionsList, error) {
	list, err := s.RoleRepository.GetRolePermissionsList()
	if err != nil {
		return nil, err
	}
	for i := range *list {
		role := &(*list)[i]
		role.PermissionIDs = make([]int, len(role.Permissions))
		for j, p := range role.Permissions {
			role.PermissionIDs[j] = p.ID
		}
	}
	return list, nil
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

func (s *RoleService) GetRolePermissionsByID(id int) (*models.RolePermissions, error) {
	role := &models.Role{
		ID: id,
	}
	rolePermissions, err := s.RoleRepository.GetRolePermissionsByID(role)
	if err != nil {
		return nil, err
	}
	rolePermissions.PermissionIDs = make([]int, len(rolePermissions.Permissions))
	for j, p := range rolePermissions.Permissions {
		rolePermissions.PermissionIDs[j] = p.ID
	}
	return rolePermissions, nil
}

func (s *RoleService) UpdateRolePermissionsByID(req *models.UpdateRolePermissionsRequest) (*models.RolePermissions, error) {
	rolePermissions := &models.RolePermissions{
		Role:          req.Role,
		PermissionIDs: req.PermissionIDs,
	}
	return s.RoleRepository.UpdateRolePermissionsByID(rolePermissions)
}
