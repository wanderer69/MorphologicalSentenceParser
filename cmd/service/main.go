package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	gw "github.com/wanderer69/MorphologicalSentenceParser/internal/gateway/grpc"
	pb "github.com/wanderer69/MorphologicalSentenceParser/pkg/server/grpc/morphological_parser"
)

func main() {
	ctx := context.Background()

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":2000")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
	signal.Notify(sigCh, syscall.SIGTERM)

	s := grpc.NewServer()
	pb.RegisterMorphologicalSentenceParserServer(s, gw.NewServer())
	// Serve gRPC server
	log.Println("Serving gRPC on connection ")
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	for {
		select {
		case <-sigCh:
			return
		case <-ctx.Done():
			return
		}
	}
}
