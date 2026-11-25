package logic

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	auth "aqua/services/passport/domain"

	"golang.org/x/crypto/bcrypt"
)

// LoginByPassword authenticates the user with username/email/phone and password.
func (s *AuthService) LoginByPassword(ctx context.Context, account string, password string, loginCtx LoginContext) (*LoginResult, error) {
	logger := s.log(ctx).With(slog.String("user.account", account))
	um, err := s.repo.GetUserByAccount(ctx, account)
	if err != nil {
		logger.Warn("failed to load user by account",
			slog.Any("error", err))
		return nil, NewAuthError()
	}

	if err := bcrypt.CompareHashAndPassword([]byte(um.Password), []byte(password)); err != nil {
		logger.Warn("password verification failed",
			slog.Int("user.uid", um.UID))
		return nil, NewAuthError()
	}
	if um.Lock != 0 {
		logger.Warn("account is locked",
			slog.Int("user.uid", um.UID))
		return nil, NewAccountLockedError()
	}

	roles, loadErr := s.loadRoleNames(ctx, um.UID)
	if loadErr != nil {
		logger.Warn("failed to load user roles for token",
			slog.Int("user.uid", um.UID),
			slog.Any("error", loadErr))
		roles = []string{}
	}
	accessToken, err := s.generateAccessToken(um, roles)
	if err != nil {
		logger.Error("failed to generate access token",
			slog.Int("user.uid", um.UID),
			slog.Any("error", err))
		return nil, NewTokenError("token generation failed")
	}
	sessionID := uuid.New().String()
	refreshToken, jti, err := s.generateRefreshToken(um, sessionID)
	if err != nil {
		logger.Error("failed to generate refresh token",
			slog.Int("user.uid", um.UID),
			slog.Any("error", err))
		return nil, NewTokenError("token generation failed")
	}
	session := &auth.Session{
		SessionID:  sessionID,
		UID:        um.UID,
		DeviceInfo: loginCtx.DeviceInfo,
	}
	createdAt := time.Now().UTC()
	session.CreatedAt = &createdAt
	if err := s.repo.StoreUserSession(ctx, session, refreshTokenLife); err != nil {
		logger.Error("failed to store refresh token",
			slog.Int("user.uid", um.UID),
			slog.Any("error", err))
		return nil, NewTokenError("token generation failed")
	}

	go func(uid int, ip string, device string, refreshJTI string) {
		bgCtx := context.WithoutCancel(ctx)
		record := auth.LoginAuditRecord{
			UID:         uid,
			LoginIP:     ip,
			DeviceInfo:  device,
			LoginMethod: "password",
			RefreshJTI:  refreshJTI,
			CreatedAt:   time.Now().UTC(),
		}
		if err := s.repo.RecordLoginLog(bgCtx, record); err != nil {
			logger.Error("failed to log user login information",
				slog.Int("user.uid", uid),
				slog.Any("error", err))
		}
	}(um.UID, loginCtx.LoginIP, loginCtx.DeviceInfo, jti)

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Exp:          time.Now().Add(accessTokenLife).Unix(),
	}, nil
}
