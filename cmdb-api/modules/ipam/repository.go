package ipam

import "cmdb-api/database"

type IPAMRepository struct{}

func NewIPAMRepository() *IPAMRepository {
	return &IPAMRepository{}
}

// Subnet CRUD
func (r *IPAMRepository) CreateSubnet(s *Subnet) error {
	return database.DB.Create(s).Error
}

func (r *IPAMRepository) GetSubnetByID(id uint) (*Subnet, error) {
	var s Subnet
	err := database.DB.First(&s, id).Error
	return &s, err
}

func (r *IPAMRepository) ListSubnets(parentID *uint) ([]Subnet, error) {
	var subnets []Subnet
	db := database.DB
	if parentID != nil {
		db = db.Where("parent_id = ?", *parentID)
	} else {
		db = db.Where("parent_id IS NULL")
	}
	err := db.Find(&subnets).Error
	return subnets, err
}

func (r *IPAMRepository) DeleteSubnet(id uint) error {
	return database.DB.Delete(&Subnet{}, id).Error
}

// IP Address
func (r *IPAMRepository) CreateIP(ip *IPAddress) error {
	return database.DB.Create(ip).Error
}

func (r *IPAMRepository) GetIPByID(id uint) (*IPAddress, error) {
	var ip IPAddress
	err := database.DB.First(&ip, id).Error
	return &ip, err
}

func (r *IPAMRepository) ListIPsBySubnet(subnetID uint) ([]IPAddress, error) {
	var ips []IPAddress
	err := database.DB.Where("subnet_id = ?", subnetID).Find(&ips).Error
	return ips, err
}

func (r *IPAMRepository) UpdateIPStatus(id uint, status string, ciID *uint, allocatedBy string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if ciID != nil {
		updates["ci_id"] = *ciID
	}
	if allocatedBy != "" {
		updates["allocated_by"] = allocatedBy
	}
	return database.DB.Model(&IPAddress{}).Where("id = ?", id).Updates(updates).Error
}
