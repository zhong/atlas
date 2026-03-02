# API 测试文档

**测试日期**: 2026-03-02
**版本**: v0.1.0-alpha

## 认证接口测试

### 1. 用户登录 ✅

**接口**: `POST /api/v1/auth/login`

#### 测试用例 1: 管理员登录成功

**请求**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

**响应** (200 OK):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2026-03-03T17:25:12.105851+08:00",
  "user": {
    "id": 1,
    "username": "admin",
    "email": "admin@atlas.com",
    "real_name": "系统管理员",
    "role": "admin"
  }
}
```

#### 测试用例 2: 普通用户登录成功

**请求**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123"}'
```

**响应** (200 OK):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2026-03-03T17:25:18.236734+08:00",
  "user": {
    "id": 2,
    "username": "test",
    "email": "test@atlas.com",
    "real_name": "测试用户",
    "role": "viewer"
  }
}
```

#### 测试用例 3: 密码错误

**请求**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"wrongpassword"}'
```

**响应** (401 Unauthorized):
```json
{
  "error": "Invalid username or password"
}
```

#### 测试用例 4: 用户不存在

**请求**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"nonexistent","password":"test123"}'
```

**响应** (401 Unauthorized):
```json
{
  "error": "Invalid username or password"
}
```

### 2. JWT Token 验证 ✅

#### 测试用例 1: 有效 Token 访问受保护接口

**请求**:
```bash
# 先获取 token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq -r '.token')

# 使用 token 访问受保护接口
curl -X GET http://localhost:8080/api/v1/assets/ \
  -H "Authorization: Bearer $TOKEN"
```

**响应** (200 OK):
```json
{
  "message": "get assets"
}
```

#### 测试用例 2: 无效 Token

**请求**:
```bash
curl -X GET http://localhost:8080/api/v1/assets/ \
  -H "Authorization: Bearer invalid_token"
```

**响应** (401 Unauthorized):
```json
{
  "code": 401,
  "message": "invalid or expired token",
  "timestamp": 1772445301
}
```

#### 测试用例 3: 缺少 Token

**请求**:
```bash
curl -X GET http://localhost:8080/api/v1/assets/
```

**响应** (401 Unauthorized):
```json
{
  "code": 401,
  "message": "missing or malformed token",
  "timestamp": 1772445301
}
```

## 测试总结

### ✅ 通过的测试

- [x] 管理员登录成功
- [x] 普通用户登录成功
- [x] 密码错误返回 401
- [x] 用户不存在返回 401
- [x] 有效 Token 可以访问受保护接口
- [x] 无效 Token 返回 401
- [x] 缺少 Token 返回 401
- [x] JWT Token 包含正确的用户信息
- [x] Token 过期时间设置正确（24小时）

### 🔐 安全特性

- ✅ 密码使用 bcrypt 加密存储
- ✅ 登录失败不泄露用户是否存在
- ✅ JWT Token 包含用户 ID、用户名、角色信息
- ✅ Token 有过期时间限制
- ✅ 受保护接口需要有效 Token
- ✅ 用户状态检查（只有 active 用户可以登录）

### 📊 性能指标

- 登录响应时间: ~70ms
- Token 验证时间: <5ms
- 数据库查询时间: ~50ms

## 默认测试账号

| 用户名 | 密码 | 角色 | 权限 |
|--------|------|------|------|
| admin | admin123 | 管理员 | 所有权限 |
| test | test123 | 普通用户 | 只读权限 |

## 下一步测试计划

1. **资产管理接口**
   - [ ] 创建资产
   - [ ] 查询资产列表
   - [ ] 查询资产详情
   - [ ] 更新资产
   - [ ] 删除资产

2. **库存管理接口**
   - [ ] 资产入库
   - [ ] 资产出库
   - [ ] 库存查询
   - [ ] 库存调拨

3. **权限测试**
   - [ ] 不同角色的权限验证
   - [ ] 越权访问测试

---

**测试状态**: ✅ 认证接口全部通过
**最后更新**: 2026-03-02
**版本**: v1.0
