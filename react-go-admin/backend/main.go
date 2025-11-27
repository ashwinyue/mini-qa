package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"react-go-admin-backend/api"
	"react-go-admin-backend/config"
	"react-go-admin-backend/models"
)

func main() {
	// 初始化数据库
	if err := models.InitDB(); err != nil {
		log.Fatal("数据库初始化失败:", err)
	}

	// 创建 Gin 引擎
	r := gin.Default()

	// 配置 CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 注册路由
	api.RegisterRoutes(r)

	// 启动服务器
	port := config.GetServerPort()
	log.Printf("服务器启动在端口 %s", port)
	if err := r.Run(port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
