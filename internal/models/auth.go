package models

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
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
	User        User   `json:"user"`
	AccessToken string `json:"access_token"`
}

type RegisterResponse struct {
	User `json:",inline"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}
