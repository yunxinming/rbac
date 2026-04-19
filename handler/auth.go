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
