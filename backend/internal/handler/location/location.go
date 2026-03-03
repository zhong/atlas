package location

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/atlas/ent"
	"github.com/your-org/atlas/ent/asset"
	"github.com/your-org/atlas/ent/location"
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
	Page        int   `query:"page"`
	PageSize    int   `query:"page_size"`
	WarehouseID int64 `query:"warehouse_id"`
	Status      string `query:"status"`
}

type CreateRequest struct {
	WarehouseID  int64  `json:"warehouse_id" validate:"required"`
	Name         string `json:"name" validate:"required"`
	Code         string `json:"code" validate:"required"`
	LocationCode string `json:"location_code" validate:"required"`
	Description  string `json:"description"`
}

type UpdateRequest struct {
	Name        *string `json:"name"`
	Status      *string `json:"status"`
	Description *string `json:"description"`
}

// List 获取库位列表
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

	query := h.client.Location.Query().WithWarehouse()

	if req.WarehouseID > 0 {
		query = query.Where(location.HasWarehouseWith())
	}
	if req.Status != "" {
		query = query.Where(location.StatusEQ(location.Status(req.Status)))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to count locations",
		})
	}

	locations, err := query.
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Order(ent.Asc(location.FieldCode)).
		All(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query locations",
		})
	}

	data := make([]fiber.Map, len(locations))
	for i, loc := range locations {
		item := fiber.Map{
			"id":            loc.ID,
			"name":          loc.Name,
			"code":          loc.Code,
			"location_code": loc.LocationCode,
			"status":        string(loc.Status),
			"created_at":    loc.CreatedAt,
		}

		if loc.Edges.Warehouse != nil {
			item["warehouse"] = fiber.Map{
				"id":   loc.Edges.Warehouse.ID,
				"name": loc.Edges.Warehouse.Name,
				"code": loc.Edges.Warehouse.Code,
			}
		}

		data[i] = item
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

// Get 获取库位详情
func (h *Handler) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid location ID",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	loc, err := h.client.Location.Query().
		Where(location.ID(int(id))).
		WithWarehouse().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Location not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query location",
		})
	}

	response := fiber.Map{
		"id":            loc.ID,
		"name":          loc.Name,
		"code":          loc.Code,
		"location_code": loc.LocationCode,
		"status":        string(loc.Status),
		"description":   loc.Description,
		"created_at":    loc.CreatedAt,
		"updated_at":    loc.UpdatedAt,
	}

	if loc.Edges.Warehouse != nil {
		response["warehouse"] = fiber.Map{
			"id":   loc.Edges.Warehouse.ID,
			"name": loc.Edges.Warehouse.Name,
			"code": loc.Edges.Warehouse.Code,
		}
	}

	return c.JSON(response)
}

// Create 创建库位
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
	exists, err := h.client.Location.Query().
		Where(location.Code(req.Code)).
		Exist(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check location code",
		})
	}
	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Location code already exists",
		})
	}

	builder := h.client.Location.Create().
		SetWarehouseID(int(req.WarehouseID)).
		SetName(req.Name).
		SetCode(req.Code).
		SetLocationCode(req.LocationCode).
		SetStatus("available").
		SetNillableDescription(&req.Description)

	loc, err := builder.Save(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create location",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":            loc.ID,
		"name":          loc.Name,
		"code":          loc.Code,
		"location_code": loc.LocationCode,
		"status":        loc.Status,
		"created_at":    loc.CreatedAt,
	})
}

// Update 更新库位
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid location ID",
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

	builder := h.client.Location.UpdateOneID(int(id))

	if req.Name != nil {
		builder.SetName(*req.Name)
	}
	if req.Status != nil {
		builder.SetStatus(location.Status(*req.Status))
	}
	if req.Description != nil {
		builder.SetDescription(*req.Description)
	}

	loc, err := builder.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Location not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update location",
		})
	}

	return c.JSON(fiber.Map{
		"id":         loc.ID,
		"name":       loc.Name,
		"code":       loc.Code,
		"status":     string(loc.Status),
		"updated_at": loc.UpdatedAt,
	})
}

// Delete 删除库位
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid location ID",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 检查是否有关联的资产
	assetCount, err := h.client.Asset.Query().
		Where(asset.HasLocationWith(location.ID(int(id)))).
		Count(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check assets",
		})
	}
	if assetCount > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot delete location with existing assets",
		})
	}

	if err := h.client.Location.DeleteOneID(int(id)).Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Location not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete location",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Location deleted successfully",
		"id":      id,
	})
}
