package core

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"

	"cmdb-api/pkg/response"
)

type CoreHandler struct {
	svc *CoreService
}

func NewCoreHandler() *CoreHandler {
	return &CoreHandler{svc: NewCoreService()}
}

// CIType handlers
func (h *CoreHandler) CreateCIType(c *gin.Context) {
	var req CreateCITypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 20001, err.Error())
		return
	}
	operator, _ := c.Get("username")
	op := ""
	if operator != nil {
		op = operator.(string)
	}
	ct, err := h.svc.CreateCIType(&req, op)
	if err != nil {
		response.Error(c, 20002, err.Error())
		return
	}
	response.Success(c, ct)
}

func (h *CoreHandler) GetCIType(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	ct, err := h.svc.GetCITypeWithAttributes(uint(id))
	if err != nil {
		response.Error(c, 404, "ci_type not found")
		return
	}
	response.Success(c, ct)
}

func (h *CoreHandler) ListCITypes(c *gin.Context) {
	cts, err := h.svc.ListCITypes()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, cts)
}

// CI handlers
func (h *CoreHandler) CreateCI(c *gin.Context) {
	var req CreateCIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 20001, err.Error())
		return
	}
	operator, _ := c.Get("username")
	op := ""
	if operator != nil {
		op = operator.(string)
	}
	userID, _ := c.Get("user_id")
	uid := uint(0)
	if userID != nil {
		uid = userID.(uint)
	}
	ci, err := h.svc.CreateCI(&req, uid, op)
	if err != nil {
		response.Error(c, 20003, err.Error())
		return
	}
	response.Success(c, ci)
}

func (h *CoreHandler) GetCI(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	ci, err := h.svc.GetCI(uint(id))
	if err != nil {
		response.Error(c, 404, "ci not found")
		return
	}

	// Parse attr_values for response
	var attrValues map[string]any
	_ = json.Unmarshal([]byte(ci.AttrValuesRaw), &attrValues)

	// Resolve reference attributes (IP addresses)
	ct, _ := h.svc.GetCITypeWithAttributes(ci.CITypeID)
	if ct != nil {
		for _, attr := range ct.Attributes {
			if !attr.IsReference || attr.RefTable != "cmdb_ipam.ip_addresses" {
				continue
			}
			val, ok := attrValues[attr.Name]
			if !ok || val == nil {
				continue
			}
			var ipID uint
			switch v := val.(type) {
			case float64:
				ipID = uint(v)
			case int:
				ipID = uint(v)
			case uint:
				ipID = v
			}
			if ipID > 0 {
				ip, _ := h.svc.ResolveIPReference(ipID)
				if ip != nil {
					attrValues[attr.Name+"_resolved"] = ip
				}
			}
		}
	}

	response.Success(c, gin.H{
		"id":          ci.ID,
		"ci_type_id":  ci.CITypeID,
		"status":      ci.Status,
		"attr_values": attrValues,
		"updated_by":  ci.UpdatedBy,
		"created_at":  ci.CreatedAt,
		"updated_at":  ci.UpdatedAt,
	})
}

func (h *CoreHandler) UpdateCI(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req CreateCIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 20001, err.Error())
		return
	}
	operator, _ := c.Get("username")
	op := ""
	if operator != nil {
		op = operator.(string)
	}
	userID, _ := c.Get("user_id")
	uid := uint(0)
	if userID != nil {
		uid = userID.(uint)
	}
	ci, err := h.svc.UpdateCI(uint(id), uid, req.AttrValues, op)
	if err != nil {
		response.Error(c, 20004, err.Error())
		return
	}
	response.Success(c, ci)
}

func (h *CoreHandler) DeleteCI(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	operator, _ := c.Get("username")
	op := ""
	if operator != nil {
		op = operator.(string)
	}
	if err := h.svc.DeleteCI(uint(id), op); err != nil {
		response.Error(c, 20005, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *CoreHandler) SearchCI(c *gin.Context) {
	var req CISearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, 20001, err.Error())
		return
	}

	builder := NewCISearchBuilder().
		WithQuery(req.Q).
		WithPagination(req.Page, req.PageSize).
		WithSort(req.Sort)

	cis, total, err := builder.Execute()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list": cis,
		"pagination": gin.H{
			"page":      req.Page,
			"page_size": req.PageSize,
			"total":     total,
		},
	})
}

func (h *CoreHandler) DashboardStats(c *gin.Context) {
	stats, err := h.svc.GetDashboardStats()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, stats)
}
