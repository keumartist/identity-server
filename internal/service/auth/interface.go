package auth

type AuthService interface {
	SignUpWithEmail(input SignUpInput) (Tokens, error)
	SignInWithEmail(input SignInInput) (Tokens, error)
}
