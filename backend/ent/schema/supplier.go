package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Supplier holds the schema definition for the Supplier entity.
type Supplier struct {
	ent.Schema
}

// Fields of the Supplier.
func (Supplier) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Comment("供应商名称"),
		field.String("code").
			Unique().
			NotEmpty().
			Comment("供应商编码"),
		field.JSON("category_specialties", []string{}).
			Optional().
			Comment("专长类别（如：网络设备、光模块等）"),
		field.String("contact_person").
			Optional().
			Comment("联系人"),
		field.String("contact_phone").
			Optional().
			Comment("联系电话"),
		field.String("contact_email").
			Optional().
			Comment("联系邮箱"),
		field.String("address").
			Optional().
			Comment("地址"),
		field.Float("rating").
			Optional().
			Default(0).
			Comment("评分（0-5）"),
		field.Enum("status").
			Values("active", "inactive", "blacklist").
			Default("active").
			Comment("状态：active-活跃, inactive-停用, blacklist-黑名单"),
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

// Edges of the Supplier.
func (Supplier) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("purchase_orders", PurchaseOrder.Type).
			Comment("采购订单"),
		edge.To("repair_vendors", RepairVendor.Type).
			Comment("维修供应商信息"),
	}
}

// Indexes of the Supplier.
func (Supplier) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status"),
		index.Fields("rating"),
	}
}
