package models

type RegisterRequest struct {
	Name     string `json:"name" `
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type JwtTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	User        *User  `json:"user"`
	AccessToken string `json:"access_token"`
}

type RegisterResponse struct {
	User *User `json:"user"`
}
