package logic

import (
	"context"
	"log/slog"
	"math/rand"

	"golang.org/x/crypto/bcrypt"

	"aqua/pkg/infra/config"
	"aqua/pkg/infra/email"
	"aqua/pkg/infra/logutil"
	"aqua/pkg/infra/models"
	"aqua/services/passport/dao"
	auth "aqua/services/passport/domain"
)

type RegistrationService struct {
	repo                    dao.AuthRepository
	emailConfigProvider     config.EmailConfigProvider
	turnstileConfigProvider config.TurnstileConfigProvider
}

func NewRegistrationService(
	repo dao.AuthRepository,
	emailConfigProvider config.EmailConfigProvider,
	turnstileConfigProvider config.TurnstileConfigProvider,
) *RegistrationService {
	if emailConfigProvider == nil {
		emailConfigProvider = func() models.EmailConfig { return models.EmailConfig{} }
	}
	if turnstileConfigProvider == nil {
		turnstileConfigProvider = func() models.TurnstileConfig { return models.TurnstileConfig{} }
	}

	return &RegistrationService{
		repo:                    repo,
		emailConfigProvider:     emailConfigProvider,
		turnstileConfigProvider: turnstileConfigProvider,
	}
}

func (s *RegistrationService) log(ctx context.Context) *slog.Logger {
	return logutil.FromContext(ctx).With(slog.String("service", "registration"))
}

// RegisterUser validates verification code and persists the user.
func (s *RegistrationService) RegisterUser(ctx context.Context, um *auth.Account, verificationCode string) error {
	logger := s.log(ctx).With(slog.String("user.email", um.Email))
	savedCode, err := s.repo.GetVerificationCode(ctx, um.Email)
	if err != nil || savedCode == "" {
		logger.Warn("verification code not found or expired",
			slog.Any("error", err))
		return NewValidationError("verification code has expired or does not exist")
	}
	if savedCode != verificationCode {
		logger.Warn("verification code mismatch")
		return NewValidationError("email verification code is incorrect")
	}

	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(um.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		logger.Error("failed to hash user password",
			slog.Any("error", hashErr))
		return hashErr
	}
	um.Password = string(hashedPassword)

	if err := s.repo.RegisterUser(ctx, um); err != nil {
		logger.Error("failed to register new user",
			slog.Any("error", err))
		return err
	}

	if err := s.repo.DeleteVerificationCode(ctx, um.Email); err != nil {
		logger.Warn("failed to delete verification code after registration",
			slog.Any("error", err))
	}

	return nil
}

// SendVerificationCode validates the anti-bot token, emails the code, and caches it.
func (s *RegistrationService) SendVerificationCode(ctx context.Context, address string, turnstileToken string) error {
	logger := s.log(ctx).With(slog.String("user.email", address))
	isValid, err := verifyTurnstileToken(ctx, s.turnstileConfigProvider(), turnstileToken, "")
	if err != nil {
		logger.Error("Turnstile verification request failed",
			slog.Any("error", err))
		return err
	}
	if !isValid {
		logger.Warn("Turnstile token validation failed")
		return NewValidationError("bot verification failed")
	}

	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	verificationCode := make([]byte, 6)
	for i := range verificationCode {
		verificationCode[i] = charset[rand.Intn(len(charset))]
	}

	if err := email.SendEmailWithConfig(s.emailConfigProvider(), address, "Email Verification", "Your verification code is: "+string(verificationCode)); err != nil {
		logger.Error("failed to send verification email",
			slog.Any("error", err))
		return err
	}

	if err := s.repo.StoreVerificationCode(ctx, address, string(verificationCode)); err != nil {
		logger.Error("failed to persist verification code",
			slog.Any("error", err))
		return err
	}

	return nil
}
