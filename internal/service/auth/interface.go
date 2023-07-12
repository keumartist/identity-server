package auth

type AuthService interface {
	SignUpWithEmail(input SignUpInput) (string, error)
	SignInWithEmail(input SignInInput) (Tokens, error)
	SignInWithGoogle(input SignInWithGoogleInput) (Tokens, error)
	VerifyEmailCode(input VerifyEmailCodeInput) error
	RefreshAccessToken(input RefreshAccessTokenInput) (Tokens, error)
}
