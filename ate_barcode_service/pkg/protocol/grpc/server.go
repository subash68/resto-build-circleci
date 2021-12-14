package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"os"
	"os/signal"

	v1 "github.com/subash68/ate/ate_category_service/pkg/api/category"
	"github.com/subash68/ate/ate_category_service/pkg/logger"
	"github.com/subash68/ate/ate_category_service/pkg/protocol/grpc/middleware"
	"google.golang.org/grpc"
)

const (
	DbCrt = "./cert/db_server.crt"
	DbKey = "./cert/db_server.key"
)

func RunServer(ctx context.Context, v1API v1.CategoryServiceServer, port string) error {
	log.Println(port)

	cert, err := tls.LoadX509KeyPair(DbCrt, DbKey)
	if err != nil {
		log.Println("Check Your Certifications")
		return err
	}

	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
		grpc.UnaryInterceptor(middleware.EnsureValidToken),
	}

	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", port))

	if err != nil {
		log.Println("Check port number interpretation")
		return err
	}

	//FIXME: panic: The unary server interceptor was already set and may not be reset.
	// opts = middleware.AddLogging(logger.Log, opts)

	//FIXME: need to look at how to set content type header in grpc response
	// header := metadata.Pairs("Content-Type", "application/json")
	// grpc.SetHeader(ctx, header)

	server := grpc.NewServer(opts...)
	v1.RegisterCategoryServiceServer(server, v1API)

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
