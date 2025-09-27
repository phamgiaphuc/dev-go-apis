package auth

import (
	"dev-go-apis/internal/lib"
	"dev-go-apis/internal/models"
	"fmt"
)

type IUserRepository interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserWithPasswordByEmail(email string) (*models.UserWithPassword, error)
	CreateUser(userWithPassword *models.UserWithPassword) (*models.User, error)
}

type AuthService struct {
	UserRepo IUserRepository
}

func NewAuthService(userRepo IUserRepository) *AuthService {
	return &AuthService{
		UserRepo: userRepo,
	}
}

func (s *AuthService) GenerateJwtTokens(userWithClaims *models.UserWithClaims) (*models.JwtTokens, error) {
	fmt.Printf("Generating tokens for user: %+v\n", userWithClaims.User)
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

func (s *AuthService) Login(body *models.LoginRequest) (*models.User, error) {
	userWithPassword, err := s.UserRepo.GetUserWithPasswordByEmail(body.Email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if err := lib.ComparePassword(userWithPassword.Password, body.Password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return &userWithPassword.User, nil
}

func (s *AuthService) Register(body *models.RegisterRequest) (*models.User, error) {
	_, err := s.UserRepo.GetUserByEmail(body.Email)
	if err == nil {
		return nil, fmt.Errorf("user already exists")
	}

	hashPassword, err := lib.HashPassword(body.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %s", err.Error())
	}

	userWithPassword := &models.UserWithPassword{
		User: models.User{
			Name:  body.Name,
			Email: body.Email,
			Image: lib.GetRandomAvatar(),
		},
		Password: hashPassword,
	}

	user, err := s.UserRepo.CreateUser(userWithPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %s", err.Error())
	}

	return user, nil
}
