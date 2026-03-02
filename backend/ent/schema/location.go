package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Location holds the schema definition for the Location entity.
type Location struct {
	ent.Schema
}

// Fields of the Location.
func (Location) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Comment("库位名称"),
		field.String("code").
			Unique().
			NotEmpty().
			Comment("库位编码"),
		field.String("location_code").
			Optional().
			Comment("位置代码（如：2号AI库、1号库等）"),
		field.Int("parent_location_id").
			Optional().
			Nillable().
			Comment("父级库位ID，支持层级结构"),
		field.Enum("status").
			Values("available", "full", "maintenance", "reserved").
			Default("available").
			Comment("状态：available-可用, full-已满, maintenance-维护中, reserved-已预留"),
		field.Text("description").
			Optional().
			Comment("描述"),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("创建时间"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("更新时间"),
	}
}

// Edges of the Location.
func (Location) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("warehouse", Warehouse.Type).
			Ref("locations").
			Unique().
			Required().
			Comment("所属仓库"),
		edge.To("assets", Asset.Type).
			Comment("库位中的资产"),
		edge.To("inventory_records_from", InventoryRecord.Type).
			Comment("从此库位出库的记录"),
		edge.To("inventory_records_to", InventoryRecord.Type).
			Comment("入库到此库位的记录"),
	}
}

// Indexes of the Location.
func (Location) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status"),
		index.Fields("location_code"),
	}
}
