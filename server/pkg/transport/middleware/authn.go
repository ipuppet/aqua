package middleware

import (
	"crypto/rsa"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/ipuppet/aqua/server/pkg/apperr"
	"github.com/ipuppet/aqua/server/pkg/transport/authctx"
	"github.com/ipuppet/aqua/server/pkg/transport/response"
	"github.com/ipuppet/aqua/server/pkg/transport/token"
)

// Authenticate verifies JWT (AuthN) and stores TokenUser in request context.
func Authenticate(publicKey *rsa.PublicKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, ok := parseBearerToken(c.GetHeader("Authorization"))
		if !ok {
			response.Fail(c, apperr.Unauthorized(""))
			return
		}

		claims, err := token.VerifyToken(publicKey, &token.AccessTokenClaims{}, tokenString)
		if err != nil || claims.Subject != token.SubjectAccessToken {
			response.Fail(c, apperr.Unauthorized(""))
			return
		}

		tokenUser := &authctx.TokenUser{UID: claims.UID, Roles: claims.Roles}
		if tokenUser.Roles == nil {
			tokenUser.Roles = []string{}
		}

		c.Request = c.Request.WithContext(authctx.WithTokenUser(c.Request.Context(), tokenUser))
		c.Next()
	}
}

func parseBearerToken(authorization string) (string, bool) {
	if authorization == "" {
		return "", false
	}
	if !strings.HasPrefix(authorization, "Bearer ") {
		return "", false
	}
	tokenString := strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
	if tokenString == "" {
		return "", false
	}
	return tokenString, true
}
