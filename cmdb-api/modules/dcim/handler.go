package dcim

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"cmdb-api/pkg/response"
)

type DCIMHandler struct {
	svc *DCIMService
}

func NewDCIMHandler() *DCIMHandler {
	return &DCIMHandler{svc: NewDCIMService()}
}

var coordRepo = &locationRepo{}

func (h *DCIMHandler) CreateIDC(c *gin.Context) {
	var req CreateIDCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, err.Error())
		return
	}
	idc, err := h.svc.CreateIDC(&req)
	if err != nil {
		response.Error(c, 40002, err.Error())
		return
	}
	response.Success(c, idc)
}

func (h *DCIMHandler) ListIDCs(c *gin.Context) {
	idcs, err := h.svc.ListIDCs()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, idcs)
}

func (h *DCIMHandler) CreateRoom(c *gin.Context) {
	var req CreateServerRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, err.Error())
		return
	}
	room, err := h.svc.CreateRoom(&req)
	if err != nil {
		response.Error(c, 40003, err.Error())
		return
	}
	response.Success(c, room)
}

func (h *DCIMHandler) CreateRack(c *gin.Context) {
	var req CreateRackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, err.Error())
		return
	}
	rack, err := h.svc.CreateRack(&req)
	if err != nil {
		response.Error(c, 40004, err.Error())
		return
	}
	response.Success(c, rack)
}

func (h *DCIMHandler) MountDevice(c *gin.Context) {
	var req MountDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, err.Error())
		return
	}
	if err := h.svc.MountDevice(&req); err != nil {
		response.Error(c, 40005, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *DCIMHandler) UnmountDevice(c *gin.Context) {
	rackID, _ := strconv.ParseUint(c.Param("rack_id"), 10, 64)
	uPos, _ := strconv.Atoi(c.Param("u_position"))
	if err := h.svc.UnmountDevice(uint(rackID), uPos); err != nil {
		response.Error(c, 40006, err.Error())
		return
	}
	response.Success(c, nil)
}

// SetCoord sets map coordinates for an IDC
func (h *DCIMHandler) SetCoord(c *gin.Context) {
	var req SetCoordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, err.Error())
		return
	}
	if req.MapProvider == "" {
		req.MapProvider = "amap"
	}
	coord := &LocationCoord{
		CIID:        req.CIID,
		Lat:         req.Lat,
		Lng:         req.Lng,
		Address:     req.Address,
		MapProvider: req.MapProvider,
	}
	if err := coordRepo.SetCoord(coord); err != nil {
		response.Error(c, 40002, err.Error())
		return
	}
	response.Success(c, coord)
}

// ListCoords returns all map coordinates
func (h *DCIMHandler) ListCoords(c *gin.Context) {
	coords, err := coordRepo.ListCoords()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, coords)
}

// GetRackLayout returns rack devices as an SVG-friendly array
func (h *DCIMHandler) GetRackLayout(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	rack, err := h.svc.repo.GetRackByID(uint(id))
	if err != nil {
		response.Error(c, 404, "rack not found")
		return
	}
	devices, err := h.svc.repo.GetRackDevices(uint(id))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, gin.H{
		"rack":    rack,
		"devices": devices,
	})
}
