package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/atlas/pkg/config"
	"github.com/your-org/atlas/pkg/jwt"
	"github.com/your-org/atlas/pkg/utils"
)

// Auth JWT认证中间件
func Auth(cfg *config.JWTConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 获取Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.Error(c, 401, "missing authorization header")
		}

		// 解析Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.Error(c, 401, "invalid authorization header format")
		}

		// 验证token
		claims, err := jwt.ParseToken(parts[1], cfg)
		if err != nil {
			return utils.Error(c, 401, "invalid or expired token")
		}

		// 将用户信息存入context
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// GetUserID 从context获取用户ID
func GetUserID(c *fiber.Ctx) int64 {
	if userID, ok := c.Locals("user_id").(int64); ok {
		return userID
	}
	return 0
}

// GetUsername 从context获取用户名
func GetUsername(c *fiber.Ctx) string {
	if username, ok := c.Locals("username").(string); ok {
		return username
	}
	return ""
}

// GetRole 从context获取角色
func GetRole(c *fiber.Ctx) string {
	if role, ok := c.Locals("role").(string); ok {
		return role
	}
	return ""
}
