package grpcserver

import (
	"context"
	"database/sql"
	"net/mail"
	"regexp"
	"strings"

	"google.golang.org/grpc/metadata"

	passportv1 "aqua/gen/go/aqua/passport/v1"
	"aqua/pkg/apperr"
	"aqua/pkg/transport/authctx"
	"aqua/pkg/types"
	authdomain "aqua/services/passport/domain"
	"aqua/services/passport/logic"
)

var usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

type PassportServer struct {
	authService         *logic.AuthService
	registrationService *logic.RegistrationService
	passwordReset       *logic.PasswordResetService
	sessionService      *logic.SessionService
	passportv1.UnimplementedPassportServiceServer
}

func NewPassportServer(
	authService *logic.AuthService,
	registrationService *logic.RegistrationService,
	passwordReset *logic.PasswordResetService,
	sessionService *logic.SessionService,
) *PassportServer {
	return &PassportServer{
		authService:         authService,
		registrationService: registrationService,
		passwordReset:       passwordReset,
		sessionService:      sessionService,
	}
}

func (s *PassportServer) LoginByPassword(ctx context.Context, req *passportv1.LoginByPasswordRequest) (*passportv1.LoginByPasswordResponse, error) {
	account := strings.TrimSpace(req.GetAccount())
	password := req.GetPassword()
	if account == "" {
		return nil, apperr.BadRequest("account is required")
	}
	if len(password) < 6 || len(password) > 50 {
		return nil, apperr.BadRequest("password must be between 6 and 50 characters")
	}

	loginCtx := extractLoginContext(ctx)
	loginInfo, err := s.authService.LoginByPassword(ctx, account, password, loginCtx)
	if err != nil {
		return nil, err
	}
	return &passportv1.LoginByPasswordResponse{
		AccessToken:  loginInfo.AccessToken,
		RefreshToken: loginInfo.RefreshToken,
		Exp:          loginInfo.Exp,
	}, nil
}

func (s *PassportServer) RefreshToken(ctx context.Context, req *passportv1.RefreshTokenRequest) (*passportv1.RefreshTokenResponse, error) {
	result, err := s.authService.VerifyAndRefreshToken(ctx, req.GetAccessToken(), req.GetRefreshToken())
	if err != nil {
		return nil, err
	}
	return &passportv1.RefreshTokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}, nil
}

func (s *PassportServer) RegisterUser(ctx context.Context, req *passportv1.RegisterUserRequest) (*passportv1.RegisterUserResponse, error) {
	if err := validateRegisterRequest(req); err != nil {
		return nil, err
	}

	phone := strings.TrimSpace(req.GetPhone())
	um := &authdomain.Account{
		Username: req.GetUsername(),
		Nickname: req.GetNickname(),
		Email:    req.GetEmail(),
		Phone: types.NullString{
			NullString: sql.NullString{
				String: phone,
				Valid:  phone != "",
			},
		},
		Avatar:   req.GetAvatar(),
		Gender:   int(req.GetGender()),
		Password: req.GetPassword(),
	}
	if err := s.registrationService.RegisterUser(ctx, um, req.GetVerificationCode()); err != nil {
		return nil, err
	}
	return &passportv1.RegisterUserResponse{}, nil
}

func (s *PassportServer) SendVerificationCode(ctx context.Context, req *passportv1.SendVerificationCodeRequest) (*passportv1.SendVerificationCodeResponse, error) {
	if _, err := mail.ParseAddress(req.GetEmail()); err != nil {
		return nil, apperr.BadRequest("invalid email")
	}
	if strings.TrimSpace(req.GetTurnstileToken()) == "" {
		return nil, apperr.BadRequest("turnstile_token is required")
	}
	if err := s.registrationService.SendVerificationCode(ctx, req.GetEmail(), req.GetTurnstileToken()); err != nil {
		return nil, err
	}
	return &passportv1.SendVerificationCodeResponse{}, nil
}

func (s *PassportServer) ForgotPassword(ctx context.Context, req *passportv1.ForgotPasswordRequest) (*passportv1.ForgotPasswordResponse, error) {
	if _, err := mail.ParseAddress(req.GetEmail()); err != nil {
		return nil, apperr.BadRequest("invalid email")
	}
	if strings.TrimSpace(req.GetTurnstileToken()) == "" {
		return nil, apperr.BadRequest("turnstile_token is required")
	}
	if err := s.passwordReset.SendPasswordResetEmail(ctx, req.GetEmail(), req.GetTurnstileToken()); err != nil {
		return nil, err
	}
	return &passportv1.ForgotPasswordResponse{}, nil
}

func (s *PassportServer) ResetPassword(ctx context.Context, req *passportv1.ResetPasswordRequest) (*passportv1.ResetPasswordResponse, error) {
	if strings.TrimSpace(req.GetToken()) == "" {
		return nil, apperr.BadRequest("token is required")
	}
	if len(req.GetPassword()) < 6 || len(req.GetPassword()) > 50 {
		return nil, apperr.BadRequest("password must be between 6 and 50 characters")
	}
	if err := s.passwordReset.ResetPassword(ctx, req.GetToken(), req.GetPassword()); err != nil {
		return nil, err
	}
	return &passportv1.ResetPasswordResponse{}, nil
}

func (s *PassportServer) ListMySessions(ctx context.Context, _ *passportv1.ListMySessionsRequest) (*passportv1.ListMySessionsResponse, error) {
	tokenUser := authctx.TokenUserFromContext(ctx)
	if tokenUser == nil {
		return nil, apperr.Unauthorized("")
	}

	sessions, err := s.sessionService.ListUserSessions(ctx, tokenUser.UID)
	if err != nil {
		return nil, err
	}
	return &passportv1.ListMySessionsResponse{Sessions: ToUserSessions(sessions)}, nil
}

func (s *PassportServer) RevokeMySession(ctx context.Context, req *passportv1.RevokeMySessionRequest) (*passportv1.RevokeMySessionResponse, error) {
	tokenUser := authctx.TokenUserFromContext(ctx)
	if tokenUser == nil {
		return nil, apperr.Unauthorized("")
	}

	sessionID := strings.TrimSpace(req.GetSessionId())
	if sessionID == "" {
		return nil, apperr.BadRequest("session_id is required")
	}
	if err := s.sessionService.RevokeUserSession(ctx, tokenUser.UID, sessionID); err != nil {
		return nil, err
	}
	return &passportv1.RevokeMySessionResponse{}, nil
}

func validateRegisterRequest(req *passportv1.RegisterUserRequest) error {
	if !usernamePattern.MatchString(req.GetUsername()) || len(req.GetUsername()) < 3 || len(req.GetUsername()) > 20 {
		return apperr.BadRequest("username must be alphanumeric and between 3 and 20 characters")
	}
	if strings.TrimSpace(req.GetNickname()) == "" {
		return apperr.BadRequest("nickname is required")
	}
	if _, err := mail.ParseAddress(req.GetEmail()); err != nil {
		return apperr.BadRequest("invalid email")
	}
	if len(req.GetPassword()) < 6 || len(req.GetPassword()) > 50 {
		return apperr.BadRequest("password must be between 6 and 50 characters")
	}
	if strings.TrimSpace(req.GetVerificationCode()) == "" {
		return apperr.BadRequest("verification_code is required")
	}
	return nil
}

func extractLoginContext(ctx context.Context) logic.LoginContext {
	loginCtx := logic.LoginContext{}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return loginCtx
	}
	if values := md.Get("x-forwarded-for"); len(values) > 0 {
		loginCtx.LoginIP = strings.TrimSpace(values[0])
	}
	if values := md.Get("x-device-info"); len(values) > 0 {
		loginCtx.DeviceInfo = strings.TrimSpace(values[0])
	}
	if loginCtx.DeviceInfo == "" {
		if values := md.Get("grpcgateway-user-agent"); len(values) > 0 {
			loginCtx.DeviceInfo = strings.TrimSpace(values[0])
		}
	}
	return loginCtx
}
