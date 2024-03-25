package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/TekClinic/API-Gateway/schemas"
	"github.com/gin-gonic/gin"
)

const TokenKey = "token"

// extractBearerToken returns bearer token that was passed in the request.
func extractBearerToken(ctx *gin.Context) (string, error) {
	header := ctx.GetHeader("Authorization")
	if header == "" {
		return "", errors.New("bearer token is missing")
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return parts[1], nil
}

// AuthRequired middleware validates that authorization token was passed in the request
// It DOESN'T check whether the token is valid. The responsibility of such check is an end-user
// The token is stored in ctx under key tokenKey.
func AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jwtToken, err := extractBearerToken(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}
		ctx.Set(TokenKey, jwtToken)
		ctx.Next()
	}
}
