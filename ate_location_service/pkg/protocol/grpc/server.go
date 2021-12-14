package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	v1 "github.com/subash68/ate/ate_location_service/pkg/api/location"
	"github.com/subash68/ate/ate_location_service/pkg/logger"
	"github.com/subash68/ate/ate_location_service/pkg/protocol/grpc/middleware"
	"google.golang.org/grpc"
)

func RunServer(ctx context.Context, v1API v1.LocationServiceServer, port string) error {
	log.Println(port)

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(middleware.EnsureValidToken),
	}

	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", port))

	if err != nil {
		log.Println("Check port number interpretation")
		return err
	}

	server := grpc.NewServer(opts...)
	v1.RegisterLocationServiceServer(server, v1API)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			logger.Log.Warn("shutting down gRPC server...")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()

	logger.Log.Info("starting gRPC server...")
	return server.Serve(listen)
}
