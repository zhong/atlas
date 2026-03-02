package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/asset-management/internal/middleware"
	"github.com/your-org/asset-management/internal/router"
	"github.com/your-org/asset-management/pkg/config"
	"github.com/your-org/asset-management/pkg/database"
	"github.com/your-org/asset-management/pkg/logger"
	redisClient "github.com/your-org/asset-management/pkg/redis"
	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	if err := logger.Init(&cfg.Log); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Get().Sync()

	logger.Info("Starting Asset Management System...")

	// 连接数据库
	dbClient, err := database.NewClient(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer dbClient.Close()
	logger.Info("Database connected successfully")

	// 连接Redis
	rdb, err := redisClient.NewClient(&cfg.Redis)
	if err != nil {
		logger.Fatal("Failed to connect to redis", zap.Error(err))
	}
	defer rdb.Close()
	logger.Info("Redis connected successfully")

	// 创建Fiber应用
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler(),
		AppName:      "Asset Management System",
	})

	// 设置路由
	router.Setup(app, dbClient, cfg)

	// 启动服务器
	go func() {
		addr := fmt.Sprintf(":%d", cfg.Server.Port)
		logger.Info("Server starting", zap.String("address", addr))
		if err := app.Listen(addr); err != nil {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
