package service

import (
	"errors"
	"hello/model"
	"hello/repository"
)

type UserService struct {
	repo     *repository.UserRepository
	roleRepo *repository.RoleRepository
}

func NewUserService() *UserService {
	return &UserService{
		repo:     repository.NewUserRepository(),
		roleRepo: repository.NewRoleRepository(),
	}
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	Status   *int8  `json:"status"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password"` // 可选，不传则不修改
	Status   *int8  `json:"status"`
}

type AssignRolesRequest struct {
	RoleIDs []uint `json:"role_ids"`
}

func (s *UserService) Create(req *CreateUserRequest) (*model.User, error) {
	exists, err := s.repo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("用户名已存在")
	}

	exists, err = s.repo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("邮箱已存在")
	}

	status := int8(1)
	if req.Status != nil {
		status = *req.Status
	}

	user := &model.User{
		Username: req.Username,
		Password: req.Password, // 密码会在BeforeCreate钩子中自动哈希
		Email:    req.Email,
		Status:   status,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Update(id uint, req *UpdateUserRequest) (*model.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 检查用户名是否被其他用户使用
	if user.Username != req.Username {
		exists, err := s.repo.ExistsByUsername(req.Username)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("用户名已存在")
		}
	}

	// 检查邮箱是否被其他用户使用
	if user.Email != req.Email {
		exists, err := s.repo.ExistsByEmail(req.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("邮箱已存在")
		}
	}

	user.Username = req.Username
	user.Email = req.Email
	if req.Password != "" {
		user.Password = req.Password // 会在BeforeUpdate钩子中自动哈希
	}
	if req.Status != nil {
		user.Status = *req.Status
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Delete(id uint) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}
	return s.repo.Delete(id)
}

func (s *UserService) GetByID(id uint) (*model.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) GetAll() ([]model.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) AssignRoles(id uint, req *AssignRolesRequest) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}
	return s.repo.AssignRoles(id, req.RoleIDs)
}

func (s *UserService) GetPermissionCodes(userID uint) ([]string, error) {
	return s.repo.GetUserPermissionCodes(userID)
}
