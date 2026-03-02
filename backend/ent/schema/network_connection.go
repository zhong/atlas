package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// NetworkConnection holds the schema definition for the NetworkConnection entity.
type NetworkConnection struct {
	ent.Schema
}

// Fields of the NetworkConnection.
func (NetworkConnection) Fields() []ent.Field {
	return []ent.Field{
		field.String("source_port").
			NotEmpty().
			Comment("源端口"),
		field.String("target_port").
			NotEmpty().
			Comment("目标端口"),
		field.Enum("connection_type").
			Values("ethernet", "infiniband", "fiber", "other").
			Default("ethernet").
			Comment("连接类型：ethernet-以太网, infiniband-IB网, fiber-光纤, other-其他"),
		field.String("speed").
			Optional().
			Comment("速率（如：25G, 100G）"),
		field.String("cable_type").
			Optional().
			Comment("线缆类型"),
		field.String("cable_length").
			Optional().
			Comment("线缆长度"),
		field.Enum("status").
			Values("active", "inactive", "maintenance").
			Default("active").
			Comment("状态"),
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

// Edges of the NetworkConnection.
func (NetworkConnection) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("source_equipment", Asset.Type).
			Ref("network_connections_source").
			Unique().
			Required().
			Comment("源设备"),
		edge.From("target_equipment", Asset.Type).
			Ref("network_connections_target").
			Unique().
			Required().
			Comment("目标设备"),
	}
}

// Indexes of the NetworkConnection.
func (NetworkConnection) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("connection_type"),
		index.Fields("status"),
	}
}
