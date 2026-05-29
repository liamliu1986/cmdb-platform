package discovery

import (
	"cmdb-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type DiscoveryHandler struct {
	svc *DiscoveryService
}

func NewDiscoveryHandler() *DiscoveryHandler {
	return &DiscoveryHandler{svc: NewDiscoveryService()}
}

func (h *DiscoveryHandler) CreateRule(c *gin.Context) {
	var req CreateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, err.Error())
		return
	}
	rule, err := h.svc.CreateRule(&req)
	if err != nil {
		response.Error(c, 40002, err.Error())
		return
	}
	response.Success(c, rule)
}

func (h *DiscoveryHandler) ListRules(c *gin.Context) {
	rules, err := h.svc.ListRules()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, rules)
}

func (h *DiscoveryHandler) ExecuteRule(c *gin.Context) {
	var req ExecuteRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, err.Error())
		return
	}
	task, err := h.svc.ExecuteRule(&req)
	if err != nil {
		response.Error(c, 40003, err.Error())
		return
	}
	response.Success(c, task)
}

func (h *DiscoveryHandler) RegisterAgent(c *gin.Context) {
	var req AgentRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, err.Error())
		return
	}
	agent, err := h.svc.RegisterAgent(&req)
	if err != nil {
		response.Error(c, 40004, err.Error())
		return
	}
	response.Success(c, gin.H{"agent_id": agent.ID, "token": agent.Token})
}

func (h *DiscoveryHandler) AgentHeartbeat(c *gin.Context) {
	token := c.GetHeader("X-Agent-Token")
	if token == "" {
		response.Error(c, 401, "missing agent token")
		return
	}
	var req AgentHeartbeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 40001, err.Error())
		return
	}
	if err := h.svc.AgentHeartbeat(token, &req); err != nil {
		response.Error(c, 401, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *DiscoveryHandler) ListAgents(c *gin.Context) {
	agents, err := h.svc.ListAgents()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, agents)
}
