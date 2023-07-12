package auth

type Tokens struct {
	IdToken      string
	AccessToken  string
	RefreshToken string
}

type SignUpInput struct {
	Email    string
	Password string
}

type SignInInput struct {
	Email    string
	Password string
}

type SignInWithGoogleInput struct {
	Code string
}

type VerifyEmailCodeInput struct {
	Code  string
	Email string
}

type RefreshAccessTokenInput struct {
	Token string
}
