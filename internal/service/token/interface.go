package token

type TokenService interface {
	GenerateAccessToken(id, email string, expireAt uint) (string, error)
	GenerateRefreshToken(id, email string, expireAt uint) (string, error)
	GenerateIdToken(id, email string, expireAt uint) (string, error)
}
