package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/Kobo-coder/miniproject2/api"
	"google.golang.org/grpc"
)

const dockerComposePort = ":50000"

var wantsAccess = false
var nextInLine = os.Args[1]
var port = flag.String("port", dockerComposePort, "The port to run the server on")
var client api.TokenServiceClient = nil

func main() {

	if len(os.Args) > 2 && os.Args[2] == "--start" {
		giveToken(nextInLine)
	}
	go worker()
	startServer()
}

func startServer() {
	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatalf("Failed ot listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	server := TokenServiceServer{}

	api.RegisterTokenServiceServer(grpcServer, &server)
	log.Printf("Token Service Server listening to %s\n", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func enter() {
	//Enter the critical section
	log.Println("Node entered critical section")
}

func resourceAccess() {
	fmt.Println("Doing God's work on critical section")
}

func exit() {
	log.Println("Node exited critical section")
	wantsAccess = false
}

func giveToken(nextNode string) {
	if client == nil {
		client = *newClient()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := client.RecieveToken(ctx, &api.Empty{})

	if err != nil {
		log.Fatalf("Dank shit my guy")
	}
}

func newClient() *api.TokenServiceClient {
	conn, err := grpc.Dial(nextInLine, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Failed to connect to next node")
	}

	client := api.NewTokenServiceClient(conn)

	return &client
}

type TokenServiceServer struct {
	api.UnimplementedTokenServiceServer
}

func worker() {
	//This is the goroutine that sleeps.
	for {
		sleep := rand.Intn(10)
		time.Sleep(time.Duration(sleep) * time.Second)
		wantsAccess = true
	}
}

func (s *TokenServiceServer) RecieveToken(context.Context, *api.Empty) (*api.Empty, error) {
	if wantsAccess {
		enter()
		resourceAccess()
		exit()
	}
	return &api.Empty{}, nil
}
