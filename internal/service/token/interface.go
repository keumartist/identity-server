package token

type TokenService interface {
	GenerateAccessToken(id, email string) (string, error)
	GenerateRefreshToken(id, email string) (string, error)
}
