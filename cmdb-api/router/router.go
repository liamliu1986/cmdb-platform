package router

import (
	"github.com/gin-gonic/gin"
	"cmdb-api/config"
	"cmdb-api/modules/auth"
)

func Setup(r *gin.Engine) {
	cfg := config.Load()
	authHandler := auth.NewAuthHandler(cfg)

	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
	}
}
