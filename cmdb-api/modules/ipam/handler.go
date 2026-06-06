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
