package asset

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/atlas/ent"
	"github.com/your-org/atlas/ent/asset"
	"github.com/your-org/atlas/ent/assettype"
	"github.com/your-org/atlas/pkg/config"
)

type Handler struct {
	client *ent.Client
	cfg    *config.Config
}

func NewHandler(client *ent.Client, cfg *config.Config) *Handler {
	return &Handler{
		client: client,
		cfg:    cfg,
	}
}

type ListRequest struct {
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
	Status   string `query:"status"`
	Category string `query:"category"`
	Keyword  string `query:"keyword"`
}

type ListResponse struct {
	Data       []*AssetInfo `json:"data"`
	Total      int          `json:"total"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	TotalPages int          `json:"total_pages"`
}

type AssetInfo struct {
	ID           int64     `json:"id"`
	AssetNo      string    `json:"asset_no"`
	Name         string    `json:"name"`
	SerialNumber string    `json:"serial_number"`
	Status       string    `json:"status"`
	Category     string    `json:"category"`
	TypeName     string    `json:"type_name"`
	Location     string    `json:"location"`
	ProjectZone  string    `json:"project_zone"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateRequest struct {
	AssetTypeID  int64  `json:"asset_type_id" validate:"required"`
	AssetNo      string `json:"asset_no" validate:"required"`
	Name         string `json:"name" validate:"required"`
	SerialNumber string `json:"serial_number"`
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	ProjectZone  string `json:"project_zone"`
	Specs        map[string]interface{} `json:"specs"`
}

// List 获取资产列表
func (h *Handler) List(c *fiber.Ctx) error {
	var req ListRequest
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	// 设置默认值
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 构建查询
	query := h.client.Asset.Query().
		WithAssetType().
		WithLocation()

	// 状态筛选
	if req.Status != "" {
		query = query.Where(asset.StatusEQ(asset.Status(req.Status)))
	}

	// 分类筛选
	if req.Category != "" {
		query = query.Where(asset.HasAssetTypeWith(
			assettype.CategoryEQ(assettype.Category(req.Category)),
		))
	}

	// 关键词搜索
	if req.Keyword != "" {
		query = query.Where(
			asset.Or(
				asset.AssetNoContains(req.Keyword),
				asset.NameContains(req.Keyword),
				asset.SnContains(req.Keyword),
			),
		)
	}

	// 获取总数
	total, err := query.Count(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to count assets",
		})
	}

	// 分页查询
	assets, err := query.
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Order(ent.Desc(asset.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query assets",
		})
	}

	// 转换为响应格式
	data := make([]*AssetInfo, len(assets))
	for i, a := range assets {
		info := &AssetInfo{
			ID:           int64(a.ID),
			AssetNo:      a.AssetNo,
			Name:         a.Name,
			SerialNumber: a.Sn,
			Status:       string(a.Status),
			ProjectZone:  string(a.ProjectZone),
			CreatedAt:    a.CreatedAt,
		}

		if a.Edges.AssetType != nil {
			info.Category = string(a.Edges.AssetType.Category)
			info.TypeName = a.Edges.AssetType.Name
		}

		if a.Edges.Location != nil {
			info.Location = a.Edges.Location.Name
		}

		data[i] = info
	}

	totalPages := (total + req.PageSize - 1) / req.PageSize

	return c.JSON(ListResponse{
		Data:       data,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	})
}

// Get 获取资产详情
func (h *Handler) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid asset ID",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a, err := h.client.Asset.Query().
		Where(asset.ID(int(id))).
		WithAssetType().
		WithLocation(func(q *ent.LocationQuery) {
			q.WithWarehouse()
		}).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Asset not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query asset",
		})
	}

	// 构建详细响应
	response := fiber.Map{
		"id":                   a.ID,
		"asset_no":             a.AssetNo,
		"name":                 a.Name,
		"serial_number":        a.Sn,
		"brand":                a.Brand,
		"model":                a.Model,
		"status":               string(a.Status),
		"project_zone":         string(a.ProjectZone),
		"borrow_status":        string(a.BorrowStatus),
		"purchase_date":        a.PurchaseDate,
		"warranty_expire_date": a.WarrantyExpireDate,
		"purchase_price":       a.PurchasePrice,
		"specs":                a.Specs,
		"notes":                a.Notes,
		"created_at":           a.CreatedAt,
		"updated_at":           a.UpdatedAt,
	}

	if a.Edges.AssetType != nil {
		response["asset_type"] = fiber.Map{
			"id":       a.Edges.AssetType.ID,
			"name":     a.Edges.AssetType.Name,
			"code":     a.Edges.AssetType.Code,
			"category": a.Edges.AssetType.Category,
		}
	}

	if a.Edges.Location != nil {
		loc := a.Edges.Location
		locationInfo := fiber.Map{
			"id":            loc.ID,
			"name":          loc.Name,
			"code":          loc.Code,
			"location_code": loc.LocationCode,
		}
		if loc.Edges.Warehouse != nil {
			locationInfo["warehouse"] = fiber.Map{
				"id":   loc.Edges.Warehouse.ID,
				"name": loc.Edges.Warehouse.Name,
				"code": loc.Edges.Warehouse.Code,
			}
		}
		response["location"] = locationInfo
	}

	return c.JSON(response)
}

// Create 创建资产
func (h *Handler) Create(c *fiber.Ctx) error {
	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 检查资产编码是否已存在
	exists, err := h.client.Asset.Query().
		Where(asset.AssetNo(req.AssetNo)).
		Exist(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check asset number",
		})
	}
	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Asset number already exists",
		})
	}

	// 创建资产
	builder := h.client.Asset.Create().
		SetAssetTypeID(int(req.AssetTypeID)).
		SetAssetNo(req.AssetNo).
		SetName(req.Name).
		SetStatus("in_stock").
		SetBorrowStatus("available")

	if req.SerialNumber != "" {
		builder.SetSn(req.SerialNumber)
	}
	if req.Brand != "" {
		builder.SetBrand(req.Brand)
	}
	if req.Model != "" {
		builder.SetModel(req.Model)
	}
	if req.ProjectZone != "" {
		builder.SetProjectZone(asset.ProjectZone(req.ProjectZone))
	}
	if req.Specs != nil {
		builder.SetSpecs(req.Specs)
	}

	a, err := builder.Save(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create asset",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":         a.ID,
		"asset_no":   a.AssetNo,
		"name":       a.Name,
		"status":     a.Status,
		"created_at": a.CreatedAt,
	})
}
