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
