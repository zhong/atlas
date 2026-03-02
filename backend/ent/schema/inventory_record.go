package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// InventoryRecord holds the schema definition for the InventoryRecord entity.
type InventoryRecord struct {
	ent.Schema
}

// Fields of the InventoryRecord.
func (InventoryRecord) Fields() []ent.Field {
	return []ent.Field{
		field.String("record_no").
			Unique().
			NotEmpty().
			Comment("记录编号"),
		field.Enum("record_type").
			Values("inbound", "outbound", "transfer", "adjust", "check", "borrow", "return_item").
			Comment("记录类型：inbound-入库, outbound-出库, transfer-调拨, adjust-调整, check-盘点, borrow-借出, return_item-归还"),
		field.Int("quantity").
			Default(1).
			Comment("数量"),
		field.String("reason").
			Optional().
			Comment("原因"),
		field.Text("note").
			Optional().
			Comment("备注"),
		field.String("from_location_name").
			Optional().
			Comment("源库位名称（冗余字段，便于查询）"),
		field.String("to_location_name").
			Optional().
			Comment("目标库位名称（冗余字段，便于查询）"),
		field.Time("transaction_date").
			Default(time.Now).
			Comment("交易日期"),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("创建时间"),
	}
}

// Edges of the InventoryRecord.
func (InventoryRecord) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("asset", Asset.Type).
			Ref("inventory_records").
			Unique().
			Required().
			Comment("关联资产"),
		edge.From("from_location", Location.Type).
			Ref("inventory_records_from").
			Unique().
			Comment("源库位"),
		edge.From("to_location", Location.Type).
			Ref("inventory_records_to").
			Unique().
			Comment("目标库位"),
		edge.From("operator", User.Type).
			Ref("inventory_records").
			Unique().
			Required().
			Comment("操作人"),
	}
}

// Indexes of the InventoryRecord.
func (InventoryRecord) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("record_type"),
		index.Fields("transaction_date"),
		index.Fields("record_no"),
	}
}
