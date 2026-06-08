package core

import (
	"encoding/json"
	"fmt"

	"cmdb-api/database"

	"gorm.io/gorm"
)

var allowedTables = map[string]bool{
	"cmdb_core.cis":         true,
	"cmdb_core.ci_types":    true,
	"cmdb_discovery.rules":  true,
	"cmdb_discovery.agents": true,
}

type CoreRepository struct {
	db *gorm.DB
}

func NewCoreRepository() *CoreRepository {
	return &CoreRepository{db: database.DB}
}

func (r *CoreRepository) WithTx(tx *gorm.DB) *CoreRepository {
	return &CoreRepository{db: tx}
}

// CIType
func (r *CoreRepository) CreateCIType(ct *CIType) error {
	return r.db.Create(ct).Error
}

func (r *CoreRepository) GetCITypeByID(id uint) (*CIType, error) {
	var ct CIType
	err := r.db.First(&ct, id).Error
	return &ct, err
}

func (r *CoreRepository) GetCITypeWithAttributes(id uint) (*CIType, error) {
	var ct CIType
	err := r.db.Preload("Attributes").First(&ct, id).Error
	return &ct, err
}

func (r *CoreRepository) GetCITypeByName(name string) (*CIType, error) {
	var ct CIType
	err := r.db.Where("name = ?", name).First(&ct).Error
	return &ct, err
}

func (r *CoreRepository) ListCITypes() ([]CIType, error) {
	var cts []CIType
	err := r.db.Find(&cts).Error
	return cts, err
}

func (r *CoreRepository) UpdateCIType(id uint, updates map[string]interface{}) error {
	return r.db.Model(&CIType{}).Where("id = ?", id).Updates(updates).Error
}

func (r *CoreRepository) DeleteCIType(id uint) error {
	return r.db.Delete(&CIType{}, id).Error
}

// Attribute
func (r *CoreRepository) CreateAttribute(attr *Attribute) error {
	return r.db.Create(attr).Error
}

func (r *CoreRepository) GetAttributeByID(id uint) (*Attribute, error) {
	var attr Attribute
	err := r.db.First(&attr, id).Error
	return &attr, err
}

// CITypeAttribute
func (r *CoreRepository) AddAttributeToCIType(cta *CITypeAttribute) error {
	return r.db.Create(cta).Error
}

func (r *CoreRepository) GetCITypeAttributes(ciTypeID uint) ([]Attribute, error) {
	var attrs []Attribute
	err := r.db.
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
	return r.db.Create(ci).Error
}

func (r *CoreRepository) GetCIByID(id uint) (*CI, error) {
	var ci CI
	err := r.db.First(&ci, id).Error
	return &ci, err
}

func (r *CoreRepository) UpdateCI(id uint, attrValues map[string]interface{}) error {
	raw, _ := json.Marshal(attrValues)
	return r.db.Model(&CI{}).Where("id = ?", id).Update("attr_values", string(raw)).Error
}

func (r *CoreRepository) DeleteCI(id uint) error {
	return r.db.Delete(&CI{}, id).Error
}

func (r *CoreRepository) ListCIsByType(ciTypeID uint, page, pageSize int) ([]CI, int64, error) {
	var cis []CI
	var total int64
	db := r.db.Where("ci_type_id = ?", ciTypeID)
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&cis).Error
	return cis, total, err
}

// OperationLog
func (r *CoreRepository) CreateOperationLog(log *OperationLog) error {
	return r.db.Create(log).Error
}

func (r *CoreRepository) CountCIsByType() ([]struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}, error) {
	var results []struct {
		Name  string `json:"name"`
		Value int64  `json:"value"`
	}
	err := r.db.Raw(`
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
	err := r.db.Raw(`
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
	err := r.db.Table(table).Count(&count).Error
	return count, err
}
