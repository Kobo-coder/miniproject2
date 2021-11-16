package main

import (
	"context"
	"flag"
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
		go giveToken()
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
	log.Println("Doing God's work on critical section . . .")
	time.Sleep(1 * time.Second)
	wantsAccess = false
}

func exit() {
	log.Println("Node exited critical section")
}

func giveToken() {
	if client == nil {
		client = *newClient(nextNode)
	}

	time.Sleep(50 * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	log.Printf("Making request to next in line")
	_, err := client.ReceiveToken(ctx, &api.Empty{})

	if err != nil {
		log.Fatalf("Dank shit my guy: %v", err)
	}
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
	rand.Seed(time.Now().UnixNano())
	for {
		sleep := rand.Intn(20)
		time.Sleep(time.Duration(sleep+10) * time.Second)
		log.Println("I want access NOW!!!")
		wantsAccess = true
	}
}

func (s *TokenServiceServer) ReceiveToken(context.Context, *api.Empty) (*api.Empty, error) {
	log.Printf("Received token")
	if wantsAccess {
		go enterCriticalSectionAndGiveToken()
	} else {
		log.Printf("Not interested in entering critical section, passing token along . . .")
		go giveToken()
	}

	return &api.Empty{}, nil
}

func enterCriticalSectionAndGiveToken() {
	enter()
	resourceAccess()
	exit()
	giveToken()
}
