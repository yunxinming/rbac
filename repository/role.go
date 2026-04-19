package repository

import (
	"hello/model"

	"gorm.io/gorm"
)

type RoleRepository struct{}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{}
}

func (r *RoleRepository) Create(role *model.Role) error {
	return DB.Create(role).Error
}

func (r *RoleRepository) Update(role *model.Role) error {
	return DB.Save(role).Error
}

func (r *RoleRepository) Delete(id uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 删除角色-权限关联
		if err := tx.Exec("DELETE FROM role_permissions WHERE role_id = ?", id).Error; err != nil {
			return err
		}
		// 删除用户-角色关联
		if err := tx.Exec("DELETE FROM user_roles WHERE role_id = ?", id).Error; err != nil {
			return err
		}
		// 删除角色
		return tx.Delete(&model.Role{}, id).Error
	})
}

func (r *RoleRepository) GetByID(id uint) (*model.Role, error) {
	var role model.Role
	err := DB.Preload("Permissions").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) GetAll() ([]model.Role, error) {
	var roles []model.Role
	err := DB.Preload("Permissions").Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) GetByCodes(codes []string) ([]model.Role, error) {
	var roles []model.Role
	err := DB.Where("code IN ?", codes).Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) AssignPermissions(roleID uint, permissionIDs []uint) error {
	// 先删除旧的关联
	if err := DB.Exec("DELETE FROM role_permissions WHERE role_id = ?", roleID).Error; err != nil {
		return err
	}
	// 创建新的关联
	for _, permID := range permissionIDs {
		if err := DB.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", roleID, permID).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *RoleRepository) ExistsByCode(code string) (bool, error) {
	var count int64
	err := DB.Model(&model.Role{}).Where("code = ?", code).Count(&count).Error
	return count > 0, err
}
