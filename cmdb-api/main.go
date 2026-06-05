package main

import (
    "github.com/gin-gonic/gin"
    "cmdb-api/config"
    "cmdb-api/database"
    "cmdb-api/modules/auth"
    "cmdb-api/modules/core"
    "cmdb-api/modules/dcim"
    "cmdb-api/modules/discovery"
    "cmdb-api/modules/ipam"
    "cmdb-api/router"
)

func main() {
    cfg := config.Load()
    database.InitPostgres(cfg)
    database.InitRedis(cfg)

    schemas := []string{"cmdb_auth", "cmdb_core", "cmdb_dcim", "cmdb_discovery", "cmdb_ipam"}
    for _, s := range schemas {
        database.DB.Exec("CREATE SCHEMA IF NOT EXISTS " + s)
    }

    auth.Migrate()
    core.Migrate()
    core.InitBuiltinCITypes()
    dcim.Migrate()
    discovery.Migrate()
    ipam.Migrate()
    r := gin.Default()
    router.Setup(r)
    r.Run(":" + cfg.ServerPort)
}
