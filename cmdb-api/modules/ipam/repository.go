package ipam

import (
	"cmdb-api/database"

	"gorm.io/gorm"
)

type IPAMRepository struct {
	db *gorm.DB
}

func NewIPAMRepository() *IPAMRepository {
	return &IPAMRepository{db: database.DB}
}

func (r *IPAMRepository) WithTx(tx *gorm.DB) *IPAMRepository {
	return &IPAMRepository{db: tx}
}

// Subnet CRUD
func (r *IPAMRepository) CreateSubnet(s *Subnet) error {
	return r.db.Create(s).Error
}

func (r *IPAMRepository) GetSubnetByID(id uint) (*Subnet, error) {
	var s Subnet
	err := r.db.First(&s, id).Error
	return &s, err
}

func (r *IPAMRepository) ListSubnets(parentID *uint) ([]Subnet, error) {
	var subnets []Subnet
	db := r.db
	if parentID != nil {
		db = db.Where("parent_id = ?", *parentID)
	} else {
		db = db.Where("parent_id IS NULL")
	}
	err := db.Find(&subnets).Error
	return subnets, err
}

func (r *IPAMRepository) DeleteSubnet(id uint) error {
	return r.db.Delete(&Subnet{}, id).Error
}

// IP Address
func (r *IPAMRepository) CreateIP(ip *IPAddress) error {
	return r.db.Create(ip).Error
}

func (r *IPAMRepository) GetIPByID(id uint) (*IPAddress, error) {
	var ip IPAddress
	err := r.db.First(&ip, id).Error
	return &ip, err
}

func (r *IPAMRepository) ListIPsBySubnet(subnetID uint) ([]IPAddress, error) {
	var ips []IPAddress
	err := r.db.Where("subnet_id = ?", subnetID).Order("id").Find(&ips).Error
	return ips, err
}

func (r *IPAMRepository) GetAvailableIPsBySubnet(subnetID uint) ([]IPAddress, error) {
	var ips []IPAddress
	err := r.db.Where("subnet_id = ? AND status = ?", subnetID, "free").Order("id").Find(&ips).Error
	return ips, err
}

func (r *IPAMRepository) UpdateIPStatus(id uint, status string, ciID *uint, allocatedBy string) error {
	updates := map[string]any{
		"status": status,
	}
	if ciID != nil {
		updates["ci_id"] = *ciID
	}
	if allocatedBy != "" {
		updates["allocated_by"] = allocatedBy
	}
	return r.db.Model(&IPAddress{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateIPStatusFull updates all fields unconditionally
func (r *IPAMRepository) UpdateIPStatusFull(id uint, status string, ciID *uint, allocatedBy string) error {
	updates := map[string]any{
		"status":       status,
		"ci_id":        ciID,
		"allocated_by": allocatedBy,
	}
	return r.db.Model(&IPAddress{}).Where("id = ?", id).Updates(updates).Error
}

// GetIPByCIID retrieves IP address by CI ID
func (r *IPAMRepository) GetIPByCIID(ciID uint) (*IPAddress, error) {
	var ip IPAddress
	err := r.db.Where("ci_id = ?", ciID).First(&ip).Error
	return &ip, err
}

// User-IP Assignment CRUD

func (r *IPAMRepository) CreateUserIP(ui *UserIPAddress) error {
	return r.db.Create(ui).Error
}

func (r *IPAMRepository) DeleteUserIP(userID, ipAddressID uint) error {
	return r.db.Where("user_id = ? AND ip_address_id = ?", userID, ipAddressID).Delete(&UserIPAddress{}).Error
}

func (r *IPAMRepository) GetUserIPByIPAddressID(ipAddressID uint) (*UserIPAddress, error) {
	var ui UserIPAddress
	err := r.db.Where("ip_address_id = ?", ipAddressID).First(&ui).Error
	return &ui, err
}

func (r *IPAMRepository) ListUserIPs(userID uint) ([]UserIPAddress, error) {
	var uis []UserIPAddress
	err := r.db.Where("user_id = ?", userID).Find(&uis).Error
	return uis, err
}

func (r *IPAMRepository) ListUserIPsWithDetails(userID uint) ([]IPAddress, error) {
	var ips []IPAddress
	err := r.db.
		Joins("JOIN cmdb_ipam.user_ip_addresses ON user_ip_addresses.ip_address_id = ip_addresses.id").
		Where("user_ip_addresses.user_id = ?", userID).
		Find(&ips).Error
	return ips, err
}
