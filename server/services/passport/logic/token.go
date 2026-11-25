package logic

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"aqua/pkg/transport/token"
	authdomain "aqua/services/passport/domain"
)

func (s *AuthService) generateAccessToken(um *authdomain.Account, roles []string) (string, error) {
	expireTime := time.Now().Add(accessTokenLife)
	if roles == nil {
		roles = []string{}
	}
	jwtClaims := token.AccessTokenClaims{
		UID:      um.UID,
		Roles:    roles,
		Username: um.Username,
		Nickname: um.Nickname,
		Avatar:   um.Avatar,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    JwtIssuer,
			Subject:   token.SubjectAccessToken,
		},
	}
	return token.SignToken(s.privateKey, jwtClaims)
}

// generateRefreshToken returns (tokenString, jti, error). Caller must StoreUserSession.
func (s *AuthService) generateRefreshToken(um *authdomain.Account, sessionID string) (string, string, error) {
	expireTime := time.Now().Add(refreshTokenLife)
	jti := uuid.New().String()
	jwtClaims := token.RefreshTokenClaims{
		UID:       um.UID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    JwtIssuer,
			Subject:   token.SubjectRefreshToken,
		},
	}
	signedToken, err := token.SignToken(s.privateKey, jwtClaims)
	if err != nil {
		return "", "", err
	}
	return signedToken, jti, nil
}

// VerifyAndRefreshToken validates tokens and refreshes them if necessary.
// Verification uses the public key derived from the private key—no user password lookup or cache required.
func (s *AuthService) VerifyAndRefreshToken(ctx context.Context, accessTokenString string, refreshTokenString string) (*TokenRefreshResult, error) {
	logger := s.log(ctx)
	publicKey := &s.privateKey.PublicKey
	accessTokenClaims, err := token.VerifyToken(publicKey, &token.AccessTokenClaims{}, accessTokenString)
	if err == nil && accessTokenClaims.Subject == token.SubjectAccessToken {
		um, userErr := s.repo.GetUserByUid(ctx, accessTokenClaims.UID)
		if userErr != nil {
			logger.Warn("failed to load user by access token",
				slog.Int("user.uid", accessTokenClaims.UID),
				slog.Any("error", userErr))
			return nil, NewAuthError()
		}
		if um.Lock != 0 {
			logger.Warn("access denied: account is locked",
				slog.Int("user.uid", um.UID))
			return nil, NewAccountLockedError()
		}
		return &TokenRefreshResult{}, nil
	}

	refreshTokenClaims, refreshErr := token.VerifyToken(publicKey, &token.RefreshTokenClaims{}, refreshTokenString)
	if refreshErr != nil || refreshTokenClaims.Subject != token.SubjectRefreshToken {
		if refreshErr == nil {
			refreshErr = errors.New("invalid refresh token subject")
		}
		logger.Error("failed to parse refresh token",
			slog.Any("error", refreshErr))
		return nil, NewTokenError("refresh token parse failed: " + refreshErr.Error())
	}
	sessionID := refreshTokenClaims.GetSessionID()
	if sessionID == "" {
		logger.Warn("refresh token missing session_id",
			slog.Int("user.uid", refreshTokenClaims.UID),
			slog.String("jti", refreshTokenClaims.TokenID()))
		return nil, NewTokenError("refresh token missing session_id")
	}
	revoked, revokeErr := token.IsRefreshTokenRevoked(ctx, s.repo.Cache(), refreshTokenClaims.UID, refreshTokenClaims.IssuedAt.Unix())
	if revokeErr != nil {
		logger.Error("failed to check refresh token revoke marker",
			slog.Int("user.uid", refreshTokenClaims.UID),
			slog.Any("error", revokeErr))
		return nil, NewTokenError("refresh token verification failed")
	}
	if revoked {
		logger.Warn("refresh token rejected by revoke marker",
			slog.Int("user.uid", refreshTokenClaims.UID),
			slog.String("jti", refreshTokenClaims.TokenID()))
		_ = s.repo.DeleteUserSession(ctx, refreshTokenClaims.UID, sessionID)
		return nil, NewTokenError("refresh token invalid or expired")
	}
	session, err := s.repo.GetUserSession(ctx, refreshTokenClaims.UID, sessionID)
	if err != nil {
		logger.Warn("refresh token not found in Redis (revoked or expired)",
			slog.Any("error", err))
		return nil, NewTokenError("refresh token verification failed")
	}
	if session.UID != refreshTokenClaims.UID {
		logger.Warn("refresh token uid mismatch",
			slog.Int("token.uid", refreshTokenClaims.UID),
			slog.Int("stored.uid", session.UID))
		return nil, NewTokenError("refresh token verification failed")
	}

	// Rotate: shorten the old token TTL to a grace period so that concurrent
	// in-flight requests can still succeed, then issue a new pair.
	session.SessionID = sessionID
	_ = s.repo.StoreUserSession(ctx, session, refreshTokenGracePeriod)

	um, userErr := s.repo.GetUserByUid(ctx, refreshTokenClaims.UID)
	if userErr != nil {
		logger.Warn("failed to load user by refresh token",
			slog.Int("user.uid", refreshTokenClaims.UID),
			slog.Any("error", userErr))
		return nil, NewAuthError()
	}
	if um.Lock != 0 {
		logger.Warn("refresh denied: account is locked",
			slog.Int("user.uid", um.UID))
		_ = s.repo.DeleteUserSession(ctx, refreshTokenClaims.UID, sessionID)
		return nil, NewAccountLockedError()
	}

	roles, loadErr := s.loadRoleNames(ctx, um.UID)
	if loadErr != nil {
		logger.Warn("failed to load user roles for token",
			slog.Int("user.uid", um.UID),
			slog.Any("error", loadErr))
		roles = []string{}
	}
	newAccessToken, err := s.generateAccessToken(um, roles)
	if err != nil {
		logger.Error("failed to generate new access token",
			slog.Int("user.uid", um.UID),
			slog.Any("error", err))
		return nil, NewTokenError("access token generation failed")
	}
	newRefreshToken, _, err := s.generateRefreshToken(um, sessionID)
	if err != nil {
		logger.Error("failed to generate new refresh token",
			slog.Int("user.uid", um.UID),
			slog.Any("error", err))
		return nil, NewTokenError("refresh token generation failed")
	}
	session.UID = um.UID
	session.SessionID = sessionID
	createdAt := time.Now().UTC()
	session.CreatedAt = &createdAt
	if err := s.repo.StoreUserSession(ctx, session, refreshTokenLife); err != nil {
		logger.Error("failed to store new refresh token",
			slog.Int("user.uid", um.UID),
			slog.Any("error", err))
		return nil, NewTokenError("refresh token save failed")
	}

	return &TokenRefreshResult{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) loadRoleNames(ctx context.Context, uid int) ([]string, error) {
	roles, err := s.rbacReader.LoadUserRoles(ctx, uid)
	if err != nil {
		return nil, err
	}
	roleNames := make([]string, 0, len(roles))
	for _, role := range roles {
		roleNames = append(roleNames, role.Code)
	}
	return roleNames, nil
}
