package auth

type Token string

type Auth interface {
	CreateToken() (Token, error)
	ValidateToken(token Token) error
}

type JWTAuth struct {
	secret string
}

func New(secret string) *JWTAuth {
	return &JWTAuth{
		secret: secret,
	}
}

func (j *JWTAuth) CreateToken() (Token, error) {
	return "", nil
}

func (j *JWTAuth) ValidateToken(token Token) error {
	return nil
}
