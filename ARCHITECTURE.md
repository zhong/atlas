# Atlas - IT资产管理系统技术架构文档

## 1. 项目概述

### 1.1 项目背景

本项目为GPU算力公司（GPU云基础设施公司）开发的IT资产全生命周期管理系统。公司拥有多个自主可控的IDC，管理大量GPU服务器、网络设备、存储设备及硬件备件。

### 1.2 核心痛点

- 当前管理方式原始：人工流程 + 邮件/即时通讯 + Excel/在线表格
- 数据不准确、不及时
- 效率低下
- 缺乏统一的管理平台

### 1.3 项目目标

- 实现设备和备件的全生命周期管理
- 提升运营效率
- 保证数据准确性和及时性
- 支持多IDC、多角色协作
- 提供可视化的DCIM能力

### 1.4 业务范围

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

### 1.5 规模指标

- 设备类型：GPU服务器、CPU服务器、交换机、存储系统、各类备件
- 设备数量：几万级别
- 备件数量：10万-百万级别
- 月入库/出库：几百到几千次
- IDC数量：几个到几十个
- 用户角色：采购、仓库、运维、IDC工程师、财务、管理层

## 2. 技术架构

### 2.1 架构选型：现代化混合架构

**核心理念**：标准化的地方用工具，复杂的地方手写

```
┌─────────────────────────────────────────────────────────────┐
│                         前端层                               │
│  Web管理后台 + 移动端H5 + 数据大屏                           │
│  React 18 + TypeScript + Ant Design                         │
└─────────────────────────────────────────────────────────────┘
                              ↓ HTTPS/WebSocket
┌─────────────────────────────────────────────────────────────┐
│                      API网关层                               │
│              Nginx (反向代理 + 负载均衡)                      │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                    核心Go应用                                │
│  ┌───────────────────────────────────────────────────┐     │
│  │ 自动生成层 (Ent框架)                               │     │
│  │ ├── 数据模型 (Schema定义)                          │     │
│  │ ├── CRUD操作 (自动生成，类型安全)                  │     │
│  │ ├── 查询构建器 (支持复杂查询)                      │     │
│  │ └── 数据库迁移 (自动生成)                          │     │
│  │                                                    │     │
│  │ 手写业务层                                          │     │
│  │ ├── 审批流程引擎 (工作流、状态机)                  │     │
│  │ ├── 库存计算服务 (事务处理、锁机制)                │     │
│  │ ├── 设备生命周期管理 (状态变更、联动)              │     │
│  │ ├── DCIM业务逻辑 (容量计算、上架规划)              │     │
│  │ └── 集成服务 (企业微信/飞书、第三方系统)           │     │
│  │                                                    │     │
│  │ API层 (Fiber框架)                                  │     │
│  │ ├── 标准CRUD API (可生成)                         │     │
│  │ ├── 复杂业务API (手写)                            │     │
│  │ ├── WebSocket (实时通知)                          │     │
│  │ └── 中间件 (认证、日志、限流)                      │     │
│  └───────────────────────────────────────────────────┘     │
│                                                             │
│  技术栈：Go 1.21+ + Fiber v2 + Ent                          │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌──────────┬──────────┬──────────┬──────────┬───────────────┐
│PostgreSQL│  Redis   │  MinIO   │ RabbitMQ │  独立服务      │
│ (主数据库)│  (缓存)  │(文件存储)│(消息队列)│ (DCIM/报表)   │
└──────────┴──────────┴──────────┴──────────┴───────────────┘
```

### 2.2 技术栈详细说明

#### 后端技术栈

| 组件 | 技术选型 | 版本 | 说明 |
|------|---------|------|------|
| 语言 | Go | 1.21+ | 高性能、编译部署方便 |
| Web框架 | Fiber | v2 | 高性能、Express风格API |
| ORM | Ent | latest | 代码生成、类型安全 |
| 数据库 | PostgreSQL | 15+ | 企业级、支持大数据量 |
| 缓存 | Redis | 7+ | 高性能缓存、分布式锁 |
| 消息队列 | RabbitMQ | 3.12+ | 可靠的消息传递 |
| 文件存储 | MinIO | latest | S3兼容、自托管 |
| 配置管理 | Viper | latest | 灵活的配置管理 |
| 日志 | Zap | latest | 高性能结构化日志 |
| 验证 | validator | v10 | 数据验证 |
| JWT | golang-jwt | v5 | 认证授权 |

#### 前端技术栈

| 组件 | 技术选型 | 版本 | 说明 |
|------|---------|------|------|
| 构建工具 | Vite | 5+ | 快速的开发构建 |
| 框架 | React | 18+ | 主流前端框架 |
| 语言 | TypeScript | 5+ | 类型安全 |
| UI库 | Ant Design | 5+ | 企业级UI组件 |
| 路由 | React Router | v6 | 标准路由方案 |
| 状态管理 | Zustand | latest | 轻量级状态管理 |
| 请求库 | Axios | latest | HTTP客户端 |
| 数据获取 | TanStack Query | v5 | 服务端状态管理 |
| 表单 | React Hook Form | latest | 高性能表单 |
| 表单验证 | Zod | latest | TypeScript优先验证 |
| 图表 | ECharts | 5+ | 数据可视化 |
| 图形 | AntV G6 | latest | 图可视化（DCIM） |
| 移动端 | Ant Design Mobile | latest | 移动端组件 |

#### 独立服务

| 服务 | 技术栈 | 说明 |
|------|--------|------|
| DCIM可视化服务 | Go + WebSocket | 实时机柜状态、3D渲染数据 |
| 报表服务 | Go + 定时任务 | 数据聚合、导出 |
| 通知服务 | Go + RabbitMQ | 企业微信/飞书通知 |

### 2.3 为什么选择这个架构？

#### 优势分析

```
1. 开发效率提升 60%
   ├── Ent自动生成CRUD代码
   ├── 类型安全减少bug
   ├── 清晰的分层架构
   └── 代码生成器加速开发

2. 完全的灵活性
   ├── 复杂业务逻辑完全手写
   ├── 性能优化空间大
   ├── 无供应商锁定
   └── 可以实现任何需求

3. 适合大规模数据
   ├── PostgreSQL支持百万级数据
   ├── Ent查询优化
   ├── Redis缓存加速
   └── 可以精细控制查询

4. 非常适合AI Agent协作
   ├── 声明式Schema，AI易理解
   ├── 生成代码模式统一
   ├── 清晰的项目结构
   └── 类型安全减少错误

5. 长期可维护
   ├── 代码量少（比传统减少56%）
   ├── 结构清晰
   ├── 测试友好
   └── 重构容易
```

## 3. 数据模型设计

### 3.1 核心实体关系图

```
┌─────────────┐      ┌──────────────┐      ┌─────────────┐
│  Supplier   │──────│PurchaseOrder │──────│   OrderItem │
│  供应商      │ 1:N  │  采购订单     │ 1:N  │  订单明细   │
└─────────────┘      └──────────────┘      └─────────────┘
                            │
                            │ 1:N
                            ↓
┌─────────────┐      ┌──────────────┐      ┌─────────────┐
│   Asset     │──────│InventoryRecord│─────│  Location   │
│   资产      │ 1:N  │  库存记录     │ N:1  │  库位       │
└─────────────┘      └──────────────┘      └─────────────┘
      │                                           │
      │ N:1                                       │ N:1
      ↓                                           ↓
┌─────────────┐                           ┌─────────────┐
│ AssetType   │                           │  Warehouse  │
│ 资产类型     │                           │   仓库      │
└─────────────┘                           └─────────────┘

┌─────────────┐      ┌──────────────┐      ┌─────────────┐
│ DataCenter  │──────│    Room      │──────│    Rack     │
│   IDC       │ 1:N  │   机房       │ 1:N  │   机柜      │
└─────────────┘      └──────────────┘      └─────────────┘
                                                  │
                                                  │ 1:N
                                                  ↓
                                           ┌─────────────┐
                                           │ RackUnit    │
                                           │   U位       │
                                           └─────────────┘
                                                  │
                                                  │ 1:1
                                                  ↓
                                           ┌─────────────┐
                                           │   Asset     │
                                           │   资产      │
                                           └─────────────┘

┌─────────────┐      ┌──────────────┐      ┌─────────────┐
│   User      │──────│  Approval    │──────│ApprovalNode │
│   用户      │ 1:N  │   审批       │ 1:N  │  审批节点   │
└─────────────┘      └──────────────┘      └─────────────┘
```

### 3.2 核心实体定义（Ent Schema）

以下是主要实体的Schema定义，完整代码将在项目中实现。

#### 资产 (Asset)

```go
// ent/schema/asset.go
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

# IT资产管理系统 - 技术架构文档（续）

## 3. 数据模型设计（续）

### 3.2 核心实体定义（续）

#### 库存记录 (InventoryRecord)

```go
type InventoryRecord struct {
    ent.Schema
}

func (InventoryRecord) Fields() []ent.Field {
    return []ent.Field{
        field.Enum("record_type").Values(
            "in",        // 入库
            "out",       // 出库
            "transfer",  // 调拨
            "adjust",    // 调整
            "check",     // 盘点
        ).Comment("记录类型"),
        field.Int("quantity").Default(1).Comment("数量"),
        field.String("reason").Optional().Comment("原因"),
        field.Text("note").Optional().Comment("备注"),
    }
}

func (InventoryRecord) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("asset", Asset.Type).Ref("inventory_records").Unique(),
        edge.To("from_location", Location.Type).Unique(),
        edge.To("to_location", Location.Type).Unique(),
        edge.To("operator", User.Type).Unique(),
    }
}
```

#### 采购订单 (PurchaseOrder)

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

#### IDC数据中心 (DataCenter)

```go
type DataCenter struct {
    ent.Schema
}

func (DataCenter) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").Unique().Comment("IDC名称"),
        field.String("code").Unique().Comment("IDC编码"),
        field.String("location").Comment("地理位置"),
        field.String("address").Optional().Comment("详细地址"),
        field.String("contact").Optional().Comment("联系人"),
        field.String("phone").Optional().Comment("联系电话"),
        field.Text("description").Optional().Comment("描述"),
    }
}

func (DataCenter) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("rooms", Room.Type),
    }
}
```

#### 机柜 (Rack)

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

#### 审批 (Approval)

```go
type Approval struct {
    ent.Schema
}

func (Approval) Fields() []ent.Field {
    return []ent.Field{
        field.String("approval_no").Unique().Comment("审批编号"),
        field.String("entity_type").Comment("实体类型：purchase_order, asset_transfer等"),
        field.Int("entity_id").Comment("实体ID"),
        field.Enum("status").Values(
            "pending",    // 待审批
            "approved",   // 已批准
            "rejected",   // 已拒绝
            "cancelled",  // 已取消
        ).Default("pending").Comment("状态"),
        field.Int("current_node").Default(0).Comment("当前节点"),
        field.Text("note").Optional().Comment("备注"),
    }
}

func (Approval) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("nodes", ApprovalNode.Type),
        edge.To("initiator", User.Type).Unique(),
    }
}
```

### 3.3 完整实体列表

| 实体 | 说明 | 优先级 |
|------|------|--------|
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

## 4. API设计

### 4.1 API设计原则

- RESTful风格
- 统一的响应格式
- 版本控制（/api/v1）
- 完善的错误处理
- OpenAPI文档

### 4.2 统一响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "timestamp": 1234567890
}
```

错误响应：
```json
{
  "code": 40001,
  "message": "资产编号已存在",
  "errors": [
    {
      "field": "asset_no",
      "message": "该资产编号已被使用"
    }
  ],
  "timestamp": 1234567890
}
```

### 4.3 核心API列表

#### 资产管理 API

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

#### 库存管理 API

```
GET    /api/v1/inventory/stock           获取库存统计
POST   /api/v1/inventory/in              入库
POST   /api/v1/inventory/out             出库
POST   /api/v1/inventory/transfer        调拨
POST   /api/v1/inventory/check           盘点
GET    /api/v1/inventory/records         获取库存记录
GET    /api/v1/inventory/alerts          获取库存预警
```

#### 采购管理 API

```
GET    /api/v1/purchase/orders           获取采购订单列表
POST   /api/v1/purchase/orders           创建采购订单
GET    /api/v1/purchase/orders/:id       获取订单详情
PUT    /api/v1/purchase/orders/:id       更新订单
POST   /api/v1/purchase/orders/:id/submit 提交审批
GET    /api/v1/purchase/suppliers        获取供应商列表
POST   /api/v1/purchase/suppliers        创建供应商
```

#### DCIM API

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

#### 审批 API

```
GET    /api/v1/approvals                 获取审批列表
GET    /api/v1/approvals/:id             获取审批详情
POST   /api/v1/approvals/:id/approve     批准
POST   /api/v1/approvals/:id/reject      拒绝
GET    /api/v1/approvals/pending         获取待审批列表
GET    /api/v1/approvals/flows           获取审批流程模板
```

#### 报表 API

```
GET    /api/v1/reports/inventory         库存报表
GET    /api/v1/reports/asset-value       资产价值报表
GET    /api/v1/reports/utilization       设备利用率报表
GET    /api/v1/reports/purchase          采购统计报表
POST   /api/v1/reports/export            导出报表
```

### 4.4 WebSocket API

```
WS     /api/v1/ws/notifications          实时通知
WS     /api/v1/ws/dcim/racks/:id         机柜实时状态
```

## 5. 项目结构

### 5.1 后端项目结构

```
asset-management-backend/
├── cmd/
│   ├── api/                    # API服务入口
│   │   └── main.go
│   ├── worker/                 # 后台任务worker
│   │   └── main.go
│   └── migrate/                # 数据库迁移工具
│       └── main.go
├── ent/
│   ├── schema/                 # Schema定义（手写）
│   │   ├── user.go
│   │   ├── asset.go
│   │   ├── inventory.go
│   │   └── ...
│   ├── generate.go             # 代码生成入口
│   └── [generated files]       # Ent自动生成的代码
├── internal/
│   ├── handler/                # HTTP处理器
│   │   ├── asset_handler.go
│   │   ├── inventory_handler.go
│   │   ├── purchase_handler.go
│   │   ├── dcim_handler.go
│   │   └── approval_handler.go
│   ├── service/                # 业务逻辑层
│   │   ├── asset_service.go
│   │   ├── inventory_service.go
│   │   ├── approval_service.go
│   │   └── dcim_service.go
│   ├── middleware/             # 中间件
│   │   ├── auth.go
│   │   ├── logger.go
│   │   ├── cors.go
│   │   └── error.go
│   ├── dto/                    # 数据传输对象
│   │   ├── request/
│   │   └── response/
│   ├── pkg/                    # 内部包
│   │   ├── workflow/           # 工作流引擎
│   │   ├── notification/       # 通知服务
│   │   ├── integration/        # 第三方集成
│   │   │   ├── wecom/          # 企业微信
│   │   │   └── feishu/         # 飞书
│   │   └── cache/              # 缓存封装
│   └── router/                 # 路由配置
│       └── router.go
├── pkg/                        # 公共库（可复用）
│   ├── config/                 # 配置管理
│   ├── database/               # 数据库连接
│   ├── redis/                  # Redis客户端
│   ├── logger/                 # 日志
│   ├── jwt/                    # JWT工具
│   └── utils/                  # 工具函数
├── api/                        # API文档
│   └── openapi.yaml
├── config/                     # 配置文件
│   ├── config.yaml
│   ├── config.dev.yaml
│   └── config.prod.yaml
├── migrations/                 # 数据库迁移（Ent生成）
├── scripts/                    # 脚本
│   ├── generate_handler.go     # Handler生成器
│   └── generate_dto.go         # DTO生成器
├── tests/                      # 测试
│   ├── integration/
│   └── e2e/
├── docker/                     # Docker相关
│   ├── Dockerfile
│   └── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile                    # 常用命令
└── README.md
```

### 5.2 前端项目结构

```
asset-management-web/
├── src/
│   ├── pages/                  # 页面组件
│   │   ├── dashboard/          # 仪表盘
│   │   ├── inventory/          # 库存管理
│   │   │   ├── AssetList.tsx
│   │   │   ├── AssetDetail.tsx
│   │   │   ├── InboundForm.tsx
│   │   │   └── OutboundForm.tsx
│   │   ├── purchase/           # 采购管理
│   │   │   ├── OrderList.tsx
│   │   │   └── OrderForm.tsx
│   │   ├── dcim/               # DCIM
│   │   │   ├── DataCenterList.tsx
│   │   │   ├── RackList.tsx
│   │   │   ├── RackView.tsx    # 机柜可视化
│   │   │   └── RoomView.tsx    # 机房平面图
│   │   ├── approval/           # 审批
│   │   │   ├── ApprovalList.tsx
│   │   │   └── ApprovalDetail.tsx
│   │   ├── report/             # 报表
│   │   └── system/             # 系统管理
│   │       ├── UserManagement.tsx
│   │       └── RoleManagement.tsx
│   ├── components/             # 通用组件
│   │   ├── common/             # 基础组件
│   │   │   ├── PageHeader/
│   │   │   ├── SearchForm/
│   │   │   └── DataTable/
│   │   └── business/           # 业务组件
│   │       ├── AssetSelector/
│   │       ├── LocationPicker/
│   │       └── ApprovalFlow/
│   ├── layouts/                # 布局组件
│   │   ├── MainLayout.tsx
│   │   └── BlankLayout.tsx
│   ├── stores/                 # Zustand状态管理
│   │   ├── authStore.ts
│   │   ├── userStore.ts
│   │   └── notificationStore.ts
│   ├── services/               # API服务
│   │   ├── api.ts              # Axios配置
│   │   ├── assetService.ts
│   │   ├── inventoryService.ts
│   │   ├── purchaseService.ts
│   │   └── dcimService.ts
│   ├── hooks/                  # 自定义Hooks
│   │   ├── useAuth.ts
│   │   ├── usePermission.ts
│   │   └── useWebSocket.ts
│   ├── types/                  # TypeScript类型
│   │   ├── asset.ts
│   │   ├── inventory.ts
│   │   └── api.ts
│   ├── utils/                  # 工具函数
│   │   ├── format.ts
│   │   ├── validate.ts
│   │   └── export.ts
│   ├── routes/                 # 路由配置
│   │   └── index.tsx
│   ├── styles/                 # 样式
│   │   └── global.css
│   ├── App.tsx
│   └── main.tsx
├── public/
├── index.html
├── vite.config.ts
├── tsconfig.json
├── package.json
└── README.md
```

## 6. 开发计划

### 6.1 阶段划分

```
阶段0：基础设施搭建（1-2周）
阶段1：核心功能开发（6-8周）
阶段2：补充功能开发（4-6周）
阶段3：优化和测试（2-4周）
总计：13-20周（3-5个月）
```

### 6.2 详细开发计划

#### 阶段0：基础设施搭建（Week 1-2）

**后端任务**：
- [ ] 初始化Go项目
- [ ] 配置Fiber框架
- [ ] 配置Ent ORM
- [ ] 定义核心Schema（15个实体）
- [ ] 运行代码生成
- [ ] 配置数据库连接
- [ ] 配置Redis连接
- [ ] 实现基础中间件（认证、日志、错误处理）
- [ ] 实现JWT认证
- [ ] 配置Viper配置管理
- [ ] 编写Makefile

**前端任务**：
- [ ] 初始化Vite + React项目
- [ ] 配置TypeScript
- [ ] 配置Ant Design
- [ ] 实现基础布局
- [ ] 配置路由
- [ ] 封装Axios
- [ ] 实现登录页面
- [ ] 配置Zustand状态管理

**基础设施**：
- [ ] 编写docker-compose.yml
- [ ] 配置PostgreSQL
- [ ] 配置Redis
- [ ] 配置MinIO
- [ ] 配置RabbitMQ
- [ ] 编写README和开发文档

**交付物**：
- 可运行的后端API（健康检查、登录接口）
- 可运行的前端（登录页面）
- 完整的开发环境（Docker Compose）

#### 阶段1：核心功能开发（Week 3-10）

**Week 3-4：库存管理模块**

后端：
- [ ] 资产CRUD API
- [ ] 入库API和业务逻辑
- [ ] 出库API和业务逻辑
- [ ] 库存查询API
- [ ] 库存统计API
- [ ] 单元测试

前端：
- [ ] 资产列表页面
- [ ] 资产详情页面
- [ ] 新增/编辑资产表单
- [ ] 入库表单
- [ ] 出库表单
- [ ] 库存查询页面

**Week 5-7：DCIM子系统**

后端：
- [ ] IDC/机房/机柜CRUD API
- [ ] 机柜U位管理API
- [ ] 设备上架/下架API
- [ ] 容量统计API
- [ ] DCIM可视化服务（独立服务）
- [ ] WebSocket实时推送

前端：
- [ ] IDC管理页面
- [ ] 机房管理页面
- [ ] 机柜列表页面
- [ ] 机柜2D可视化组件
- [ ] 上架规划工具
- [ ] 实时状态更新

**Week 8-10：采购和审批流程**

后端：
- [ ] 采购订单CRUD API
- [ ] 供应商管理API
- [ ] 审批流程引擎
- [ ] 审批API
- [ ] 通知服务（独立服务）
- [ ] 企业微信/飞书集成

前端：
- [ ] 采购订单列表
- [ ] 采购订单表单
- [ ] 供应商管理
- [ ] 审批列表
- [ ] 审批详情和操作
- [ ] 审批流程配置

#### 阶段2：补充功能开发（Week 11-16）

**Week 11-12：设备生命周期管理**
- [ ] 设备分配功能
- [ ] 设备调拨功能
- [ ] 维修管理
- [ ] 报废流程
- [ ] 状态变更记录

**Week 13-14：报表和分析**
- [ ] 报表服务（独立服务）
- [ ] 库存报表
- [ ] 资产价值报表
- [ ] 设备利用率报表
- [ ] 采购统计报表
- [ ] 数据大屏

**Week 15-16：移动端**
- [ ] 移动端H5页面
- [ ] 扫码入库功能
- [ ] 扫码出库功能
- [ ] 现场拍照上传
- [ ] 签名功能
- [ ] PWA离线支持

#### 阶段3：优化和测试（Week 17-20）

**Week 17-18：性能优化**
- [ ] 数据库索引优化
- [ ] Redis缓存策略
- [ ] API性能测试
- [ ] 前端性能优化
- [ ] 大数据量测试

**Week 19：安全和权限**
- [ ] RBAC权限系统完善
- [ ] 数据权限（行级权限）
- [ ] 操作日志审计
- [ ] 安全加固
- [ ] 渗透测试

**Week 20：集成测试和部署**
- [ ] 端到端测试
- [ ] 用户验收测试
- [ ] 部署文档
- [ ] 运维文档
- [ ] 培训材料

### 6.3 AI Agent协作策略

每个功能模块的开发流程：

```
1. 需求澄清（AI + 你）
   - 确认功能范围
   - 确认数据模型
   - 确认API接口

2. Schema定义（AI主导，你Review）
   - AI生成Schema代码
   - 你Review并调整
   - 运行代码生成

3. 后端开发（AI主导）
   - AI生成Handler代码
   - AI实现业务逻辑
   - AI编写单元测试
   - 你Review和测试

4. 前端开发（AI主导）
   - AI生成页面框架
   - AI实现交互逻辑
   - AI对接API
   - 你Review和测试

5. 集成测试（AI + 你）
   - 端到端测试
   - 修复bug
   - 性能测试

6. 部署（你主导，AI辅助）
   - 构建Docker镜像
   - 更新部署配置
   - 执行部署
```

## 7. 部署架构

### 7.1 开发环境

```
Docker Compose本地开发环境：
├── PostgreSQL (端口5432)
├── Redis (端口6379)
├── MinIO (端口9000, 9001)
├── RabbitMQ (端口5672, 15672)
└── 后端API (端口8080)
```

### 7.2 生产环境

```
┌─────────────────────────────────────────────────────────┐
│                      负载均衡                            │
│                   Nginx / HAProxy                       │
└─────────────────────────────────────────────────────────┘
                          ↓
┌──────────────┬──────────────┬──────────────────────────┐
│  API服务器1   │  API服务器2   │  API服务器N              │
│  (Docker)    │  (Docker)    │  (Docker)                │
└──────────────┴──────────────┴──────────────────────────┘
                          ↓
┌──────────────┬──────────────┬──────────────────────────┐
│ PostgreSQL   │    Redis     │  MinIO        RabbitMQ   │
│ (主从复制)    │  (集群)      │  (集群)       (集群)      │
└──────────────┴──────────────┴──────────────────────────┘
```

### 7.3 部署清单

**服务器要求**：
- API服务器：2核4G内存（最小），建议4核8G
- 数据库服务器：4核16G内存，SSD存储
- Redis服务器：2核4G内存
- 文件存储：根据需求配置

**软件要求**：
- Docker 20+
- Docker Compose 2+
- PostgreSQL 15+
- Redis 7+
- Nginx 1.20+

## 8. 监控和运维

### 8.1 监控指标

- API响应时间
- 数据库查询性能
- 缓存命中率
- 错误率
- 并发用户数
- 资源使用率（CPU、内存、磁盘）

### 8.2 日志管理

- 应用日志（Zap结构化日志）
- 访问日志（Nginx）
- 错误日志
- 审计日志（操作记录）

### 8.3 备份策略

- 数据库：每日全量备份 + 实时增量备份
- 文件存储：定期备份
- 配置文件：版本控制

## 9. 安全设计

### 9.1 认证授权

- JWT Token认证
- RBAC权限模型
- 数据权限（行级权限）
- API访问控制

### 9.2 数据安全

- 敏感数据加密
- HTTPS传输
- SQL注入防护
- XSS防护
- CSRF防护

### 9.3 审计

- 操作日志记录
- 登录日志
- 数据变更追踪

## 10. 下一步行动

现在开始搭建项目脚手架：

1. ✅ 创建技术架构文档
2. ⏭️ 配置Docker开发环境
3. ⏭️ 搭建后端项目脚手架
4. ⏭️ 定义核心数据模型Schema
5. ⏭️ 搭建前端项目脚手架

---

文档版本：v1.0
创建日期：2026-03-02
最后更新：2026-03-02
