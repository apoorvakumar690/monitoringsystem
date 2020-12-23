package main

import (
	"context"
	"log"

	"google.golang.org/grpc"

	"stash.bms.bz/bms/monitoringsystem"
	"stash.bms.bz/bms/monitoringsystem/_example/apmgrpc/health"
	"stash.bms.bz/bms/monitoringsystem/apmgrpc"
)

func main() {
	apm, err := monitoringsystem.New(monitoringsystem.Elastic, true, monitoringsystem.Option{
		ElasticServiceName: "OTE-IN-GRPC-CLIENT-MONITORING",
	})
	// Set up a connection to the server.
	conn, err := grpc.Dial(
		":15001",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			apmgrpc.UnaryClientInterceptor(apmgrpc.WithAPM(apm)),
		),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	// Make a echo client and send gRPC.
	h := health.NewHealthClient(conn)
	response, err := h.Check(context.Background(), &health.Request{})
	if err != nil {
		log.Fatalf("Health check error: %v", err)
	}
	log.Println("Status : ", response.Status)
}
