package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/your-org/atlas/ent"
	"github.com/your-org/atlas/pkg/config"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

func main() {
	// 加载配置
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 构建数据库连接字符串
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	// 连接数据库
	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	log.Println("Starting to seed database...")

	// 1. 创建权限
	log.Println("\n📝 Creating permissions...")
	permissions := createPermissions(ctx, client)
	log.Printf("✅ Created %d permissions", len(permissions))

	// 2. 创建角色
	log.Println("\n👥 Creating roles...")
	roles := createRoles(ctx, client, permissions)
	log.Printf("✅ Created %d roles", len(roles))

	// 3. 创建用户
	log.Println("\n🧑 Creating users...")
	users := createUsers(ctx, client, roles)
	log.Printf("✅ Created %d users", len(users))

	// 4. 创建仓库
	log.Println("\n🏢 Creating warehouses...")
	warehouses := createWarehouses(ctx, client)
	log.Printf("✅ Created %d warehouses", len(warehouses))

	// 5. 创建库位
	log.Println("\n📦 Creating locations...")
	locations := createLocations(ctx, client, warehouses)
	log.Printf("✅ Created %d locations", len(locations))

	// 6. 创建资产类型
	log.Println("\n🏷️  Creating asset types...")
	assetTypes := createAssetTypes(ctx, client)
	log.Printf("✅ Created %d asset types", len(assetTypes))

	// 7. 创建供应商
	log.Println("\n🏭 Creating suppliers...")
	suppliers := createSuppliers(ctx, client)
	log.Printf("✅ Created %d suppliers", len(suppliers))

	// 8. 创建数据中心
	log.Println("\n🏢 Creating data centers...")
	dataCenters := createDataCenters(ctx, client)
	log.Printf("✅ Created %d data centers", len(dataCenters))

	log.Println("\n✅ Database seeding completed successfully!")
	log.Println("\n📊 Summary:")
	log.Printf("  - Permissions: %d", len(permissions))
	log.Printf("  - Roles: %d", len(roles))
	log.Printf("  - Users: %d", len(users))
	log.Printf("  - Warehouses: %d", len(warehouses))
	log.Printf("  - Locations: %d", len(locations))
	log.Printf("  - Asset Types: %d", len(assetTypes))
	log.Printf("  - Suppliers: %d", len(suppliers))
	log.Printf("  - Data Centers: %d", len(dataCenters))

	log.Println("\n🎉 You can now start the API server!")
}

func createPermissions(ctx context.Context, client *ent.Client) []*ent.Permission {
	permissions := []struct {
		name        string
		code        string
		resource    string
		action      string
		description string
	}{
		{"查看资产", "asset:read", "asset", "read", "查看资产信息"},
		{"创建资产", "asset:create", "asset", "create", "创建新资产"},
		{"更新资产", "asset:update", "asset", "update", "更新资产信息"},
		{"删除资产", "asset:delete", "asset", "delete", "删除资产"},
		{"查看库存", "inventory:read", "inventory", "read", "查看库存信息"},
		{"管理库存", "inventory:manage", "inventory", "manage", "管理库存（入库/出库）"},
		{"查看采购", "purchase:read", "purchase", "read", "查看采购订单"},
		{"创建采购", "purchase:create", "purchase", "create", "创建采购订单"},
		{"审批采购", "purchase:approve", "purchase", "approve", "审批采购订单"},
		{"查看用户", "user:read", "user", "read", "查看用户信息"},
		{"管理用户", "user:manage", "user", "manage", "管理用户"},
	}

	var result []*ent.Permission
	for _, p := range permissions {
		perm, err := client.Permission.Create().
			SetName(p.name).
			SetCode(p.code).
			SetResource(p.resource).
			SetAction(p.action).
			SetDescription(p.description).
			Save(ctx)
		if err != nil {
			log.Printf("Warning: Failed to create permission %s: %v", p.code, err)
			continue
		}
		result = append(result, perm)
	}
	return result
}

func createRoles(ctx context.Context, client *ent.Client, permissions []*ent.Permission) []*ent.Role {
	// 管理员角色 - 所有权限
	admin, err := client.Role.Create().
		SetName("管理员").
		SetCode("admin").
		SetDescription("系统管理员，拥有所有权限").
		SetSortOrder(1).
		AddPermissions(permissions...).
		Save(ctx)
	if err != nil {
		log.Fatalf("Failed to create admin role: %v", err)
	}

	// 仓库管理员 - 库存相关权限
	warehousePerms := filterPermissions(permissions, []string{
		"asset:read", "asset:create", "asset:update",
		"inventory:read", "inventory:manage",
	})
	warehouse, err := client.Role.Create().
		SetName("仓库管理员").
		SetCode("warehouse_admin").
		SetDescription("管理仓库和库存").
		SetSortOrder(2).
		AddPermissions(warehousePerms...).
		Save(ctx)
	if err != nil {
		log.Fatalf("Failed to create warehouse role: %v", err)
	}

	// 采购员 - 采购相关权限
	purchasePerms := filterPermissions(permissions, []string{
		"asset:read", "purchase:read", "purchase:create",
	})
	purchaser, err := client.Role.Create().
		SetName("采购员").
		SetCode("purchaser").
		SetDescription("负责采购订单管理").
		SetSortOrder(3).
		AddPermissions(purchasePerms...).
		Save(ctx)
	if err != nil {
		log.Fatalf("Failed to create purchaser role: %v", err)
	}

	// 普通用户 - 只读权限
	viewerPerms := filterPermissions(permissions, []string{
		"asset:read", "inventory:read", "purchase:read",
	})
	viewer, err := client.Role.Create().
		SetName("普通用户").
		SetCode("viewer").
		SetDescription("只能查看信息").
		SetSortOrder(4).
		AddPermissions(viewerPerms...).
		Save(ctx)
	if err != nil {
		log.Fatalf("Failed to create viewer role: %v", err)
	}

	return []*ent.Role{admin, warehouse, purchaser, viewer}
}

func createUsers(ctx context.Context, client *ent.Client, roles []*ent.Role) []*ent.User {
	// 密码加密
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

	// 管理员用户
	admin, err := client.User.Create().
		SetUsername("admin").
		SetPassword(string(hashedPassword)).
		SetEmail("admin@atlas.com").
		SetRealName("系统管理员").
		SetDepartment("IT部门").
		SetStatus("active").
		AddRoles(roles[0]). // admin role
		Save(ctx)
	if err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	// 测试用户
	testPassword, _ := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
	test, err := client.User.Create().
		SetUsername("test").
		SetPassword(string(testPassword)).
		SetEmail("test@atlas.com").
		SetRealName("测试用户").
		SetDepartment("测试部门").
		SetStatus("active").
		AddRoles(roles[3]). // viewer role
		Save(ctx)
	if err != nil {
		log.Fatalf("Failed to create test user: %v", err)
	}

	return []*ent.User{admin, test}
}

func createWarehouses(ctx context.Context, client *ent.Client) []*ent.Warehouse {
	warehouses := []struct {
		name          string
		code          string
		warehouseType string
		location      string
		address       string
	}{
		{"2号AI库", "WH-AI-02", "idc", "青岛", "山东省青岛市"},
		{"1号库", "WH-01", "warehouse", "北京", "北京市海淀区"},
		{"3号库", "WH-03", "warehouse", "上海", "上海市浦东新区"},
		{"5号基地库", "WH-05", "idc", "广州", "广东省广州市"},
		{"北京小库房", "WH-BJ-SMALL", "warehouse", "北京", "北京市朝阳区"},
	}

	var result []*ent.Warehouse
	for _, w := range warehouses {
		wh, err := client.Warehouse.Create().
			SetName(w.name).
			SetCode(w.code).
			SetWarehouseType(w.warehouseType).
			SetLocation(w.location).
			SetAddress(w.address).
			SetStatus("active").
			Save(ctx)
		if err != nil {
			log.Printf("Warning: Failed to create warehouse %s: %v", w.code, err)
			continue
		}
		result = append(result, wh)
	}
	return result
}

func createLocations(ctx context.Context, client *ent.Client, warehouses []*ent.Warehouse) []*ent.Location {
	var result []*ent.Location

	// 为每个仓库创建几个库位
	for i, wh := range warehouses {
		for j := 1; j <= 3; j++ {
			loc, err := client.Location.Create().
				SetName(fmt.Sprintf("%s-区域%d", wh.Name, j)).
				SetCode(fmt.Sprintf("%s-LOC-%02d", wh.Code, j)).
				SetLocationCode(fmt.Sprintf("%s-LOC-%02d", wh.Code, j)).
				SetWarehouse(wh).
				SetStatus("available").
				Save(ctx)
			if err != nil {
				log.Printf("Warning: Failed to create location: %v", err)
				continue
			}
			result = append(result, loc)
		}

		// 只为前两个仓库创建详细信息
		if i >= 2 {
			break
		}
	}

	return result
}

func createAssetTypes(ctx context.Context, client *ent.Client) []*ent.AssetType {
	assetTypes := []struct {
		name        string
		code        string
		category    string
		description string
	}{
		{"GPU服务器", "GPU-SERVER", "server", "配备GPU的服务器"},
		{"CPU服务器", "CPU-SERVER", "server", "标准CPU服务器"},
		{"交换机-25G", "SWITCH-25G", "switch", "25G以太网交换机"},
		{"交换机-100G", "SWITCH-100G", "switch", "100G以太网交换机"},
		{"网卡-25G", "NIC-25G", "network_card", "25G网卡"},
		{"网卡-100G", "NIC-100G", "network_card", "100G网卡"},
		{"存储设备", "STORAGE", "storage", "存储系统"},
	}

	var result []*ent.AssetType
	for _, at := range assetTypes {
		assetType, err := client.AssetType.Create().
			SetName(at.name).
			SetCode(at.code).
			SetCategory(at.category).
			SetDescription(at.description).
			SetStatus("active").
			Save(ctx)
		if err != nil {
			log.Printf("Warning: Failed to create asset type %s: %v", at.code, err)
			continue
		}
		result = append(result, assetType)
	}
	return result
}

func createSuppliers(ctx context.Context, client *ent.Client) []*ent.Supplier {
	suppliers := []struct {
		name        string
		code        string
		specialties []string
	}{
		{"山石网科", "HILLSTONE", []string{"防火墙", "网络安全设备"}},
		{"新华三", "H3C", []string{"交换机", "路由器"}},
		{"华云光电", "HUAYUN", []string{"光模块", "光纤"}},
		{"四通", "SITONG", []string{"GPU维修", "服务器维修"}},
		{"超融核", "CHAORONG", []string{"GPU维修"}},
	}

	var result []*ent.Supplier
	for _, s := range suppliers {
		supplier, err := client.Supplier.Create().
			SetName(s.name).
			SetCode(s.code).
			SetCategorySpecialties(s.specialties).
			SetStatus("active").
			Save(ctx)
		if err != nil {
			log.Printf("Warning: Failed to create supplier %s: %v", s.code, err)
			continue
		}
		result = append(result, supplier)
	}
	return result
}

func createDataCenters(ctx context.Context, client *ent.Client) []*ent.DataCenter {
	dataCenters := []struct {
		name     string
		code     string
		location string
		address  string
	}{
		{"青岛数据中心", "DC-QD", "青岛", "山东省青岛市"},
		{"北京数据中心", "DC-BJ", "北京", "北京市海淀区"},
	}

	var result []*ent.DataCenter
	for _, dc := range dataCenters {
		dataCenter, err := client.DataCenter.Create().
			SetName(dc.name).
			SetCode(dc.code).
			SetLocation(dc.location).
			SetAddress(dc.address).
			SetStatus("active").
			Save(ctx)
		if err != nil {
			log.Printf("Warning: Failed to create data center %s: %v", dc.code, err)
			continue
		}
		result = append(result, dataCenter)
	}
	return result
}

// 辅助函数：根据code过滤权限
func filterPermissions(permissions []*ent.Permission, codes []string) []*ent.Permission {
	codeMap := make(map[string]bool)
	for _, code := range codes {
		codeMap[code] = true
	}

	var result []*ent.Permission
	for _, p := range permissions {
		if codeMap[p.Code] {
			result = append(result, p)
		}
	}
	return result
}
