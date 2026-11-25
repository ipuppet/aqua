package grpcserver

import (
	"context"

	"google.golang.org/protobuf/types/known/structpb"

	userv1 "github.com/ipuppet/aqua/server/gen/go/aqua/user/v1"
	"github.com/ipuppet/aqua/server/pkg/apperr"
	"github.com/ipuppet/aqua/server/pkg/transport/authctx"
	user "github.com/ipuppet/aqua/server/services/user/domain"
	"github.com/ipuppet/aqua/server/services/user/logic"
)

type UserServer struct {
	userService *logic.UserService
	userv1.UnimplementedUserServiceServer
}

func NewUserServer(userService *logic.UserService) *UserServer {
	return &UserServer{userService: userService}
}

func (s *UserServer) GetUserInfo(ctx context.Context, req *userv1.GetUserInfoRequest) (*userv1.User, error) {
	info, err := s.userService.GetUserInfo(ctx, int(req.Uid))
	if err != nil {
		return nil, apperr.FromDAO(err)
	}
	return ToPBUser(*info), nil
}

func (s *UserServer) GetSelf(ctx context.Context, req *userv1.GetSelfRequest) (*userv1.User, error) {
	tokenUser := authctx.TokenUserFromContext(ctx)
	if tokenUser == nil {
		return nil, apperr.Unauthorized("")
	}
	info, err := s.userService.GetUserInfo(ctx, tokenUser.UID)
	if err != nil {
		return nil, apperr.FromDAO(err)
	}
	return ToPBUser(*info), nil
}

func (s *UserServer) UpdateSelf(ctx context.Context, req *userv1.UpdateSelfRequest) (*userv1.UpdateSelfResponse, error) {
	tokenUser := authctx.TokenUserFromContext(ctx)
	if tokenUser == nil {
		return nil, apperr.Unauthorized("")
	}
	gender := int(req.Gender)
	updates := &user.UserUpdate{
		Nickname: &req.Nickname,
		Email:    &req.Email,
		Phone:    &req.Phone,
		Avatar:   &req.Avatar,
		Gender:   &gender,
	}
	err := s.userService.UpdateSelf(ctx, tokenUser.UID, updates)
	return &userv1.UpdateSelfResponse{}, err
}

func (s *UserServer) GetMenu(ctx context.Context, req *userv1.GetMenuRequest) (*userv1.GetMenuResponse, error) {
	tokenUser := authctx.TokenUserFromContext(ctx)
	if tokenUser == nil {
		return nil, apperr.Unauthorized("")
	}
	menus, err := s.userService.GetUserMenu(ctx, tokenUser.UID)
	if err != nil {
		return nil, err
	}

	var pbMenus []*structpb.Struct
	for _, m := range menus {
		str, err := structpb.NewStruct(m)
		if err == nil {
			pbMenus = append(pbMenus, str)
		}
	}
	return &userv1.GetMenuResponse{Menus: pbMenus}, nil
}

func (s *UserServer) GetStructure(ctx context.Context, req *userv1.GetStructureRequest) (*userv1.GetStructureResponse, error) {
	tokenUser := authctx.TokenUserFromContext(ctx)
	if tokenUser == nil {
		return nil, apperr.Unauthorized("")
	}
	strMap, err := s.userService.GetUserStructure(ctx)
	if err != nil {
		return nil, err
	}
	str, _ := structpb.NewStruct(strMap)
	return &userv1.GetStructureResponse{Structure: str}, nil
}
