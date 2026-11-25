package logic

import (
	"github.com/google/wire"
)

// Set is logic providers.
var Set = wire.NewSet(
	NewAuthService,
	NewRegistrationService,
	NewPasswordResetService,
	NewSessionService,
)
