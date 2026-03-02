package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// OrderItem holds the schema definition for the OrderItem entity.
type OrderItem struct {
	ent.Schema
}

// Fields of the OrderItem.
func (OrderItem) Fields() []ent.Field {
	return []ent.Field{
		field.String("category").
			NotEmpty().
			Comment("类别"),
		field.Text("model_config").
			Comment("型号/配置要求"),
		field.String("warranty_period").
			Optional().
			Comment("质保期（如：3年）"),
		field.String("delivery_timeline").
			Optional().
			Comment("货期（如：PO+7）"),
		field.Int("required_qty").
			Default(0).
			Comment("需求数量"),
		field.Int("spare_qty").
			Default(0).
			Comment("备件需求数量"),
		field.Int("available_stock").
			Default(0).
			Comment("可用库存"),
		field.Int("purchase_qty").
			Default(0).
			Comment("采购数量（计算：需求数量+备件需求-可用库存）"),
		field.Float("unit_price").
			Optional().
			Default(0).
			Comment("单价"),
		field.Float("total_price").
			Optional().
			Default(0).
			Comment("总价"),
		field.String("purpose").
			Optional().
			Comment("用途"),
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

// Edges of the OrderItem.
func (OrderItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("purchase_order", PurchaseOrder.Type).
			Ref("items").
			Unique().
			Required().
			Comment("所属采购订单"),
	}
}
