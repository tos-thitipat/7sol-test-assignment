package main

import (
	"pie_fire_dine/config"
	"pie_fire_dine/logger"
	"pie_fire_dine/server"
)

func setup() {
	config.Load()
	logger.SetupLogger(config.GetLogger())
}

func main() {
	setup()
	server.Start()
}
