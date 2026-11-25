package logic

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ipuppet/aqua/server/pkg/infra/logutil"
	"github.com/ipuppet/aqua/server/pkg/infra/models"
	auth "github.com/ipuppet/aqua/server/services/passport/domain"
)

func TestRegisterUser(t *testing.T) {
	// 1. Arrange: Setup In-Memory database (Given)
	repo := newMemoryAuthRepo()

	// Pre-populate some existing data
	existingEmail := "existing@example.com"
	repo.AddUser(&auth.Account{
		Email: existingEmail,
	})

	// Add valid verification code
	verificationKey := "test@example.com"
	validCode := "123456"
	_ = repo.StoreVerificationCode(context.Background(), verificationKey, validCode)

	svc := NewRegistrationService(
		repo,
		func() models.EmailConfig { return models.EmailConfig{} },
		func() models.TurnstileConfig { return models.TurnstileConfig{} },
	)
	ctx := logutil.WithLogger(context.Background(), slog.Default())

	tests := []struct {
		name             string
		email            string
		password         string
		verificationCode string
		expectError      bool
	}{
		{
			name:             "Successful Registration",
			email:            "test@example.com",
			password:         "StrongPass1!",
			verificationCode: validCode,
			expectError:      false,
		},
		{
			name:             "Incorrect Verification Code",
			email:            "test2@example.com",
			password:         "StrongPass1!",
			verificationCode: "000000",
			expectError:      true,
		},
		{
			name:             "Existing Email Refused",
			email:            existingEmail,
			password:         "StrongPass1!",
			verificationCode: validCode,
			expectError:      true, // Validation should throw an error due to unique constraint or similar
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Pre-Seed verification code for individual test case if it corresponds to an actual test email
			if !tc.expectError {
				// Ensuring our code is definitely matching the email in reality
				// (though our logic handles email to key translation in RegisterUser)
				_ = repo.StoreVerificationCode(ctx, tc.email, tc.verificationCode)
			}

			um := &auth.Account{
				Email:    tc.email,
				Password: tc.password,
			}

			// 2. Act: Execute register function (When)
			err := svc.RegisterUser(ctx, um, tc.verificationCode)

			// 3. Assert: Validate creation state (Then)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify DB state actually changed
				savedUser, getErr := repo.GetUserByAccount(ctx, tc.email)
				assert.NoError(t, getErr)
				assert.NotNil(t, savedUser)
				assert.Equal(t, tc.email, savedUser.Email)

				// Ensure code gets deleted after registration
				_, codeErr := repo.GetVerificationCode(ctx, tc.email)
				assert.Error(t, codeErr, "Expected verification code to be deleted after registration")
			}
		})
	}
}
