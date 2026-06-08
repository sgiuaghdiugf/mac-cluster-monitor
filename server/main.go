package main

import (
	"log"
	"mac-cluster-monitor/server/config"
	"mac-cluster-monitor/server/database"
	"mac-cluster-monitor/server/handlers"
	"mac-cluster-monitor/server/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Printf("警告: 加载配置文件失败,使用默认配置: %v", err)
	}

	// 设置gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库
	db, err := database.Init(cfg)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	defer db.Close()

	// 创建处理器
	authHandler := handlers.NewAuthHandler(db, cfg)
	deviceHandler := handlers.NewDeviceHandler(db)
	userHandler := handlers.NewUserHandler(db)

	// 设置定时任务: 检查离线设备
	go func() {
		ticker := time.NewTicker(time.Duration(cfg.Monitor.CheckInterval) * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			deviceHandler.CheckOfflineDevices()
		}
	}()

	// 设置定时任务: 更新盈利数据
	go func() {
		interval := time.Duration(cfg.Profit.Interval) * time.Second
		if interval < 1 {
			interval = 5 * time.Second
		}
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			deviceHandler.UpdateProfit(cfg.Profit.Increment)
		}
	}()

	// 设置路由
	r := gin.Default()

	// CORS配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 公开API
	api := r.Group("/api")
	{
		// 登录 (不支持注册)
		api.POST("/login", authHandler.Login)

		// 设备心跳 (不需要认证)
		api.POST("/heartbeat", deviceHandler.Heartbeat)
	}

	// 需要认证的API
	authorized := api.Group("/")
	authorized.Use(middleware.AuthMiddleware(cfg))
	{
		// 设备管理 (普通用户只能看自己的)
		authorized.GET("/devices", deviceHandler.GetDevices)
		authorized.GET("/devices/stats", deviceHandler.GetDeviceStats)

		// 普通用户Dashboard
		authorized.GET("/user/dashboard", userHandler.GetUserDashboard)
	}

	// 管理员专用API
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(cfg))
	admin.Use(middleware.AdminMiddleware())
	{
		// 用户管理
		admin.GET("/users", userHandler.GetUsers)
		admin.POST("/users", userHandler.CreateUser)
		admin.PUT("/users/:id", userHandler.UpdateUser)
		admin.DELETE("/users/:id", userHandler.DeleteUser)

		// 设备管理
		admin.DELETE("/devices/:id", deviceHandler.DeleteDevice)
		admin.GET("/devices/unbound", deviceHandler.GetUnboundDevices)
		admin.PUT("/devices/:id/bind", userHandler.BindDevice)
		admin.PUT("/devices/:id/unbind", userHandler.UnbindDevice)
	}

	// 启动服务器
	log.Printf("服务器启动在端口 %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
