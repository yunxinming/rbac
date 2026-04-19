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
