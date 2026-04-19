package service

import (
	"errors"
	"hello/model"
	"hello/pkg/jwt"
	"hello/repository"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(),
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string             `json:"token"`
	User  *model.User        `json:"user"`
	Menus []model.Permission `json:"menus"`
}

func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if !user.CheckPassword(req.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	// 提取角色编码
	var roleCodes []string
	for _, role := range user.Roles {
		roleCodes = append(roleCodes, role.Code)
	}

	// 生成Token
	token, err := jwt.GenerateToken(user.ID, user.Username, roleCodes)
	if err != nil {
		return nil, errors.New("生成Token失败")
	}

	// 获取用户菜单
	permRepo := repository.NewPermissionRepository()
	menus, err := permRepo.GetUserMenus(user.ID)
	if err != nil {
		menus = []model.Permission{}
	}

	// 清空密码
	user.Password = ""

	return &LoginResponse{
		Token: token,
		User:  user,
		Menus: menus,
	}, nil
}

func (s *AuthService) Logout() error {
	// JWT无状态，客户端删除Token即可
	// 如需服务端失效，可在此实现Token黑名单
	return nil
}
