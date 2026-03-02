# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Atlas** - GPU算力公司的IT资产全生命周期管理系统，支持设备采购、入库、库存管理、出库、DCIM可视化、设备生命周期管理等功能。

Atlas（阿特拉斯）源自希腊神话中的擎天神，寓意全面掌控和管理企业的IT资产与基础设施。

### Business Scope

```
设备生命周期管理：
├── 采购管理：采购需求、订单跟踪、供应商管理
├── 入库管理：扫码入库、批量入库、质检
├── 库存管理：实时库存、库位管理、盘点、预警
├── 出库管理：领用申请、审批、扫码出库
├── DCIM：IDC管理、机柜可视化、U位管理、上架规划
├── 设备分配：设备分配、调拨
├── 维修管理：故障报修、维修跟踪、备件更换
└── 报废处理：报废申请、审批、资产处置
```

### Scale Indicators

- Device types: GPU servers, CPU servers, switches, storage systems, various components
- Device quantity: Tens of thousands
- Component quantity: 100K-1M level
- Monthly in/out: Hundreds to thousands of transactions
- IDC count: Several to dozens
- User roles: Procurement, warehouse, operations, IDC engineers, finance, management

## Technology Stack

**Backend:**
- Go 1.21+ with Fiber v2 web framework
- Ent ORM (code generation + type safety)
- PostgreSQL 15+ database
- Redis 7+ for caching
- RabbitMQ for message queue
- MinIO for file storage

**Frontend (planned):**
- Vite 5+ + React 18 + TypeScript 5
- Ant Design 5 UI library
- Zustand for state management
- TanStack Query (React Query v5)

## Development Commands

### Backend (from `backend/` directory)

```bash
make install      # Install Go dependencies
make generate     # Generate Ent code from schemas
make build        # Build the application
make run          # Run the API server
make test         # Run tests with coverage
make lint         # Run golangci-lint
make clean        # Clean build artifacts and logs
```

### Docker Infrastructure (from project root)

```bash
docker-compose up -d      # Start all services (PostgreSQL, Redis, MinIO, RabbitMQ)
docker-compose down       # Stop all services
docker-compose ps         # Check service status
docker-compose logs -f    # View logs
```

### Service Access

- PostgreSQL: localhost:5432 (admin/admin123)
- Redis: localhost:6379 (password: redis123)
- MinIO Console: http://localhost:9001 (minioadmin/minioadmin123)
- RabbitMQ Management: http://localhost:15672 (admin/admin123)
- Adminer (DB tool): http://localhost:8081

## Architecture

### Hybrid Architecture Philosophy

**Core principle**: Use tools for standardized parts, write code for complex parts.

```
Frontend (React + TypeScript)
    ↓
API Gateway (Nginx)
    ↓
Go Application (Fiber + Ent)
├── Auto-generated layer (Ent)
│   ├── Data models (Schema definitions)
│   ├── CRUD operations (type-safe)
│   └── Query builders
├── Hand-written business layer
│   ├── Approval workflow engine
│   ├── Inventory calculation service
│   ├── Device lifecycle management
│   └── DCIM business logic
└── API layer (Fiber)
    ↓
PostgreSQL + Redis + MinIO + RabbitMQ
```

### Why This Architecture?

1. **Development efficiency +60%**
   - Ent auto-generates CRUD code
   - Type safety reduces bugs
   - Clear layered architecture
   - Code generators accelerate development

2. **Complete flexibility**
   - Complex business logic fully hand-written
   - Large optimization space
   - No vendor lock-in
   - Can implement any requirement

3. **Suitable for large-scale data**
   - PostgreSQL supports millions of records
   - Ent query optimization
   - Redis cache acceleration
   - Fine-grained query control

4. **AI Agent friendly**
   - Declarative schemas, easy for AI to understand
   - Generated code follows consistent patterns
   - Clear project structure
   - Type safety reduces errors

### Backend Project Structure

```
backend/
├── cmd/api/main.go              # API server entry point
├── ent/
│   ├── schema/                  # Schema definitions (hand-written)
│   ├── generate.go              # Code generation entry
│   └── [generated files]        # Ent auto-generated code
├── internal/
│   ├── handler/                 # HTTP handlers
│   ├── service/                 # Business logic layer
│   ├── middleware/              # Middleware (auth, logger, error)
│   ├── dto/                     # Data transfer objects
│   └── router/                  # Route configuration
├── pkg/                         # Reusable packages
│   ├── config/                  # Configuration management (Viper)
│   ├── database/                # Database connection
│   ├── redis/                   # Redis client
│   ├── logger/                  # Structured logging (Zap)
│   ├── jwt/                     # JWT utilities
│   └── utils/                   # Utility functions
└── config/                      # Configuration files
```

## Development Workflow

### Adding New Entities

1. Create schema file in `backend/ent/schema/`
2. Run `make generate` to generate Ent code
3. Create corresponding Handler and Service
4. Register routes in Router
5. Write tests

### Ent Schema Pattern

All schemas follow this pattern:

```go
type EntityName struct {
    ent.Schema
}

func (EntityName) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").Comment("描述"),
        field.Enum("status").Values("active", "inactive").Default("active"),
        field.Time("created_at").Default(time.Now),
    }
}

func (EntityName) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("related", RelatedEntity.Type),
    }
}
```

### API Response Format

All APIs use unified response format:

```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "timestamp": 1234567890
}
```

Error response:

```json
{
  "code": 40001,
  "message": "错误描述",
  "errors": [{"field": "field_name", "message": "详细错误"}],
  "timestamp": 1234567890
}
```

## Core Data Models

### Complete Entity List

| Entity | Description | Priority |
|--------|-------------|----------|
| User | 用户 | P0 |
| Role | 角色 | P0 |
| Permission | 权限 | P0 |
| AssetType | 资产类型 | P0 |
| Asset | 资产 | P0 |
| Location | 库位 | P0 |
| Warehouse | 仓库 | P0 |
| InventoryRecord | 库存记录 | P0 |
| Supplier | 供应商 | P1 |
| PurchaseOrder | 采购订单 | P1 |
| OrderItem | 订单明细 | P1 |
| DataCenter | 数据中心 | P1 |
| Room | 机房 | P1 |
| Rack | 机柜 | P1 |
| RackUnit | U位 | P1 |
| Approval | 审批 | P1 |
| ApprovalNode | 审批节点 | P1 |
| ApprovalFlow | 审批流程模板 | P1 |
| MaintenanceRecord | 维修记录 | P2 |
| TransferRecord | 调拨记录 | P2 |
| StockCheck | 盘点记录 | P2 |
| Notification | 通知 | P2 |
| OperationLog | 操作日志 | P2 |

### Key Entity Examples

#### Asset (资产)

```go
type Asset struct {
    ent.Schema
}

func (Asset) Fields() []ent.Field {
    return []ent.Field{
        field.String("asset_no").Unique().Comment("资产编号"),
        field.String("asset_type").Comment("资产类型：gpu_server, cpu_server, switch, storage, component"),
        field.String("brand").Optional().Comment("品牌"),
        field.String("model").Optional().Comment("型号"),
        field.String("sn").Unique().Optional().Comment("序列号"),
        field.Enum("status").Values(
            "in_stock",      // 在库
            "deployed",      // 已部署
            "maintenance",   // 维修中
            "retired",       // 已报废
        ).Default("in_stock").Comment("状态"),
        field.JSON("specs", map[string]interface{}{}).Optional().Comment("规格参数"),
        field.Time("purchase_date").Optional().Comment("采购日期"),
        field.Time("warranty_expire_date").Optional().Comment("保修到期日"),
        field.Float("purchase_price").Optional().Comment("采购价格"),
        field.Text("description").Optional().Comment("描述"),
    }
}

func (Asset) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("asset_type_ref", AssetType.Type).Ref("assets").Unique(),
        edge.To("location", Location.Type).Unique(),
        edge.To("rack_unit", RackUnit.Type).Unique(),
        edge.To("inventory_records", InventoryRecord.Type),
        edge.To("maintenance_records", MaintenanceRecord.Type),
    }
}
```

#### PurchaseOrder (采购订单)

```go
type PurchaseOrder struct {
    ent.Schema
}

func (PurchaseOrder) Fields() []ent.Field {
    return []ent.Field{
        field.String("order_no").Unique().Comment("订单编号"),
        field.Float("total_amount").Comment("总金额"),
        field.Enum("status").Values(
            "draft",      // 草稿
            "pending",    // 待审批
            "approved",   // 已批准
            "rejected",   // 已拒绝
            "ordered",    // 已下单
            "received",   // 已收货
            "completed",  // 已完成
            "cancelled",  // 已取消
        ).Default("draft").Comment("状态"),
        field.Time("order_date").Optional().Comment("下单日期"),
        field.Time("expected_date").Optional().Comment("预计到货日期"),
        field.Text("note").Optional().Comment("备注"),
    }
}

func (PurchaseOrder) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("supplier", Supplier.Type).Unique(),
        edge.To("items", OrderItem.Type),
        edge.To("creator", User.Type).Unique(),
        edge.To("approval", Approval.Type).Unique(),
    }
}
```

#### Rack (机柜)

```go
type Rack struct {
    ent.Schema
}

func (Rack) Fields() []ent.Field {
    return []ent.Field{
        field.String("rack_no").Comment("机柜编号"),
        field.Int("total_units").Default(42).Comment("总U数"),
        field.Float("power_capacity").Optional().Comment("电力容量(kW)"),
        field.Float("power_used").Default(0).Comment("已用电力(kW)"),
        field.Int("weight_capacity").Optional().Comment("承重(kg)"),
        field.Enum("status").Values(
            "available",   // 可用
            "full",        // 已满
            "maintenance", // 维护中
            "reserved",    // 已预留
        ).Default("available").Comment("状态"),
        field.JSON("position", map[string]interface{}{}).Optional().Comment("位置坐标"),
    }
}

func (Rack) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("room", Room.Type).Ref("racks").Unique(),
        edge.To("units", RackUnit.Type),
    }
}
```

## Key Implementation Notes

### Ent Code Generation

- Always run `make generate` after modifying schemas
- Generated code is in `ent/` directory (do not edit manually)
- Migrations are auto-generated by Ent

### Business Logic Layer

Complex business logic should be in `internal/service/`:
- Approval workflow engine (state machine)
- Inventory calculations (transactions, locking)
- Device lifecycle management (state transitions)
- DCIM capacity planning

### Middleware Stack

Request flow: Logger → Auth → Error Handler → Business Logic

### Configuration

- Uses Viper for configuration management
- Config file: `backend/config/config.yaml`
- Environment-specific configs: `config.dev.yaml`, `config.prod.yaml`

## Testing

```bash
make test                    # Run all tests
go test -v ./...            # Run tests with verbose output
go test -v ./internal/service/...  # Test specific package
```

## API Design

### RESTful Principles

- Unified response format (see above)
- Version control: `/api/v1/`
- Complete error handling
- OpenAPI documentation

### Core API Endpoints

#### Asset Management

```
GET    /api/v1/assets              获取资产列表
POST   /api/v1/assets              创建资产
GET    /api/v1/assets/:id          获取资产详情
PUT    /api/v1/assets/:id          更新资产
DELETE /api/v1/assets/:id          删除资产
POST   /api/v1/assets/batch        批量导入资产
GET    /api/v1/assets/:id/history  获取资产历史记录
POST   /api/v1/assets/:id/transfer 资产调拨
```

#### Inventory Management

```
GET    /api/v1/inventory/stock           获取库存统计
POST   /api/v1/inventory/in              入库
POST   /api/v1/inventory/out             出库
POST   /api/v1/inventory/transfer        调拨
POST   /api/v1/inventory/check           盘点
GET    /api/v1/inventory/records         获取库存记录
GET    /api/v1/inventory/alerts          获取库存预警
```

#### Procurement Management

```
GET    /api/v1/purchase/orders           获取采购订单列表
POST   /api/v1/purchase/orders           创建采购订单
GET    /api/v1/purchase/orders/:id       获取订单详情
PUT    /api/v1/purchase/orders/:id       更新订单
POST   /api/v1/purchase/orders/:id/submit 提交审批
GET    /api/v1/purchase/suppliers        获取供应商列表
POST   /api/v1/purchase/suppliers        创建供应商
```

#### DCIM

```
GET    /api/v1/dcim/datacenters          获取IDC列表
POST   /api/v1/dcim/datacenters          创建IDC
GET    /api/v1/dcim/rooms                获取机房列表
GET    /api/v1/dcim/racks                获取机柜列表
GET    /api/v1/dcim/racks/:id            获取机柜详情
GET    /api/v1/dcim/racks/:id/units      获取机柜U位占用情况
POST   /api/v1/dcim/racks/:id/mount      设备上架
POST   /api/v1/dcim/racks/:id/unmount    设备下架
GET    /api/v1/dcim/capacity             获取容量统计
```

#### Approval

```
GET    /api/v1/approvals                 获取审批列表
GET    /api/v1/approvals/:id             获取审批详情
POST   /api/v1/approvals/:id/approve     批准
POST   /api/v1/approvals/:id/reject      拒绝
GET    /api/v1/approvals/pending         获取待审批列表
GET    /api/v1/approvals/flows           获取审批流程模板
```

#### Reports

```
GET    /api/v1/reports/inventory         库存报表
GET    /api/v1/reports/asset-value       资产价值报表
GET    /api/v1/reports/utilization       设备利用率报表
GET    /api/v1/reports/purchase          采购统计报表
POST   /api/v1/reports/export            导出报表
```

#### WebSocket

```
WS     /api/v1/ws/notifications          实时通知
WS     /api/v1/ws/dcim/racks/:id         机柜实时状态
```

## Current Development Stage

**Stage 0: Infrastructure Setup (Week 1-2)** - 80% complete

Completed:
- ✅ Technical architecture design
- ✅ Docker environment configuration
- ✅ Backend scaffolding
- ✅ Basic middleware (auth, logger, error)
- ✅ Configuration management
- ✅ Database and Redis connections

Next steps:
- Define core data model schemas (P0 entities)
- Setup frontend scaffolding
- Implement first feature module (inventory management)

## Development Plan (4 Stages, 13-20 Weeks)

### Stage 0: Infrastructure Setup (Week 1-2) - CURRENT

**Backend:**
- [x] Initialize Go project with Fiber and Ent
- [x] Configure database and Redis connections
- [x] Implement basic middleware (auth, logger, error)
- [x] Implement JWT authentication
- [x] Write Makefile
- [ ] Define core schemas (15 entities)
- [ ] Run code generation

**Frontend:**
- [ ] Initialize Vite + React + TypeScript
- [ ] Configure Ant Design
- [ ] Implement basic layout and routing
- [ ] Implement login page
- [ ] Configure Zustand state management

**Infrastructure:**
- [x] Write docker-compose.yml
- [x] Configure PostgreSQL, Redis, MinIO, RabbitMQ
- [x] Write README and development docs

**Deliverables:**
- Runnable backend API (health check, login)
- Runnable frontend (login page)
- Complete development environment (Docker Compose)

### Stage 1: Core Features (Week 3-10)

**Week 3-4: Inventory Management Module**
- Backend: Asset CRUD API, in/out stock API, inventory query/stats
- Frontend: Asset list/detail pages, in/out stock forms, inventory query
- Unit tests

**Week 5-7: DCIM Subsystem**
- Backend: IDC/Room/Rack CRUD, rack unit management, mount/unmount API, capacity stats
- Backend: DCIM visualization service (independent), WebSocket real-time push
- Frontend: IDC/Room/Rack management pages, 2D rack visualization, mount planning tool

**Week 8-10: Procurement and Approval Workflow**
- Backend: Purchase order CRUD, supplier management, approval workflow engine, notification service
- Backend: WeChat Work/Feishu integration
- Frontend: Purchase order list/form, supplier management, approval list/detail, approval flow config

### Stage 2: Supplementary Features (Week 11-16)

**Week 11-12: Device Lifecycle Management**
- Device allocation, transfer, maintenance management, retirement process
- State change tracking

**Week 13-14: Reports and Analytics**
- Report service (independent)
- Inventory, asset value, utilization, procurement reports
- Data dashboard

**Week 15-16: Mobile**
- Mobile H5 pages
- Scan in/out stock, photo upload, signature
- PWA offline support

### Stage 3: Optimization and Testing (Week 17-20)

**Week 17-18: Performance Optimization**
- Database index optimization, Redis cache strategy
- API performance testing, frontend optimization
- Large data volume testing

**Week 19: Security and Permissions**
- RBAC permission system refinement
- Data permissions (row-level)
- Operation log audit, security hardening

**Week 20: Integration Testing and Deployment**
- End-to-end testing, user acceptance testing
- Deployment and operations documentation
- Training materials

## Important Conventions

1. **Schema-first development**: Define Ent schemas before implementing features
2. **Layered architecture**: Handler → Service → Repository (Ent)
3. **Type safety**: Leverage Ent's type-safe query builders
4. **Error handling**: Use unified error response format
5. **Logging**: Use structured logging with Zap
6. **Authentication**: JWT-based with middleware
7. **API versioning**: All APIs under `/api/v1/`

## AI Agent Collaboration Strategy

### Feature Module Development Flow

```
1. Requirements Clarification (AI + You)
   - Confirm feature scope
   - Confirm data models
   - Confirm API interfaces

2. Schema Definition (AI leads, You review)
   - AI generates Schema code
   - You review and adjust
   - Run code generation

3. Backend Development (AI leads)
   - AI generates Handler code
   - AI implements business logic
   - AI writes unit tests
   - You review and test

4. Frontend Development (AI leads)
   - AI generates page framework
   - AI implements interaction logic
   - AI integrates APIs
   - You review and test

5. Integration Testing (AI + You)
   - End-to-end testing
   - Bug fixes
   - Performance testing

6. Deployment (You lead, AI assists)
   - Build Docker images
   - Update deployment config
   - Execute deployment
```

### Best Practices for AI Collaboration

- **Incremental development**: Small commits, always compilable and testable
- **Learn from existing code**: Study 3 similar features before implementing new ones
- **Identify common patterns**: Follow project conventions
- **Staged implementation**: Break complex work into 3-5 stages
- **Quality standards**: Every commit must compile, pass tests, include new tests

## Security Design

### Authentication & Authorization

- JWT Token authentication
- RBAC permission model
- Data permissions (row-level)
- API access control

### Data Security

- Sensitive data encryption
- HTTPS transmission
- SQL injection protection
- XSS protection
- CSRF protection

### Audit

- Operation log recording
- Login logs
- Data change tracking

## Documentation

- [ARCHITECTURE.md](./ARCHITECTURE.md) - Complete technical architecture (1000+ lines)
- [PROJECT_SUMMARY.md](./PROJECT_SUMMARY.md) - Project progress summary
- [QUICKSTART.md](./QUICKSTART.md) - Quick start guide
- [backend/README.md](./backend/README.md) - Backend development guide
