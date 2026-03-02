package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// AssetType holds the schema definition for the AssetType entity.
type AssetType struct {
	ent.Schema
}

// Fields of the AssetType.
func (AssetType) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Comment("资产类型名称"),
		field.String("code").
			Unique().
			NotEmpty().
			Comment("资产类型代码"),
		field.Enum("category").
			Values("server", "switch", "network_card", "storage", "component", "other").
			Comment("类别：server-服务器, switch-交换机, network_card-网卡, storage-存储, component-组件, other-其他"),
		field.Text("description").
			Optional().
			Comment("描述"),
		field.JSON("default_specs", map[string]interface{}{}).
			Optional().
			Comment("默认规格参数"),
		field.Int("sort_order").
			Default(0).
			Comment("排序"),
		field.Enum("status").
			Values("active", "inactive").
			Default("active").
			Comment("状态"),
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

// Edges of the AssetType.
func (AssetType) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("assets", Asset.Type).
			Comment("此类型的资产"),
	}
}
