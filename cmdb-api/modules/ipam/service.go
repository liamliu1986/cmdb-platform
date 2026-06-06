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
	_, ipNet, err := net.ParseCIDR(req.CIDR)
	if err != nil {
		return nil, errors.New("invalid CIDR format")
	}

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

	// Auto-generate IP addresses in the subnet
	if err := s.initIPPool(subnet.ID, ipNet); err != nil {
		_ = err
	}

	return subnet, nil
}

func (s *IPAMService) initIPPool(subnetID uint, ipNet *net.IPNet) error {
	ip := ipNet.IP.Mask(ipNet.Mask)
	incIP(ip)

	count := 0
	maxIPs := 254
	for ipNet.Contains(ip) {
		nextIP := make(net.IP, len(ip))
		copy(nextIP, ip)
		incIP(nextIP)
		if !ipNet.Contains(nextIP) {
			break
		}

		ipAddr := &IPAddress{
			SubnetID: subnetID,
			IP:       ip.String(),
			Status:   "free",
		}
		if err := s.repo.CreateIP(ipAddr); err != nil {
			return err
		}

		incIP(ip)
		count++
		if count >= maxIPs {
			break
		}
	}
	return nil
}

func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
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

func (s *IPAMService) AllocateIPByID(ipID uint, operator string) (*IPAddress, error) {
	ip, err := s.repo.GetIPByID(ipID)
	if err != nil {
		return nil, errors.New("ip not found")
	}
	if ip.Status != "free" {
		return nil, errors.New("ip is not available")
	}
	if err := s.repo.UpdateIPStatus(ip.ID, "allocated", nil, operator); err != nil {
		return nil, err
	}
	ip.Status = "allocated"
	ip.AllocatedBy = operator
	return ip, nil
}

func (s *IPAMService) GetAvailableIPsBySubnet(subnetID uint) ([]IPAddress, error) {
	return s.repo.GetAvailableIPsBySubnet(subnetID)
}
