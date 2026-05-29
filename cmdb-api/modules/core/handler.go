package core

import (
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
	ct, err := h.svc.GetCIType(uint(id))
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
	ci, err := h.svc.CreateCI(&req, op)
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
	response.Success(c, ci)
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
	ci, err := h.svc.UpdateCI(uint(id), req.AttrValues, op)
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
