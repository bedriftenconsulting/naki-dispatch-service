package main

import (
	"fmt"
	"log"

	"github.com/naki/dispatch-service/conf"
	"github.com/naki/dispatch-service/database"
	"github.com/naki/dispatch-service/functions/api_functions"
	"github.com/naki/dispatch-service/routers"
)

func main() {
	if err := conf.Load(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := database.Connect(); err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	api_functions.InitRedis()

	api_functions.StartKafkaConsumers()

	r := routers.SetupRouter()

	addr := fmt.Sprintf(":%s", conf.AppConfig.Port)
	log.Printf("dispatch service starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
