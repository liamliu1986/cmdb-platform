package core

import (
	"gorm.io/gorm"
	"time"
)

// CIType 配置项类型
type CIType struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"size:32;uniqueIndex;not null" json:"name"`
	Alias        string         `gorm:"size:32;not null" json:"alias"`
	UniqueAttrID uint           `gorm:"not null" json:"unique_attr_id"`
	Icon         string         `gorm:"size:255" json:"icon"`
	Enabled      bool           `gorm:"default:true" json:"enabled"`
	IsBuiltin    bool           `gorm:"default:false" json:"is_builtin"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CIType) TableName() string { return "cmdb_core.ci_types" }

// Attribute 属性定义
type Attribute struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:32;uniqueIndex;not null" json:"name"`
	Alias        string    `gorm:"size:32;not null" json:"alias"`
	ValueType    string    `gorm:"size:16;not null" json:"value_type"`
	IsChoice     bool      `gorm:"default:false" json:"is_choice"`
	IsList       bool      `gorm:"default:false" json:"is_list"`
	IsUnique     bool      `gorm:"default:false" json:"is_unique"`
	IsIndex      bool      `gorm:"default:false" json:"is_index"`
	IsPassword   bool      `gorm:"default:false" json:"is_password"`
	IsComputed   bool      `gorm:"default:false" json:"is_computed"`
	ComputeExpr  string    `gorm:"type:text" json:"compute_expr,omitempty"`
	DefaultValue string    `gorm:"type:jsonb" json:"default_value,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

func (Attribute) TableName() string { return "cmdb_core.attributes" }

// CITypeAttribute CIType 与属性关联
type CITypeAttribute struct {
	ID          uint `gorm:"primaryKey" json:"id"`
	CITypeID    uint `gorm:"not null;index" json:"ci_type_id"`
	AttributeID uint `gorm:"not null;index" json:"attribute_id"`
	Order       int  `gorm:"default:0" json:"order"`
	IsRequired  bool `gorm:"default:false" json:"is_required"`
	DefaultShow bool `gorm:"default:true" json:"default_show"`
}

func (CITypeAttribute) TableName() string { return "cmdb_core.ci_type_attributes" }

// RelationType 关系类型
type RelationType struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"size:16;uniqueIndex;not null" json:"name"`
	Description string `gorm:"size:255" json:"description"`
}

func (RelationType) TableName() string { return "cmdb_core.relation_types" }

// CITypeRelation CIType 间关系模板
type CITypeRelation struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	ParentCITypeID uint   `gorm:"not null;index" json:"parent_ci_type_id"`
	ChildCITypeID  uint   `gorm:"not null;index" json:"child_ci_type_id"`
	RelationTypeID uint   `gorm:"not null" json:"relation_type_id"`
	Constraint     string `gorm:"size:16;default:'one2many'" json:"constraint"`
}

func (CITypeRelation) TableName() string { return "cmdb_core.ci_type_relations" }

// CI 配置项实例 - 使用 JSONB 存储属性值
type CI struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	CITypeID        uint           `gorm:"not null;index" json:"ci_type_id"`
	Status          string         `gorm:"size:16;default:'active'" json:"status"`
	AttrValuesRaw   string         `gorm:"column:attr_values;type:jsonb;not null;default:'{}'" json:"-"`
	IsAutoDiscovery bool           `gorm:"default:false" json:"is_auto_discovery"`
	UpdatedBy       string         `gorm:"size:64" json:"updated_by"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CI) TableName() string { return "cmdb_core.cis" }

// CIRelation CI 实例间关系
type CIRelation struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	FirstCIID      uint      `gorm:"not null;index" json:"first_ci_id"`
	SecondCIID     uint      `gorm:"not null;index" json:"second_ci_id"`
	RelationTypeID uint      `gorm:"not null" json:"relation_type_id"`
	CreatedAt      time.Time `json:"created_at"`
}

func (CIRelation) TableName() string { return "cmdb_core.ci_relations" }

// OperationLog 操作审计日志
type OperationLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TargetType string    `gorm:"size:32;not null" json:"target_type"`
	TargetID   uint      `gorm:"not null" json:"target_id"`
	Action     string    `gorm:"size:32;not null" json:"action"`
	Operator   string    `gorm:"size:64;not null" json:"operator"`
	OldValue   string    `gorm:"type:jsonb" json:"old_value,omitempty"`
	NewValue   string    `gorm:"type:jsonb" json:"new_value,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

func (OperationLog) TableName() string { return "cmdb_core.operation_logs" }
