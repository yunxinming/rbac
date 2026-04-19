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
