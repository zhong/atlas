package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Role holds the schema definition for the Role entity.
type Role struct {
	ent.Schema
}

// Fields of the Role.
func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique().
			NotEmpty().
			Comment("角色名称"),
		field.String("code").
			Unique().
			NotEmpty().
			Comment("角色代码"),
		field.String("description").
			Optional().
			Comment("角色描述"),
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

// Edges of the Role.
func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("users", User.Type).
			Ref("roles").
			Comment("拥有此角色的用户"),
		edge.To("permissions", Permission.Type).
			Comment("角色权限"),
	}
}
