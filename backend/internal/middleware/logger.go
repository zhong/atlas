package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/asset-management/pkg/logger"
	"go.uber.org/zap"
)

// Logger 日志中间件
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// 处理请求
		err := c.Next()

		// 记录日志
		logger.Info("HTTP Request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", time.Since(start)),
			zap.String("ip", c.IP()),
			zap.String("user_agent", c.Get("User-Agent")),
		)

		return err
	}
}
