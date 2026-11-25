package logic

import (
	"context"
	"log/slog"

	"github.com/ipuppet/aqua/server/pkg/infra/logutil"
	"github.com/ipuppet/aqua/server/pkg/rbac"
	"github.com/ipuppet/aqua/server/pkg/transport/token"
	"github.com/ipuppet/aqua/server/services/user/dao"
	user "github.com/ipuppet/aqua/server/services/user/domain"
)

type UserManagementService struct {
	repo dao.UserManagementRepository
}

// NewUserManagementService constructs UserManagementService.
func NewUserManagementService(repo dao.UserManagementRepository) *UserManagementService {
	return &UserManagementService{repo: repo}
}

func (s *UserManagementService) log(ctx context.Context) *slog.Logger {
	return logutil.FromContext(ctx)
}

func (s *UserManagementService) GetUsers(ctx context.Context, page int, pageSize int) ([]user.UserModelResponse, error) {
	offset, limit := normalizePagination(page, pageSize)
	users, err := s.repo.GetUsers(ctx, offset, limit)
	if err != nil {
		s.log(ctx).Error("failed to get users",
			slog.Int("page", page),
			slog.Int("page_size", pageSize),
			slog.Any("error", err))
		return nil, err
	}
	return ToUserModelResponses(users), nil
}

func (s *UserManagementService) SearchUsers(ctx context.Context, kw string, page int, pageSize int) ([]user.UserModelResponse, error) {
	offset, limit := normalizePagination(page, pageSize)
	users, err := s.repo.SearchUsers(ctx, kw, offset, limit)
	if err != nil {
		s.log(ctx).Error("failed to search users",
			slog.String("search.keyword", kw),
			slog.Int("page", page),
			slog.Int("page_size", pageSize),
			slog.Any("error", err))
		return nil, err
	}
	return ToUserModelResponses(users), nil
}

func (s *UserManagementService) UpdateUser(ctx context.Context, uid int, updates *user.UserUpdate) error {
	if err := s.repo.UpdateUser(ctx, uid, updates); err != nil {
		s.log(ctx).Error("failed to update user",
			slog.Int("user.uid", uid),
			slog.Any("error", err))
		return err
	}
	_ = s.repo.Cache().Del(ctx, rbac.CacheKeyUserInfo(uid))
	if err := token.RevokeUserRefreshTokens(ctx, s.repo.Cache(), uid); err != nil {
		s.log(ctx).Warn("failed to revoke user refresh tokens after profile update",
			slog.Int("user.uid", uid),
			slog.Any("error", err))
	}
	return nil
}
