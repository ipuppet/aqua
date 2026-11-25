package logic

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/ipuppet/aqua/server/pkg/infra/config"
	"github.com/ipuppet/aqua/server/pkg/infra/email"
	"github.com/ipuppet/aqua/server/pkg/infra/filestore"
	"github.com/ipuppet/aqua/server/pkg/infra/logutil"
	"github.com/ipuppet/aqua/server/pkg/infra/models"
	"github.com/ipuppet/aqua/server/pkg/transport/token"
	"github.com/ipuppet/aqua/server/services/passport/dao"
)

type PasswordResetService struct {
	repo                    dao.AuthRepository
	emailConfigProvider     config.EmailConfigProvider
	turnstileConfigProvider config.TurnstileConfigProvider
	store                   *filestore.Client
}

func NewPasswordResetService(
	repo dao.AuthRepository,
	emailConfigProvider config.EmailConfigProvider,
	turnstileConfigProvider config.TurnstileConfigProvider,
	store *filestore.Client,
) *PasswordResetService {
	if emailConfigProvider == nil {
		emailConfigProvider = func() models.EmailConfig { return models.EmailConfig{} }
	}
	if turnstileConfigProvider == nil {
		turnstileConfigProvider = func() models.TurnstileConfig { return models.TurnstileConfig{} }
	}

	return &PasswordResetService{
		repo:                    repo,
		emailConfigProvider:     emailConfigProvider,
		turnstileConfigProvider: turnstileConfigProvider,
		store:                   store,
	}
}

func (s *PasswordResetService) log(ctx context.Context) *slog.Logger {
	return logutil.FromContext(ctx).With(slog.String("service", "password_reset"))
}

func generateResetToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *PasswordResetService) SendPasswordResetEmail(ctx context.Context, emailAddress, turnstileToken string) error {
	logger := s.log(ctx).With(slog.String("user.email", emailAddress))

	isValid, err := verifyTurnstileToken(ctx, s.turnstileConfigProvider(), turnstileToken, "")
	if err != nil {
		logger.Error("Turnstile verification request failed", "error", err)
		return err
	}
	if !isValid {
		logger.Warn("Turnstile token validation failed")
		return NewValidationError("bot verification failed")
	}

	// Verify user exists
	user, err := s.repo.GetUserByAccount(ctx, emailAddress)
	if err != nil {
		// Log the error but returning generic error to prevent email enumeration
		logger.Warn("Password reset requested for non-existent/invalid email", "email", emailAddress, "error", err)
		return nil // Pretend email was sent
	}

	token := generateResetToken()

	if err := s.repo.StorePasswordResetToken(ctx, token, user.Email); err != nil {
		logger.Error("Failed to cache reset token", "error", err)
		return fmt.Errorf("could not generate reset token")
	}

	// Prepare email

	templateContent, readErr := s.store.GetContent("email", "password_reset.html")
	var tmpl *template.Template
	var parseErr error
	if readErr == nil {
		tmpl, parseErr = template.New("email").Parse(templateContent)
	} else {
		parseErr = readErr
	}

	var body string
	if parseErr != nil {
		logger.Error("Failed to read/parse email template", "error", parseErr)
		return parseErr
	}
	var buf bytes.Buffer
	data := struct {
		Nickname string
		Username string
		Token    string
		Year     int
	}{
		Nickname: user.Nickname,
		Username: user.Username,
		Token:    token,
		Year:     time.Now().Year(),
	}
	if execErr := tmpl.Execute(&buf, data); execErr != nil {
		logger.Error("Failed to execute email template", "error", execErr)
		return fmt.Errorf("failed to generate email body")
	}
	body = buf.String()

	err = email.SendEmailWithConfig(s.emailConfigProvider(), user.Email, "Reset Your Password", body)
	if err != nil {
		logger.Error("Failed to send reset email", "error", err)
		return fmt.Errorf("failed to send reset email")
	}

	logger.Info("Password reset email sent", "email", user.Email)
	return nil
}

func (s *PasswordResetService) ResetPassword(ctx context.Context, resetToken, newPassword string) error {
	logger := s.log(ctx)

	emailAddress, err := s.repo.GetPasswordResetEmail(ctx, resetToken)
	if err != nil {
		logger.Error("Invalid or expired reset token", "error", err)
		return errors.New("invalid or expired reset token")
	}

	userModel, err := s.repo.GetUserByAccount(ctx, emailAddress)
	if err != nil {
		logger.Error("Failed to fetch user by linked email", "error", err, "email", emailAddress)
		return errors.New("user not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to hash new password", "error", err)
		return fmt.Errorf("internal server error")
	}

	if err := s.repo.UpdateUserPassword(ctx, userModel.UID, string(hashedPassword)); err != nil {
		logger.Error("Failed to update password", "error", err)
		return fmt.Errorf("failed to reset password")
	}

	if err := token.RevokeUserRefreshTokens(ctx, s.repo.Cache(), userModel.UID); err != nil {
		logger.Warn("Failed to write refresh revoke marker after password reset",
			"uid", userModel.UID,
			"error", err)
	}

	// Consume the token
	if err := s.repo.DeletePasswordResetToken(ctx, resetToken); err != nil {
		logger.Warn("Failed to invalidate reset token", "token", resetToken, "error", err)
	}

	logger.Info("Password reset successful", "uid", userModel.UID)
	return nil
}
