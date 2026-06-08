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

func (s *IPAMService) GetIPByID(id uint) (*IPAddress, error) {
	return s.repo.GetIPByID(id)
}

func (s *IPAMService) GetAvailableIPsBySubnet(subnetID uint) ([]IPAddress, error) {
	return s.repo.GetAvailableIPsBySubnet(subnetID)
}

// --- User-IP Assignment (Layer 1) ---

func (s *IPAMService) AssignIPToUser(ipID, userID uint, assignedBy string) error {
	ip, err := s.repo.GetIPByID(ipID)
	if err != nil {
		return errors.New("ip not found")
	}
	if ip.Status != "free" {
		return errors.New("ip is not available for assignment")
	}
	ui := &UserIPAddress{
		UserID:      userID,
		IPAddressID: ipID,
		AssignedBy:  assignedBy,
	}
	if err := s.repo.CreateUserIP(ui); err != nil {
		return err
	}
	return s.repo.UpdateIPStatusFull(ipID, "reserved", nil, "")
}

func (s *IPAMService) UnassignIPFromUser(ipID, userID uint) error {
	ip, err := s.repo.GetIPByID(ipID)
	if err != nil {
		return errors.New("ip not found")
	}
	if ip.Status == "allocated" {
		return errors.New("cannot unassign ip currently allocated to a ci")
	}
	if err := s.repo.DeleteUserIP(userID, ipID); err != nil {
		return err
	}
	return s.repo.UpdateIPStatusFull(ipID, "free", nil, "")
}

func (s *IPAMService) GetUserAssignedIPs(userID uint) ([]IPAddress, error) {
	return s.repo.ListUserIPsWithDetails(userID)
}

func (s *IPAMService) IsIPAssignedToUser(ipID, userID uint) (bool, error) {
	ui, err := s.repo.GetUserIPByIPAddressID(ipID)
	if err != nil {
		return false, nil // not found = not assigned
	}
	return ui.UserID == userID, nil
}

// --- IP-CI Binding (Layer 2) ---

func (s *IPAMService) AllocateIPForCI(ipID, ciID, userID uint, operator string) error {
	ip, err := s.repo.GetIPByID(ipID)
	if err != nil {
		return errors.New("ip not found")
	}
	if ip.Status != "reserved" {
		return errors.New("ip must be reserved to user before allocating to ci")
	}
	if ip.CIID != nil && *ip.CIID != 0 && *ip.CIID != ciID {
		return errors.New("ip is already allocated to another ci")
	}
	// Verify ownership
	assigned, err := s.IsIPAssignedToUser(ipID, userID)
	if err != nil {
		return err
	}
	if !assigned {
		return errors.New("ip is not assigned to this user")
	}
	return s.repo.UpdateIPStatusFull(ipID, "allocated", &ciID, operator)
}

func (s *IPAMService) ReleaseIPFromCI(ipID uint) error {
	ip, err := s.repo.GetIPByID(ipID)
	if err != nil {
		return errors.New("ip not found")
	}
	if ip.Status != "allocated" {
		return errors.New("ip is not allocated to a ci")
	}
	// Return to reserved status (user pool)
	return s.repo.UpdateIPStatusFull(ipID, "reserved", nil, "")
}

// ReleaseIP rejects if the IP is bound to a CI.
func (s *IPAMService) ReleaseIP(ipID uint) error {
	ip, err := s.repo.GetIPByID(ipID)
	if err != nil {
		return errors.New("ip not found")
	}
	if ip.CIID != nil && *ip.CIID != 0 {
		return errors.New("cannot release ip bound to a ci")
	}
	return s.repo.UpdateIPStatusFull(ipID, "free", nil, "")
}
