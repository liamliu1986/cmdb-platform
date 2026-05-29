package dcim

import (
	"cmdb-api/database"
	"gorm.io/gorm"
	"time"
)

// LocationCoord 地图坐标
type LocationCoord struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CIID        uint           `gorm:"uniqueIndex;not null" json:"ci_id"` // IDC CI instance ID
	Lat         float64        `gorm:"not null" json:"lat"`               // 纬度
	Lng         float64        `gorm:"not null" json:"lng"`               // 经度
	Address     string         `gorm:"size:255" json:"address"`
	MapProvider string         `gorm:"size:16;default:'amap'" json:"map_provider"` // amap/baidu
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (LocationCoord) TableName() string { return "cmdb_dcim.location_coords" }

// locationRepository
type locationRepo struct{}

func (r *locationRepo) SetCoord(coord *LocationCoord) error {
	var existing LocationCoord
	err := database.DB.Where("ci_id = ?", coord.CIID).First(&existing).Error
	if err == nil {
		return database.DB.Model(&LocationCoord{}).Where("ci_id = ?", coord.CIID).Updates(map[string]interface{}{
			"lat": coord.Lat, "lng": coord.Lng, "address": coord.Address,
		}).Error
	}
	return database.DB.Create(coord).Error
}

func (r *locationRepo) GetCoord(ciID uint) (*LocationCoord, error) {
	var coord LocationCoord
	err := database.DB.Where("ci_id = ?", ciID).First(&coord).Error
	return &coord, err
}

func (r *locationRepo) ListCoords() ([]LocationCoord, error) {
	var coords []LocationCoord
	err := database.DB.Find(&coords).Error
	return coords, err
}
