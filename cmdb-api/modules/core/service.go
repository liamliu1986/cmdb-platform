package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"cmdb-api/modules/ipam"
)

type CoreService struct {
	repo    *CoreRepository
	ipamSvc *ipam.IPAMService
}

func NewCoreService() *CoreService {
	return &CoreService{
		repo:    NewCoreRepository(),
		ipamSvc: ipam.NewIPAMService(),
	}
}

func parseUintFromInterface(v any) uint {
	switch val := v.(type) {
	case float64:
		return uint(val)
	case int:
		return uint(val)
	case uint:
		return val
	case string:
		id, _ := strconv.ParseUint(val, 10, 64)
		return uint(id)
	default:
		return 0
	}
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

func (s *CoreService) ResolveIPReference(ipID uint) (*ipam.IPAddress, error) {
	return s.ipamSvc.GetIPByID(ipID)
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
		IsReference: req.IsReference,
		RefTable:    req.RefTable,
	}
	if err := s.repo.CreateAttribute(attr); err != nil {
		return nil, err
	}
	return attr, nil
}

func (s *CoreService) CreateCI(req *CreateCIRequest, userID uint, operator string) (*CI, error) {
	ct, err := s.repo.GetCITypeWithAttributes(req.CITypeID)
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

	// Process reference attributes before creating CI
	for _, attr := range ct.Attributes {
		if !attr.IsReference || attr.RefTable != "cmdb_ipam.ip_addresses" {
			continue
		}
		val, ok := req.AttrValues[attr.Name]
		if !ok || val == nil {
			continue
		}
		ipID := parseUintFromInterface(val)
		if ipID == 0 {
			continue
		}
		// Validate IP belongs to user
		assigned, err := s.ipamSvc.IsIPAssignedToUser(ipID, userID)
		if err != nil {
			return nil, err
		}
		if !assigned {
			return nil, fmt.Errorf("ip for attribute '%s' is not assigned to user", attr.Name)
		}
		// Validate IP is reserved (not already allocated)
		ip, err := s.ipamSvc.GetIPByID(ipID)
		if err != nil {
			return nil, fmt.Errorf("ip for attribute '%s' not found", attr.Name)
		}
		if ip.Status != "reserved" {
			return nil, fmt.Errorf("ip for attribute '%s' is not available", attr.Name)
		}
		req.AttrValues[attr.Name] = ipID
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

	// After CI created, bind IPs
	for _, attr := range ct.Attributes {
		if !attr.IsReference || attr.RefTable != "cmdb_ipam.ip_addresses" {
			continue
		}
		val, ok := req.AttrValues[attr.Name]
		if !ok || val == nil {
			continue
		}
		ipID := parseUintFromInterface(val)
		if ipID > 0 {
			_ = s.ipamSvc.AllocateIPForCI(ipID, ci.ID, userID, operator)
		}
	}

	s.logOperation("ci", ci.ID, "create", operator, nil, ci)
	return ci, nil
}

func (s *CoreService) GetCI(id uint) (*CI, error) {
	return s.repo.GetCIByID(id)
}

func (s *CoreService) UpdateCI(id, userID uint, attrValues map[string]any, operator string) (*CI, error) {
	oldCI, err := s.repo.GetCIByID(id)
	if err != nil {
		return nil, err
	}

	ct, err := s.repo.GetCITypeWithAttributes(oldCI.CITypeID)
	if err != nil {
		return nil, err
	}

	// Parse old attr_values
	var oldAttrValues map[string]any
	_ = json.Unmarshal([]byte(oldCI.AttrValuesRaw), &oldAttrValues)

	// Process reference attributes
	for _, attr := range ct.Attributes {
		if !attr.IsReference || attr.RefTable != "cmdb_ipam.ip_addresses" {
			continue
		}

		oldVal := oldAttrValues[attr.Name]
		newVal, hasNew := attrValues[attr.Name]

		// Release old IP if changed or removed
		if oldVal != nil {
			oldIPID := parseUintFromInterface(oldVal)
			if oldIPID > 0 && (!hasNew || newVal == nil || parseUintFromInterface(newVal) != oldIPID) {
				_ = s.ipamSvc.ReleaseIPFromCI(oldIPID)
			}
		}

		// Allocate new IP
		if hasNew && newVal != nil {
			newIPID := parseUintFromInterface(newVal)
			if newIPID > 0 {
				assigned, _ := s.ipamSvc.IsIPAssignedToUser(newIPID, userID)
				if !assigned {
					return nil, fmt.Errorf("ip for attribute '%s' is not assigned to user", attr.Name)
				}
				if err := s.ipamSvc.AllocateIPForCI(newIPID, id, userID, operator); err != nil {
					return nil, fmt.Errorf("failed to allocate ip for attribute '%s': %w", attr.Name, err)
				}
			}
		}
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

	// Parse attr_values and release referenced IPs
	var attrValues map[string]any
	_ = json.Unmarshal([]byte(ci.AttrValuesRaw), &attrValues)

	ct, _ := s.repo.GetCITypeWithAttributes(ci.CITypeID)
	for _, attr := range ct.Attributes {
		if !attr.IsReference || attr.RefTable != "cmdb_ipam.ip_addresses" {
			continue
		}
		val, ok := attrValues[attr.Name]
		if !ok || val == nil {
			continue
		}
		ipID := parseUintFromInterface(val)
		if ipID > 0 {
			_ = s.ipamSvc.ReleaseIPFromCI(ipID)
		}
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
