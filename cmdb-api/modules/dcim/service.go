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

func (s *DCIMService) MountDevice(req *MountDeviceRequest) error {
	rack, err := s.repo.GetRackByID(req.RackID)
	if err != nil {
		return errors.New("rack not found")
	}
	if req.UPosition < 1 || req.UPosition > rack.TotalU {
		return errors.New("invalid U position")
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
