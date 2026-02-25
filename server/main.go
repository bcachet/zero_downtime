// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/bcachet/zero_downtime/helloworld"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

var (
	port  int
	delay int
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	time.Sleep(time.Duration(delay) * time.Second)
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	flag.IntVar(&port, "port", 50051, "The server port")
	flag.IntVar(&delay, "delay", 0, "Delay before returning the response")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})

	// Create a channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Register service in Consul
	config := api.DefaultConfig()
	config.Address = "consul:8500"
	consul, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating Consul client: %v", err)
	}
	name := "blue"
	if port == 50001 {
		name = "green"
	}
	err = consul.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      fmt.Sprintf("greeter-%s", name),
		Name:    "greeter",
		Tags:    []string{name},
		Port:    port,
		Address: fmt.Sprintf("greeter-server-%s", name),
	})
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	// Run the server in a separate goroutine
	go func() {
		log.Printf("Server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Block until a signal is received
	<-stop
	log.Println("Shutting down server gracefully...")

	// Gracefully stop the server
	s.GracefulStop()

	log.Println("Server stopped")
}
