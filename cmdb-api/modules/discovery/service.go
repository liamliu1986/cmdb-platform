package discovery

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
)

type DiscoveryService struct {
	repo *DiscoveryRepository
}

func NewDiscoveryService() *DiscoveryService {
	return &DiscoveryService{repo: NewDiscoveryRepository()}
}

func (s *DiscoveryService) generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *DiscoveryService) CreateRule(req *CreateRuleRequest) (*DiscoveryRule, error) {
	config, _ := json.Marshal(req.Config)
	rule := &DiscoveryRule{
		Name:     req.Name,
		Type:     req.Type,
		Config:   string(config),
		Schedule: req.Schedule,
		Enabled:  true,
	}
	if err := s.repo.CreateRule(rule); err != nil {
		return nil, err
	}
	return rule, nil
}

func (s *DiscoveryService) ListRules() ([]DiscoveryRule, error) {
	return s.repo.ListRules()
}

func (s *DiscoveryService) ExecuteRule(req *ExecuteRuleRequest) (*DiscoveryTask, error) {
	rule, err := s.repo.GetRuleByID(req.RuleID)
	if err != nil {
		return nil, errors.New("rule not found")
	}
	task := &DiscoveryTask{
		RuleID: rule.ID,
		Status: "pending",
	}
	if err := s.repo.CreateTask(task); err != nil {
		return nil, err
	}
	// In a real implementation, this would enqueue the task for async execution
	return task, nil
}

func (s *DiscoveryService) RegisterAgent(req *AgentRegisterRequest) (*Agent, error) {
	labels, _ := json.Marshal(req.Labels)
	agent := &Agent{
		Name:    req.Name,
		Token:   s.generateToken(),
		IP:      req.IP,
		OS:      req.OS,
		Arch:    req.Arch,
		Version: req.Version,
		Labels:  string(labels),
		Status:  "online",
	}
	if err := s.repo.CreateAgent(agent); err != nil {
		return nil, err
	}
	return agent, nil
}

func (s *DiscoveryService) AgentHeartbeat(token string, req *AgentHeartbeatRequest) error {
	agent, err := s.repo.GetAgentByToken(token)
	if err != nil {
		return errors.New("invalid agent token")
	}
	return s.repo.UpdateAgentHeartbeat(agent.ID, req.Status)
}

func (s *DiscoveryService) ListAgents() ([]Agent, error) {
	return s.repo.ListAgents()
}
