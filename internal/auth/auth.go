package auth

import (
	"fmt"
	"github.com/Ocean1342/ya-gofermart/internal/storage"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Token string

type Auth interface {
	CreateToken(user *storage.User) (Token, error)
	GetUserFromToken(token Token) (*storage.User, error)
	CreateHash(inside string) (string, error)
	CompareHashAndPassword(password, hash string) bool
}

type JWTAuth struct {
	secret string
	ttl    time.Duration
}

type Claims struct {
	jwt.RegisteredClaims
	User *storage.User
}

func New(secret string, ttl time.Duration) *JWTAuth {
	return &JWTAuth{
		secret: secret,
		ttl:    ttl,
	}
}

func (j *JWTAuth) CreateToken(user *storage.User) (Token, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ttl)),
		},
		User: user,
	})
	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}
	return Token(tokenString), nil
}

func (j *JWTAuth) GetUserFromToken(token Token) (*storage.User, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(string(token), claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		logrus.Errorf("could not parse token. err:`%s`, token:`%s`", token)
		return nil, fmt.Errorf("could not parse token. err:`%s`", err)
	}
	return claims.User, nil
}

// TODO: вопрос -  стоит ли закладывать возможность изменения алгоритма хеширования на уровне БД?
func (j *JWTAuth) CreateHash(inside string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(inside), bcrypt.DefaultCost)
	return string(bytes), err
}
func (j *JWTAuth) CompareHashAndPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
