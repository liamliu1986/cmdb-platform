package integration

import (
	"os"

	"github.com/gin-gonic/gin"
	"cmdb-api/pkg/response"
)

type IntegrationHandler struct {
	prometheus *PrometheusClient
	elk        *ELKClient
	emailSvc   *EmailService
	emailCfg   *EmailConfig
}

func NewIntegrationHandler() *IntegrationHandler {
	// Default configs - override via environment in production
	emailCfg := &EmailConfig{
		SMTPHost: getEnv("SMTP_HOST", ""),
		SMTPPort: getEnv("SMTP_PORT", "587"),
		Username: getEnv("SMTP_USER", ""),
		Password: getEnv("SMTP_PASS", ""),
		From:     getEnv("SMTP_FROM", "cmdb@example.com"),
	}
	return &IntegrationHandler{
		prometheus: NewPrometheusClient(getEnv("PROMETHEUS_URL", "http://localhost:9090")),
		elk:        NewELKClient(getEnv("ELASTICSEARCH_URL", "http://localhost:9200"), "logstash-*"),
		emailSvc:   NewEmailService(emailCfg),
		emailCfg:   emailCfg,
	}
}

// PrometheusQuery handles instant Prometheus queries
func (h *IntegrationHandler) PrometheusQuery(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		response.Error(c, 50001, "query parameter is required")
		return
	}
	result, err := h.prometheus.QueryInstant(query)
	if err != nil {
		response.Error(c, 50002, err.Error())
		return
	}
	response.Success(c, result)
}

// ELKSearch handles ELK log searches
func (h *IntegrationHandler) ELKSearch(c *gin.Context) {
	var req struct {
		Query string `json:"query" binding:"required"`
		Size  int    `json:"size"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	result, err := h.elk.SearchLogs(req.Query, req.Size)
	if err != nil {
		response.Error(c, 50003, err.Error())
		return
	}
	response.Success(c, result)
}

// SendTestEmail handles test email sending
func (h *IntegrationHandler) SendTestEmail(c *gin.Context) {
	var req struct {
		To      string `json:"to" binding:"required"`
		Subject string `json:"subject" binding:"required"`
		Body    string `json:"body" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 50001, err.Error())
		return
	}
	if err := h.emailSvc.SendAlert([]string{req.To}, req.Subject, req.Body); err != nil {
		response.Error(c, 50004, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "email sent"})
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
