package logic

import (
	"context"
	"errors"
	"log/slog"
	"slices"

	"aqua/pkg/infra/filestore"
	"aqua/pkg/infra/logutil"
	"aqua/pkg/rbac"
	"aqua/services/user/dao"
	user "aqua/services/user/domain"
)

// UserService is the logic implementation.
type UserService struct {
	repo       dao.UserRepository
	rbacReader rbac.Reader
	store      *filestore.Client
}

// NewUserService constructs a UserService backed by the provided repository and RBAC reader.
func NewUserService(repo dao.UserRepository, reader rbac.Reader, store *filestore.Client) *UserService {
	return &UserService{repo: repo, rbacReader: reader, store: store}
}

func (s *UserService) log(ctx context.Context) *slog.Logger {
	return logutil.FromContext(ctx)
}

func (s *UserService) GetUserInfo(ctx context.Context, uid int) (*user.UserModelResponse, error) {
	info, err := s.repo.GetUserInfo(ctx, uid)
	if err != nil {
		s.log(ctx).Error("failed to get user info", "uid", uid, "error", err)
		return nil, err
	}
	return ToUserModelResponse(info), nil
}

func (s *UserService) UpdateSelf(ctx context.Context, uid int, updates *user.UserUpdate) error {
	if uid <= 0 {
		return errors.New("invalid uid")
	}

	if err := s.repo.UpdateUser(ctx, uid, updates); err != nil {
		s.log(ctx).Error("failed to update self", "uid", uid, "error", err)
		return err
	}
	// Invalidate user cache so subsequent token / get info calls return fresh data
	if err := s.repo.Cache().Del(ctx, rbac.CacheKeyUserInfo(uid)); err != nil {
		s.log(ctx).Warn("failed to invalidate user cache", "uid", uid, "error", err)
	}
	return nil
}

func (s *UserService) GetUserMenu(ctx context.Context, uid int) ([]map[string]any, error) {
	menus := []map[string]any{}
	if err := s.store.GetJSON("user", "menu.json", &menus); err != nil {
		s.log(ctx).Error("failed to load menu from file", "error", err)
		return nil, err
	}

	permissionsModel, err := s.rbacReader.LoadUserPermissions(ctx, uid)
	if err != nil {
		s.log(ctx).Error("failed to load user permissions", "uid", uid, "error", err)
		return nil, err
	}

	permissions := make([]string, 0, len(permissionsModel))
	for _, permission := range permissionsModel {
		permissions = append(permissions, permission.Code)
	}

	// 检查是否有全部权限
	if slices.Contains(permissions, "ALL") {
		return menus, nil
	}

	// 根据权限过滤菜单（即使没有权限，也应显示 module=public 的项）
	result := []map[string]any{}
	for _, menu := range menus {
		// 先处理不需要权限的公共菜单
		if permission, ok := menu["permission"].(string); ok && permission == "public" {
			result = append(result, menu)
			continue
		}
		for _, permission := range permissions {
			if menu["permission"] == permission {
				result = append(result, menu)
				break // 避免重复添加
			}
		}
	}
	return result, nil
}

func (s *UserService) GetUserStructure(ctx context.Context) (map[string]any, error) {
	data := map[string]any{}
	if err := s.store.GetJSON("user", "structure.json", &data); err != nil {
		s.log(ctx).Error("failed to load user structure from file", "error", err)
		return nil, err
	}
	return data, nil
}
