package ipam

import (
	"gorm.io/gorm"
	"time"
)

// Subnet 子网
type Subnet struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ParentID  *uint          `gorm:"index" json:"parent_id"`
	CIDR      string         `gorm:"size:43;not null" json:"cidr"` // IPv6 max length
	Name      string         `gorm:"size:64;not null" json:"name"`
	VlanID    string         `gorm:"size:16" json:"vlan_id"`
	Status    string         `gorm:"size:16;default:'active'" json:"status"` // active/inactive
	Path      string         `gorm:"size:255" json:"path"`                   // hierarchical path like "10.0.0.0/8.10.0.0/16"
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Subnet) TableName() string { return "cmdb_ipam.subnets" }

// IPAddress IP 地址
type IPAddress struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	SubnetID    uint           `gorm:"not null;index" json:"subnet_id"`
	IP          string         `gorm:"size:43;not null;index" json:"ip"`
	Status      string         `gorm:"size:16;default:'free'" json:"status"` // free/allocated/reserved/disabled
	CIID        *uint          `gorm:"index" json:"ci_id"`                   // associated CI instance
	AllocatedBy string         `gorm:"size:64" json:"allocated_by"`
	AllocatedAt *time.Time     `json:"allocated_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (IPAddress) TableName() string { return "cmdb_ipam.ip_addresses" }

// IPAMHistory IPAM 操作历史
type IPAMHistory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SubnetID  *uint     `json:"subnet_id"`
	IPID      *uint     `json:"ip_id"`
	Action    string    `gorm:"size:32;not null" json:"action"` // create/allocate/release/reserve
	Operator  string    `gorm:"size:64" json:"operator"`
	Detail    string    `gorm:"type:text" json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}

func (IPAMHistory) TableName() string { return "cmdb_ipam.histories" }

// UserIPAddress 用户被分配的IP地址（用户个人IP池）
type UserIPAddress struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	IPAddressID uint           `gorm:"not null;uniqueIndex" json:"ip_address_id"`
	AssignedBy  string         `gorm:"size:64" json:"assigned_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (UserIPAddress) TableName() string { return "cmdb_ipam.user_ip_addresses" }
