package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// DataCenter holds the schema definition for the DataCenter entity.
type DataCenter struct {
	ent.Schema
}

// Fields of the DataCenter.
func (DataCenter) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Comment("IDC名称"),
		field.String("code").
			Unique().
			NotEmpty().
			Comment("IDC编码"),
		field.String("location").
			NotEmpty().
			Comment("地理位置"),
		field.String("address").
			Optional().
			Comment("详细地址"),
		field.String("contact").
			Optional().
			Comment("联系人"),
		field.String("phone").
			Optional().
			Comment("联系电话"),
		field.Enum("status").
			Values("active", "inactive", "maintenance").
			Default("active").
			Comment("状态"),
		field.Text("description").
			Optional().
			Comment("描述"),
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

// Edges of the DataCenter.
func (DataCenter) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("rooms", Room.Type).
			Comment("机房"),
	}
}
