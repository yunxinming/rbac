package model

import (
	"time"

	"gorm.io/gorm"
)

// PermissionType 权限类型
type PermissionType string

const (
	PermissionTypeDirectory PermissionType = "directory" // 目录
	PermissionTypeMenu      PermissionType = "menu"      // 菜单
	PermissionTypeOperation PermissionType = "operation" // 操作
)

type Permission struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Name       string         `gorm:"size:64;not null" json:"name"`
	Code       string         `gorm:"size:128;uniqueIndex;not null" json:"code"`
	Type       PermissionType `gorm:"size:20;not null;default:operation" json:"type"`
	ParentID   uint           `gorm:"default:0;index" json:"parent_id"`
	Path       string         `gorm:"size:128" json:"path"`        // 前端路由路径
	Icon       string         `gorm:"size:64" json:"icon"`         // 图标
	Sort       int            `gorm:"default:0" json:"sort"`       // 排序
	ApiPath    string         `gorm:"size:255" json:"api_path"`    // API路径
	ApiMethod  string         `gorm:"size:10" json:"api_method"`   // API方法
	Component  string         `gorm:"size:128" json:"component"`   // 前端组件路径
	Status     int8           `gorm:"default:1" json:"status"`     // 状态：1启用 0禁用
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Children   []Permission   `gorm:"-" json:"children,omitempty"` // 子权限（非数据库字段）
}

func (Permission) TableName() string {
	return "permissions"
}
