package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type CoreService struct {
	repo *CoreRepository
}

func NewCoreService() *CoreService {
	return &CoreService{repo: NewCoreRepository()}
}

func (s *CoreService) CreateCIType(req *CreateCITypeRequest, operator string) (*CIType, error) {
	if req.Alias == "" {
		req.Alias = req.Name
	}
	ct := &CIType{
		Name:         req.Name,
		Alias:        req.Alias,
		UniqueAttrID: req.UniqueAttrID,
		Icon:         req.Icon,
	}
	if err := s.repo.CreateCIType(ct); err != nil {
		return nil, err
	}
	s.logOperation("ci_type", ct.ID, "create", operator, nil, ct)
	return ct, nil
}

func (s *CoreService) GetCIType(id uint) (*CIType, error) {
	return s.repo.GetCITypeByID(id)
}

func (s *CoreService) GetCITypeWithAttributes(id uint) (*CIType, error) {
	return s.repo.GetCITypeWithAttributes(id)
}

func (s *CoreService) ListCITypes() ([]CIType, error) {
	return s.repo.ListCITypes()
}

func (s *CoreService) CreateAttribute(req *CreateAttributeRequest) (*Attribute, error) {
	attr := &Attribute{
		Name:        req.Name,
		Alias:       req.Alias,
		ValueType:   req.ValueType,
		IsChoice:    req.IsChoice,
		IsList:      req.IsList,
		IsUnique:    req.IsUnique,
		IsIndex:     req.IsIndex,
		IsPassword:  req.IsPassword,
		IsComputed:  req.IsComputed,
		ComputeExpr: req.ComputeExpr,
	}
	if err := s.repo.CreateAttribute(attr); err != nil {
		return nil, err
	}
	return attr, nil
}

func (s *CoreService) CreateCI(req *CreateCIRequest, operator string) (*CI, error) {
	ct, err := s.repo.GetCITypeByID(req.CITypeID)
	if err != nil {
		return nil, errors.New("ci_type not found")
	}
	if ct.UniqueAttrID > 0 {
		attr, _ := s.repo.GetAttributeByID(ct.UniqueAttrID)
		if attr != nil {
			uniqueVal, ok := req.AttrValues[attr.Name]
			if !ok || uniqueVal == nil || uniqueVal == "" {
				return nil, fmt.Errorf("unique attribute '%s' is required", attr.Name)
			}
		}
	}
	raw, _ := json.Marshal(req.AttrValues)
	ci := &CI{
		CITypeID:      req.CITypeID,
		AttrValuesRaw: string(raw),
		UpdatedBy:     operator,
	}
	if err := s.repo.CreateCI(ci); err != nil {
		return nil, err
	}
	s.logOperation("ci", ci.ID, "create", operator, nil, ci)
	return ci, nil
}

func (s *CoreService) GetCI(id uint) (*CI, error) {
	return s.repo.GetCIByID(id)
}

func (s *CoreService) UpdateCI(id uint, attrValues map[string]interface{}, operator string) (*CI, error) {
	oldCI, err := s.repo.GetCIByID(id)
	if err != nil {
		return nil, err
	}
	if err := s.repo.UpdateCI(id, attrValues); err != nil {
		return nil, err
	}
	newCI, _ := s.repo.GetCIByID(id)
	s.logOperation("ci", id, "update", operator, oldCI, newCI)
	return newCI, nil
}

func (s *CoreService) DeleteCI(id uint, operator string) error {
	ci, err := s.repo.GetCIByID(id)
	if err != nil {
		return err
	}
	if err := s.repo.DeleteCI(id); err != nil {
		return err
	}
	s.logOperation("ci", id, "delete", operator, ci, nil)
	return nil
}

func (s *CoreService) logOperation(targetType string, targetID uint, action, operator string, oldVal, newVal interface{}) {
	oldBytes, _ := json.Marshal(oldVal)
	newBytes, _ := json.Marshal(newVal)
	log := &OperationLog{
		TargetType: targetType,
		TargetID:   targetID,
		Action:     action,
		Operator:   operator,
		OldValue:   string(oldBytes),
		NewValue:   string(newBytes),
		CreatedAt:  time.Now(),
	}
	s.repo.CreateOperationLog(log)
}

type DashboardStats struct {
	TotalCI     int64 `json:"total_ci"`
	TotalCIType int64 `json:"total_citype"`
	TotalRule   int64 `json:"total_rule"`
	TotalAgent  int64 `json:"total_agent"`
	CIByType    []struct {
		Name  string `json:"name"`
		Value int64  `json:"value"`
	} `json:"ci_by_type"`
	CIByStatus []struct {
		Status string `json:"status"`
		Value  int64  `json:"value"`
	} `json:"ci_by_status"`
}

func (s *CoreService) GetDashboardStats() (*DashboardStats, error) {
	totalCI, err := s.repo.CountTotal("cmdb_core.cis")
	if err != nil {
		return nil, fmt.Errorf("count cis: %w", err)
	}
	totalCIType, err := s.repo.CountTotal("cmdb_core.ci_types")
	if err != nil {
		return nil, fmt.Errorf("count citypes: %w", err)
	}
	totalRule, err := s.repo.CountTotal("cmdb_discovery.rules")
	if err != nil {
		return nil, fmt.Errorf("count rules: %w", err)
	}
	totalAgent, err := s.repo.CountTotal("cmdb_discovery.agents")
	if err != nil {
		return nil, fmt.Errorf("count agents: %w", err)
	}
	ciByType, err := s.repo.CountCIsByType()
	if err != nil {
		return nil, fmt.Errorf("count ci by type: %w", err)
	}
	ciByStatus, err := s.repo.CountCIsByStatus()
	if err != nil {
		return nil, fmt.Errorf("count ci by status: %w", err)
	}

	return &DashboardStats{
		TotalCI:     totalCI,
		TotalCIType: totalCIType,
		TotalRule:   totalRule,
		TotalAgent:  totalAgent,
		CIByType:    ciByType,
		CIByStatus:  ciByStatus,
	}, nil
}
