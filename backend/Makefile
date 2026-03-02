.PHONY: help install generate build run test lint clean docker-up docker-down

# 默认目标
help:
	@echo "Atlas - IT Asset Management System"
	@echo ""
	@echo "Usage:"
	@echo "  make install      - Install dependencies"
	@echo "  make generate     - Generate Ent code"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests"
	@echo "  make lint         - Run linter"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make docker-up    - Start Docker services"
	@echo "  make docker-down  - Stop Docker services"

# 安装依赖
install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# 生成Ent代码
generate:
	@echo "Generating Ent code..."
	go generate ./ent

# 构建应用
build:
	@echo "Building application..."
	go build -o bin/api cmd/api/main.go

# 运行应用
run:
	@echo "Running application..."
	go run cmd/api/main.go

# 运行测试
test:
	@echo "Running tests..."
	go test -v -cover ./...

# 代码检查
lint:
	@echo "Running linter..."
	golangci-lint run

# 清理
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf logs/

# 启动Docker服务
docker-up:
	@echo "Starting Docker services..."
	cd .. && docker-compose up -d

# 停止Docker服务
docker-down:
	@echo "Stopping Docker services..."
	cd .. && docker-compose down

# 查看Docker服务状态
docker-ps:
	@echo "Docker services status..."
	cd .. && docker-compose ps

# 查看日志
docker-logs:
	@echo "Docker services logs..."
	cd .. && docker-compose logs -f
