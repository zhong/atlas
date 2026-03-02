package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// PurchaseOrder holds the schema definition for the PurchaseOrder entity.
type PurchaseOrder struct {
	ent.Schema
}

// Fields of the PurchaseOrder.
func (PurchaseOrder) Fields() []ent.Field {
	return []ent.Field{
		field.String("order_no").
			Unique().
			NotEmpty().
			Comment("订单编号"),
		field.String("project").
			Optional().
			Comment("项目名称"),
		field.Float("total_amount").
			Default(0).
			Comment("总金额"),
		field.Enum("status").
			Values("draft", "pending", "approved", "rejected", "ordered", "received", "completed", "cancelled").
			Default("draft").
			Comment("状态：draft-草稿, pending-待审批, approved-已批准, rejected-已拒绝, ordered-已下单, received-已收货, completed-已完成, cancelled-已取消"),
		field.Time("order_date").
			Optional().
			Nillable().
			Comment("下单日期"),
		field.Time("expected_date").
			Optional().
			Nillable().
			Comment("预计到货日期"),
		field.String("delivery_location").
			Optional().
			Comment("收货地址"),
		field.Text("note").
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

// Edges of the PurchaseOrder.
func (PurchaseOrder) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("supplier", Supplier.Type).
			Ref("purchase_orders").
			Unique().
			Required().
			Comment("供应商"),
		edge.To("items", OrderItem.Type).
			Comment("订单明细"),
		edge.From("creator", User.Type).
			Ref("created_purchase_orders").
			Unique().
			Required().
			Comment("创建人"),
		edge.To("approval", Approval.Type).
			Unique().
			Comment("审批记录"),
	}
}

// Indexes of the PurchaseOrder.
func (PurchaseOrder) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("order_no"),
		index.Fields("status"),
		index.Fields("order_date"),
		index.Fields("project"),
	}
}
