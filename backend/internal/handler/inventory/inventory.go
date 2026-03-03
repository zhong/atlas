package inventory

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/your-org/atlas/ent"
	"github.com/your-org/atlas/ent/asset"
	"github.com/your-org/atlas/ent/inventoryrecord"
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

type InboundRequest struct {
	AssetID    int64   `json:"asset_id" validate:"required"`
	LocationID int64   `json:"location_id" validate:"required"`
	Quantity   int     `json:"quantity" validate:"required,min=1"`
	Notes      string  `json:"notes"`
	OperatorID int64   `json:"operator_id" validate:"required"`
}

type OutboundRequest struct {
	AssetID    int64  `json:"asset_id" validate:"required"`
	Quantity   int    `json:"quantity" validate:"required,min=1"`
	Recipient  string `json:"recipient" validate:"required"`
	Purpose    string `json:"purpose"`
	Notes      string `json:"notes"`
	OperatorID int64  `json:"operator_id" validate:"required"`
}

type TransferRequest struct {
	AssetID        int64  `json:"asset_id" validate:"required"`
	FromLocationID int64  `json:"from_location_id" validate:"required"`
	ToLocationID   int64  `json:"to_location_id" validate:"required"`
	Quantity       int    `json:"quantity" validate:"required,min=1"`
	Notes          string `json:"notes"`
	OperatorID     int64  `json:"operator_id" validate:"required"`
}

type StockQueryRequest struct {
	LocationID int64  `query:"location_id"`
	AssetType  string `query:"asset_type"`
	Status     string `query:"status"`
}

// Inbound 资产入库
func (h *Handler) Inbound(c *fiber.Ctx) error {
	var req InboundRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 开始事务
	tx, err := h.client.Tx(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to start transaction",
		})
	}

	// 检查资产是否存在
	a, err := tx.Asset.Query().
		Where(asset.ID(int(req.AssetID))).
		Only(ctx)
	if err != nil {
		tx.Rollback()
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Asset not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query asset",
		})
	}

	// 检查库位是否存在
	loc, err := tx.Location.Query().
		Where(location.ID(int(req.LocationID))).
		Only(ctx)
	if err != nil {
		tx.Rollback()
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Location not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query location",
		})
	}

	// 更新资产状态和库位
	_, err = tx.Asset.UpdateOne(a).
		SetStatus("in_stock").
		SetLocation(loc).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update asset",
		})
	}

	// 创建入库记录
	recordNo := fmt.Sprintf("IN-%s-%d", time.Now().Format("20060102150405"), req.AssetID)
	record, err := tx.InventoryRecord.Create().
		SetRecordNo(recordNo).
		SetAsset(a).
		SetToLocation(loc).
		SetToLocationName(loc.Name).
		SetRecordType(inventoryrecord.RecordTypeInbound).
		SetQuantity(req.Quantity).
		SetOperator(tx.User.GetX(ctx, int(req.OperatorID))).
		SetNillableNote(&req.Notes).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create inventory record",
		})
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":          record.ID,
		"asset_id":    a.ID,
		"asset_no":    a.AssetNo,
		"location":    loc.Name,
		"quantity":    record.Quantity,
		"record_type": string(record.RecordType),
		"created_at":  record.CreatedAt,
	})
}

// Outbound 资产出库
func (h *Handler) Outbound(c *fiber.Ctx) error {
	var req OutboundRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 开始事务
	tx, err := h.client.Tx(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to start transaction",
		})
	}

	// 检查资产是否存在且在库
	a, err := tx.Asset.Query().
		Where(
			asset.ID(int(req.AssetID)),
			asset.StatusEQ("in_stock"),
		).
		WithLocation().
		Only(ctx)
	if err != nil {
		tx.Rollback()
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Asset not found or not in stock",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query asset",
		})
	}

	// 更新资产状态
	_, err = tx.Asset.UpdateOne(a).
		SetStatus("deployed").
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update asset",
		})
	}

	// 创建出库记录
	recordNo := fmt.Sprintf("OUT-%s-%d", time.Now().Format("20060102150405"), req.AssetID)
	builder := tx.InventoryRecord.Create().
		SetRecordNo(recordNo).
		SetAsset(a).
		SetRecordType(inventoryrecord.RecordTypeOutbound).
		SetQuantity(req.Quantity).
		SetOperator(tx.User.GetX(ctx, int(req.OperatorID)))

	if a.Edges.Location != nil {
		builder.SetFromLocation(a.Edges.Location).
			SetFromLocationName(a.Edges.Location.Name)
	}
	if req.Purpose != "" {
		builder.SetReason(req.Purpose)
	}
	if req.Notes != "" {
		builder.SetNote(req.Notes)
	}

	record, err := builder.Save(ctx)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create inventory record",
		})
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":          record.ID,
		"asset_id":    a.ID,
		"asset_no":    a.AssetNo,
		"quantity":    record.Quantity,
		"record_type": string(record.RecordType),
		"created_at":  record.CreatedAt,
	})
}

// Transfer 资产调拨
func (h *Handler) Transfer(c *fiber.Ctx) error {
	var req TransferRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 开始事务
	tx, err := h.client.Tx(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to start transaction",
		})
	}

	// 检查资产
	a, err := tx.Asset.Query().
		Where(asset.ID(int(req.AssetID))).
		Only(ctx)
	if err != nil {
		tx.Rollback()
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Asset not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query asset",
		})
	}

	// 检查目标库位
	toLoc, err := tx.Location.Query().
		Where(location.ID(int(req.ToLocationID))).
		Only(ctx)
	if err != nil {
		tx.Rollback()
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Target location not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query location",
		})
	}

	// 更新资产库位
	_, err = tx.Asset.UpdateOne(a).
		SetLocation(toLoc).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update asset",
		})
	}

	// 创建调拨记录
	fromLoc, _ := tx.Location.Query().
		Where(location.ID(int(req.FromLocationID))).
		Only(ctx)

	recordNo := fmt.Sprintf("TRF-%s-%d", time.Now().Format("20060102150405"), req.AssetID)
	record, err := tx.InventoryRecord.Create().
		SetRecordNo(recordNo).
		SetAsset(a).
		SetFromLocation(fromLoc).
		SetFromLocationName(fromLoc.Name).
		SetToLocation(toLoc).
		SetToLocationName(toLoc.Name).
		SetRecordType(inventoryrecord.RecordTypeTransfer).
		SetQuantity(req.Quantity).
		SetOperator(tx.User.GetX(ctx, int(req.OperatorID))).
		SetNillableNote(&req.Notes).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create inventory record",
		})
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":          record.ID,
		"asset_id":    a.ID,
		"asset_no":    a.AssetNo,
		"to_location": toLoc.Name,
		"quantity":    record.Quantity,
		"record_type": string(record.RecordType),
		"created_at":  record.CreatedAt,
	})
}

// GetStock 查询库存统计
func (h *Handler) GetStock(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 统计各状态资产数量
	totalCount, _ := h.client.Asset.Query().Count(ctx)
	inStockCount, _ := h.client.Asset.Query().
		Where(asset.StatusEQ("in_stock")).
		Count(ctx)
	deployedCount, _ := h.client.Asset.Query().
		Where(asset.StatusEQ("deployed")).
		Count(ctx)
	maintenanceCount, _ := h.client.Asset.Query().
		Where(asset.StatusEQ("maintenance")).
		Count(ctx)
	retiredCount, _ := h.client.Asset.Query().
		Where(asset.StatusEQ("retired")).
		Count(ctx)

	return c.JSON(fiber.Map{
		"total":       totalCount,
		"in_stock":    inStockCount,
		"deployed":    deployedCount,
		"maintenance": maintenanceCount,
		"retired":     retiredCount,
		"timestamp":   time.Now(),
	})
}

// GetRecords 查询库存记录
func (h *Handler) GetRecords(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 20)
	if pageSize > 100 {
		pageSize = 100
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 构建查询
	query := h.client.InventoryRecord.Query().
		WithAsset().
		WithFromLocation().
		WithToLocation().
		WithOperator()

	// 获取总数
	total, err := query.Count(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to count records",
		})
	}

	// 分页查询
	records, err := query.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order(ent.Desc(inventoryrecord.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query records",
		})
	}

	// 转换为响应格式
	data := make([]fiber.Map, len(records))
	for i, r := range records {
		item := fiber.Map{
			"id":          r.ID,
			"record_type": string(r.RecordType),
			"quantity":    r.Quantity,
			"reason":      r.Reason,
			"note":        r.Note,
			"created_at":  r.CreatedAt,
		}

		if r.Edges.Asset != nil {
			item["asset"] = fiber.Map{
				"id":       r.Edges.Asset.ID,
				"asset_no": r.Edges.Asset.AssetNo,
				"name":     r.Edges.Asset.Name,
			}
		}

		if r.Edges.FromLocation != nil {
			item["from_location"] = fiber.Map{
				"id":   r.Edges.FromLocation.ID,
				"name": r.Edges.FromLocation.Name,
				"code": r.Edges.FromLocation.Code,
			}
		}

		if r.Edges.ToLocation != nil {
			item["to_location"] = fiber.Map{
				"id":   r.Edges.ToLocation.ID,
				"name": r.Edges.ToLocation.Name,
				"code": r.Edges.ToLocation.Code,
			}
		}

		if r.Edges.Operator != nil {
			item["operator"] = fiber.Map{
				"id":        r.Edges.Operator.ID,
				"username":  r.Edges.Operator.Username,
				"real_name": r.Edges.Operator.RealName,
			}
		}

		data[i] = item
	}

	totalPages := (total + pageSize - 1) / pageSize

	return c.JSON(fiber.Map{
		"data":        data,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	})
}
