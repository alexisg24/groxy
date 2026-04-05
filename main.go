package main

import (
	config "github.com/alexisg24/groxy/core/config"
	healthcheck "github.com/alexisg24/groxy/core/health-check"
	"github.com/alexisg24/groxy/core/server"
)

func main() {
	config.Init()
	healthcheck.InitializeHealthCheck(config.GlobalConfig.GetAllItems())
	server.StartHttpServer()
}
