package main

import (
	"flag"
	"log"

	"github.com/feliksik/restgrpc"
	"github.com/feliksik/restgrpc/example/compute"
	"github.com/feliksik/restgrpc/example/sayhello"
)

func main() {
	flag.Parse()

	proxy := restgrpc.NewGateway()

	proxy.AddService(restgrpc.NewEndpoint(
		"Compute",
		"localhost:50051",
		compute.RegisterComputeServiceHandler),
	)

	proxy.AddService(restgrpc.NewEndpoint(
		"SayHello",
		"localhost:50052",
		sayhello.RegisterSayHelloHandler,
	))

	if err := proxy.Bind(":8080"); err != nil {
		log.Fatal(err)
	}
}
