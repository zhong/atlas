package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/your-org/asset-management/pkg/logger"
	"github.com/your-org/asset-management/pkg/utils"
	"go.uber.org/zap"
)

// ErrorHandler 全局错误处理中间件
func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// 默认500错误
		code := fiber.StatusInternalServerError
		message := "Internal Server Error"

		// 如果是Fiber错误，使用其状态码和消息
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
			message = e.Message
		}

		// 记录错误日志
		logger.Error("Request Error",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", code),
			zap.Error(err),
		)

		// 返回错误响应
		return c.Status(code).JSON(utils.Response{
			Code:    code,
			Message: message,
		})
	}
}
