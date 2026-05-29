package dcim

type CreateIDCRequest struct {
	Name     string `json:"name" binding:"required"`
	Address  string `json:"address"`
	Contact  string `json:"contact"`
	RegionID uint   `json:"region_id"`
}

type CreateServerRoomRequest struct {
	Name  string `json:"name" binding:"required"`
	Floor string `json:"floor"`
	IDCID uint   `json:"idc_id" binding:"required"`
}

type CreateRackRequest struct {
	Name   string `json:"name" binding:"required"`
	TotalU int    `json:"total_u"`
	RoomID uint   `json:"room_id" binding:"required"`
}

type MountDeviceRequest struct {
	RackID     uint `json:"rack_id" binding:"required"`
	UPosition  int  `json:"u_position" binding:"required"`
	DeviceCIID uint `json:"device_ci_id" binding:"required"`
}
