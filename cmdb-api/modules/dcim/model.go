package dcim

import (
	"gorm.io/gorm"
	"time"
)

// IDC 数据中心
type IDC struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:64;not null" json:"name"`
	Address   string         `gorm:"size:255" json:"address"`
	Contact   string         `gorm:"size:64" json:"contact"`
	RegionID  uint           `gorm:"index" json:"region_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (IDC) TableName() string { return "cmdb_dcim.idcs" }

// ServerRoom 机房
type ServerRoom struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:64;not null" json:"name"`
	Floor     string         `gorm:"size:16" json:"floor"`
	IDCID     uint           `gorm:"not null;index" json:"idc_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ServerRoom) TableName() string { return "cmdb_dcim.server_rooms" }

// Rack 机柜
type Rack struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:64;not null" json:"name"`
	TotalU    int            `gorm:"not null;default:42" json:"total_u"`
	RoomID    uint           `gorm:"not null;index" json:"room_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Rack) TableName() string { return "cmdb_dcim.racks" }

// RackLayout 机柜U位占用
type RackLayout struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	RackID     uint      `gorm:"not null;index" json:"rack_id"`
	UPosition  int       `gorm:"not null" json:"u_position"`           // U位编号，从1开始
	DeviceCIID uint      `gorm:"not null;index" json:"device_ci_id"`   // 设备CI的ID
	Direction  string    `gorm:"size:16;default:'front'" json:"direction"` // front/back
	CreatedAt  time.Time `json:"created_at"`
}

func (RackLayout) TableName() string { return "cmdb_dcim.rack_layouts" }
