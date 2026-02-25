package main

import (
	"context"
	"flag"
	"log"

	_ "github.com/mbobakov/grpc-consul-resolver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/bcachet/zero_downtime/helloworld"
)

var (
	name string
)

func main() {
	flag.StringVar(&name, "name", "", "Name to be used against greeter service")
	flag.Parse()
	conn, err := grpc.NewClient(
		"consul://consul:8500/greeter",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = conn.Close()
	}()

	// Create gRPC client
	client := pb.NewGreeterClient(conn)
	log.Printf("Sending SayHello request with name %s", name)
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp.Message)
}
