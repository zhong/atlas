package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Asset holds the schema definition for the Asset entity.
type Asset struct {
	ent.Schema
}

// Fields of the Asset.
func (Asset) Fields() []ent.Field {
	return []ent.Field{
		field.String("asset_no").
			Unique().
			NotEmpty().
			Comment("资产编号"),
		field.String("name").
			NotEmpty().
			Comment("资产名称"),
		field.String("brand").
			Optional().
			Comment("品牌"),
		field.String("model").
			Optional().
			Comment("型号"),
		field.String("sn").
			Unique().
			Optional().
			Comment("序列号"),
		field.Enum("project_zone").
			Values("ai_cloud", "hpc_cloud", "industry_cloud", "other").
			Optional().
			Comment("项目分区：ai_cloud-AI云, hpc_cloud-超算云, industry_cloud-行业云, other-其他"),
		field.Enum("status").
			Values("in_stock", "deployed", "borrowed", "in_transit", "maintenance", "retired").
			Default("in_stock").
			Comment("状态：in_stock-在库, deployed-已部署, borrowed-已借出, in_transit-运输中, maintenance-维修中, retired-已报废"),
		field.Enum("borrow_status").
			Values("available", "borrowed", "reserved").
			Default("available").
			Optional().
			Comment("借用状态：available-可借, borrowed-已借出, reserved-已预约"),
		field.Int("borrowed_by_id").
			Optional().
			Nillable().
			Comment("借用人ID"),
		field.Time("borrowed_at").
			Optional().
			Nillable().
			Comment("借出时间"),
		field.Time("expected_return_at").
			Optional().
			Nillable().
			Comment("预计归还时间"),
		field.JSON("specs", map[string]interface{}{}).
			Optional().
			Comment("规格参数（JSON格式）"),
		field.Time("purchase_date").
			Optional().
			Nillable().
			Comment("采购日期"),
		field.Time("warranty_expire_date").
			Optional().
			Nillable().
			Comment("保修到期日"),
		field.Float("purchase_price").
			Optional().
			Comment("采购价格"),
		field.Text("description").
			Optional().
			Comment("描述"),
		field.Text("notes").
			Optional().
			Comment("备注（记录设备生命周期事件）"),
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

// Edges of the Asset.
func (Asset) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("asset_type", AssetType.Type).
			Ref("assets").
			Unique().
			Required().
			Comment("资产类型"),
		edge.From("location", Location.Type).
			Ref("assets").
			Unique().
			Comment("当前库位"),
		edge.From("rack_unit", RackUnit.Type).
			Ref("asset").
			Unique().
			Comment("机柜U位（如果已上架）"),
		edge.To("inventory_records", InventoryRecord.Type).
			Comment("库存变动记录"),
		edge.To("repair_tickets", RepairTicket.Type).
			Comment("维修工单"),
		edge.To("network_connections_source", NetworkConnection.Type).
			Comment("作为源设备的网络连接"),
		edge.To("network_connections_target", NetworkConnection.Type).
			Comment("作为目标设备的网络连接"),
		edge.To("ip_addresses", IPAddress.Type).
			Comment("IP地址"),
	}
}

// Indexes of the Asset.
func (Asset) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("asset_no"),
		index.Fields("sn"),
		index.Fields("status"),
		index.Fields("project_zone"),
		index.Fields("borrow_status"),
		index.Fields("brand", "model"),
	}
}
