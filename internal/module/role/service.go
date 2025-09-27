package role

import "dev-go-apis/internal/models"

type IRoleRepository interface {
	GetRoleList() ([]models.RoleList, error)
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
