package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// RackUnit holds the schema definition for the RackUnit entity.
type RackUnit struct {
	ent.Schema
}

// Fields of the RackUnit.
func (RackUnit) Fields() []ent.Field {
	return []ent.Field{
		field.Int("unit_number").
			Comment("U位编号（1-42）"),
		field.Int("height").
			Default(1).
			Comment("占用高度（U数）"),
		field.Enum("status").
			Values("available", "occupied", "reserved").
			Default("available").
			Comment("状态：available-可用, occupied-已占用, reserved-已预留"),
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

// Edges of the RackUnit.
func (RackUnit) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("rack", Rack.Type).
			Ref("units").
			Unique().
			Required().
			Comment("所属机柜"),
		edge.To("asset", Asset.Type).
			Unique().
			Comment("安装的设备"),
	}
}

// Indexes of the RackUnit.
func (RackUnit) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("unit_number"),
		index.Fields("status"),
	}
}
