package service

import (
	"errors"
	"hello/model"
	"hello/repository"
)

type RoleService struct {
	repo            *repository.RoleRepository
	permissionRepo  *repository.PermissionRepository
}

func NewRoleService() *RoleService {
	return &RoleService{
		repo:           repository.NewRoleRepository(),
		permissionRepo: repository.NewPermissionRepository(),
	}
}

type CreateRoleRequest struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
}

type UpdateRoleRequest struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
}

type AssignPermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids"`
}

func (s *RoleService) Create(req *CreateRoleRequest) (*model.Role, error) {
	exists, err := s.repo.ExistsByCode(req.Code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("角色编码已存在")
	}

	role := &model.Role{
		Name: req.Name,
		Code: req.Code,
	}
	if err := s.repo.Create(role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *RoleService) Update(id uint, req *UpdateRoleRequest) (*model.Role, error) {
	role, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("角色不存在")
	}

	if role.Code != req.Code {
		exists, err := s.repo.ExistsByCode(req.Code)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("角色编码已存在")
		}
	}

	role.Name = req.Name
	role.Code = req.Code
	if err := s.repo.Update(role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *RoleService) Delete(id uint) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("角色不存在")
	}
	return s.repo.Delete(id)
}

func (s *RoleService) GetByID(id uint) (*model.Role, error) {
	return s.repo.GetByID(id)
}

func (s *RoleService) GetAll() ([]model.Role, error) {
	return s.repo.GetAll()
}

func (s *RoleService) AssignPermissions(id uint, req *AssignPermissionsRequest) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("角色不存在")
	}
	return s.repo.AssignPermissions(id, req.PermissionIDs)
}
