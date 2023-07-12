package auth

type SignUpWithEmailRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInWithEmailRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInWithGoogleRequest struct {
	Code string `json:"code"`
}

type VerifyEmailRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type RefreshAccessTokenRequest struct {
	Token string
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	IDToken      string `json:"idToken"`
}
