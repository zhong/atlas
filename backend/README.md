# Atlas Backend

Atlas IT资产管理系统后端服务

## 技术栈

- Go 1.21+
- Fiber v2 (Web框架)
- Ent (ORM)
- PostgreSQL 15+
- Redis 7+
- RabbitMQ 3+

## 项目结构

```
backend/
├── cmd/                    # 应用入口
│   ├── api/               # API服务
│   ├── worker/            # 后台任务
│   └── migrate/           # 数据库迁移
├── ent/                   # Ent ORM
│   ├── schema/           # Schema定义
│   └── [generated]       # 生成的代码
├── internal/              # 内部代码
│   ├── handler/          # HTTP处理器
│   ├── service/          # 业务逻辑
│   ├── middleware/       # 中间件
│   ├── dto/              # 数据传输对象
│   ├── pkg/              # 内部包
│   └── router/           # 路由配置
├── pkg/                   # 公共库
│   ├── config/           # 配置管理
│   ├── database/         # 数据库
│   ├── redis/            # Redis
│   ├── logger/           # 日志
│   ├── jwt/              # JWT
│   └── utils/            # 工具函数
├── config/                # 配置文件
├── tests/                 # 测试
└── Makefile              # 常用命令
```

## 快速开始

### 1. 安装依赖

```bash
make install
```

### 2. 启动基础设施服务

```bash
make docker-up
```

### 3. 生成Ent代码

```bash
make generate
```

### 4. 运行应用

```bash
make run
```

应用将在 http://localhost:8080 启动

## 开发命令

```bash
# 安装依赖
make install

# 生成Ent代码
make generate

# 构建应用
make build

# 运行应用
make run

# 运行测试
make test

# 代码检查
make lint

# 清理构建产物
make clean

# Docker相关
make docker-up      # 启动Docker服务
make docker-down    # 停止Docker服务
make docker-ps      # 查看服务状态
make docker-logs    # 查看日志
```

## API文档

启动服务后访问：
- Swagger UI: http://localhost:8080/swagger
- OpenAPI Spec: http://localhost:8080/api/openapi.yaml

## 配置

配置文件位于 `config/config.yaml`，可以通过环境变量覆盖：

```bash
export ASSET_MGMT_SERVER_PORT=8080
export ASSET_MGMT_DATABASE_HOST=localhost
```

## 测试

```bash
# 运行所有测试
make test

# 运行特定包的测试
go test -v ./internal/service/...

# 运行测试并生成覆盖率报告
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 代码生成

### Ent Schema

1. 在 `ent/schema/` 目录下创建新的schema文件
2. 运行 `make generate` 生成代码
3. Ent会自动生成CRUD操作、查询构建器等代码

示例：
```go
// ent/schema/user.go
package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
)

type User struct {
    ent.Schema
}

func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("username").Unique(),
        field.String("email").Unique(),
    }
}
```

## 部署

### Docker构建

```bash
docker build -t asset-management-api:latest .
```

### 运行容器

```bash
docker run -d \
  -p 8080:8080 \
  -e ASSET_MGMT_DATABASE_HOST=postgres \
  asset-management-api:latest
```

## 贡献指南

1. Fork项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建Pull Request

## 许可证

MIT License
