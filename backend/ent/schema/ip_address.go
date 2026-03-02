package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// IPAddress holds the schema definition for the IPAddress entity.
type IPAddress struct {
	ent.Schema
}

// Fields of the IPAddress.
func (IPAddress) Fields() []ent.Field {
	return []ent.Field{
		field.String("ip_address").
			NotEmpty().
			Comment("IP地址"),
		field.String("subnet").
			Optional().
			Comment("子网掩码/CIDR"),
		field.String("gateway").
			Optional().
			Comment("网关"),
		field.Enum("ip_type").
			Values("internal", "public", "management").
			Default("internal").
			Comment("IP类型：internal-内网, public-公网, management-管理网"),
		field.Enum("status").
			Values("allocated", "available", "reserved").
			Default("available").
			Comment("状态：allocated-已分配, available-可用, reserved-已预留"),
		field.String("vlan").
			Optional().
			Comment("VLAN"),
		field.Text("notes").
			Optional().
			Comment("备注"),
		field.Time("allocated_at").
			Optional().
			Nillable().
			Comment("分配时间"),
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

// Edges of the IPAddress.
func (IPAddress) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("asset", Asset.Type).
			Ref("ip_addresses").
			Unique().
			Comment("关联设备"),
	}
}

// Indexes of the IPAddress.
func (IPAddress) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("ip_address").Unique(),
		index.Fields("ip_type"),
		index.Fields("status"),
	}
}
