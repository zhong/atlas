package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Approval holds the schema definition for the Approval entity.
type Approval struct {
	ent.Schema
}

// Fields of the Approval.
func (Approval) Fields() []ent.Field {
	return []ent.Field{
		field.String("approval_no").
			Unique().
			NotEmpty().
			Comment("审批编号"),
		field.String("entity_type").
			NotEmpty().
			Comment("实体类型：purchase_order, asset_transfer等"),
		field.Int("entity_id").
			Comment("实体ID"),
		field.Enum("status").
			Values("pending", "approved", "rejected", "cancelled").
			Default("pending").
			Comment("状态：pending-待审批, approved-已批准, rejected-已拒绝, cancelled-已取消"),
		field.Int("current_node").
			Default(0).
			Comment("当前节点"),
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

// Edges of the Approval.
func (Approval) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("nodes", ApprovalNode.Type).
			Comment("审批节点"),
		edge.From("initiator", User.Type).
			Ref("initiated_approvals").
			Unique().
			Required().
			Comment("发起人"),
		edge.From("purchase_order", PurchaseOrder.Type).
			Ref("approval").
			Unique().
			Comment("关联的采购订单"),
	}
}

// Indexes of the Approval.
func (Approval) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("approval_no"),
		index.Fields("status"),
		index.Fields("entity_type", "entity_id"),
	}
}
