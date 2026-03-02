package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Room holds the schema definition for the Room entity.
type Room struct {
	ent.Schema
}

// Fields of the Room.
func (Room) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Comment("机房名称"),
		field.String("code").
			NotEmpty().
			Comment("机房编码"),
		field.String("floor").
			Optional().
			Comment("楼层"),
		field.Float("area").
			Optional().
			Comment("面积（平方米）"),
		field.Float("power_capacity").
			Optional().
			Comment("电力容量（kW）"),
		field.Float("power_used").
			Default(0).
			Comment("已用电力（kW）"),
		field.Enum("status").
			Values("active", "inactive", "maintenance").
			Default("active").
			Comment("状态"),
		field.JSON("layout_data", map[string]interface{}{}).
			Optional().
			Comment("平面图数据（JSON格式）"),
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

// Edges of the Room.
func (Room) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("data_center", DataCenter.Type).
			Ref("rooms").
			Unique().
			Required().
			Comment("所属IDC"),
		edge.To("racks", Rack.Type).
			Comment("机柜"),
	}
}
