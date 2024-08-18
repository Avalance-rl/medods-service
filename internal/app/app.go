package app

import (
	"medods-service/internal/cache"
	"medods-service/internal/config"
	"medods-service/internal/database"
	"medods-service/internal/services"
	"net"
)



func init() {
	config.Init()
	database.Init()
	cache.InitRedis()
}


func Run() {
	r := services.SetupRouter()
	r.Run(net.JoinHostPort(config.CFG.HTTPServer.Address, config.CFG.HTTPServer.Port))
}

