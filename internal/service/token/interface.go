package token

type TokenService interface {
	GenerateAccessToken(id, email string, expirationInSeconds uint) (string, error)
	GenerateRefreshToken(id, email string, expirationInSeconds uint) (string, error)
	GenerateIdToken(id, email string, expirationInSeconds uint) (string, error)
}
