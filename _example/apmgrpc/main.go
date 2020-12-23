package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"stash.bms.bz/bms/monitoringsystem"
	"stash.bms.bz/bms/monitoringsystem/_example/apmgrpc/health"
	"stash.bms.bz/bms/monitoringsystem/apmgrpc"
)

// H monitors the underlying services that iam-broker relies on.
type H struct{}

// NewHealth new health instance
func NewHealth() *H {
	return new(H)
}

// Check reports whether the service is up or down.
func (*H) Check(ctx context.Context, _ *health.Request) (*health.Response, error) {
	return &health.Response{
		Status: 1,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":15001")
	if err != nil {
		log.Fatalln(err)
	}
	apm, _ := monitoringsystem.New(monitoringsystem.Elastic, true, monitoringsystem.Option{
		ElasticServiceName: "OTE-IN-GRPC-MONITORING",
	})
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			apmgrpc.UnaryServerInterceptor(apmgrpc.WithAPM(apm)),
		),
	)
	// Register gRPC server
	health.RegisterHealthServer(server, NewHealth())
	log.Printf("Server started listening on 15001")
	// Start serving requests
	if err = server.Serve(listener); err != nil {
		log.Fatalln("Server", err)
	}
}
