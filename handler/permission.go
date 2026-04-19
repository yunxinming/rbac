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

// Tree 获取权限树
func (h *PermissionHandler) Tree(c *gin.Context) {
	tree, err := h.permissionService.GetTree()
	if err != nil {
		response.InternalError(c, "获取权限树失败")
		return
	}
	response.Success(c, tree)
}
