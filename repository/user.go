package repository

import (
	"hello/model"

	"gorm.io/gorm"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(user *model.User) error {
	return DB.Create(user).Error
}

func (r *UserRepository) Update(user *model.User) error {
	return DB.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 删除用户-角色关联
		if err := tx.Exec("DELETE FROM user_roles WHERE user_id = ?", id).Error; err != nil {
			return err
		}
		// 删除用户
		return tx.Delete(&model.User{}, id).Error
	})
}

func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := DB.Preload("Roles.Permissions").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := DB.Preload("Roles.Permissions").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetAll() ([]model.User, error) {
	var users []model.User
	err := DB.Preload("Roles").Find(&users).Error
	return users, err
}

func (r *UserRepository) GetByIDs(ids []uint) ([]model.User, error) {
	var users []model.User
	err := DB.Find(&users, ids).Error
	return users, err
}

func (r *UserRepository) AssignRoles(userID uint, roleIDs []uint) error {
	// 先删除旧的关联
	if err := DB.Exec("DELETE FROM user_roles WHERE user_id = ?", userID).Error; err != nil {
		return err
	}
	// 创建新的关联
	for _, roleID := range roleIDs {
		if err := DB.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", userID, roleID).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := DB.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := DB.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// GetUserPermissionCodes 获取用户的所有权限码
func (r *UserRepository) GetUserPermissionCodes(userID uint) ([]string, error) {
	var codes []string
	err := DB.Table("permissions").
		Select("DISTINCT permissions.code").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Pluck("permissions.code", &codes).Error
	return codes, err
}
