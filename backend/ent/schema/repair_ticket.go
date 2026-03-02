package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// RepairTicket holds the schema definition for the RepairTicket entity.
type RepairTicket struct {
	ent.Schema
}

// Fields of the RepairTicket.
func (RepairTicket) Fields() []ent.Field {
	return []ent.Field{
		field.String("ticket_number").
			Unique().
			NotEmpty().
			Comment("工单编号"),
		field.String("gpu_model").
			Optional().
			Comment("GPU型号"),
		field.String("serial_number").
			Optional().
			Comment("序列号"),
		field.Text("fault_description").
			NotEmpty().
			Comment("故障描述"),
		field.Enum("fault_type").
			Values("hardware", "software", "performance", "other").
			Default("hardware").
			Comment("故障类型：hardware-硬件, software-软件, performance-性能, other-其他"),
		field.Enum("severity").
			Values("low", "medium", "high", "critical").
			Default("medium").
			Comment("严重程度：low-低, medium-中, high-高, critical-紧急"),
		field.Enum("status").
			Values("reported", "diagnosed", "in_repair", "testing", "resolved", "closed").
			Default("reported").
			Comment("状态：reported-已报修, diagnosed-已诊断, in_repair-维修中, testing-测试中, resolved-已解决, closed-已关闭"),
		field.Int("assigned_to_id").
			Optional().
			Nillable().
			Comment("分配给（工程师ID）"),
		field.Time("reported_at").
			Default(time.Now).
			Comment("报修时间"),
		field.Time("resolved_at").
			Optional().
			Nillable().
			Comment("解决时间"),
		field.Time("sla_deadline").
			Optional().
			Nillable().
			Comment("SLA截止时间"),
		field.Text("resolution_notes").
			Optional().
			Comment("解决方案/备注"),
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

// Edges of the RepairTicket.
func (RepairTicket) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("asset", Asset.Type).
			Ref("repair_tickets").
			Unique().
			Required().
			Comment("关联设备"),
		edge.From("vendor", RepairVendor.Type).
			Ref("repair_tickets").
			Unique().
			Comment("维修供应商"),
		edge.From("reported_by", User.Type).
			Ref("repair_tickets").
			Unique().
			Required().
			Comment("报修人"),
	}
}

// Indexes of the RepairTicket.
func (RepairTicket) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("ticket_number"),
		index.Fields("status"),
		index.Fields("severity"),
		index.Fields("reported_at"),
	}
}
