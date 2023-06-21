package auth

type SignUpWithEmailRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInWithEmailRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	IDToken      string `json:"idToken"`
}
