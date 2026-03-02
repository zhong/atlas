# Atlas - IT资产管理系统

**Atlas** 是为GPU算力公司打造的IT资产全生命周期管理系统，支持设备采购、入库、库存管理、出库、DCIM可视化、设备生命周期管理等功能。

Atlas（阿特拉斯）源自希腊神话中的擎天神，寓意全面掌控和管理企业的IT资产与基础设施。

## 技术栈

### 后端
- **语言**: Go 1.21+
- **Web框架**: Fiber v2
- **ORM**: Ent (代码生成 + 类型安全)
- **数据库**: PostgreSQL 15+
- **缓存**: Redis 7+
- **消息队列**: RabbitMQ 3+
- **文件存储**: MinIO

### 前端
- **构建工具**: Vite 5+
- **框架**: React 18+
- **语言**: TypeScript 5+
- **UI库**: Ant Design 5+
- **状态管理**: Zustand
- **数据获取**: TanStack Query (React Query v5)

## 快速开始

### 前置要求

- Docker 20+
- Docker Compose 2+
- Go 1.21+ (后端开发)
- Node.js 18+ (前端开发)

### 启动开发环境

1. 克隆项目
```bash
git clone <repository-url>
cd hello
```

2. 启动基础设施服务
```bash
docker-compose up -d
```

3. 验证服务状态
```bash
docker-compose ps
```

### 服务访问地址

| 服务 | 地址 | 用户名 | 密码 |
|------|------|--------|------|
| PostgreSQL | localhost:5432 | admin | admin123 |
| Redis | localhost:6379 | - | redis123 |
| MinIO API | localhost:9000 | minioadmin | minioadmin123 |
| MinIO Console | http://localhost:9001 | minioadmin | minioadmin123 |
| RabbitMQ | localhost:5672 | admin | admin123 |
| RabbitMQ Management | http://localhost:15672 | admin | admin123 |
| Adminer (数据库管理) | http://localhost:8081 | - | - |

### 后端开发

```bash
cd backend

# 安装依赖
go mod download

# 生成Ent代码
go generate ./ent

# 运行数据库迁移
go run cmd/migrate/main.go

# 启动API服务
go run cmd/api/main.go
```

### 前端开发

```bash
cd web

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

## 项目结构

```
.
├── backend/                 # 后端Go项目
│   ├── cmd/                # 应用入口
│   ├── ent/                # Ent Schema和生成代码
│   ├── internal/           # 内部代码
│   ├── pkg/                # 公共库
│   └── config/             # 配置文件
├── web/                    # 前端React项目
│   ├── src/
│   │   ├── pages/         # 页面组件
│   │   ├── components/    # 通用组件
│   │   ├── services/      # API服务
│   │   └── stores/        # 状态管理
│   └── public/
├── docker-compose.yml      # Docker Compose配置
├── scripts/                # 脚本文件
├── ARCHITECTURE.md         # 技术架构文档
└── README.md              # 本文件
```

## 开发计划

详见 [ARCHITECTURE.md](./ARCHITECTURE.md) 第6章节。

### 当前阶段：阶段0 - 基础设施搭建

- [x] 创建技术架构文档
- [x] 配置Docker开发环境
- [ ] 搭建后端项目脚手架
- [ ] 定义核心数据模型Schema
- [ ] 搭建前端项目脚手架

## 常用命令

### Docker

```bash
# 启动所有服务
docker-compose up -d

# 停止所有服务
docker-compose down

# 查看服务状态
docker-compose ps

# 查看服务日志
docker-compose logs -f [service-name]

# 重启服务
docker-compose restart [service-name]

# 清理所有数据（危险操作）
docker-compose down -v
```

### 后端

```bash
# 代码生成
make generate

# 运行测试
make test

# 代码检查
make lint

# 构建
make build

# 运行
make run
```

### 前端

```bash
# 开发
npm run dev

# 构建
npm run build

# 预览构建结果
npm run preview

# 代码检查
npm run lint

# 类型检查
npm run type-check
```

## 文档

- [技术架构文档](./ARCHITECTURE.md) - 完整的技术架构设计
- [API文档](./backend/api/openapi.yaml) - OpenAPI规范
- [开发指南](./docs/development.md) - 开发规范和最佳实践
- [部署文档](./docs/deployment.md) - 部署指南

## 贡献指南

1. Fork项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启Pull Request

## 许可证

[MIT License](LICENSE)

## 联系方式

项目负责人 - [@your-name](mailto:your-email@example.com)

项目链接: [https://github.com/your-org/asset-management](https://github.com/your-org/asset-management)
