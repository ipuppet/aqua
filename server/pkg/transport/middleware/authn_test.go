package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"aqua/pkg/transport/authctx"
	"aqua/pkg/transport/token"
)

func TestAuthenticate_UnauthorizedWithoutBearer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	r := gin.New()
	r.Use(Authenticate(&priv.PublicKey))
	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthenticate_WithValidTokenSetsPrincipal(t *testing.T) {
	gin.SetMode(gin.TestMode)
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	token, err := token.SignToken(priv, &token.AccessTokenClaims{
		UID:   42,
		Roles: []string{"blog.manager"},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(600 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   token.SubjectAccessToken,
		},
	})
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	r := gin.New()
	r.Use(Authenticate(&priv.PublicKey))
	r.GET("/", func(c *gin.Context) {
		tokenUser := authctx.TokenUserFromContext(c.Request.Context())
		if tokenUser == nil {
			t.Fatalf("token user should not be nil")
		}
		if tokenUser.UID != 42 {
			t.Fatalf("expected uid 42, got %d", tokenUser.UID)
		}
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
	}
}
