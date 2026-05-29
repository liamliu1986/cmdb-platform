package discovery

import (
	"gorm.io/gorm"
	"time"
)

// DiscoveryRule 发现规则
type DiscoveryRule struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:64;not null" json:"name"`
	Type      string         `gorm:"size:16;not null" json:"type"`                   // cloud/server/network
	Config    string         `gorm:"type:jsonb;not null;default:'{}'" json:"config"` // plugin config (AK/SK, region, etc.)
	Schedule  string         `gorm:"size:32" json:"schedule"`                        // cron expression or "manual"
	Enabled   bool           `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (DiscoveryRule) TableName() string { return "cmdb_discovery.rules" }

// DiscoveryTask 发现任务
type DiscoveryTask struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	RuleID        uint       `gorm:"not null;index" json:"rule_id"`
	Status        string     `gorm:"size:16;default:'pending'" json:"status"` // pending/running/success/failed
	ResultSummary string     `gorm:"type:text" json:"result_summary"`
	StartedAt     *time.Time `json:"started_at"`
	CompletedAt   *time.Time `json:"completed_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

func (DiscoveryTask) TableName() string { return "cmdb_discovery.tasks" }

// DiscoveryResult 发现结果
type DiscoveryResult struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TaskID      uint      `gorm:"not null;index" json:"task_id"`
	CITypeName  string    `gorm:"size:32" json:"ci_type_name"` // e.g., "AliyunECS"
	UniqueKey   string    `gorm:"size:128" json:"unique_key"`  // for deduplication
	RawData     string    `gorm:"type:jsonb" json:"raw_data"`
	MatchedCIID *uint     `json:"matched_ci_id"`               // matched existing CI or nil
	Action      string    `gorm:"size:16" json:"action"`       // created/updated/ignored
	CreatedAt   time.Time `json:"created_at"`
}

func (DiscoveryResult) TableName() string { return "cmdb_discovery.results" }

// Agent Agent 注册信息
type Agent struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"size:64;not null" json:"name"`
	Token         string         `gorm:"size:128;not null" json:"-"` // auth token
	IP            string         `gorm:"size:43" json:"ip"`
	OS            string         `gorm:"size:16" json:"os"`          // linux/windows/darwin
	Arch          string         `gorm:"size:16" json:"arch"`        // amd64/arm64
	Version       string         `gorm:"size:16" json:"version"`
	Labels        string         `gorm:"type:jsonb" json:"labels"`
	LastHeartbeat *time.Time     `json:"last_heartbeat"`
	Status        string         `gorm:"size:16;default:'offline'" json:"status"` // online/offline
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Agent) TableName() string { return "cmdb_discovery.agents" }
