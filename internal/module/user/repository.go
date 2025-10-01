package user

import (
	"context"
	"database/sql"
	"time"

	"dev-go-apis/internal/models"

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

func (repo *UserRepository) UpdateUserById(id string, user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE "users"
		SET
			name = $1,
			email = $2,
			email_verified = $3,
			image = $4,
			is_banned = $5,
			updated_at = NOW()
		WHERE id = $6
	`
	if _, err := repo.DBClient.ExecContext(ctx, query, user.Name, user.Email, user.EmailVerified, user.Image, user.IsBanned, id); err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO "users" (name, email, image, email_verified)
		VALUES ($1, $2, $3)
		RETURNING *, (SELECT name FROM roles WHERE id = role_id) AS role
	`
	if err := repo.DBClient.GetContext(ctx, user, query, user.Name, user.Email, user.Image, user.EmailVerified); err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *UserRepository) CreateAccount(account *models.Account) (*models.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO "accounts" (user_id, account_id, provider_id, password)
		VALUES ($1, $2, $3, $4)
		RETURNING *
	`
	if err := repo.DBClient.GetContext(ctx, account, query, account.UserID, account.AccountID, account.ProviderId, account.Password); err != nil {
		return nil, err
	}

	return account, nil
}

func (repo *UserRepository) CreateUserAccount(userWithAccount *models.UserWithAccount) (*models.UserWithAccount, error) {
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

	userQuery := `
		INSERT INTO "users" (name, email, image, email_verified)
		VALUES ($1, $2, $3, $4)
		RETURNING *, (SELECT name FROM roles WHERE id = role_id) AS role;
	`
	if err = tx.GetContext(ctx, &userWithAccount.User, userQuery, userWithAccount.User.Name, userWithAccount.User.Email, userWithAccount.User.Image, userWithAccount.User.EmailVerified); err != nil {
		return nil, err
	}

	if userWithAccount.Account.AccountID == "" {
		userWithAccount.Account.AccountID = userWithAccount.User.ID.String()
	}

	accountQuery := `
		INSERT INTO "accounts" (user_id, account_id, provider_id, password)
		VALUES ($1, $2, $3, $4)
		RETURNING *;
	`
	if err = tx.GetContext(ctx, &userWithAccount.Account, accountQuery, userWithAccount.User.ID, userWithAccount.Account.AccountID, userWithAccount.Account.ProviderId, userWithAccount.Account.Password); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return userWithAccount, nil
}

func (repo *UserRepository) GetUser(user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := repo.DBClient.GetContext(ctx, user, `
		SELECT u.*, r.name AS role
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1 OR u.email = $2;
	`, user.ID, user.Email); err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *UserRepository) GetAccount(account *models.Account) (*models.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := repo.DBClient.GetContext(ctx, account, `
		SELECT *
		FROM accounts a
		WHERE a.user_id = $1 AND a.provider_id = $2
	`, account.UserID, account.ProviderId); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return account, nil
}

func (repo *UserRepository) GetUserPermissionsByRoleID(roleID int) ([]models.Permission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	permissions := []models.Permission{}
	if err := repo.DBClient.SelectContext(ctx, &permissions, `
		SELECT p.*
		FROM permissions p
		LEFT JOIN role_permissions rp ON rp.permission_id = p.id AND rp.role_id = $1
	`, roleID); err != nil {
		return nil, err
	}

	return permissions, nil
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
