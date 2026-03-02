# Atlas 快速参考

本文档提供Atlas项目的快速查找信息。

## 📁 项目结构

```
atlas/
├── backend/              # 后端Go项目
│   ├── cmd/             # 应用入口
│   ├── ent/             # Ent ORM（Schema + 生成代码）
│   ├── internal/        # 内部代码
│   ├── pkg/             # 公共库
│   └── config/          # 配置文件
├── docs/                # 文档
│   ├── references/      # Excel参考文件
│   └── DATA_MODEL.md    # 数据模型文档
├── scripts/             # 脚本
└── docker-compose.yml   # Docker配置
```

## 🔗 重要链接

- **GitHub仓库**: https://github.com/zhong/atlas
- **本地路径**: `/Users/chenzhong/Developer/hello`

## 📚 文档导航

| 文档 | 用途 | 路径 |
|------|------|------|
| README.md | 项目总览 | `./README.md` |
| ARCHITECTURE.md | 完整技术架构（1000+行） | `./ARCHITECTURE.md` |
| CLAUDE.md | Claude Code工作指南 | `./CLAUDE.md` |
| PROJECT_SUMMARY.md | 项目进度总结 | `./PROJECT_SUMMARY.md` |
| QUICKSTART.md | 快速开始指南 | `./QUICKSTART.md` |
| CHANGELOG.md | 开发日志 | `./CHANGELOG.md` |
| DATA_MODEL.md | 数据模型设计 | `./docs/DATA_MODEL.md` |

## 🚀 常用命令

### Docker服务

```bash
# 启动所有服务
docker-compose up -d

# 停止所有服务
docker-compose down

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f [service-name]
```

### 后端开发

```bash
cd backend

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

# 清理
make clean
```

### Git操作

```bash
# 查看状态
git status

# 提交更改
git add .
git commit -m "message"
git push

# 查看日志
git log --oneline
```

## 🗄️ 数据库连接

| 服务 | 地址 | 用户名 | 密码 |
|------|------|--------|------|
| PostgreSQL | localhost:5432 | admin | admin123 |
| Redis | localhost:6379 | - | redis123 |
| MinIO Console | http://localhost:9001 | minioadmin | minioadmin123 |
| RabbitMQ Management | http://localhost:15672 | admin | admin123 |
| Adminer | http://localhost:8081 | - | - |

## 📊 数据模型

### 21个实体

**P0（核心）**:
- User, Role, Permission
- Warehouse, Location
- AssetType, Asset
- InventoryRecord

**P1（扩展）**:
- Supplier, PurchaseOrder, OrderItem
- DataCenter, Room, Rack, RackUnit
- Approval, ApprovalNode
- NetworkConnection, IPAddress
- RepairVendor, RepairTicket

详见 `docs/DATA_MODEL.md`

## 🎯 当前状态

**版本**: v0.1.0-alpha
**阶段**: 阶段0完成（100%）

**已完成**:
- ✅ 技术架构设计
- ✅ Docker环境
- ✅ 后端脚手架
- ✅ 数据模型（21个实体）
- ✅ 代码生成（76文件，22,637行）
- ✅ Excel分析（5文件）

**下一步**:
- 🔄 数据库迁移
- 🔄 前端脚手架
- 🔄 库存管理模块

## 📈 统计数据

- **代码量**: 34,000+行
- **文件数**: 115+个
- **提交数**: 3次
- **开发时间**: 1天
- **实体数**: 21个
- **生成代码**: 76文件，22,637行

## 🔍 快速查找

### 查找Schema定义
```bash
ls backend/ent/schema/
cat backend/ent/schema/asset.go
```

### 查看生成的代码
```bash
ls backend/ent/
cat backend/ent/asset/asset.go
```

### 查看Excel参考文件
```bash
ls docs/references/
```

### 查看配置
```bash
cat backend/config/config.yaml
cat docker-compose.yml
```

## 🛠️ 故障排查

### Docker服务无法启动
```bash
# 检查端口占用
lsof -i :5432  # PostgreSQL
lsof -i :6379  # Redis
lsof -i :9000  # MinIO

# 清理并重启
docker-compose down -v
docker-compose up -d
```

### 后端编译错误
```bash
cd backend

# 清理并重新生成
make clean
make install
make generate
```

### 查看日志
```bash
# Docker服务日志
docker-compose logs -f postgres
docker-compose logs -f redis

# 应用日志
tail -f backend/logs/app.log
```

## 📞 获取帮助

1. 查看文档：`README.md`, `ARCHITECTURE.md`, `CLAUDE.md`
2. 查看开发日志：`CHANGELOG.md`
3. 查看数据模型：`docs/DATA_MODEL.md`
4. 查看Excel分析：`docs/references/`

---

**最后更新**: 2026-03-02
**版本**: v1.0
