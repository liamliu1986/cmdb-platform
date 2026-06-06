package router

import (
	"github.com/gin-gonic/gin"
	"cmdb-api/config"
	"cmdb-api/middleware"
	"cmdb-api/modules/auth"
	"cmdb-api/modules/core"
	"cmdb-api/modules/dcim"
	"cmdb-api/modules/discovery"
	"cmdb-api/modules/integration"
	"cmdb-api/modules/ipam"
)

func Setup(r *gin.Engine) {
	cfg := config.Load()
	authHandler := auth.NewAuthHandler(cfg)
	coreHandler := core.NewCoreHandler()
	ipamHandler := ipam.NewIPAMHandler()
	dcimHandler := dcim.NewDCIMHandler()
	discoveryHandler := discovery.NewDiscoveryHandler()
	integrationHandler := integration.NewIntegrationHandler()
	jwtMiddleware := middleware.JWTAuth(cfg)

	api := r.Group("/api/v1")
	{
		// Public
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)

		// Agent endpoints (public, use X-Agent-Token)
		api.POST("/discovery/agents/register", discoveryHandler.RegisterAgent)
		api.POST("/discovery/agents/heartbeat", discoveryHandler.AgentHeartbeat)

		// Protected
		authorized := api.Group("", jwtMiddleware)
		{
			// Stats
			authorized.GET("/stats/dashboard", coreHandler.DashboardStats)

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

			// IPAM
			authorized.POST("/ipam/subnets", ipamHandler.CreateSubnet)
			authorized.GET("/ipam/subnets", ipamHandler.ListSubnets)
			authorized.GET("/ipam/subnets/:id/ips", ipamHandler.ListIPsBySubnet)
			authorized.GET("/ipam/subnets/:id/ips/available", ipamHandler.ListAvailableIPs)
			authorized.GET("/ipam/ips/:id", ipamHandler.GetIP)
			authorized.POST("/ipam/ips/allocate", ipamHandler.AllocateIP)
			authorized.POST("/ipam/ips/:id/allocate-by-id", ipamHandler.AllocateIPByID)
			authorized.POST("/ipam/ips/:id/release", ipamHandler.ReleaseIP)

			// IPAM User-IP Assignment
			authorized.POST("/ipam/users/:user_id/ips", ipamHandler.AssignIPToUser)
			authorized.DELETE("/ipam/users/:user_id/ips/:ip_address_id", ipamHandler.UnassignIPFromUser)
			authorized.GET("/ipam/users/:user_id/ips", ipamHandler.GetUserAssignedIPs)

			// Auth permissions
			authorized.GET("/auth/permissions/subnets", authHandler.GetUserSubnetPermissions)

			// DCIM
			authorized.POST("/dcim/idcs", dcimHandler.CreateIDC)
			authorized.GET("/dcim/idcs", dcimHandler.ListIDCs)
			authorized.POST("/dcim/rooms", dcimHandler.CreateRoom)
			authorized.GET("/dcim/rooms", dcimHandler.ListRooms)
			authorized.POST("/dcim/racks", dcimHandler.CreateRack)
			authorized.GET("/dcim/racks", dcimHandler.ListRacks)
			authorized.GET("/dcim/racks/:id", dcimHandler.GetRack)
			authorized.POST("/dcim/racks/mount", dcimHandler.MountDevice)
			authorized.DELETE("/dcim/racks/:rack_id/devices/:u_position", dcimHandler.UnmountDevice)

			// DCIM coordinates
			authorized.POST("/dcim/coords", dcimHandler.SetCoord)
			authorized.GET("/dcim/coords", dcimHandler.ListCoords)
			authorized.GET("/dcim/racks/:id/layout", dcimHandler.GetRackLayout)
			authorized.GET("/dcim/racks/:id/capacity", dcimHandler.GetRackCapacity)

			// Discovery
			authorized.POST("/discovery/rules", discoveryHandler.CreateRule)
			authorized.GET("/discovery/rules", discoveryHandler.ListRules)
			authorized.POST("/discovery/rules/execute", discoveryHandler.ExecuteRule)
			authorized.GET("/discovery/agents", discoveryHandler.ListAgents)

			// Integration
			authorized.GET("/integration/prometheus/query", integrationHandler.PrometheusQuery)
			authorized.POST("/integration/elk/search", integrationHandler.ELKSearch)
			authorized.POST("/integration/email/test", integrationHandler.SendTestEmail)
		}
	}
}
