package repository

import (
	"hello/config"
	"hello/model"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	var err error

	// 确保数据目录存在
	dbDir := filepath.Dir(config.AppConfig.DBPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return err
	}

	DB, err = gorm.Open(sqlite.Open(config.AppConfig.DBPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	// 自动迁移
	err = DB.AutoMigrate(&model.User{}, &model.Role{}, &model.Permission{})
	if err != nil {
		return err
	}

	return nil
}

// InitData 初始化基础数据（超级管理员）
func InitData() error {
	// 检查是否已存在admin用户
	var count int64
	DB.Model(&model.User{}).Where("username = ?", "admin").Count(&count)
	if count > 0 {
		return nil
	}

	// 创建初始权限（树形结构）
	permissions := []model.Permission{
		// 目录
		{Name: "系统管理", Code: "system", Type: model.PermissionTypeDirectory, Icon: "Setting", Sort: 1},

		// 菜单
		{Name: "用户管理", Code: "system:user", Type: model.PermissionTypeMenu, ParentID: 1, Path: "/users", Icon: "User", Component: "user/List", Sort: 1},
		{Name: "角色管理", Code: "system:role", Type: model.PermissionTypeMenu, ParentID: 1, Path: "/roles", Icon: "Avatar", Component: "role/List", Sort: 2},
		{Name: "权限管理", Code: "system:permission", Type: model.PermissionTypeMenu, ParentID: 1, Path: "/permissions", Icon: "Key", Component: "permission/List", Sort: 3},

		// 用户操作
		{Name: "查看用户", Code: "system:user:read", Type: model.PermissionTypeOperation, ParentID: 2, ApiPath: "/api/users", ApiMethod: "GET", Sort: 1},
		{Name: "创建用户", Code: "system:user:create", Type: model.PermissionTypeOperation, ParentID: 2, ApiPath: "/api/users", ApiMethod: "POST", Sort: 2},
		{Name: "更新用户", Code: "system:user:update", Type: model.PermissionTypeOperation, ParentID: 2, ApiPath: "/api/users/:id", ApiMethod: "PUT", Sort: 3},
		{Name: "删除用户", Code: "system:user:delete", Type: model.PermissionTypeOperation, ParentID: 2, ApiPath: "/api/users/:id", ApiMethod: "DELETE", Sort: 4},

		// 角色操作
		{Name: "查看角色", Code: "system:role:read", Type: model.PermissionTypeOperation, ParentID: 3, ApiPath: "/api/roles", ApiMethod: "GET", Sort: 1},
		{Name: "创建角色", Code: "system:role:create", Type: model.PermissionTypeOperation, ParentID: 3, ApiPath: "/api/roles", ApiMethod: "POST", Sort: 2},
		{Name: "更新角色", Code: "system:role:update", Type: model.PermissionTypeOperation, ParentID: 3, ApiPath: "/api/roles/:id", ApiMethod: "PUT", Sort: 3},
		{Name: "删除角色", Code: "system:role:delete", Type: model.PermissionTypeOperation, ParentID: 3, ApiPath: "/api/roles/:id", ApiMethod: "DELETE", Sort: 4},

		// 权限操作
		{Name: "查看权限", Code: "system:permission:read", Type: model.PermissionTypeOperation, ParentID: 4, ApiPath: "/api/permissions", ApiMethod: "GET", Sort: 1},
		{Name: "创建权限", Code: "system:permission:create", Type: model.PermissionTypeOperation, ParentID: 4, ApiPath: "/api/permissions", ApiMethod: "POST", Sort: 2},
		{Name: "更新权限", Code: "system:permission:update", Type: model.PermissionTypeOperation, ParentID: 4, ApiPath: "/api/permissions/:id", ApiMethod: "PUT", Sort: 3},
		{Name: "删除权限", Code: "system:permission:delete", Type: model.PermissionTypeOperation, ParentID: 4, ApiPath: "/api/permissions/:id", ApiMethod: "DELETE", Sort: 4},
	}

	for i := range permissions {
		DB.FirstOrCreate(&permissions[i], model.Permission{Code: permissions[i].Code})
	}

	// 创建超级管理员角色（包含所有权限）
	adminRole := model.Role{
		Name:        "超级管理员",
		Code:        "admin",
		Permissions: permissions,
	}
	if err := DB.Create(&adminRole).Error; err != nil {
		return err
	}

	// 创建超级管理员用户
	adminUser := model.User{
		Username: "admin",
		Password: "admin123",
		Email:    "admin@example.com",
		Status:   1,
		Roles:    []model.Role{adminRole},
	}
	if err := DB.Create(&adminUser).Error; err != nil {
		return err
	}

	return nil
}
