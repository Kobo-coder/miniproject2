package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/Kobo-coder/miniproject2/api"
)

var wantsAccess = false
var nextInLine = os.Args[1]

func main() {
	if len(os.Args) > 2 && os.Args[2] == "--start" {
		giveToken(nextInLine)
	}
	go worker()
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

func giveToken() {

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
