package user

import (
	"context"
	"dev-go-apis/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	DBClient *sqlx.DB
}

func NewUserRepository(dbClient *sqlx.DB) *UserRepository {
	return &UserRepository{
		DBClient: dbClient,
	}
}

func (repo *UserRepository) CreateUser(userWithPassword *models.UserWithPassword) (*models.User, error) {
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

	user := &models.User{}
	userQuery := `
		INSERT INTO "users" (name, email, image)
		VALUES ($1, $2, $3)
		RETURNING *, (SELECT name FROM roles WHERE id = role_id) AS role;
	`
	if err = tx.GetContext(ctx, user, userQuery, userWithPassword.Name, userWithPassword.Email, userWithPassword.Image); err != nil {
		return nil, err
	}

	account := &models.Account{}
	accountQuery := `
		INSERT INTO "accounts" (user_id, account_id, provider_id, password)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, user_id, account_id, provider_id, password, created_at, updated_at;
	`
	if err = tx.GetContext(ctx, account, accountQuery, user.ID, user.ID, models.ProviderCredential, userWithPassword.Password); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user := models.User{}
	if err := repo.DBClient.GetContext(ctx, &user, `
		SELECT u.*, r.name AS role
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.email = $1
	`, email); err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepository) GetUserWithPasswordByEmail(email string) (*models.UserWithPassword, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userWithPasssword := &models.UserWithPassword{}
	if err := repo.DBClient.GetContext(ctx, userWithPasssword, `
		SELECT u.*, r.name AS role, a.password AS password
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		LEFT JOIN accounts a ON u.id = a.user_id
		WHERE u.email = $1
	`, email); err != nil {
		return nil, err
	}

	return userWithPasssword, nil
}

func (repo *UserRepository) GetUserByID(id uuid.UUID) (*models.UserWithAccounts, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user := models.User{}
	if err := repo.DBClient.GetContext(ctx, &user, `
		SELECT u.*, r.name AS role
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1;
	`, id.String()); err != nil {
		return nil, err
	}

	accounts := []models.Account{}
	if err := repo.DBClient.SelectContext(ctx, &accounts, `
		SELECT * FROM "accounts" WHERE user_id = $1
	`, user.ID); err != nil {
		return nil, err
	}

	return &models.UserWithAccounts{User: user, Accounts: accounts}, nil
}
