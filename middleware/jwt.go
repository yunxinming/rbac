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
