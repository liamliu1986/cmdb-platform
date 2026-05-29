package dcim

import "cmdb-api/database"

type DCIMRepository struct{}

func NewDCIMRepository() *DCIMRepository {
	return &DCIMRepository{}
}

// IDC
func (r *DCIMRepository) CreateIDC(idc *IDC) error {
	return database.DB.Create(idc).Error
}

func (r *DCIMRepository) ListIDCs() ([]IDC, error) {
	var idcs []IDC
	err := database.DB.Find(&idcs).Error
	return idcs, err
}

// ServerRoom
func (r *DCIMRepository) CreateRoom(room *ServerRoom) error {
	return database.DB.Create(room).Error
}

func (r *DCIMRepository) ListRoomsByIDC(idcID uint) ([]ServerRoom, error) {
	var rooms []ServerRoom
	err := database.DB.Where("idc_id = ?", idcID).Find(&rooms).Error
	return rooms, err
}

// Rack
func (r *DCIMRepository) CreateRack(rack *Rack) error {
	return database.DB.Create(rack).Error
}

func (r *DCIMRepository) ListRacksByRoom(roomID uint) ([]Rack, error) {
	var racks []Rack
	err := database.DB.Where("room_id = ?", roomID).Find(&racks).Error
	return racks, err
}

func (r *DCIMRepository) GetRackByID(id uint) (*Rack, error) {
	var rack Rack
	err := database.DB.First(&rack, id).Error
	return &rack, err
}

// RackLayout
func (r *DCIMRepository) MountDevice(layout *RackLayout) error {
	return database.DB.Create(layout).Error
}

func (r *DCIMRepository) UnmountDevice(rackID uint, uPosition int) error {
	return database.DB.Where("rack_id = ? AND u_position = ?", rackID, uPosition).Delete(&RackLayout{}).Error
}

func (r *DCIMRepository) GetRackDevices(rackID uint) ([]RackLayout, error) {
	var layouts []RackLayout
	err := database.DB.Where("rack_id = ?", rackID).Order("u_position").Find(&layouts).Error
	return layouts, err
}
