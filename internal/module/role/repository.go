package role

import (
	"context"
	"dev-go-apis/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type RoleRepository struct {
	DBClient *sqlx.DB
}

func NewRoleRepository(dbClient *sqlx.DB) *RoleRepository {
	return &RoleRepository{
		DBClient: dbClient,
	}
}

func (repo *RoleRepository) DeleteRolePermissions(roleIds *models.RoleIDs) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := repo.DBClient.ExecContext(ctx, `DELETE FROM roles WHERE id = ANY($1)`, pq.Array(*roleIds))
	if err != nil {
		return err
	}

	return nil
}

func (repo *RoleRepository) UpdateRolePermissionsByID(rolePermissions *models.RolePermissions) (*models.RolePermissions, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := repo.DBClient.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	role := &models.Role{}
	roleQuery := `
		UPDATE roles
		SET
				name = $2,
				description = $3
		WHERE id = $1
		RETURNING id, name, description;
	`
	if err = tx.GetContext(ctx, role, roleQuery, rolePermissions.ID, rolePermissions.Name, rolePermissions.Description); err != nil {
		return nil, err
	}

	deleteQuery := `
    DELETE FROM role_permissions 
    WHERE role_id = $1 AND permission_id <> ALL($2)
	`
	_, err = tx.ExecContext(ctx, deleteQuery, rolePermissions.Role.ID, pq.Array(rolePermissions.PermissionIDs))
	if err != nil {
		return nil, err
	}

	if len(rolePermissions.PermissionIDs) > 0 {
		insertQuery := `
			INSERT INTO role_permissions (role_id, permission_id)
			SELECT $1, p_id
			FROM unnest($2::int[]) AS p_id
			ON CONFLICT DO NOTHING
    `
		_, err = tx.ExecContext(
			ctx,
			insertQuery,
			rolePermissions.ID,
			pq.Array(rolePermissions.PermissionIDs),
		)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return rolePermissions, nil
}

func (repo *RoleRepository) GetRolePermissionsByID(role *models.Role) (*models.RolePermissions, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rolePermissions := &models.RolePermissions{}
	if err := repo.DBClient.GetContext(ctx, rolePermissions, `
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
		WHERE r.id = $1
		GROUP BY r.id, r.name, r.description
		ORDER BY r.id;
	`, role.ID); err != nil {
		return nil, err
	}

	return rolePermissions, nil
}

func (repo *RoleRepository) CreateRole(role *models.Role) (*models.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := repo.DBClient.GetContext(ctx, role, `
		INSERT INTO "roles" (name, description)
		VALUES ($1, $2)
		RETURNING *;
	`, role.Name, role.Description); err != nil {
		return nil, err
	}

	return role, nil
}

func (repo *RoleRepository) GetRolePermissionsList() (*models.RolePermissionsList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rolePermissionsList := &models.RolePermissionsList{}
	if err := repo.DBClient.SelectContext(ctx, rolePermissionsList, `
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

	return rolePermissionsList, nil
}
