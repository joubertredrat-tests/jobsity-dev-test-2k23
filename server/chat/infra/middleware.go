package infra

import (
	"joubertredrat-tests/jobsity-dev-test-2k23/chat/domain"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	HEADER_CONTENT_TYPE  = "Content-Type"
	HEADER_AUTHORIZATION = "Authorization"
)

func JSONBodyMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.GetHeader(HEADER_CONTENT_TYPE) != "application/json" {
			ctx.Abort()
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "This route expects Content-Type with application/json"})
		}
	}
}

func JwtCheckMiddleware(tokenService domain.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader(HEADER_AUTHORIZATION)
		if token == "" {
			forbidden(ctx)
			return
		}

		jwtToken := strings.Split(token, " ")
		if len(jwtToken) != 2 {
			forbidden(ctx)
			return
		}
		jwtTokenString := jwtToken[1]
		if jwtTokenString == "" {
			forbidden(ctx)
			return
		}

		user, err := tokenService.Check(ctx, domain.NewUserToken(jwtTokenString))
		if err != nil {
			forbidden(ctx)
			return
		}

		ctx.Set("userAuth", user)
	}
}

func forbidden(ctx *gin.Context) {
	ctx.Abort()
	ctx.JSON(http.StatusForbidden, gin.H{"error": "You need to be authenticated to use this route"})
}
