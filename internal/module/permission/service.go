package permission

import "dev-go-apis/internal/models"

type IPermissionRepository interface {
	GetPermissionList() ([]models.PermissionList, error)
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
