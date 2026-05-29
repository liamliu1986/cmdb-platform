package discovery

type CreateRuleRequest struct {
	Name     string                 `json:"name" binding:"required"`
	Type     string                 `json:"type" binding:"required,oneof=cloud server network"`
	Config   map[string]interface{} `json:"config"`
	Schedule string                 `json:"schedule"`
}

type ExecuteRuleRequest struct {
	RuleID uint `json:"rule_id" binding:"required"`
}

type AgentRegisterRequest struct {
	Name    string                 `json:"name" binding:"required"`
	IP      string                 `json:"ip"`
	OS      string                 `json:"os"`
	Arch    string                 `json:"arch"`
	Version string                 `json:"version"`
	Labels  map[string]interface{} `json:"labels"`
}

type AgentHeartbeatRequest struct {
	Status string `json:"status"` // running
}
