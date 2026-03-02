package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Warehouse holds the schema definition for the Warehouse entity.
type Warehouse struct {
	ent.Schema
}

// Fields of the Warehouse.
func (Warehouse) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Comment("仓库名称"),
		field.String("code").
			Unique().
			NotEmpty().
			Comment("仓库编码"),
		field.Enum("warehouse_type").
			Values("idc", "warehouse", "office").
			Default("warehouse").
			Comment("仓库类型：idc-数据中心, warehouse-仓库, office-办公室"),
		field.String("location").
			NotEmpty().
			Comment("地理位置"),
		field.String("address").
			Optional().
			Comment("详细地址"),
		field.String("contact").
			Optional().
			Comment("联系人"),
		field.String("phone").
			Optional().
			Comment("联系电话"),
		field.Enum("status").
			Values("active", "inactive").
			Default("active").
			Comment("状态"),
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

// Edges of the Warehouse.
func (Warehouse) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("locations", Location.Type).
			Comment("仓库库位"),
	}
}

// Indexes of the Warehouse.
func (Warehouse) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("warehouse_type"),
		index.Fields("status"),
	}
}
