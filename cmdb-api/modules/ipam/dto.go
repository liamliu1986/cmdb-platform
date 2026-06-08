package ipam

type CreateSubnetRequest struct {
	ParentID *uint  `json:"parent_id"`
	CIDR     string `json:"cidr" binding:"required"`
	Name     string `json:"name" binding:"required"`
	VlanID   string `json:"vlan_id"`
}

type AllocateIPRequest struct {
	SubnetID uint `json:"subnet_id" binding:"required"`
}

type ReleaseIPRequest struct {
	IPID uint `json:"ip_id" binding:"required"`
}

type AssignIPToUserRequest struct {
	IPAddressID uint `json:"ip_address_id" binding:"required"`
}
