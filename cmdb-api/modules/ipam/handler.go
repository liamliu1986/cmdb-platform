package ipam

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"cmdb-api/pkg/response"
)

type IPAMHandler struct {
	svc             *IPAMService
	resourceCreator func(subnetID uint, subnetName string) error
}

func NewIPAMHandler() *IPAMHandler {
	return &IPAMHandler{svc: NewIPAMService()}
}

func (h *IPAMHandler) SetResourceCreator(fn func(subnetID uint, subnetName string) error) {
	h.resourceCreator = fn
}

func (h *IPAMHandler) CreateSubnet(c *gin.Context) {
	var req CreateSubnetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 30001, err.Error())
		return
	}
	subnet, err := h.svc.CreateSubnet(&req)
	if err != nil {
		response.Error(c, 30002, err.Error())
		return
	}
	if h.resourceCreator != nil {
		_ = h.resourceCreator(subnet.ID, subnet.Name)
	}
	response.Success(c, subnet)
}

func (h *IPAMHandler) ListSubnets(c *gin.Context) {
	var parentID *uint
	if pid := c.Query("parent_id"); pid != "" {
		id, _ := strconv.ParseUint(pid, 10, 64)
		uid := uint(id)
		parentID = &uid
	}
	subnets, err := h.svc.ListSubnets(parentID)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, subnets)
}

func (h *IPAMHandler) AllocateIP(c *gin.Context) {
	var req AllocateIPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 30001, err.Error())
		return
	}
	operator, _ := c.Get("username")
	op := ""
	if operator != nil {
		op = operator.(string)
	}
	ip, err := h.svc.AllocateIP(&req, op)
	if err != nil {
		response.Error(c, 30003, err.Error())
		return
	}
	response.Success(c, ip)
}

func (h *IPAMHandler) ReleaseIP(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.ReleaseIP(uint(id)); err != nil {
		response.Error(c, 30004, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *IPAMHandler) ListIPsBySubnet(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	status := c.Query("status")
	var ips []IPAddress
	var err error
	if status == "free" {
		ips, err = h.svc.GetAvailableIPsBySubnet(uint(id))
	} else {
		ips, err = h.svc.repo.ListIPsBySubnet(uint(id))
	}
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, ips)
}

func (h *IPAMHandler) AllocateIPByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	operator, _ := c.Get("username")
	op := ""
	if operator != nil {
		op = operator.(string)
	}
	ip, err := h.svc.AllocateIPByID(uint(id), op)
	if err != nil {
		response.Error(c, 30003, err.Error())
		return
	}
	response.Success(c, ip)
}

func (h *IPAMHandler) GetIP(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	ip, err := h.svc.repo.GetIPByID(uint(id))
	if err != nil {
		response.Error(c, 404, "ip not found")
		return
	}
	response.Success(c, ip)
}

func (h *IPAMHandler) ListAvailableIPs(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	ips, err := h.svc.GetAvailableIPsBySubnet(uint(id))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, ips)
}

// --- User-IP Assignment Handlers ---

func (h *IPAMHandler) AssignIPToUser(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)
	var req AssignIPToUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 30001, err.Error())
		return
	}
	operator, _ := c.Get("username")
	op := ""
	if operator != nil {
		op = operator.(string)
	}
	if err := h.svc.AssignIPToUser(req.IPAddressID, uint(userID), op); err != nil {
		response.Error(c, 30005, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *IPAMHandler) UnassignIPFromUser(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)
	ipAddressID, _ := strconv.ParseUint(c.Param("ip_address_id"), 10, 64)
	if err := h.svc.UnassignIPFromUser(uint(ipAddressID), uint(userID)); err != nil {
		response.Error(c, 30006, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *IPAMHandler) GetUserAssignedIPs(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)
	ips, err := h.svc.GetUserAssignedIPs(uint(userID))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, ips)
}
