package token

type TokenType int

const (
	AccessToken TokenType = iota
	IdToken
	RefreshToken
)

type GenerateTokenInput struct {
	Id                  string
	Email               string
	ExpirationInSeconds uint
	TokenType           TokenType
}

type VerifyTokenInput struct {
	Token     string
	TokenType TokenType
}

type GetUserIDFromTokenInput struct {
	Token     string
	TokenType TokenType
}
