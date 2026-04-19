package service

import (
	"errors"
	"hello/model"
	"hello/repository"
)

type PermissionService struct {
	repo *repository.PermissionRepository
}

func NewPermissionService() *PermissionService {
	return &PermissionService{
		repo: repository.NewPermissionRepository(),
	}
}

type CreatePermissionRequest struct {
	Name      string `json:"name" binding:"required"`
	Code      string `json:"code" binding:"required"`
	Type      string `json:"type" binding:"required,oneof=directory menu operation"`
	ParentID  uint   `json:"parent_id"`
	Path      string `json:"path"`
	Icon      string `json:"icon"`
	Sort      int    `json:"sort"`
	ApiPath   string `json:"api_path"`
	ApiMethod string `json:"api_method"`
	Component string `json:"component"`
	Status    *int8  `json:"status"`
}

type UpdatePermissionRequest struct {
	Name      string `json:"name" binding:"required"`
	Code      string `json:"code" binding:"required"`
	Type      string `json:"type" binding:"required,oneof=directory menu operation"`
	ParentID  uint   `json:"parent_id"`
	Path      string `json:"path"`
	Icon      string `json:"icon"`
	Sort      int    `json:"sort"`
	ApiPath   string `json:"api_path"`
	ApiMethod string `json:"api_method"`
	Component string `json:"component"`
	Status    *int8  `json:"status"`
}

func (s *PermissionService) Create(req *CreatePermissionRequest) (*model.Permission, error) {
	exists, err := s.repo.ExistsByCode(req.Code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("权限编码已存在")
	}

	status := int8(1)
	if req.Status != nil {
		status = *req.Status
	}

	permission := &model.Permission{
		Name:      req.Name,
		Code:      req.Code,
		Type:      model.PermissionType(req.Type),
		ParentID:  req.ParentID,
		Path:      req.Path,
		Icon:      req.Icon,
		Sort:      req.Sort,
		ApiPath:   req.ApiPath,
		ApiMethod: req.ApiMethod,
		Component: req.Component,
		Status:    status,
	}
	if err := s.repo.Create(permission); err != nil {
		return nil, err
	}
	return permission, nil
}

func (s *PermissionService) Update(id uint, req *UpdatePermissionRequest) (*model.Permission, error) {
	permission, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("权限不存在")
	}

	if permission.Code != req.Code {
		exists, err := s.repo.ExistsByCode(req.Code)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("权限编码已存在")
		}
	}

	permission.Name = req.Name
	permission.Code = req.Code
	permission.Type = model.PermissionType(req.Type)
	permission.ParentID = req.ParentID
	permission.Path = req.Path
	permission.Icon = req.Icon
	permission.Sort = req.Sort
	permission.ApiPath = req.ApiPath
	permission.ApiMethod = req.ApiMethod
	permission.Component = req.Component
	if req.Status != nil {
		permission.Status = *req.Status
	}

	if err := s.repo.Update(permission); err != nil {
		return nil, err
	}
	return permission, nil
}

func (s *PermissionService) Delete(id uint) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("权限不存在")
	}
	return s.repo.Delete(id)
}

func (s *PermissionService) GetByID(id uint) (*model.Permission, error) {
	return s.repo.GetByID(id)
}

func (s *PermissionService) GetAll() ([]model.Permission, error) {
	return s.repo.GetAll()
}

// GetTree 获取权限树
func (s *PermissionService) GetTree() ([]model.Permission, error) {
	return s.repo.GetAllWithChildren()
}
