package database

import (
	"context"
	"fmt"
	"time"

	"github.com/your-org/atlas/ent"
	"github.com/your-org/atlas/pkg/config"

	_ "github.com/lib/pq"
)

// NewClient 创建数据库客户端
func NewClient(cfg *config.DatabaseConfig) (*ent.Client, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 配置连接池
	db := client.Driver().(*ent.Driver).DB()
	db.SetMaxOpenConns(cfg.MaxOpen)
	db.SetMaxIdleConns(cfg.MaxIdle)
	db.SetConnMaxLifetime(time.Hour)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Schema.Create(ctx); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return client, nil
}
