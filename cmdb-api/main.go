package main

import (
    "github.com/gin-gonic/gin"
    "cmdb-api/config"
    "cmdb-api/database"
    "cmdb-api/modules/core"
    "cmdb-api/router"
)

func main() {
    cfg := config.Load()
    database.InitPostgres(cfg)
    database.InitRedis(cfg)
    core.Migrate()
    core.InitBuiltinCITypes()
    r := gin.Default()
    router.Setup(r)
    r.Run(":" + cfg.ServerPort)
}
