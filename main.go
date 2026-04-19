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

	// CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

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
				users.GET("", middleware.RequirePermission("system:user:read"), userHandler.List)
				users.GET("/:id", middleware.RequirePermission("system:user:read"), userHandler.Get)
				users.POST("", middleware.RequirePermission("system:user:create"), userHandler.Create)
				users.PUT("/:id", middleware.RequirePermission("system:user:update"), userHandler.Update)
				users.DELETE("/:id", middleware.RequirePermission("system:user:delete"), userHandler.Delete)
				users.PUT("/:id/roles", middleware.RequirePermission("system:user:update"), userHandler.AssignRoles)
			}

			// 角色管理
			roles := protected.Group("/roles")
			{
				roles.GET("", middleware.RequirePermission("system:role:read"), roleHandler.List)
				roles.GET("/:id", middleware.RequirePermission("system:role:read"), roleHandler.Get)
				roles.POST("", middleware.RequirePermission("system:role:create"), roleHandler.Create)
				roles.PUT("/:id", middleware.RequirePermission("system:role:update"), roleHandler.Update)
				roles.DELETE("/:id", middleware.RequirePermission("system:role:delete"), roleHandler.Delete)
				roles.PUT("/:id/permissions", middleware.RequirePermission("system:role:update"), roleHandler.AssignPermissions)
			}

			// 权限管理
			permissions := protected.Group("/permissions")
			{
				permissions.GET("", middleware.RequirePermission("system:permission:read"), permissionHandler.List)
				permissions.GET("/tree", middleware.RequirePermission("system:permission:read"), permissionHandler.Tree)
				permissions.GET("/:id", middleware.RequirePermission("system:permission:read"), permissionHandler.Get)
				permissions.POST("", middleware.RequirePermission("system:permission:create"), permissionHandler.Create)
				permissions.PUT("/:id", middleware.RequirePermission("system:permission:update"), permissionHandler.Update)
				permissions.DELETE("/:id", middleware.RequirePermission("system:permission:delete"), permissionHandler.Delete)
			}
		}
	}
}
