package token

import "github.com/golang-jwt/jwt/v5"

const (
	SubjectAccessToken  = "UserAccessToken"
	SubjectRefreshToken = "UserRefreshToken"
)

// AccessTokenClaims represents the payload embedded in an access token.
type AccessTokenClaims struct {
	UID      int      `json:"uid"`
	Roles    []string `json:"roles"`
	Username string   `json:"un"`
	Nickname string   `json:"nk"`
	Avatar   string   `json:"av"`
	jwt.RegisteredClaims
}

// RefreshTokenClaims represents the payload embedded in a refresh token.
// It keeps jwt.RegisteredClaims for compatibility; jti is stored in StandardClaims.Id.
type RefreshTokenClaims struct {
	UID       int    `json:"uid"`
	SessionID string `json:"session_id,omitempty"`
	jwt.RegisteredClaims
}

// ClaimsWithUID represents a token claim that provides user identification.
type ClaimsWithUID interface {
	jwt.Claims
	GetUID() int
}

// GetUID returns the user identifier bound to the access token.
func (c *AccessTokenClaims) GetUID() int {
	return c.UID
}

// GetUID returns the user identifier bound to the refresh token.
func (c *RefreshTokenClaims) GetUID() int {
	return c.UID
}

// TokenID returns the refresh token jti.
func (c *RefreshTokenClaims) TokenID() string {
	return c.ID
}

// GetSessionID returns stable session id for refresh token rotation.
func (c *RefreshTokenClaims) GetSessionID() string {
	return c.SessionID
}
