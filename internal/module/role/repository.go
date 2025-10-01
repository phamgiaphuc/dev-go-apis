package role

import (
	"context"
	"log"
	"time"

	"dev-go-apis/internal/models"

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

func (repo *RoleRepository) DeleteRole(roleIds *models.RoleIDs) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := repo.DBClient.ExecContext(ctx, `DELETE FROM roles WHERE id = ANY($1)`, pq.Array(*roleIds))
	if err != nil {
		return err
	}

	return nil
}

func (repo *RoleRepository) UpdateRoleById(roleWithPermissions *models.RoleWithPermissions) (*models.RoleWithPermissions, error) {
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
	if err = tx.GetContext(ctx, role, roleQuery, roleWithPermissions.ID, roleWithPermissions.Name, roleWithPermissions.Description); err != nil {
		return nil, err
	}

	deleteQuery := `
    DELETE FROM role_permissions 
    WHERE role_id = $1 AND permission_id <> ALL($2)
	`
	_, err = tx.ExecContext(ctx, deleteQuery, roleWithPermissions.Role.ID, pq.Array(roleWithPermissions.PermissionIDs))
	if err != nil {
		return nil, err
	}

	if len(roleWithPermissions.PermissionIDs) > 0 {
		insertQuery := `
			INSERT INTO role_permissions (role_id, permission_id)
			SELECT $1, p_id
			FROM unnest($2::int[]) AS p_id
			ON CONFLICT DO NOTHING
    `
		_, err = tx.ExecContext(
			ctx,
			insertQuery,
			roleWithPermissions.ID,
			pq.Array(roleWithPermissions.PermissionIDs),
		)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return roleWithPermissions, nil
}

func (repo *RoleRepository) GetRoleById(roleWithPermissions *models.RoleWithPermissions) (*models.RoleWithPermissions, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := repo.DBClient.GetContext(ctx, roleWithPermissions, `
		SELECT r.*,
		COALESCE(
			ARRAY_AGG(
				p.id
			) FILTER (WHERE p.id IS NOT NULL), '{}'
		) AS permission_ids
		FROM roles r
		LEFT JOIN role_permissions rp ON rp.role_id = r.id
		LEFT JOIN permissions p ON p.id = rp.permission_id
		WHERE r.id = $1
		GROUP BY r.id, r.name, r.description
		ORDER BY r.id;
	`, roleWithPermissions.ID); err != nil {
		return nil, err
	}

	return roleWithPermissions, nil
}

func (repo *RoleRepository) CreateRole(roleWithPermissions *models.RoleWithPermissions) (*models.RoleWithPermissions, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := repo.DBClient.GetContext(ctx, roleWithPermissions, `
		INSERT INTO "roles" (name, description)
		VALUES ($1, $2)
		RETURNING *;
	`, roleWithPermissions.Name, roleWithPermissions.Description); err != nil {
		return nil, err
	}

	return roleWithPermissions, nil
}

func (repo *RoleRepository) GetRoleList() (*models.RoleList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	roleList := &models.RoleList{}
	if err := repo.DBClient.SelectContext(ctx, roleList, `
		SELECT r.*,
		COALESCE(
			ARRAY_AGG(
				p.id
			) FILTER (WHERE p.id IS NOT NULL), '{}'
		) AS permission_ids
		FROM roles r
		LEFT JOIN role_permissions rp ON rp.role_id = r.id
		LEFT JOIN permissions p ON p.id = rp.permission_id
		GROUP BY r.id, r.name, r.description
		ORDER BY r.id;
	`); err != nil {
		log.Printf("%v\n", err)
		return nil, err
	}

	return roleList, nil
}
