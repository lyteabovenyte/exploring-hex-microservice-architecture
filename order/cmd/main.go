package main

import (
	"github.com/lyteabovenyte/microservices-main/order/config"
	"github.com/lyteabovenyte/microservices-main/order/internal/adapters/db"
	"github.com/lyteabovenyte/microservices-main/order/internal/adapters/grpc"
	"github.com/lyteabovenyte/microservices-main/order/internal/application/core/api"
	"log"
)

func main() {
	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("failed to connect to the database. Error: %v", err)
	}

	application := api.NewApplication(dbAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
