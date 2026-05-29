package ipam

import (
	"errors"
	"net"
)

type IPAMService struct {
	repo *IPAMRepository
}

func NewIPAMService() *IPAMService {
	return &IPAMService{repo: NewIPAMRepository()}
}

func (s *IPAMService) CreateSubnet(req *CreateSubnetRequest) (*Subnet, error) {
	// Validate CIDR
	_, ipNet, err := net.ParseCIDR(req.CIDR)
	if err != nil {
		return nil, errors.New("invalid CIDR format")
	}
	_ = ipNet

	subnet := &Subnet{
		ParentID: req.ParentID,
		CIDR:     req.CIDR,
		Name:     req.Name,
		VlanID:   req.VlanID,
		Status:   "active",
	}
	if err := s.repo.CreateSubnet(subnet); err != nil {
		return nil, err
	}
	return subnet, nil
}

func (s *IPAMService) ListSubnets(parentID *uint) ([]Subnet, error) {
	return s.repo.ListSubnets(parentID)
}

func (s *IPAMService) AllocateIP(req *AllocateIPRequest, operator string) (*IPAddress, error) {
	// Find a free IP in the subnet
	ips, err := s.repo.ListIPsBySubnet(req.SubnetID)
	if err != nil {
		return nil, err
	}
	for _, ip := range ips {
		if ip.Status == "free" {
			if err := s.repo.UpdateIPStatus(ip.ID, "allocated", nil, operator); err != nil {
				return nil, err
			}
			ip.Status = "allocated"
			ip.AllocatedBy = operator
			return &ip, nil
		}
	}
	return nil, errors.New("no free IP available in subnet")
}

func (s *IPAMService) ReleaseIP(ipID uint) error {
	return s.repo.UpdateIPStatus(ipID, "free", nil, "")
}
