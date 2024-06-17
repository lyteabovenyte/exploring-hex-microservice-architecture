package grpc

import (
	"fmt"
	"google.golang.org/grpc/reflection"
	"log"
	"net"

	"github.com/huseyinbabal/microservices-proto/golang/order"
	"github.com/lyteabovenyte/microservices-main/order/config"
	"github.com/lyteabovenyte/microservices-main/order/internal/ports"
	"google.golang.org/grpc"
)

type Adapter struct {
	api                            ports.APIPort // core application dependency
	port                           int           // port to serve grpc on
	order.UnimplementedOrderServer               // for forward compatibility
}

func NewAdapter(api ports.APIPort, port int) *Adapter {
	return &Adapter{api: api, port: port}
}

func (a Adapter) Run() {
	var err error

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatalf("failed to listen to port %d, error: %v", a.port, err)
	}

	grpcServer := grpc.NewServer()
	order.RegisterOrderServer(grpcServer, a)
	// enable reflection to make grpcurl
	if config.GetEnv() == "development" {
		reflection.Register(grpcServer)
	}
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve grpc on port")
	}
}
