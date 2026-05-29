package router

import (
	"github.com/gin-gonic/gin"
	"cmdb-api/config"
	"cmdb-api/middleware"
	"cmdb-api/modules/auth"
	"cmdb-api/modules/core"
)

func Setup(r *gin.Engine) {
	cfg := config.Load()
	authHandler := auth.NewAuthHandler(cfg)
	coreHandler := core.NewCoreHandler()
	jwtMiddleware := middleware.JWTAuth(cfg)

	api := r.Group("/api/v1")
	{
		// Public
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)

		// Protected
		authorized := api.Group("", jwtMiddleware)
		{
			// CIType
			authorized.POST("/citypes", coreHandler.CreateCIType)
			authorized.GET("/citypes", coreHandler.ListCITypes)
			authorized.GET("/citypes/:id", coreHandler.GetCIType)

			// CI
			authorized.POST("/ci", coreHandler.CreateCI)
			authorized.GET("/ci/:id", coreHandler.GetCI)
			authorized.PUT("/ci/:id", coreHandler.UpdateCI)
			authorized.DELETE("/ci/:id", coreHandler.DeleteCI)
			authorized.GET("/ci/s", coreHandler.SearchCI)
		}
	}
}
