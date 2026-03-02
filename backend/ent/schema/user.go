package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").
			Unique().
			NotEmpty().
			Comment("用户名"),
		field.String("password").
			Sensitive().
			NotEmpty().
			Comment("密码（加密存储）"),
		field.String("email").
			Unique().
			NotEmpty().
			Comment("邮箱"),
		field.String("phone").
			Optional().
			Comment("电话"),
		field.String("real_name").
			NotEmpty().
			Comment("真实姓名"),
		field.String("department").
			Optional().
			Comment("部门"),
		field.Enum("status").
			Values("active", "inactive", "locked").
			Default("active").
			Comment("状态：active-活跃, inactive-停用, locked-锁定"),
		field.Time("last_login_at").
			Optional().
			Nillable().
			Comment("最后登录时间"),
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

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("roles", Role.Type).
			Comment("用户角色"),
		edge.To("created_purchase_orders", PurchaseOrder.Type).
			Comment("创建的采购订单"),
		edge.To("inventory_records", InventoryRecord.Type).
			Comment("操作的库存记录"),
		edge.To("initiated_approvals", Approval.Type).
			Comment("发起的审批"),
		edge.To("repair_tickets", RepairTicket.Type).
			Comment("报修工单"),
	}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email"),
		index.Fields("status"),
		index.Fields("department"),
	}
}
