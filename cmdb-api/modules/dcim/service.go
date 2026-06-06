package dcim

import "errors"

type DCIMService struct {
	repo *DCIMRepository
}

func NewDCIMService() *DCIMService {
	return &DCIMService{repo: NewDCIMRepository()}
}

func (s *DCIMService) CreateIDC(req *CreateIDCRequest) (*IDC, error) {
	idc := &IDC{
		Name:     req.Name,
		Address:  req.Address,
		Contact:  req.Contact,
		RegionID: req.RegionID,
	}
	if err := s.repo.CreateIDC(idc); err != nil {
		return nil, err
	}
	return idc, nil
}

func (s *DCIMService) ListIDCs() ([]IDC, error) {
	return s.repo.ListIDCs()
}

func (s *DCIMService) CreateRoom(req *CreateServerRoomRequest) (*ServerRoom, error) {
	room := &ServerRoom{
		Name:  req.Name,
		Floor: req.Floor,
		IDCID: req.IDCID,
	}
	if err := s.repo.CreateRoom(room); err != nil {
		return nil, err
	}
	return room, nil
}

func (s *DCIMService) ListRoomsByIDC(idcID uint) ([]ServerRoom, error) {
	return s.repo.ListRoomsByIDC(idcID)
}

func (s *DCIMService) CreateRack(req *CreateRackRequest) (*Rack, error) {
	rack := &Rack{
		Name:   req.Name,
		TotalU: req.TotalU,
		RoomID: req.RoomID,
	}
	if rack.TotalU <= 0 {
		rack.TotalU = 42
	}
	if err := s.repo.CreateRack(rack); err != nil {
		return nil, err
	}
	return rack, nil
}

func (s *DCIMService) ListRacksByRoom(roomID uint) ([]Rack, error) {
	return s.repo.ListRacksByRoom(roomID)
}

func (s *DCIMService) GetRackByID(id uint) (*Rack, error) {
	return s.repo.GetRackByID(id)
}

func (s *DCIMService) GetRackLayout(rackID uint) (*Rack, []RackLayout, error) {
	rack, err := s.repo.GetRackByID(rackID)
	if err != nil {
		return nil, nil, errors.New("rack not found")
	}
	devices, err := s.repo.GetRackDevices(rackID)
	if err != nil {
		return nil, nil, err
	}
	return rack, devices, nil
}

func (s *DCIMService) GetRackCapacity(rackID uint) (map[string]any, error) {
	rack, devices, err := s.GetRackLayout(rackID)
	if err != nil {
		return nil, err
	}
	occupied := len(devices)
	return map[string]any{
		"total_u":    rack.TotalU,
		"occupied":   occupied,
		"available":  rack.TotalU - occupied,
		"usage_rate": float64(occupied) / float64(rack.TotalU),
	}, nil
}

func (s *DCIMService) MountDevice(req *MountDeviceRequest) error {
	rack, err := s.repo.GetRackByID(req.RackID)
	if err != nil {
		return errors.New("rack not found")
	}
	if req.UPosition < 1 || req.UPosition > rack.TotalU {
		return errors.New("invalid U position")
	}
	existing, err := s.repo.GetRackDevices(req.RackID)
	if err != nil {
		return err
	}
	for _, d := range existing {
		if d.UPosition == req.UPosition {
			return errors.New("U position already occupied")
		}
	}
	layout := &RackLayout{
		RackID:     req.RackID,
		UPosition:  req.UPosition,
		DeviceCIID: req.DeviceCIID,
	}
	return s.repo.MountDevice(layout)
}

func (s *DCIMService) UnmountDevice(rackID uint, uPosition int) error {
	return s.repo.UnmountDevice(rackID, uPosition)
}
