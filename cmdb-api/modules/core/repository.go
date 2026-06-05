package core

import (
	"encoding/json"
	"fmt"

	"cmdb-api/database"
)

var allowedTables = map[string]bool{
	"cmdb_core.cis":         true,
	"cmdb_core.ci_types":    true,
	"cmdb_discovery.rules":  true,
	"cmdb_discovery.agents": true,
}

type CoreRepository struct{}

func NewCoreRepository() *CoreRepository {
	return &CoreRepository{}
}

// CIType
func (r *CoreRepository) CreateCIType(ct *CIType) error {
	return database.DB.Create(ct).Error
}

func (r *CoreRepository) GetCITypeByID(id uint) (*CIType, error) {
	var ct CIType
	err := database.DB.First(&ct, id).Error
	return &ct, err
}

func (r *CoreRepository) GetCITypeWithAttributes(id uint) (*CIType, error) {
	var ct CIType
	err := database.DB.Preload("Attributes").First(&ct, id).Error
	return &ct, err
}

func (r *CoreRepository) GetCITypeByName(name string) (*CIType, error) {
	var ct CIType
	err := database.DB.Where("name = ?", name).First(&ct).Error
	return &ct, err
}

func (r *CoreRepository) ListCITypes() ([]CIType, error) {
	var cts []CIType
	err := database.DB.Find(&cts).Error
	return cts, err
}

func (r *CoreRepository) UpdateCIType(id uint, updates map[string]interface{}) error {
	return database.DB.Model(&CIType{}).Where("id = ?", id).Updates(updates).Error
}

func (r *CoreRepository) DeleteCIType(id uint) error {
	return database.DB.Delete(&CIType{}, id).Error
}

// Attribute
func (r *CoreRepository) CreateAttribute(attr *Attribute) error {
	return database.DB.Create(attr).Error
}

func (r *CoreRepository) GetAttributeByID(id uint) (*Attribute, error) {
	var attr Attribute
	err := database.DB.First(&attr, id).Error
	return &attr, err
}

// CITypeAttribute
func (r *CoreRepository) AddAttributeToCIType(cta *CITypeAttribute) error {
	return database.DB.Create(cta).Error
}

func (r *CoreRepository) GetCITypeAttributes(ciTypeID uint) ([]Attribute, error) {
	var attrs []Attribute
	err := database.DB.
		Joins("JOIN cmdb_core.ci_type_attributes ON ci_type_attributes.attribute_id = attributes.id").
		Where("ci_type_attributes.ci_type_id = ?", ciTypeID).
		Order("ci_type_attributes.order").
		Find(&attrs).Error
	return attrs, err
}

// CI
func (r *CoreRepository) CreateCI(ci *CI) error {
	raw, _ := json.Marshal(ci.AttrValuesRaw)
	// Note: AttrValuesRaw is the JSONB field. We store the raw JSON string.
	// The actual attr values should be passed as a map and marshaled.
	// For now, just create with the raw string.
	_ = raw
	return database.DB.Create(ci).Error
}

func (r *CoreRepository) GetCIByID(id uint) (*CI, error) {
	var ci CI
	err := database.DB.First(&ci, id).Error
	return &ci, err
}

func (r *CoreRepository) UpdateCI(id uint, attrValues map[string]interface{}) error {
	raw, _ := json.Marshal(attrValues)
	return database.DB.Model(&CI{}).Where("id = ?", id).Update("attr_values", string(raw)).Error
}

func (r *CoreRepository) DeleteCI(id uint) error {
	return database.DB.Delete(&CI{}, id).Error
}

func (r *CoreRepository) ListCIsByType(ciTypeID uint, page, pageSize int) ([]CI, int64, error) {
	var cis []CI
	var total int64
	db := database.DB.Where("ci_type_id = ?", ciTypeID)
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&cis).Error
	return cis, total, err
}

// OperationLog
func (r *CoreRepository) CreateOperationLog(log *OperationLog) error {
	return database.DB.Create(log).Error
}

func (r *CoreRepository) CountCIsByType() ([]struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}, error) {
	var results []struct {
		Name  string `json:"name"`
		Value int64  `json:"value"`
	}
	err := database.DB.Raw(`
		SELECT ct.name, COUNT(c.id) as value
		FROM cmdb_core.ci_types ct
		LEFT JOIN cmdb_core.cis c ON c.ci_type_id = ct.id
		GROUP BY ct.id, ct.name
		ORDER BY value DESC
	`).Scan(&results).Error
	return results, err
}

func (r *CoreRepository) CountCIsByStatus() ([]struct {
	Status string `json:"status"`
	Value  int64  `json:"value"`
}, error) {
	var results []struct {
		Status string `json:"status"`
		Value  int64  `json:"value"`
	}
	err := database.DB.Raw(`
		SELECT status, COUNT(*) as value
		FROM cmdb_core.cis
		GROUP BY status
		ORDER BY value DESC
	`).Scan(&results).Error
	return results, err
}

func (r *CoreRepository) CountTotal(table string) (int64, error) {
	if !allowedTables[table] {
		return 0, fmt.Errorf("invalid table name: %s", table)
	}
	var count int64
	err := database.DB.Table(table).Count(&count).Error
	return count, err
}
