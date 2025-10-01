package permission

import (
	"context"
	"time"

	"dev-go-apis/internal/models"

	"github.com/jmoiron/sqlx"
)

type PermissionRepository struct {
	DBClient *sqlx.DB
}

func NewPermissionRepository(dbClient *sqlx.DB) *PermissionRepository {
	return &PermissionRepository{
		DBClient: dbClient,
	}
}

func (repo *PermissionRepository) GetPermissionList() ([]models.PermissionList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	permissionList := []models.PermissionList{}
	if err := repo.DBClient.SelectContext(ctx, &permissionList, `
		SELECT 
			pg.id,
			pg.name,
			pg.description,
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
		FROM permission_groups pg
		LEFT JOIN permissions p ON p.group_id = pg.id
		GROUP BY pg.id, pg.name, pg.description
		ORDER BY pg.id;
	`); err != nil {
		return nil, err
	}

	return permissionList, nil
}
