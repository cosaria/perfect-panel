package catalog

import "time"

const (
	RevisionName          = "0003_catalog_node_relations"
	nodeGroupTypeTag      = "tag"
	nodeGroupTypeExplicit = "explicit_nodes"
	tagGroupCodePrefix    = "tag:"
)

type NodeGroup struct {
	Id        int64     `gorm:"primaryKey"`
	Code      string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_node_groups_code"`
	Name      string    `gorm:"type:varchar(255);not null;default:''"`
	Type      string    `gorm:"type:varchar(64);not null;index:idx_node_groups_type"`
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func (NodeGroup) TableName() string {
	return "node_groups"
}

type NodeGroupNode struct {
	Id          int64     `gorm:"primaryKey"`
	NodeGroupId int64     `gorm:"not null;uniqueIndex:idx_node_group_nodes_unique;index:idx_node_group_nodes_group"`
	NodeId      int64     `gorm:"not null;uniqueIndex:idx_node_group_nodes_unique;index:idx_node_group_nodes_node"`
	CreatedAt   time.Time `gorm:"<-:create"`
	UpdatedAt   time.Time
}

func (NodeGroupNode) TableName() string {
	return "node_group_nodes"
}

type PlanNodeGroupRule struct {
	Id          int64     `gorm:"primaryKey"`
	SubscribeId int64     `gorm:"not null;uniqueIndex:idx_plan_node_group_rules_unique;index:idx_plan_node_group_rules_subscribe"`
	NodeGroupId int64     `gorm:"not null;uniqueIndex:idx_plan_node_group_rules_unique;index:idx_plan_node_group_rules_group"`
	RuleType    string    `gorm:"type:varchar(64);not null;index:idx_plan_node_group_rules_type"`
	RuleValue   string    `gorm:"type:varchar(255);not null;default:''"`
	CreatedAt   time.Time `gorm:"<-:create"`
	UpdatedAt   time.Time
}

func (PlanNodeGroupRule) TableName() string {
	return "plan_node_group_rules"
}
