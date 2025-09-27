package permission

import "dev-go-apis/internal/models"

type IPermissionRepository interface {
	GetPermissionList() ([]models.PermissionList, error)
	CreatePermissionGroup(permissionGroup *models.PermissionGroup) (*models.PermissionGroup, error)
}

type PermissionService struct {
	PermissionRepo IPermissionRepository
}

func NewPermissionService(permissionRepo IPermissionRepository) *PermissionService {
	return &PermissionService{
		PermissionRepo: permissionRepo,
	}
}

func (s *PermissionService) GetPermissionList() ([]models.PermissionList, error) {
	return s.PermissionRepo.GetPermissionList()
}

func (s *PermissionService) CreatePermissionGroup(body *models.CreatePermissionGroupRequest) (*models.PermissionGroup, error) {
	permissionGroup := &models.PermissionGroup{
		Name:        body.Name,
		Description: body.Description,
	}
	return s.PermissionRepo.CreatePermissionGroup(permissionGroup)
}
