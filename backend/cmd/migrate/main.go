package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/your-org/atlas/ent"
	"github.com/your-org/atlas/pkg/config"

	_ "github.com/lib/pq"
)

func main() {
	// 加载配置
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 构建数据库连接字符串
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	// 连接数据库
	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// 运行迁移
	log.Println("Running database migrations...")
	if err := client.Schema.Create(
		ctx,
		// 使用 WithDropColumn 和 WithDropIndex 选项来处理 schema 变更
		// migrate.WithDropColumn(true),
		// migrate.WithDropIndex(true),
	); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	log.Println("✅ Database migrations completed successfully!")

	// 显示创建的表
	log.Println("\nCreated tables:")
	tables := []string{
		"users", "roles", "permissions",
		"warehouses", "locations", "asset_types", "assets",
		"inventory_records", "suppliers", "purchase_orders", "order_items",
		"data_centers", "rooms", "racks", "rack_units",
		"approvals", "approval_nodes",
		"network_connections", "ip_addresses",
		"repair_vendors", "repair_tickets",
	}

	for _, table := range tables {
		log.Printf("  - %s", table)
	}

	log.Println("\n✅ All done! You can now run the seed script to populate initial data.")
	os.Exit(0)
}
