package warehouse

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/atlas/ent"
	"github.com/your-org/atlas/ent/location"
	"github.com/your-org/atlas/ent/warehouse"
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
	Type     string `query:"type"`
}

type CreateRequest struct {
	Name          string `json:"name" validate:"required"`
	Code          string `json:"code" validate:"required"`
	WarehouseType string `json:"warehouse_type" validate:"required"`
	Location      string `json:"location" validate:"required"`
	Address       string `json:"address"`
	Contact       string `json:"contact"`
}

type UpdateRequest struct {
	Name     *string `json:"name"`
	Location *string `json:"location"`
	Address  *string `json:"address"`
	Contact  *string `json:"contact"`
	Status   *string `json:"status"`
}

// List 获取仓库列表
func (h *Handler) List(c *fiber.Ctx) error {
	var req ListRequest
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := h.client.Warehouse.Query()

	if req.Status != "" {
		query = query.Where(warehouse.StatusEQ(warehouse.Status(req.Status)))
	}
	if req.Type != "" {
		query = query.Where(warehouse.WarehouseTypeEQ(warehouse.WarehouseType(req.Type)))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to count warehouses",
		})
	}

	warehouses, err := query.
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Order(ent.Asc(warehouse.FieldCode)).
		All(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query warehouses",
		})
	}

	data := make([]fiber.Map, len(warehouses))
	for i, w := range warehouses {
		data[i] = fiber.Map{
			"id":             w.ID,
			"name":           w.Name,
			"code":           w.Code,
			"warehouse_type": string(w.WarehouseType),
			"location":       w.Location,
			"address":        w.Address,
			"contact":        w.Contact,
			"status":         w.Status,
			"created_at":     w.CreatedAt,
		}
	}

	totalPages := (total + req.PageSize - 1) / req.PageSize

	return c.JSON(fiber.Map{
		"data":        data,
		"total":       total,
		"page":        req.Page,
		"page_size":   req.PageSize,
		"total_pages": totalPages,
	})
}

// Get 获取仓库详情
func (h *Handler) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid warehouse ID",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	w, err := h.client.Warehouse.Query().
		Where(warehouse.ID(int(id))).
		WithLocations().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Warehouse not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query warehouse",
		})
	}

	response := fiber.Map{
		"id":             w.ID,
		"name":           w.Name,
		"code":           w.Code,
		"warehouse_type": string(w.WarehouseType),
		"location":       w.Location,
		"address":        w.Address,
		"contact":        w.Contact,
		"status":         w.Status,
		"created_at":     w.CreatedAt,
		"updated_at":     w.UpdatedAt,
	}

	if len(w.Edges.Locations) > 0 {
		locations := make([]fiber.Map, len(w.Edges.Locations))
		for i, loc := range w.Edges.Locations {
			locations[i] = fiber.Map{
				"id":            loc.ID,
				"name":          loc.Name,
				"code":          loc.Code,
				"location_code": loc.LocationCode,
				"status":        loc.Status,
			}
		}
		response["locations"] = locations
		response["location_count"] = len(locations)
	}

	return c.JSON(response)
}

// Create 创建仓库
func (h *Handler) Create(c *fiber.Ctx) error {
	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 检查编码是否已存在
	exists, err := h.client.Warehouse.Query().
		Where(warehouse.Code(req.Code)).
		Exist(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check warehouse code",
		})
	}
	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Warehouse code already exists",
		})
	}

	w, err := h.client.Warehouse.Create().
		SetName(req.Name).
		SetCode(req.Code).
		SetWarehouseType(warehouse.WarehouseType(req.WarehouseType)).
		SetLocation(req.Location).
		SetNillableAddress(&req.Address).
		SetNillableContact(&req.Contact).
		SetStatus("active").
		Save(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create warehouse",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":             w.ID,
		"name":           w.Name,
		"code":           w.Code,
		"warehouse_type": string(w.WarehouseType),
		"location":       w.Location,
		"status":         w.Status,
		"created_at":     w.CreatedAt,
	})
}

// Update 更新仓库
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid warehouse ID",
		})
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	builder := h.client.Warehouse.UpdateOneID(int(id))

	if req.Name != nil {
		builder.SetName(*req.Name)
	}
	if req.Location != nil {
		builder.SetLocation(*req.Location)
	}
	if req.Address != nil {
		builder.SetAddress(*req.Address)
	}
	if req.Contact != nil {
		builder.SetContact(*req.Contact)
	}
	if req.Status != nil {
		builder.SetStatus(warehouse.Status(*req.Status))
	}

	w, err := builder.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Warehouse not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update warehouse",
		})
	}

	return c.JSON(fiber.Map{
		"id":         w.ID,
		"name":       w.Name,
		"code":       w.Code,
		"location":   w.Location,
		"status":     w.Status,
		"updated_at": w.UpdatedAt,
	})
}

// Delete 删除仓库
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid warehouse ID",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 检查是否有关联的库位
	locationCount, err := h.client.Location.Query().
		Where(location.HasWarehouseWith(warehouse.ID(int(id)))).
		Count(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check locations",
		})
	}
	if locationCount > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot delete warehouse with existing locations",
		})
	}

	if err := h.client.Warehouse.DeleteOneID(int(id)).Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Warehouse not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete warehouse",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Warehouse deleted successfully",
		"id":      id,
	})
}
