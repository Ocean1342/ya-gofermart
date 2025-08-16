package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Token string

type Auth interface {
	CreateToken() (Token, error)
	ValidateToken(token Token) error
	CreateHash(inside string) (string, error)
	CompareHashAndPassword(password, hash string) bool
}

type JWTAuth struct {
	secret string
	ttl    time.Duration
}

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

func New(secret string, ttl time.Duration) *JWTAuth {
	return &JWTAuth{
		secret: secret,
	}
}

func (j *JWTAuth) CreateToken() (Token, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ttl)),
		},
		// собственное утверждение
		UserID: 1,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return Token(tokenString), nil
}

func (j *JWTAuth) ValidateToken(token Token) error {
	return nil
}

// TODO: стоит ли закладывать возможность изменения алгоритма хеширования на уровне БД?
func (j *JWTAuth) CreateHash(inside string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(inside), bcrypt.DefaultCost)
	return string(bytes), err
	return inside, nil
}
func (j *JWTAuth) CompareHashAndPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
