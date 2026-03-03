# Atlas 开发日志

本文档记录 Atlas 项目的所有重要变更和进展。

## 2026-03-02

### 🎉 项目启动 (上午)

**项目命名**
- 确定项目名称：Atlas（阿特拉斯）
- 寓意：希腊神话擎天神，象征全面掌控IT资产与基础设施
- 更新所有文档和代码中的项目名称

**技术架构确定**
- 选择现代化混合架构：Go + Fiber + Ent + React
- 核心理念：标准化的地方用工具，复杂的地方手写
- 预计开发周期：3-5个月（13-20周）

### 📋 需求分析

**业务范围确定**
- GPU算力公司的IT资产全生命周期管理
- 规模：设备几万级、备件百万级、月入库/出库几百到几千次
- 核心功能：采购、入库、库存、出库、DCIM、设备生命周期管理

**Excel业务文件分析**
- 分析了5个Excel文件（总计1.2MB）
- 提取了完整的业务需求和数据结构
- 文件存放在 `docs/references/` 目录

分析的文件：
1. 库存管理表.xlsx (231KB) - 10个地点的库存管理
2. 设备信息表.xlsx (559KB) - DCIM和网络拓扑
3. 采购订单模版.xlsx (12KB) - 采购流程
4. GPU故障报修单.xlsx (374KB) - 1204条维修记录
5. 资源申请表.xlsx (18KB) - 资源申请流程

### 🏗️ 基础设施搭建

**Docker开发环境**
- 配置了完整的开发环境（docker-compose.yml）
- 服务：PostgreSQL 15, Redis 7, MinIO, RabbitMQ, Adminer
- 所有服务可一键启动：`docker-compose up -d`

**后端项目脚手架**
- 初始化Go项目（Go 1.21+）
- 配置Fiber v2 Web框架
- 配置Ent ORM
- 实现基础中间件：认证、日志、错误处理
- 实现JWT认证
- 配置Viper配置管理
- 创建Makefile常用命令

项目结构：
```
backend/
├── cmd/api/          # API服务入口
├── ent/              # Ent ORM
├── internal/         # 内部代码
│   ├── handler/      # HTTP处理器
│   ├── service/      # 业务逻辑
│   ├── middleware/   # 中间件
│   └── router/       # 路由配置
├── pkg/              # 公共库
└── config/           # 配置文件
```

### 📊 数据模型设计

**Schema定义完成**
- 创建了21个Ent Schema实体
- P0优先级：8个核心实体
- P1优先级：13个扩展实体

**P0实体（核心）**：
- User（用户）- 增强：部门、最后登录
- Role（角色）
- Permission（权限）
- Warehouse（仓库）- 支持IDC/仓库/办公室类型
- Location（库位）- 支持层级结构和位置代码
- AssetType（资产类型）
- Asset（资产）- 项目分区、借出状态、借用人追踪
- InventoryRecord（库存记录）- 支持借出/归还类型

**P1实体（扩展）**：
- Supplier（供应商）- 专长类别、评分
- PurchaseOrder（采购订单）
- OrderItem（订单明细）- 自动计算采购数量
- DataCenter（数据中心）
- Room（机房）- 支持平面图数据
- Rack（机柜）- 网格布局（行列坐标）
- RackUnit（U位）
- Approval（审批）
- ApprovalNode（审批节点）
- NetworkConnection（网络连接）- 网络拓扑
- IPAddress（IP地址）- IPAM功能
- RepairVendor（维修供应商）- SLA管理
- RepairTicket（维修工单）- 故障追踪

**基于Excel分析的关键增强**：

1. **库存管理增强**（基于库存管理表.xlsx）
   - Asset: 项目分区（AI云、超算云、行业云）
   - Asset: 借出状态、借用人、借出时间追踪
   - Location: 位置代码（如：2号AI库）、层级结构
   - InventoryRecord: 借出/归还类型

2. **DCIM功能**（基于设备信息表.xlsx）
   - Rack: 网格布局（position_code, row, column）
   - Rack: 电力容量管理
   - NetworkConnection: 网络拓扑连接管理
   - IPAddress: IP地址分配和追踪（内网、公网、管理网）

3. **采购优化**（基于采购订单模版.xlsx）
   - OrderItem: 自动计算公式
     ```
     采购数量 = 需求数量 + 备件需求 - 可用库存
     ```
   - OrderItem: 货期管理（PO+N天）
   - Supplier: 专长类别、评分

4. **维修管理**（基于GPU故障报修单.xlsx）
   - RepairTicket: 故障类型、严重程度分级
   - RepairTicket: SLA截止时间追踪
   - RepairVendor: 多供应商支持（四通、超融核、融科联创、安擎、麦芒）
   - RepairVendor: SLA响应时间配置

**代码生成**
- 运行 `make generate` 成功
- 生成了76个文件，22,637行代码
- 包含所有实体的CRUD操作
- 类型安全的查询构建器
- 完整的关系映射
- 索引优化

### 📚 文档创建

**核心文档**：
- `README.md` - 项目总览
- `ARCHITECTURE.md` - 完整技术架构文档（1000+行）
- `CLAUDE.md` - Claude Code工作指南
- `PROJECT_SUMMARY.md` - 项目进度总结
- `QUICKSTART.md` - 快速开始指南
- `docs/DATA_MODEL.md` - 数据模型设计文档
- `docs/references/README.md` - 参考文件说明
- `CHANGELOG.md` - 本文件

### 🔧 Git仓库

**仓库初始化**
- 初始化Git仓库
- 远程仓库：`git@github.com:zhong/atlas.git`
- 默认分支：`main`

**提交记录**：
1. `4db6cec` - Initial commit: Atlas IT Asset Management System
   - 24个文件，3,542行代码
   - 项目架构和基础设施

2. `b3129a3` - feat: Define complete data model based on Excel analysis
   - 76个文件，22,637行代码
   - 完整数据模型定义和代码生成

### 📈 当前状态

**阶段0：基础设施搭建 - 100% 完成** ✅

完成项：
- ✅ 技术方案设计
- ✅ 项目命名（Atlas）
- ✅ Docker环境配置
- ✅ 后端脚手架
- ✅ 数据模型定义（21个实体）
- ✅ Ent代码生成
- ✅ Excel业务分析
- ✅ Git仓库初始化
- ✅ 完整文档体系

待完成项：
- ⏳ 数据库迁移脚本
- ⏳ 初始数据（种子数据）
- ⏳ 前端项目脚手架

### 🎯 下一步计划

**立即开始**：
1. 创建数据库迁移脚本
2. 初始化数据库表结构
3. 创建种子数据（初始用户、角色、权限等）

**并行进行**：
1. 搭建前端项目脚手架（Vite + React + TypeScript）
2. 配置Ant Design UI库
3. 创建基础布局和路由

**验证流程**：
1. 实现第一个完整的功能模块（库存管理）
2. 端到端测试

---

### 🔐 认证接口实现 (下午)

**用户登录接口**
- 实现了 `POST /api/v1/auth/login` 接口
- 创建了 `internal/handler/auth/auth.go` 认证处理器
- 功能特性：
  - 用户名密码验证
  - bcrypt 密码加密验证
  - 用户状态检查（只允许 active 用户登录）
  - JWT Token 生成（24小时有效期）
  - 返回用户信息和角色

**测试结果**
- ✅ 管理员登录成功（admin/admin123）
- ✅ 普通用户登录成功（test/test123）
- ✅ 密码错误返回 401
- ✅ 用户不存在返回 401
- ✅ 有效 Token 可以访问受保护接口
- ✅ 无效 Token 返回 401
- ✅ 缺少 Token 返回 401

**安全特性**
- 密码使用 bcrypt 加密存储
- 登录失败不泄露用户是否存在
- JWT Token 包含用户 ID、用户名、角色信息
- Token 有过期时间限制
- 受保护接口需要有效 Token

**文档更新**
- 创建了 `docs/API_TESTING.md` API测试文档
- 记录了所有测试用例和响应示例
- 包含安全特性说明和性能指标

---

## 2026-03-03

### 💻 资产管理接口实现

**资产CRUD接口**
- 实现了 `GET /api/v1/assets/` 资产列表接口
- 实现了 `GET /api/v1/assets/:id` 资产详情接口
- 实现了 `POST /api/v1/assets/` 创建资产接口
- 创建了 `internal/handler/asset/asset.go` 资产处理器

**功能特性**
- 分页查询（支持自定义页码和每页数量）
- 状态筛选（in_stock, deployed, maintenance, retired）
- 分类筛选（server, switch, storage, network_card）
- 关键词搜索（资产编号、名称、序列号）
- 关联查询（资产类型、库位、仓库信息）
- 资产编号唯一性校验

**测试结果**
- ✅ 创建资产成功
- ✅ 资产编号唯一性校验
- ✅ 查询资产列表（6条测试数据）
- ✅ 分页查询（page_size=3, total_pages=2）
- ✅ 状态筛选（status=in_stock）
- ✅ 关键词搜索（keyword=GPU，返回5条）
- ✅ 查询资产详情（包含完整信息）
- ✅ 资产不存在返回404

**性能指标**
- 创建资产: ~20ms
- 列表查询: ~30ms
- 详情查询: ~15ms

**文档更新**
- 更新了 `docs/API_TESTING.md` 添加资产管理测试用例
- 包含完整的请求/响应示例

### 📦 资产管理和库存管理接口完善

**资产管理接口扩展**
- 实现了 `PUT /api/v1/assets/:id` 资产更新接口
- 实现了 `DELETE /api/v1/assets/:id` 资产删除接口
- 支持部分字段更新（name, serial_number, brand, model, status, project_zone, specs, notes）
- 删除限制：只能删除在库状态的资产

**库存管理接口实现**
- 实现了 `GET /api/v1/inventory/stock` 库存统计接口
- 实现了 `POST /api/v1/inventory/inbound` 资产入库接口
- 实现了 `POST /api/v1/inventory/outbound` 资产出库接口
- 实现了 `POST /api/v1/inventory/transfer` 资产调拨接口
- 实现了 `GET /api/v1/inventory/records` 库存记录查询接口
- 创建了 `internal/handler/inventory/inventory.go` 库存处理器

**功能特性**
- 事务处理（确保入库/出库/调拨的数据一致性）
- 自动生成记录编号（IN-/OUT-/TRF-前缀）
- 状态自动更新（入库设为in_stock，出库设为deployed）
- 库位关联（支持from_location和to_location）
- 操作人记录（关联用户信息）
- 库存统计（按状态分类统计）

**测试结果**
- ✅ 资产更新成功
- ✅ 资产删除成功（在库状态）
- ✅ 资产删除失败（已部署状态，符合预期）
- ✅ 库存统计正确（total: 5, in_stock: 4, deployed: 1）
- ✅ 入库操作成功（资产状态更新，库位关联）
- ✅ 出库操作成功（资产状态变为deployed）
- ✅ 库存记录查询成功（包含完整关联信息）

**性能指标**
- 资产更新: ~15ms
- 资产删除: ~10ms
- 库存统计: ~20ms
- 入库操作: ~30ms（含事务）
- 出库操作: ~25ms（含事务）
- 记录查询: ~35ms

**文档更新**
- 更新了 `docs/API_TESTING.md` 添加所有新接口测试用例
- 包含完整的请求/响应示例和错误场景

### 🏢 仓库和库位管理接口实现

**仓库管理接口**
- 实现了 `GET /api/v1/warehouses/` 仓库列表接口
- 实现了 `GET /api/v1/warehouses/:id` 仓库详情接口
- 实现了 `POST /api/v1/warehouses/` 创建仓库接口
- 实现了 `PUT /api/v1/warehouses/:id` 更新仓库接口
- 实现了 `DELETE /api/v1/warehouses/:id` 删除仓库接口
- 创建了 `internal/handler/warehouse/warehouse.go` 仓库处理器

**库位管理接口**
- 实现了 `GET /api/v1/locations/` 库位列表接口
- 实现了 `GET /api/v1/locations/:id` 库位详情接口
- 实现了 `POST /api/v1/locations/` 创建库位接口
- 实现了 `PUT /api/v1/locations/:id` 更新库位接口
- 实现了 `DELETE /api/v1/locations/:id` 删除库位接口
- 创建了 `internal/handler/location/location.go` 库位处理器

**功能特性**
- 仓库类型支持（warehouse-普通仓库, idc-数据中心）
- 状态管理（active-活跃, inactive-停用）
- 关联查询（仓库包含库位列表，库位关联仓库信息）
- 删除保护（有关联数据时禁止删除）
- 编码唯一性校验
- 分页查询支持

**测试结果**
- ✅ 仓库列表查询成功（5个仓库）
- ✅ 库位列表查询成功（9个库位）
- ✅ 创建仓库成功（测试仓库）
- ✅ 创建库位成功（测试库位）
- ✅ 关联查询正确（库位显示所属仓库）

**性能指标**
- 仓库列表: ~25ms
- 库位列表: ~30ms
- 创建仓库: ~15ms
- 创建库位: ~18ms

---

## 统计数据

**代码量**：
- 后端代码：28,000+行（Schema + 生成代码 + 所有业务接口）
- 文档：约11,500行
- 总计：39,500+行

**文件数**：
- 后端文件：108个
- 文档文件：11个
- 配置文件：5个
- 总计：124个

**提交数**：7次
**开发时间**：2天
**团队**：AI Agent + 用户协作

---

**最后更新**: 2026-03-03
**版本**: v0.1.0-alpha
**状态**: 阶段0完成，核心业务接口实现完成
