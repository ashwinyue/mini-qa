package api

import (
	"react-go-admin-backend/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册路由
func RegisterRoutes(r *gin.Engine) {
	// API 路由组
	api := r.Group("/api")

	// 认证路由（无需 token）
	authCtrl := NewAuthController()
	auth := api.Group("/auth")
	{
		auth.POST("/login", authCtrl.Login)
		auth.POST("/logout", authCtrl.Logout)
		auth.GET("/profile", middleware.AuthMiddleware(), authCtrl.GetProfile)
		auth.POST("/change-password", middleware.AuthMiddleware(), authCtrl.ChangePassword)
	}

	// 需要认证的路由
	authorized := api.Group("")
	authorized.Use(middleware.AuthMiddleware())
	{
		// 用户管理
		userCtrl := NewUserController()
		users := authorized.Group("/users")
		{
			users.GET("", userCtrl.GetList)
			users.GET("/:id", userCtrl.GetDetail)
			users.POST("", userCtrl.Create)
			users.PUT("/:id", userCtrl.Update)
			users.DELETE("/:id", userCtrl.Delete)
		}

		// 角色管理
		roleCtrl := NewRoleController()
		roles := authorized.Group("/roles")
		{
			roles.GET("", roleCtrl.GetList)
			roles.GET("/:id", roleCtrl.GetDetail)
			roles.POST("", roleCtrl.Create)
			roles.PUT("/:id", roleCtrl.Update)
			roles.DELETE("/:id", roleCtrl.Delete)
		}

	}
}
