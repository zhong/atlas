# 数据库迁移和初始化指南

本文档说明如何初始化Atlas数据库并填充初始数据。

## 📋 前置条件

1. Docker和Docker Compose已安装
2. Go 1.21+已安装
3. 已克隆项目并进入backend目录

## 🚀 快速开始

### 1. 启动Docker服务

```bash
# 在项目根目录
docker-compose up -d

# 等待服务启动（约30秒）
docker-compose ps

# 应该看到所有服务都是 Up 状态
```

### 2. 运行数据库迁移

```bash
cd backend

# 方式1：使用Makefile（推荐）
make migrate

# 方式2：直接运行
go run cmd/migrate/main.go
```

**输出示例**：
```
Running database migrations...
✅ Database migrations completed successfully!

Created tables:
  - users
  - roles
  - permissions
  - warehouses
  - locations
  - asset_types
  - assets
  - inventory_records
  - suppliers
  - purchase_orders
  - order_items
  - data_centers
  - rooms
  - racks
  - rack_units
  - approvals
  - approval_nodes
  - network_connections
  - ip_addresses
  - repair_vendors
  - repair_tickets

✅ All done! You can now run the seed script to populate initial data.
```

### 3. 填充初始数据

```bash
# 方式1：使用Makefile（推荐）
make seed

# 方式2：直接运行
go run cmd/seed/main.go
```

**输出示例**：
```
Starting to seed database...

📝 Creating permissions...
✅ Created 11 permissions

👥 Creating roles...
✅ Created 4 roles

🧑 Creating users...
✅ Created 2 users

🏢 Creating warehouses...
✅ Created 5 warehouses

📦 Creating locations...
✅ Created 15 locations

🏷️  Creating asset types...
✅ Created 7 asset types

🏭 Creating suppliers...
✅ Created 5 suppliers

🏢 Creating data centers...
✅ Created 2 data centers

✅ Database seeding completed successfully!

📊 Summary:
  - Permissions: 11
  - Roles: 4
  - Users: 2
  - Warehouses: 5
  - Locations: 15
  - Asset Types: 7
  - Suppliers: 5
  - Data Centers: 2

🎉 You can now start the API server!
```

### 4. 一键重置数据库

```bash
# 同时运行迁移和种子数据
make reset-db
```

## 📊 初始数据详情

### 用户账号

| 用户名 | 密码 | 角色 | 邮箱 |
|--------|------|------|------|
| admin | admin123 | 管理员 | admin@atlas.com |
| test | test123 | 普通用户 | test@atlas.com |

### 角色和权限

**管理员（admin）**
- 拥有所有权限

**仓库管理员（warehouse_admin）**
- 查看/创建/更新资产
- 查看/管理库存

**采购员（purchaser）**
- 查看资产
- 查看/创建采购订单

**普通用户（viewer）**
- 查看资产
- 查看库存
- 查看采购订单

### 仓库

基于Excel分析创建的仓库：
- 2号AI库（青岛）
- 1号库（北京）
- 3号库（上海）
- 5号基地库（广州）
- 北京小库房（北京）

每个仓库包含3个库位。

### 资产类型

- GPU服务器
- CPU服务器
- 交换机-25G
- 交换机-100G
- 网卡-25G
- 网卡-100G
- 存储设备

### 供应商

基于Excel分析创建的供应商：
- 山石网科（防火墙、网络安全设备）
- 新华三（交换机、路由器）
- 华云光电（光模块、光纤）
- 四通（GPU维修、服务器维修）
- 超融核（GPU维修）

### 数据中心

- 青岛数据中心
- 北京数据中心

## 🔍 验证数据库

### 使用Adminer（Web界面）

1. 访问 http://localhost:8081
2. 登录信息：
   - 系统：PostgreSQL
   - 服务器：postgres
   - 用户名：admin
   - 密码：admin123
   - 数据库：asset_management

### 使用psql命令行

```bash
# 连接数据库
docker exec -it asset-mgmt-postgres psql -U admin -d asset_management

# 查看所有表
\dt

# 查看用户
SELECT username, email, real_name FROM users;

# 查看角色
SELECT name, code, description FROM roles;

# 退出
\q
```

### 使用SQL查询

```sql
-- 查看权限数量
SELECT COUNT(*) FROM permissions;

-- 查看角色及其权限数量
SELECT r.name, COUNT(rp.permission_id) as perm_count
FROM roles r
LEFT JOIN role_permissions rp ON r.id = rp.role_id
GROUP BY r.id, r.name;

-- 查看用户及其角色
SELECT u.username, u.real_name, r.name as role_name
FROM users u
JOIN user_roles ur ON u.id = ur.user_id
JOIN roles r ON ur.role_id = r.id;

-- 查看仓库和库位
SELECT w.name as warehouse, COUNT(l.id) as location_count
FROM warehouses w
LEFT JOIN locations l ON w.id = l.warehouse_id
GROUP BY w.id, w.name;
```

## 🛠️ 故障排查

### 问题1：连接数据库失败

**错误**：`Failed to connect to database`

**解决方案**：
```bash
# 检查Docker服务状态
docker-compose ps

# 如果服务未运行，启动它们
docker-compose up -d

# 等待30秒让服务完全启动
sleep 30

# 重试迁移
make migrate
```

### 问题2：表已存在

**错误**：`table "users" already exists`

**解决方案**：
```bash
# 方式1：删除并重建数据库
docker-compose down -v
docker-compose up -d
sleep 30
make reset-db

# 方式2：手动删除表
docker exec -it asset-mgmt-postgres psql -U admin -d asset_management -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
make reset-db
```

### 问题3：种子数据已存在

**错误**：`duplicate key value violates unique constraint`

**解决方案**：
种子数据脚本不是幂等的，如果数据已存在会报错。需要先清空数据库：

```bash
# 完全重置
docker-compose down -v
docker-compose up -d
sleep 30
make reset-db
```

### 问题4：导入路径错误

**错误**：`package github.com/your-org/asset-management/xxx not found`

**解决方案**：
```bash
# 确保go.mod中的模块名正确
cat go.mod | head -1
# 应该显示: module github.com/your-org/atlas

# 如果不正确，更新go.mod
go mod edit -module=github.com/your-org/atlas
go mod tidy
```

## 📝 下一步

数据库初始化完成后，你可以：

1. **启动API服务器**
   ```bash
   make run
   ```

2. **测试API**
   ```bash
   # 健康检查
   curl http://localhost:8080/health

   # 登录（获取JWT token）
   curl -X POST http://localhost:8080/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username":"admin","password":"admin123"}'
   ```

3. **开始开发**
   - 实现资产管理API
   - 实现库存管理API
   - 搭建前端项目

## 🔄 重新开始

如果需要完全重新开始：

```bash
# 1. 停止并删除所有Docker容器和数据
docker-compose down -v

# 2. 重新启动服务
docker-compose up -d

# 3. 等待服务启动
sleep 30

# 4. 重置数据库
cd backend
make reset-db

# 5. 启动API服务
make run
```

---

**最后更新**: 2026-03-02
**版本**: v1.0
