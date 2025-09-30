package auth

import (
	"dev-go-apis/internal/lib"
	"dev-go-apis/internal/models"
	"fmt"
	"log"
	"time"
)

type IUserRepository interface {
	CreateUserAccount(userWithAccount *models.UserWithAccount) (*models.UserWithAccount, error)
	CreateUser(user *models.User) (*models.User, error)
	CreateAccount(account *models.Account) (*models.Account, error)
	GetUser(user *models.User) (*models.User, error)
	GetAccount(account *models.Account) (*models.Account, error)
	UpdateUserById(id string, user *models.User) (*models.User, error)
}

type ICacheRepository interface {
	SetValue(string, interface{}, time.Duration) error
	GetValue(string) (string, error)
	DeleteValue(keys ...string) (bool, error)
}

type AuthService struct {
	UserRepo  IUserRepository
	CacheRepo ICacheRepository
}

func NewAuthService(userRepo IUserRepository, cacheRepo ICacheRepository) *AuthService {
	return &AuthService{
		UserRepo:  userRepo,
		CacheRepo: cacheRepo,
	}
}

func (s *AuthService) CheckOAuthState(state string) (bool, error) {
	_, err := s.CacheRepo.GetValue(state)
	if err != nil {
		return false, fmt.Errorf("failed to get oauth state: %s", err.Error())
	}
	return true, nil
}

func (s *AuthService) GenerateOAuthState(prefix string) (string, error) {
	state := fmt.Sprintf("%s%s", prefix, lib.GenerateOTP(uint32(8)))
	if err := s.CacheRepo.SetValue(state, 1, time.Minute*5); err != nil {
		return "", fmt.Errorf("failed to generate oauth state: %s", err.Error())
	}
	return state, nil
}

func (s *AuthService) GenerateJwtTokens(userWithClaims *models.UserWithClaims) (*models.JwtTokens, error) {
	accessToken, err := lib.SignAccessToken(userWithClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %s", err.Error())
	}

	refreshToken, err := lib.SignRefreshToken(userWithClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %s", err.Error())
	}

	return &models.JwtTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) LoginWithGoogle(googleUserInfo *models.GoogleUserInfo) (*models.User, error) {
	user := &models.User{
		Name:          fmt.Sprintf("%s %s", googleUserInfo.FamilyName, googleUserInfo.GivenName),
		Email:         googleUserInfo.Email,
		Image:         googleUserInfo.Picture,
		EmailVerified: googleUserInfo.VerifiedEmail,
	}

	account := &models.Account{
		ProviderId: models.ProviderGoogle,
		AccountID:  googleUserInfo.ID,
	}

	_, err := s.UserRepo.GetUser(user)
	if err != nil {
		log.Printf("%v\n", err.Error())
		userWithAccount, err := s.UserRepo.CreateUserAccount(&models.UserWithAccount{
			User:    *user,
			Account: *account,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %s", err.Error())
		}

		return &userWithAccount.User, nil
	}

	if !user.EmailVerified {
		user.EmailVerified = googleUserInfo.VerifiedEmail
		_, err := s.UserRepo.UpdateUserById(user.ID.String(), user)
		if err != nil {
			return nil, fmt.Errorf("failed to update user: %s", err.Error())
		}
	}

	account.UserID = user.ID
	result, err := s.UserRepo.GetAccount(account)
	if err != nil {
		return nil, err
	}
	if result != nil {
		return user, nil
	}

	_, err = s.UserRepo.CreateAccount(account)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %s", err.Error())
	}

	return user, nil
}

func (s *AuthService) Login(body *models.LoginRequest) (*models.User, error) {
	user, err := s.UserRepo.GetUser(&models.User{
		Email: body.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	account, err := s.UserRepo.GetAccount(&models.Account{
		UserID:     user.ID,
		ProviderId: models.ProviderCredential,
	})
	if err != nil {
		return nil, fmt.Errorf("account not found")
	}
	if err := lib.ComparePassword(account.Password, body.Password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

func (s *AuthService) Register(body *models.RegisterRequest) (*models.User, error) {
	hashPassword, err := lib.HashPassword(body.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %s", err.Error())
	}

	user := &models.User{
		Name:  body.Name,
		Email: body.Email,
		Image: lib.GetRandomAvatar(),
	}

	account := &models.Account{
		ProviderId: models.ProviderCredential,
		Password:   hashPassword,
	}

	_, err = s.UserRepo.GetUser(user)
	if err == nil {
		account.AccountID = user.ID.String()
		account.UserID = user.ID
		_, err := s.UserRepo.GetAccount(account)
		if err == nil {
			return nil, fmt.Errorf("user already exists")
		}
		return user, nil
	}

	userWithAccount, err := s.UserRepo.CreateUserAccount(&models.UserWithAccount{
		User:    *user,
		Account: *account,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register: %s", err.Error())
	}

	return &userWithAccount.User, nil
}
