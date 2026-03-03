package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/your-org/atlas/ent"
	"github.com/your-org/atlas/internal/handler/asset"
	"github.com/your-org/atlas/internal/handler/auth"
	"github.com/your-org/atlas/internal/middleware"
	"github.com/your-org/atlas/pkg/config"
)

// Setup 设置路由
func Setup(app *fiber.App, client *ent.Client, cfg *config.Config) {
	// 全局中间件
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	app.Use(middleware.Logger())

	// 健康检查
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// API v1路由组
	v1 := app.Group("/api/v1")

	// 初始化处理器
	authHandler := auth.NewHandler(client, cfg)
	assetHandler := asset.NewHandler(client, cfg)

	// 公开路由（不需要认证）
	public := v1.Group("")
	{
		public.Post("/auth/login", authHandler.Login)
	}

	// 需要认证的路由
	protected := v1.Group("", middleware.Auth(&cfg.JWT))
	{
		// 资产管理
		assets := protected.Group("/assets")
		{
			assets.Get("/", assetHandler.List)
			assets.Get("/:id", assetHandler.Get)
			assets.Post("/", assetHandler.Create)
		}

		// 库存管理
		inventory := protected.Group("/inventory")
		{
			inventory.Get("/stock", func(c *fiber.Ctx) error {
				return c.JSON(fiber.Map{"message": "get inventory"})
			})
		}

		// 采购管理
		purchase := protected.Group("/purchase")
		{
			purchase.Get("/orders", func(c *fiber.Ctx) error {
				return c.JSON(fiber.Map{"message": "get purchase orders"})
			})
		}

		// DCIM
		dcim := protected.Group("/dcim")
		{
			dcim.Get("/datacenters", func(c *fiber.Ctx) error {
				return c.JSON(fiber.Map{"message": "get datacenters"})
			})
		}

		// 审批
		approvals := protected.Group("/approvals")
		{
			approvals.Get("/", func(c *fiber.Ctx) error {
				return c.JSON(fiber.Map{"message": "get approvals"})
			})
		}
	}
}
