package auth

import (
	"github.com/gin-gonic/gin"
	"cmdb-api/config"
	"cmdb-api/pkg/response"
)

type AuthHandler struct {
	svc *AuthService
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{svc: NewAuthService(cfg)}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 10001, err.Error())
		return
	}
	user, err := h.svc.Register(&req)
	if err != nil {
		response.Error(c, 10002, err.Error())
		return
	}
	response.Success(c, gin.H{"id": user.ID, "username": user.Username})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 10001, err.Error())
		return
	}
	resp, err := h.svc.Login(&req)
	if err != nil {
		response.Error(c, 10003, err.Error())
		return
	}
	response.Success(c, resp)
}

func (h *AuthHandler) GetUserSubnetPermissions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, 401, "unauthorized")
		return
	}
	ids, err := h.svc.GetUserPermittedSubnets(userID.(uint))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, gin.H{"subnet_ids": ids})
}
