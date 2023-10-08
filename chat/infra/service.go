package infra

import (
	"context"
	"errors"
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

var errInvalidJwtToken = errors.New("Invalid JWT token")

type Claims struct {
	UserName  string `json:"userName"`
	UserEmail string `json:"userEmail"`
	jwt.RegisteredClaims
}

type TokenServiceJWT struct {
	logger               *logrus.Logger
	jwtSecretKey         []byte
	tokenExpirationHours uint
}

func NewTokenServiceJWT(l *logrus.Logger, k string, t uint) domain.TokenService {
	return TokenServiceJWT{
		logger:               l,
		jwtSecretKey:         []byte(k),
		tokenExpirationHours: t,
	}
}

func (s TokenServiceJWT) Generate(ctx context.Context, user domain.User) (domain.UserToken, error) {
	expirationTime := time.Now().Add(time.Duration(s.tokenExpirationHours) * time.Hour)

	claims := &Claims{
		UserName:  user.Name,
		UserEmail: user.Email,
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

func (s TokenServiceJWT) Check(ctx context.Context, userToken domain.UserToken) (domain.User, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(userToken.AccessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecretKey, nil
	})
	if err != nil {
		s.logger.Error(err)
		return domain.User{}, errInvalidJwtToken
	}

	if !token.Valid {
		return domain.User{}, errInvalidJwtToken
	}

	return domain.NewUser("", claims.UserName, claims.UserEmail, "")
}
