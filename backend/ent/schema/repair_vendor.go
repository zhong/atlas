package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// RepairVendor holds the schema definition for the RepairVendor entity.
type RepairVendor struct {
	ent.Schema
}

// Fields of the RepairVendor.
func (RepairVendor) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Comment("供应商名称（如：四通、超融核、融科联创、安擎、麦芒）"),
		field.JSON("supported_models", []string{}).
			Optional().
			Comment("支持的设备型号"),
		field.Int("sla_hours").
			Default(24).
			Comment("SLA响应时间（小时）"),
		field.String("contact_person").
			Optional().
			Comment("联系人"),
		field.String("contact_phone").
			Optional().
			Comment("联系电话"),
		field.String("contact_email").
			Optional().
			Comment("联系邮箱"),
		field.Enum("status").
			Values("active", "inactive").
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

// Edges of the RepairVendor.
func (RepairVendor) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("supplier", Supplier.Type).
			Ref("repair_vendors").
			Unique().
			Comment("关联的供应商"),
		edge.To("repair_tickets", RepairTicket.Type).
			Comment("维修工单"),
	}
}
