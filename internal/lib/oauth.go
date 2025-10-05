package lib

import (
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	StateDuration = time.Minute * 3
)

var (
	GoogleStatePrefix = "G"
	GoogleScopes      = []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	}
	GoogleOAuthUserInfoURl = "https://www.googleapis.com/oauth2/v2/userinfo"
	GoogleOAuthConfig      = &oauth2.Config{
		RedirectURL:  GOOGLE_REDIRECT_URL,
		ClientID:     GOOGLE_CLIENT_ID,
		ClientSecret: GOOGLE_CLIENT_SECRET,
		Scopes:       GoogleScopes,
		Endpoint:     google.Endpoint,
	}
)
