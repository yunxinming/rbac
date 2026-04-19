# RBAC权限管理系统实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 构建一个独立的RBAC权限管理API服务，提供用户、角色、权限的管理和JWT认证功能。

**Architecture:** 采用分层架构：Handler -> Service -> Repository -> Model。Gin中间件处理JWT认证和权限校验，GORM管理SQLite数据库，首次启动自动初始化超级管理员。

**Tech Stack:** Gin, JWT (golang-jwt/jwt/v5), SQLite, GORM, bcrypt

---

## 文件结构

```
hello/
├── main.go                      # 入口，路由注册，启动服务
├── config/
│   └── config.go                # 配置管理（环境变量读取）
├── model/
│   ├── user.go                  # User模型 + GORM钩子
│   ├── role.go                  # Role模型
│   └── permission.go            # Permission模型
├── handler/
│   ├── auth.go                  # 登录/登出处理器
│   ├── user.go                  # 用户CRUD处理器
│   ├── role.go                  # 角色CRUD处理器
│   └── permission.go            # 权限CRUD处理器
├── middleware/
│   ├── jwt.go                   # JWT认证中间件
│   └── permission.go            # 权限校验中间件
├── service/
│   ├── auth.go                  # 登录业务逻辑
│   ├── user.go                  # 用户业务逻辑
│   ├── role.go                  # 角色业务逻辑
│   └── permission.go            # 权限业务逻辑
├── repository/
│   ├── user.go                  # 用户数据访问
│   ├── role.go                  # 角色数据访问
│   └── permission.go            # 权限数据访问
├── pkg/
│   ├── jwt/
│   │   └── jwt.go               # JWT生成/解析工具
│   └── response/
│       └── response.go          # 统一JSON响应工具
├── data/
│   └── rbac.db                  # SQLite数据库文件（自动生成）
└── go.mod
```

---

## Task 1: 项目初始化与依赖安装

**Files:**
- Modify: `go.mod`
- Create: `go.sum` (自动生成)

- [ ] **Step 1: 安装依赖**

```bash
go get github.com/gin-gonic/gin@latest
go get github.com/golang-jwt/jwt/v5@latest
go get gorm.io/gorm@latest
go get gorm.io/driver/sqlite@latest
go get golang.org/x/crypto@latest
```

- [ ] **Step 2: 验证依赖安装**

```bash
go mod tidy
cat go.mod
```

Expected: go.mod 包含 gin, jwt/v5, gorm, sqlite, crypto 依赖

---

## Task 2: 配置管理模块

**Files:**
- Create: `config/config.go`

- [ ] **Step 1: 创建配置模块**

```go
package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port      string
	JWTSecret string
	DBPath    string
}

var AppConfig *Config

func LoadConfig() {
	port := getEnv("PORT", "8080")
	jwtSecret := getEnv("JWT_SECRET", "default-secret-change-in-production")
	dbPath := getEnv("DB_PATH", "./data/rbac.db")

	AppConfig = &Config{
		Port:      port,
		JWTSecret: jwtSecret,
		DBPath:    dbPath,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
```

- [ ] **Step 2: 验证配置模块编译**

```bash
go build ./config/...
```

Expected: 无错误输出

---

## Task 3: 数据模型定义

**Files:**
- Create: `model/user.go`
- Create: `model/role.go`
- Create: `model/permission.go`

- [ ] **Step 1: 创建 Permission 模型**

```go
package model

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:64;not null" json:"name"`
	Code      string         `gorm:"size:64;uniqueIndex;not null" json:"code"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Permission) TableName() string {
	return "permissions"
}
```

- [ ] **Step 2: 创建 Role 模型**

```go
package model

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:64;not null" json:"name"`
	Code        string         `gorm:"size:64;uniqueIndex;not null" json:"code"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Permissions []Permission   `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

func (Role) TableName() string {
	return "roles"
}
```

- [ ] **Step 3: 创建 User 模型**

```go
package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"size:128;not null" json:"-"`
	Email     string         `gorm:"size:128;uniqueIndex;not null" json:"email"`
	Status    int8           `gorm:"default:1" json:"status"` // 1:启用 0:禁用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Roles     []Role         `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

func (User) TableName() string {
	return "users"
}

// SetPassword 设置密码（bcrypt哈希）
func (u *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// BeforeCreate GORM钩子：创建前检查密码是否已哈希
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password != "" && len(u.Password) < 60 {
		u.SetPassword(u.Password)
	}
	return nil
}

// BeforeUpdate GORM钩子：更新前检查密码是否变更
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("Password") && len(u.Password) < 60 {
		u.SetPassword(u.Password)
	}
	return nil
}
```

- [ ] **Step 4: 验证模型编译**

```bash
go build ./model/...
```

Expected: 无错误输出

---

## Task 4: 数据库初始化与自动迁移

**Files:**
- Create: `data/` 目录
- Create: `repository/db.go`

- [ ] **Step 1: 创建数据库初始化模块**

```go
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

	// 创建初始权限
	permissions := []model.Permission{
		{Name: "查看用户", Code: "user:read"},
		{Name: "创建用户", Code: "user:create"},
		{Name: "更新用户", Code: "user:update"},
		{Name: "删除用户", Code: "user:delete"},
		{Name: "查看角色", Code: "role:read"},
		{Name: "创建角色", Code: "role:create"},
		{Name: "更新角色", Code: "role:update"},
		{Name: "删除角色", Code: "role:delete"},
		{Name: "查看权限", Code: "permission:read"},
		{Name: "创建权限", Code: "permission:create"},
		{Name: "更新权限", Code: "permission:update"},
		{Name: "删除权限", Code: "permission:delete"},
	}
	for i := range permissions {
		DB.FirstOrCreate(&permissions[i], model.Permission{Code: permissions[i].Code})
	}

	// 创建超级管理员角色
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
```

- [ ] **Step 2: 验证编译**

```bash
go build ./repository/...
```

Expected: 无错误输出

---

## Task 5: 统一响应工具

**Files:**
- Create: `pkg/response/response.go`

- [ ] **Step 1: 创建响应工具**

```go
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 响应码定义
const (
	CodeSuccess         = 0
	CodeBadRequest      = 40000
	CodeUnauthorized    = 40001
	CodeForbidden       = 40003
	CodeNotFound        = 40004
	CodeInternalError   = 50000
	CodeInvalidPassword = 41001
	CodeUserDisabled    = 41002
	CodeUserExists      = 41003
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// BadRequest 参数错误
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, CodeBadRequest, message)
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, CodeUnauthorized, message)
}

// Forbidden 禁止访问
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, CodeForbidden, message)
}

// NotFound 资源不存在
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, CodeNotFound, message)
}

// InternalError 内部错误
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, CodeInternalError, message)
}
```

- [ ] **Step 2: 验证编译**

```bash
go build ./pkg/response/...
```

Expected: 无错误输出

---

## Task 6: JWT工具

**Files:**
- Create: `pkg/jwt/jwt.go`

- [ ] **Step 1: 创建JWT工具**

```go
package jwt

import (
	"errors"
	"hello/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint     `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

var (
	ErrTokenExpired     = errors.New("token已过期")
	ErrTokenInvalid     = errors.New("token无效")
	ErrTokenMalformed   = errors.New("token格式错误")
	ErrTokenNotValidYet = errors.New("token尚未生效")
)

// GenerateToken 生成JWT Token
func GenerateToken(userID uint, username string, roles []string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "rbac-system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWTSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotValidYet
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}
```

- [ ] **Step 2: 验证编译**

```bash
go build ./pkg/jwt/...
```

Expected: 无错误输出

---

## Task 7: Repository层 - 数据访问

**Files:**
- Create: `repository/user.go`
- Create: `repository/role.go`
- Create: `repository/permission.go`

- [ ] **Step 1: 创建 Permission Repository**

```go
package repository

import (
	"hello/model"

	"gorm.io/gorm"
)

type PermissionRepository struct{}

func NewPermissionRepository() *PermissionRepository {
	return &PermissionRepository{}
}

func (r *PermissionRepository) Create(permission *model.Permission) error {
	return DB.Create(permission).Error
}

func (r *PermissionRepository) Update(permission *model.Permission) error {
	return DB.Save(permission).Error
}

func (r *PermissionRepository) Delete(id uint) error {
	return DB.Delete(&model.Permission{}, id).Error
}

func (r *PermissionRepository) GetByID(id uint) (*model.Permission, error) {
	var permission model.Permission
	err := DB.First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionRepository) GetAll() ([]model.Permission, error) {
	var permissions []model.Permission
	err := DB.Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepository) GetByCodes(codes []string) ([]model.Permission, error) {
	var permissions []model.Permission
	err := DB.Where("code IN ?", codes).Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepository) ExistsByCode(code string) (bool, error) {
	var count int64
	err := DB.Model(&model.Permission{}).Where("code = ?", code).Count(&count).Error
	return count > 0, err
}
```

- [ ] **Step 2: 创建 Role Repository**

```go
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
```

- [ ] **Step 3: 创建 User Repository**

```go
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
```

- [ ] **Step 4: 验证编译**

```bash
go build ./repository/...
```

Expected: 无错误输出

---

## Task 8: Service层 - 业务逻辑

**Files:**
- Create: `service/auth.go`
- Create: `service/user.go`
- Create: `service/role.go`
- Create: `service/permission.go`

- [ ] **Step 1: 创建 Permission Service**

```go
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
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
}

type UpdatePermissionRequest struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
}

func (s *PermissionService) Create(req *CreatePermissionRequest) (*model.Permission, error) {
	exists, err := s.repo.ExistsByCode(req.Code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("权限编码已存在")
	}

	permission := &model.Permission{
		Name: req.Name,
		Code: req.Code,
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

	// 检查新编码是否已被其他权限使用
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
```

- [ ] **Step 2: 创建 Role Service**

```go
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
```

- [ ] **Step 3: 创建 User Service**

```go
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
```

- [ ] **Step 4: 创建 Auth Service**

```go
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
	Token string      `json:"token"`
	User  *model.User `json:"user"`
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

	// 清空密码
	user.Password = ""

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) Logout() error {
	// JWT无状态，客户端删除Token即可
	// 如需服务端失效，可在此实现Token黑名单
	return nil
}
```

- [ ] **Step 5: 验证编译**

```bash
go build ./service/...
```

Expected: 无错误输出

---

## Task 9: 中间件 - JWT认证

**Files:**
- Create: `middleware/jwt.go`

- [ ] **Step 1: 创建JWT认证中间件**

```go
package middleware

import (
	"hello/pkg/jwt"
	"hello/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserID   = "userID"
	ContextUsername = "username"
	ContextRoles    = "roles"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未提供认证Token")
			c.Abort()
			return
		}

		// 解析Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Token格式错误，应为: Bearer <token>")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			response.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		// 将用户信息存入Context
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUsername, claims.Username)
		c.Set(ContextRoles, claims.Roles)

		c.Next()
	}
}

// GetUserID 从Context获取用户ID
func GetUserID(c *gin.Context) uint {
	if userID, exists := c.Get(ContextUserID); exists {
		return userID.(uint)
	}
	return 0
}

// GetUsername 从Context获取用户名
func GetUsername(c *gin.Context) string {
	if username, exists := c.Get(ContextUsername); exists {
		return username.(string)
	}
	return ""
}

// GetRoles 从Context获取角色列表
func GetRoles(c *gin.Context) []string {
	if roles, exists := c.Get(ContextRoles); exists {
		return roles.([]string)
	}
	return nil
}

// IsAdmin 检查是否为管理员
func IsAdmin(c *gin.Context) bool {
	roles := GetRoles(c)
	for _, role := range roles {
		if role == "admin" {
			return true
		}
	}
	return false
}
```

- [ ] **Step 2: 验证编译**

```bash
go build ./middleware/...
```

Expected: 无错误输出

---

## Task 10: 中间件 - 权限校验

**Files:**
- Create: `middleware/permission.go`

- [ ] **Step 1: 创建权限校验中间件**

```go
package middleware

import (
	"hello/pkg/response"
	"hello/repository"

	"github.com/gin-gonic/gin"
)

// RequirePermission 校验单个权限
func RequirePermission(permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 管理员直接放行
		if IsAdmin(c) {
			c.Next()
			return
		}

		userID := GetUserID(c)
		if userID == 0 {
			response.Unauthorized(c, "用户未认证")
			c.Abort()
			return
		}

		// 获取用户权限
		userRepo := repository.NewUserRepository()
		permissionCodes, err := userRepo.GetUserPermissionCodes(userID)
		if err != nil {
			response.InternalError(c, "获取权限失败")
			c.Abort()
			return
		}

		// 检查是否拥有所需权限
		hasPermission := false
		for _, code := range permissionCodes {
			if code == permissionCode {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			response.Forbidden(c, "无权限访问")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission 校验多个权限（满足任一即可）
func RequireAnyPermission(permissionCodes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 管理员直接放行
		if IsAdmin(c) {
			c.Next()
			return
		}

		userID := GetUserID(c)
		if userID == 0 {
			response.Unauthorized(c, "用户未认证")
			c.Abort()
			return
		}

		userRepo := repository.NewUserRepository()
		userPermissionCodes, err := userRepo.GetUserPermissionCodes(userID)
		if err != nil {
			response.InternalError(c, "获取权限失败")
			c.Abort()
			return
		}

		// 检查是否拥有任一权限
		hasPermission := false
		permissionSet := make(map[string]bool)
		for _, code := range userPermissionCodes {
			permissionSet[code] = true
		}
		for _, code := range permissionCodes {
			if permissionSet[code] {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			response.Forbidden(c, "无权限访问")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAllPermissions 校验多个权限（需全部满足）
func RequireAllPermissions(permissionCodes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 管理员直接放行
		if IsAdmin(c) {
			c.Next()
			return
		}

		userID := GetUserID(c)
		if userID == 0 {
			response.Unauthorized(c, "用户未认证")
			c.Abort()
			return
		}

		userRepo := repository.NewUserRepository()
		userPermissionCodes, err := userRepo.GetUserPermissionCodes(userID)
		if err != nil {
			response.InternalError(c, "获取权限失败")
			c.Abort()
			return
		}

		// 检查是否拥有全部权限
		permissionSet := make(map[string]bool)
		for _, code := range userPermissionCodes {
			permissionSet[code] = true
		}

		for _, code := range permissionCodes {
			if !permissionSet[code] {
				response.Forbidden(c, "无权限访问")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
```

- [ ] **Step 2: 验证编译**

```bash
go build ./middleware/...
```

Expected: 无错误输出

---

## Task 11: Handler层 - 认证接口

**Files:**
- Create: `handler/auth.go`

- [ ] **Step 1: 创建认证处理器**

```go
package handler

import (
	"hello/pkg/response"
	"hello/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(),
	}
}

// Login 登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	result, err := h.authService.Login(&req)
	if err != nil {
		response.Error(c, 401, 41000, err.Error())
		return
	}

	response.Success(c, result)
}

// Logout 登出
func (h *AuthHandler) Logout(c *gin.Context) {
	// JWT无状态，客户端删除Token即可
	response.Success(c, gin.H{
		"message": "登出成功",
	})
}
```

- [ ] **Step 2: 验证编译**

```bash
go build ./handler/...
```

Expected: 无错误输出

---

## Task 12: Handler层 - 用户接口

**Files:**
- Create: `handler/user.go`

- [ ] **Step 1: 创建用户处理器**

```go
package handler

import (
	"hello/pkg/response"
	"hello/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

// List 获取用户列表
func (h *UserHandler) List(c *gin.Context) {
	users, err := h.userService.GetAll()
	if err != nil {
		response.InternalError(c, "获取用户列表失败")
		return
	}
	response.Success(c, users)
}

// Get 获取用户详情
func (h *UserHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, user)
}

// Create 创建用户
func (h *UserHandler) Create(c *gin.Context) {
	var req service.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user, err := h.userService.Create(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, user)
}

// Update 更新用户
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	var req service.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user, err := h.userService.Update(uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, user)
}

// Delete 删除用户
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	if err := h.userService.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "删除成功",
	})
}

// AssignRoles 分配角色
func (h *UserHandler) AssignRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	var req service.AssignRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.userService.AssignRoles(uint(id), &req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "分配角色成功",
	})
}
```

- [ ] **Step 2: 验证编译**

```bash
go build ./handler/...
```

Expected: 无错误输出

---

## Task 13: Handler层 - 角色接口

**Files:**
- Create: `handler/role.go`

- [ ] **Step 1: 创建角色处理器**

```go
package handler

import (
	"hello/pkg/response"
	"hello/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleService *service.RoleService
}

func NewRoleHandler() *RoleHandler {
	return &RoleHandler{
		roleService: service.NewRoleService(),
	}
}

// List 获取角色列表
func (h *RoleHandler) List(c *gin.Context) {
	roles, err := h.roleService.GetAll()
	if err != nil {
		response.InternalError(c, "获取角色列表失败")
		return
	}
	response.Success(c, roles)
}

// Get 获取角色详情
func (h *RoleHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	role, err := h.roleService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "角色不存在")
		return
	}

	response.Success(c, role)
}

// Create 创建角色
func (h *RoleHandler) Create(c *gin.Context) {
	var req service.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	role, err := h.roleService.Create(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, role)
}

// Update 更新角色
func (h *RoleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	var req service.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	role, err := h.roleService.Update(uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, role)
}

// Delete 删除角色
func (h *RoleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	if err := h.roleService.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "删除成功",
	})
}

// AssignPermissions 分配权限
func (h *RoleHandler) AssignPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	var req service.AssignPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.roleService.AssignPermissions(uint(id), &req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "分配权限成功",
	})
}
```

- [ ] **Step 2: 验证编译**

```bash
go build ./handler/...
```

Expected: 无错误输出

---

## Task 14: Handler层 - 权限接口

**Files:**
- Create: `handler/permission.go`

- [ ] **Step 1: 创建权限处理器**

```go
package handler

import (
	"hello/pkg/response"
	"hello/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	permissionService *service.PermissionService
}

func NewPermissionHandler() *PermissionHandler {
	return &PermissionHandler{
		permissionService: service.NewPermissionService(),
	}
}

// List 获取权限列表
func (h *PermissionHandler) List(c *gin.Context) {
	permissions, err := h.permissionService.GetAll()
	if err != nil {
		response.InternalError(c, "获取权限列表失败")
		return
	}
	response.Success(c, permissions)
}

// Get 获取权限详情
func (h *PermissionHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的权限ID")
		return
	}

	permission, err := h.permissionService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "权限不存在")
		return
	}

	response.Success(c, permission)
}

// Create 创建权限
func (h *PermissionHandler) Create(c *gin.Context) {
	var req service.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	permission, err := h.permissionService.Create(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, permission)
}

// Update 更新权限
func (h *PermissionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的权限ID")
		return
	}

	var req service.UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	permission, err := h.permissionService.Update(uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, permission)
}

// Delete 删除权限
func (h *PermissionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的权限ID")
		return
	}

	if err := h.permissionService.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "删除成功",
	})
}
```

- [ ] **Step 2: 验证编译**

```bash
go build ./handler/...
```

Expected: 无错误输出

---

## Task 15: 主入口与路由注册

**Files:**
- Modify: `main.go`

- [ ] **Step 1: 创建主入口文件**

```go
package main

import (
	"fmt"
	"hello/config"
	"hello/handler"
	"hello/middleware"
	"hello/repository"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	config.LoadConfig()

	// 初始化数据库
	if err := repository.InitDB(); err != nil {
		panic(fmt.Sprintf("初始化数据库失败: %v", err))
	}

	// 初始化基础数据
	if err := repository.InitData(); err != nil {
		panic(fmt.Sprintf("初始化基础数据失败: %v", err))
	}

	// 创建Gin引擎
	r := gin.Default()

	// 全局中间件
	r.Use(gin.Recovery())

	// 注册路由
	setupRoutes(r)

	// 启动服务
	addr := fmt.Sprintf(":%s", config.AppConfig.Port)
	fmt.Printf("RBAC服务启动在 %s\n", addr)
	if err := r.Run(addr); err != nil {
		panic(fmt.Sprintf("启动服务失败: %v", err))
	}
}

func setupRoutes(r *gin.Engine) {
	// Handler实例
	authHandler := handler.NewAuthHandler()
	userHandler := handler.NewUserHandler()
	roleHandler := handler.NewRoleHandler()
	permissionHandler := handler.NewPermissionHandler()

	// API路由组
	api := r.Group("/api")
	{
		// 认证接口（无需认证）
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", middleware.JWTAuth(), authHandler.Logout)
		}

		// 需要认证的接口
		protected := api.Group("")
		protected.Use(middleware.JWTAuth())
		{
			// 用户管理
			users := protected.Group("/users")
			{
				users.GET("", middleware.RequirePermission("user:read"), userHandler.List)
				users.GET("/:id", middleware.RequirePermission("user:read"), userHandler.Get)
				users.POST("", middleware.RequirePermission("user:create"), userHandler.Create)
				users.PUT("/:id", middleware.RequirePermission("user:update"), userHandler.Update)
				users.DELETE("/:id", middleware.RequirePermission("user:delete"), userHandler.Delete)
				users.PUT("/:id/roles", middleware.RequirePermission("user:update"), userHandler.AssignRoles)
			}

			// 角色管理
			roles := protected.Group("/roles")
			{
				roles.GET("", middleware.RequirePermission("role:read"), roleHandler.List)
				roles.GET("/:id", middleware.RequirePermission("role:read"), roleHandler.Get)
				roles.POST("", middleware.RequirePermission("role:create"), roleHandler.Create)
				roles.PUT("/:id", middleware.RequirePermission("role:update"), roleHandler.Update)
				roles.DELETE("/:id", middleware.RequirePermission("role:delete"), roleHandler.Delete)
				roles.PUT("/:id/permissions", middleware.RequirePermission("role:update"), roleHandler.AssignPermissions)
			}

			// 权限管理
			permissions := protected.Group("/permissions")
			{
				permissions.GET("", middleware.RequirePermission("permission:read"), permissionHandler.List)
				permissions.GET("/:id", middleware.RequirePermission("permission:read"), permissionHandler.Get)
				permissions.POST("", middleware.RequirePermission("permission:create"), permissionHandler.Create)
				permissions.PUT("/:id", middleware.RequirePermission("permission:update"), permissionHandler.Update)
				permissions.DELETE("/:id", middleware.RequirePermission("permission:delete"), permissionHandler.Delete)
			}
		}
	}
}
```

- [ ] **Step 2: 编译整个项目**

```bash
go build -o rbac-server .
```

Expected: 无错误输出，生成 rbac-server 可执行文件

---

## Task 16: 运行测试

**Files:**
- 无新文件

- [ ] **Step 1: 启动服务**

```bash
JWT_SECRET=your-secret-key ./rbac-server
```

Expected: 输出 "RBAC服务启动在 :8080"

- [ ] **Step 2: 测试登录接口**

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

Expected: 返回包含token的JSON响应

- [ ] **Step 3: 测试获取用户列表**

```bash
curl -X GET http://localhost:8080/api/users \
  -H "Authorization: Bearer <上一步获取的token>"
```

Expected: 返回用户列表JSON

---

## Task 17: 清理团队

**Files:**
- 无

- [ ] **Step 1: 关闭团队成员**

通知所有团队成员工作已完成。

---

## 实现顺序总结

1. **基础设施** (Task 1-6): 依赖、配置、模型、数据库、工具函数
2. **数据层** (Task 7): Repository层
3. **业务层** (Task 8): Service层
4. **中间件** (Task 9-10): JWT认证、权限校验
5. **接口层** (Task 11-14): Handler层
6. **集成** (Task 15-16): 路由注册、运行测试
