# Atlas - 快速开始指南

## 当前项目状态

✅ 技术架构文档已完成
✅ Docker开发环境已配置
✅ 后端项目脚手架已搭建
✅ 数据模型Schema已定义（21个实体）
✅ Ent代码已生成（76个文件，22,637行代码）
✅ Excel业务文件已分析
✅ Git仓库已初始化并推送到GitHub
⏳ 数据库迁移待创建
⏳ 前端项目待搭建

## 立即开始开发

### 1. 启动基础设施服务

```bash
# 在项目根目录
docker-compose up -d

# 验证服务状态
docker-compose ps

# 应该看到以下服务运行中：
# - asset-mgmt-postgres (PostgreSQL)
# - asset-mgmt-redis (Redis)
# - asset-mgmt-minio (MinIO)
# - asset-mgmt-rabbitmq (RabbitMQ)
# - asset-mgmt-adminer (数据库管理工具)
```

### 2. 访问服务

| 服务 | 地址 | 用户名 | 密码 |
|------|------|--------|------|
| PostgreSQL | localhost:5432 | admin | admin123 |
| Redis | localhost:6379 | - | redis123 |
| MinIO Console | http://localhost:9001 | minioadmin | minioadmin123 |
| RabbitMQ Management | http://localhost:15672 | admin | admin123 |
| Adminer | http://localhost:8081 | - | - |

### 3. 后端开发准备

```bash
cd backend

# 安装Go依赖
make install

# 此时会看到一些错误，因为还没有定义Schema
# 这是正常的，下一步我们将定义Schema
```

## 下一步：数据库迁移和初始数据

数据模型已经定义完成，现在需要：

1. **创建数据库迁移脚本**
2. **初始化数据库表结构**
3. **创建初始数据（种子数据）**

### 数据模型概览

已定义21个实体：
- **P0**: User, Role, Permission, Warehouse, Location, AssetType, Asset, InventoryRecord
- **P1**: Supplier, PurchaseOrder, OrderItem, DataCenter, Room, Rack, RackUnit, Approval, ApprovalNode, NetworkConnection, IPAddress, RepairVendor, RepairTicket

详见 `docs/DATA_MODEL.md` 获取完整的数据模型设计文档。

### 示例：查看生成的代码

```bash
cd backend

# 查看生成的实体
ls -la ent/

# 查看Asset实体
cat ent/asset/asset.go
```

### 运行后端服务

```bash
cd backend

# 运行服务
make run

# 服务将在 http://localhost:8080 启动
```

### 测试API

```bash
# 健康检查
curl http://localhost:8080/health

# 应该返回：
# {"status":"ok"}
```

## 开发工作流

### 添加新的实体

1. 在 `backend/ent/schema/` 创建新的schema文件
2. 运行 `make generate` 生成代码
3. 创建对应的Handler、Service
4. 在Router中注册路由
5. 测试API

### 常用命令

```bash
# 后端
cd backend
make help          # 查看所有命令
make install       # 安装依赖
make generate      # 生成Ent代码
make run           # 运行服务
make test          # 运行测试
make lint          # 代码检查

# Docker
docker-compose up -d      # 启动服务
docker-compose down       # 停止服务
docker-compose ps         # 查看状态
docker-compose logs -f    # 查看日志
```

## 文档导航

- [README.md](./README.md) - 项目总览
- [ARCHITECTURE.md](./ARCHITECTURE.md) - 完整技术架构文档
- [PROJECT_SUMMARY.md](./PROJECT_SUMMARY.md) - 项目进度总结
- [backend/README.md](./backend/README.md) - 后端开发文档

## 需要帮助？

如果遇到问题：

1. 检查Docker服务是否正常运行：`docker-compose ps`
2. 查看服务日志：`docker-compose logs -f [service-name]`
3. 检查配置文件：`backend/config/config.yaml`
4. 查看应用日志：`backend/logs/app.log`

## AI Agent协作提示

当你准备好继续开发时，告诉我：

- "定义User Schema" - 我会帮你创建User实体
- "定义所有P0 Schema" - 我会创建所有优先级P0的实体
- "搭建前端项目" - 我会开始搭建React前端
- "实现资产管理API" - 我会实现第一个完整的功能模块

让我们一起高效地完成这个项目！🚀
