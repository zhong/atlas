package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Rack holds the schema definition for the Rack entity.
type Rack struct {
	ent.Schema
}

// Fields of the Rack.
func (Rack) Fields() []ent.Field {
	return []ent.Field{
		field.String("rack_no").
			NotEmpty().
			Comment("机柜编号"),
		field.String("position_code").
			Optional().
			Comment("位置代码（如：A01, M13）"),
		field.String("row").
			Optional().
			Comment("行（如：A, B, M）"),
		field.Int("column").
			Optional().
			Comment("列（如：01, 13）"),
		field.Int("total_units").
			Default(42).
			Comment("总U数"),
		field.Float("power_capacity").
			Optional().
			Comment("电力容量（kW）"),
		field.Float("power_used").
			Default(0).
			Comment("已用电力（kW）"),
		field.Int("weight_capacity").
			Optional().
			Comment("承重（kg）"),
		field.Enum("status").
			Values("available", "full", "maintenance", "reserved").
			Default("available").
			Comment("状态：available-可用, full-已满, maintenance-维护中, reserved-已预留"),
		field.JSON("position", map[string]interface{}{}).
			Optional().
			Comment("位置坐标（用于平面图显示）"),
		field.Text("notes").
			Optional().
			Comment("备注"),
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

// Edges of the Rack.
func (Rack) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("room", Room.Type).
			Ref("racks").
			Unique().
			Required().
			Comment("所属机房"),
		edge.To("units", RackUnit.Type).
			Comment("U位"),
	}
}

// Indexes of the Rack.
func (Rack) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status"),
		index.Fields("position_code"),
		index.Fields("row", "column"),
	}
}
