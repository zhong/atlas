package auth

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/atlas/ent"
	"github.com/your-org/atlas/ent/user"
	"github.com/your-org/atlas/pkg/config"
	"github.com/your-org/atlas/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	client *ent.Client
	cfg    *config.Config
}

func NewHandler(client *ent.Client, cfg *config.Config) *Handler {
	return &Handler{
		client: client,
		cfg:    cfg,
	}
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      UserInfo  `json:"user"`
}

type UserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	RealName string `json:"real_name"`
	Role     string `json:"role"`
}

// Login 用户登录
func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// 查询用户
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u, err := h.client.User.Query().
		Where(user.Username(req.Username)).
		WithRoles().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid username or password",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}

	// 检查用户状态
	if u.Status != "active" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "User account is not active",
		})
	}

	// 获取用户角色
	roles := u.Edges.Roles
	var roleCode string
	if len(roles) > 0 {
		roleCode = roles[0].Code
	}

	// 生成JWT token
	token, err := jwt.GenerateToken(int64(u.ID), u.Username, roleCode, &h.cfg.JWT)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// 返回响应
	expiresAt := time.Now().Add(time.Hour * time.Duration(h.cfg.JWT.ExpireTime))
	return c.JSON(LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: UserInfo{
			ID:       int64(u.ID),
			Username: u.Username,
			Email:    u.Email,
			RealName: u.RealName,
			Role:     roleCode,
		},
	})
}
