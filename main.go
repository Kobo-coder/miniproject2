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
var port = flag.String("port", dockerComposePort, "The port to run the server on")
var client api.TokenServiceClient = nil
var nextNode = os.Args[1]

func main() {
	if len(os.Args) > 2 && os.Args[2] == "--start" {
		go giveToken(nextNode)
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
		client = *newClient(nextNode)
	}

	time.Sleep(time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	log.Printf("Making request to next in line")
	_, err := client.ReceiveToken(ctx, &api.Empty{})

	if err != nil {
		log.Fatalf("Dank shit my guy: %v", err)
	}

	log.Printf("I passed the spliff")
}

func newClient(nextNode string) *api.TokenServiceClient {
	conn, err := grpc.Dial(nextNode, grpc.WithInsecure(), grpc.WithBlock())

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

func (s *TokenServiceServer) ReceiveToken(context.Context, *api.Empty) (*api.Empty, error) {
	log.Printf("Received token")
	if wantsAccess {
		enter()
		resourceAccess()
		exit()
	}

	go giveToken(nextNode)

	return &api.Empty{}, nil
}
