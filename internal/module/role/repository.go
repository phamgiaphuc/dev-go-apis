package role

import (
	"context"
	"dev-go-apis/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type RoleRepository struct {
	DBClient *sqlx.DB
}

func NewRoleRepository(dbClient *sqlx.DB) *RoleRepository {
	return &RoleRepository{
		DBClient: dbClient,
	}
}

func (repo *RoleRepository) GetRoleList() ([]models.RoleList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	roleList := []models.RoleList{}
	if err := repo.DBClient.SelectContext(ctx, &roleList, `
		SELECT r.*,
		COALESCE(
			json_agg(
				json_build_object(
					'id', p.id,
					'name', p.name,
					'description', p.description,
					'group_id', p.group_id
				)
			) FILTER (WHERE p.id IS NOT NULL), '[]'
		) AS permissions
		FROM roles r
		LEFT JOIN role_permissions rp ON rp.role_id = r.id
		LEFT JOIN permissions p ON p.id = rp.permission_id
		GROUP BY r.id, r.name, r.description
		ORDER BY r.id;
	`); err != nil {
		return nil, err
	}

	return roleList, nil
}
