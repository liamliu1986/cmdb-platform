package discovery

import (
	"cmdb-api/database"
	"time"
)

type DiscoveryRepository struct{}

func NewDiscoveryRepository() *DiscoveryRepository {
	return &DiscoveryRepository{}
}

// Rule CRUD
func (r *DiscoveryRepository) CreateRule(rule *DiscoveryRule) error {
	return database.DB.Create(rule).Error
}

func (r *DiscoveryRepository) GetRuleByID(id uint) (*DiscoveryRule, error) {
	var rule DiscoveryRule
	err := database.DB.First(&rule, id).Error
	return &rule, err
}

func (r *DiscoveryRepository) ListRules() ([]DiscoveryRule, error) {
	var rules []DiscoveryRule
	err := database.DB.Find(&rules).Error
	return rules, err
}

func (r *DiscoveryRepository) UpdateRule(id uint, updates map[string]interface{}) error {
	return database.DB.Model(&DiscoveryRule{}).Where("id = ?", id).Updates(updates).Error
}

func (r *DiscoveryRepository) DeleteRule(id uint) error {
	return database.DB.Delete(&DiscoveryRule{}, id).Error
}

// Task
func (r *DiscoveryRepository) CreateTask(task *DiscoveryTask) error {
	return database.DB.Create(task).Error
}

func (r *DiscoveryRepository) GetTaskByID(id uint) (*DiscoveryTask, error) {
	var task DiscoveryTask
	err := database.DB.First(&task, id).Error
	return &task, err
}

func (r *DiscoveryRepository) UpdateTaskStatus(id uint, status string, summary string) error {
	updates := map[string]interface{}{"status": status}
	if summary != "" {
		updates["result_summary"] = summary
	}
	return database.DB.Model(&DiscoveryTask{}).Where("id = ?", id).Updates(updates).Error
}

func (r *DiscoveryRepository) ListTasks(ruleID uint) ([]DiscoveryTask, error) {
	var tasks []DiscoveryTask
	err := database.DB.Where("rule_id = ?", ruleID).Order("created_at DESC").Find(&tasks).Error
	return tasks, err
}

// Result
func (r *DiscoveryRepository) CreateResult(result *DiscoveryResult) error {
	return database.DB.Create(result).Error
}

// Agent
func (r *DiscoveryRepository) CreateAgent(agent *Agent) error {
	return database.DB.Create(agent).Error
}

func (r *DiscoveryRepository) GetAgentByToken(token string) (*Agent, error) {
	var agent Agent
	err := database.DB.Where("token = ?", token).First(&agent).Error
	return &agent, err
}

func (r *DiscoveryRepository) UpdateAgentHeartbeat(id uint, status string) error {
	now := time.Now()
	return database.DB.Model(&Agent{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_heartbeat": now,
		"status":         status,
	}).Error
}

func (r *DiscoveryRepository) ListAgents() ([]Agent, error) {
	var agents []Agent
	err := database.DB.Find(&agents).Error
	return agents, err
}
