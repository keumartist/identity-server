package token

type TokenService interface {
	GenerateToken(input GenerateTokenInput) (string, error)
	VerifyToken(input VerifyTokenInput) (bool, string, string, error)
}
