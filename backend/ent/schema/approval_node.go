package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ApprovalNode holds the schema definition for the ApprovalNode entity.
type ApprovalNode struct {
	ent.Schema
}

// Fields of the ApprovalNode.
func (ApprovalNode) Fields() []ent.Field {
	return []ent.Field{
		field.Int("node_order").
			Comment("节点顺序"),
		field.String("approver_name").
			NotEmpty().
			Comment("审批人姓名"),
		field.String("approver_email").
			Optional().
			Comment("审批人邮箱"),
		field.Enum("status").
			Values("pending", "approved", "rejected", "skipped").
			Default("pending").
			Comment("状态：pending-待审批, approved-已批准, rejected-已拒绝, skipped-已跳过"),
		field.Text("comment").
			Optional().
			Comment("审批意见"),
		field.Time("approved_at").
			Optional().
			Nillable().
			Comment("审批时间"),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("创建时间"),
	}
}

// Edges of the ApprovalNode.
func (ApprovalNode) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("approval", Approval.Type).
			Ref("nodes").
			Unique().
			Required().
			Comment("所属审批"),
	}
}
