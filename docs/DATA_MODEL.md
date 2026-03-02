# Atlas 数据模型设计文档

本文档基于Excel参考文件分析结果，详细说明Atlas系统的数据模型设计。

## 📊 数据模型概览

### 已实现的实体（21个）

| 实体 | 说明 | 优先级 | 状态 |
|------|------|--------|------|
| User | 用户 | P0 | ✅ |
| Role | 角色 | P0 | ✅ |
| Permission | 权限 | P0 | ✅ |
| Warehouse | 仓库 | P0 | ✅ |
| Location | 库位 | P0 | ✅ |
| AssetType | 资产类型 | P0 | ✅ |
| Asset | 资产 | P0 | ✅ |
| InventoryRecord | 库存记录 | P0 | ✅ |
| Supplier | 供应商 | P1 | ✅ |
| PurchaseOrder | 采购订单 | P1 | ✅ |
| OrderItem | 订单明细 | P1 | ✅ |
| DataCenter | 数据中心 | P1 | ✅ |
| Room | 机房 | P1 | ✅ |
| Rack | 机柜 | P1 | ✅ |
| RackUnit | U位 | P1 | ✅ |
| Approval | 审批 | P1 | ✅ |
| ApprovalNode | 审批节点 | P1 | ✅ |
| NetworkConnection | 网络连接 | P1 | ✅ |
| IPAddress | IP地址 | P1 | ✅ |
| RepairVendor | 维修供应商 | P1 | ✅ |
| RepairTicket | 维修工单 | P1 | ✅ |

## 🎯 基于Excel分析的关键设计决策

### 1. 库存管理增强（基于 库存管理表.xlsx）

**发现：**
- 10个不同地点的库存管理
- 项目分区：AI云、超算云、行业云
- 设备借出/归还流程
- 详细的生命周期记录

**实现：**

#### Asset（资产）增强字段
```go
project_zone: Enum (ai_cloud, hpc_cloud, industry_cloud, other)
status: Enum (in_stock, deployed, borrowed, in_transit, maintenance, retired)
borrow_status: Enum (available, borrowed, reserved)
borrowed_by_id: Integer
borrowed_at: Timestamp
expected_return_at: Timestamp
notes: Text  // 记录设备生命周期事件
```

#### Location（库位）增强字段
```go
location_code: String  // 如：2号AI库、1号库
parent_location_id: Integer  // 支持层级结构
```

#### InventoryRecord（库存记录）增强
```go
record_type: Enum (in, out, transfer, adjust, check, borrow, return)
from_location_name: String  // 冗余字段，便于查询
to_location_name: String
```

### 2. DCIM功能（基于 设备信息表.xlsx）

**发现：**
- 机柜网格布局（A-N行，01-13列）
- 网络拓扑连接（Leaf-Spine架构）
- IP地址管理（内网、公网、管理网）
- 光纤连接管理

**实现：**

#### Rack（机柜）
```go
position_code: String  // 如：A01, M13
row: String  // 如：A, B, M
column: Integer  // 如：01, 13
position: JSON  // 位置坐标，用于平面图显示
power_capacity: Float
power_used: Float
```

#### NetworkConnection（网络连接）- 新增
```go
source_equipment_id: FK -> Asset
source_port: String
target_equipment_id: FK -> Asset
target_port: String
connection_type: Enum (ethernet, infiniband, fiber, other)
speed: String  // 如：25G, 100G
cable_type: String
cable_length: String
```

#### IPAddress（IP地址管理）- 新增
```go
ip_address: String (unique)
subnet: String
gateway: String
ip_type: Enum (internal, public, management)
status: Enum (allocated, available, reserved)
vlan: String
```

### 3. 采购流程优化（基于 采购订单模版.xlsx）

**发现：**
- 采购数量自动计算：需求数量 + 备件需求 - 可用库存
- 货期格式：PO+N天
- 供应商专长管理

**实现：**

#### OrderItem（订单明细）增强
```go
required_qty: Integer  // 需求数量
spare_qty: Integer  // 备件需求
available_stock: Integer  // 可用库存
purchase_qty: Integer  // 采购数量（自动计算）
warranty_period: String  // 质保期
delivery_timeline: String  // 货期（如：PO+7）
```

#### Supplier（供应商）
```go
category_specialties: JSON  // 专长类别数组
rating: Float  // 评分（0-5）
```

### 4. 维修管理（基于 GPU故障报修单.xlsx）

**发现：**
- 1204条维修记录
- 5个不同供应商（四通、超融核、融科联创、安擎、麦芒）
- SLA管理需求

**实现：**

#### RepairTicket（维修工单）- 新增
```go
ticket_number: String (unique)
asset_id: FK -> Asset
gpu_model: String
serial_number: String
fault_description: Text
fault_type: Enum (hardware, software, performance, other)
severity: Enum (low, medium, high, critical)
status: Enum (reported, diagnosed, in_repair, testing, resolved, closed)
vendor_id: FK -> RepairVendor
reported_by_id: FK -> User
assigned_to_id: FK -> User
reported_at: Timestamp
resolved_at: Timestamp
sla_deadline: Timestamp
```

#### RepairVendor（维修供应商）- 新增
```go
name: String  // 如：四通、超融核
supported_models: JSON  // 支持的设备型号
sla_hours: Integer  // SLA响应时间
```

## 📐 实体关系图

### 核心关系

```
User ──┬── Role ── Permission
       ├── InventoryRecord
       ├── PurchaseOrder
       ├── Approval
       └── RepairTicket

Warehouse ── Location ── Asset ──┬── AssetType
                                 ├── InventoryRecord
                                 ├── RackUnit ── Rack ── Room ── DataCenter
                                 ├── NetworkConnection (source/target)
                                 ├── IPAddress
                                 └── RepairTicket

PurchaseOrder ──┬── Supplier
                ├── OrderItem
                └── Approval ── ApprovalNode

RepairTicket ──┬── Asset
               ├── RepairVendor ── Supplier
               └── User (reported_by)
```

## 🔑 关键字段说明

### 枚举类型定义

#### Asset.status
- `in_stock`: 在库
- `deployed`: 已部署
- `borrowed`: 已借出
- `in_transit`: 运输中
- `maintenance`: 维修中
- `retired`: 已报废

#### Asset.project_zone
- `ai_cloud`: AI云
- `hpc_cloud`: 超算云（HPC Cloud）
- `industry_cloud`: 行业云
- `other`: 其他

#### InventoryRecord.record_type
- `in`: 入库
- `out`: 出库
- `transfer`: 调拨
- `adjust`: 调整
- `check`: 盘点
- `borrow`: 借出
- `return`: 归还

#### PurchaseOrder.status
- `draft`: 草稿
- `pending`: 待审批
- `approved`: 已批准
- `rejected`: 已拒绝
- `ordered`: 已下单
- `received`: 已收货
- `completed`: 已完成
- `cancelled`: 已取消

#### RepairTicket.status
- `reported`: 已报修
- `diagnosed`: 已诊断
- `in_repair`: 维修中
- `testing`: 测试中
- `resolved`: 已解决
- `closed`: 已关闭

## 💡 业务规则实现

### 1. 采购数量自动计算
```go
// 在 OrderItem 创建/更新时
purchase_qty = required_qty + spare_qty - available_stock
if purchase_qty < 0 {
    purchase_qty = 0
}
```

### 2. 设备借出流程
```go
// 借出时
asset.status = "borrowed"
asset.borrow_status = "borrowed"
asset.borrowed_by_id = user_id
asset.borrowed_at = now()
asset.expected_return_at = expected_date

// 创建库存记录
inventory_record.record_type = "borrow"
```

### 3. SLA计算
```go
// 创建维修工单时
repair_ticket.sla_deadline = reported_at + vendor.sla_hours
```

### 4. 机柜容量检查
```go
// 设备上架前检查
available_units = rack.total_units - count(occupied_units)
if equipment.height_u > available_units {
    return error("机柜空间不足")
}
```

## 📝 索引策略

### 高频查询字段索引

```go
// Asset
index: asset_no (unique)
index: sn (unique)
index: status
index: project_zone
index: borrow_status
index: brand, model (composite)

// InventoryRecord
index: record_type
index: transaction_date
index: record_no (unique)

// PurchaseOrder
index: order_no (unique)
index: status
index: order_date
index: project

// RepairTicket
index: ticket_number (unique)
index: status
index: severity
index: reported_at

// IPAddress
index: ip_address (unique)
index: ip_type
index: status

// Rack
index: position_code
index: row, column (composite)
index: status
```

## 🚀 下一步

1. **运行数据库迁移**
   ```bash
   cd backend
   make generate  # ✅ 已完成
   # 接下来需要创建迁移脚本
   ```

2. **实现业务逻辑层**
   - AssetService: 资产管理、借出/归还
   - InventoryService: 库存管理、调拨
   - PurchaseService: 采购订单、自动计算
   - DCIMService: 机柜管理、上架规划
   - RepairService: 维修工单、SLA跟踪

3. **创建API接口**
   - 按照CLAUDE.md中定义的API列表实现

4. **前端开发**
   - 资产列表和详情页
   - 库存管理界面
   - 机柜可视化
   - 采购订单表单

## 📚 参考文件

- `docs/references/库存管理表.xlsx` - 库存管理业务流程
- `docs/references/设备信息表.xlsx` - DCIM和网络拓扑
- `docs/references/采购订单模版.xlsx` - 采购流程
- `docs/references/GPU故障报修单.xlsx` - 维修管理
- `docs/references/资源申请表.xlsx` - 资源申请流程

---

**文档版本**: v1.0
**创建日期**: 2026-03-02
**最后更新**: 2026-03-02
**状态**: Schema定义完成，代码生成成功 ✅
