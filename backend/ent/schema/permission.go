package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Permission holds the schema definition for the Permission entity.
type Permission struct {
	ent.Schema
}

// Fields of the Permission.
func (Permission) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique().
			NotEmpty().
			Comment("权限名称"),
		field.String("code").
			Unique().
			NotEmpty().
			Comment("权限代码"),
		field.String("resource").
			NotEmpty().
			Comment("资源"),
		field.String("action").
			NotEmpty().
			Comment("操作：create, read, update, delete"),
		field.String("description").
			Optional().
			Comment("权限描述"),
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

// Edges of the Permission.
func (Permission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("roles", Role.Type).
			Ref("permissions").
			Comment("拥有此权限的角色"),
	}
}
