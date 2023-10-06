package infra

import (
	"context"
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type TokenServiceJWT struct {
	logger               *logrus.Logger
	jwtSecretKey         []byte
	tokenExpirationHours uint
}

func NewTokenServiceJWT(l *logrus.Logger, s []byte, t uint) domain.TokenService {
	return TokenServiceJWT{
		logger:               l,
		jwtSecretKey:         s,
		tokenExpirationHours: t,
	}
}

func (s TokenServiceJWT) Generate(ctx context.Context, user domain.User) (domain.UserToken, error) {
	expirationTime := time.Now().Add(time.Duration(s.tokenExpirationHours) * time.Hour)

	claims := &Claims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecretKey)
	if err != nil {
		return domain.UserToken{}, nil
	}

	return domain.NewUserToken(tokenString), nil
}
