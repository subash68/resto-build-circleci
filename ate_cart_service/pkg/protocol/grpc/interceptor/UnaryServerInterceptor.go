package interceptor

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Get the metadata from the incoming context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, fmt.Errorf("couldn't parse incoming context metadata")
		}
		// Retrieve the client OS, this will be empty if it does not exist
		os := md.Get("client-os")
		// Get the client IP Address
		ip, err := getClientIP(ctx)
		if err != nil {
			return nil, err
		}
		// Populate the EdgeLocation type with the IP and OS
		// req.(*api.EdgeLocation).IpAddress = ip
		// req.(*api.EdgeLocation).OperatingSystem = os[0]

		h, err := handler(ctx, req)
		log.Printf("server interceptor hit: hydrating type with OS: '%v' and IP: '%v'", os[0], ip)
		return h, err
	}
}

// GetClientIP inspects the context to retrieve the ip address of the client
func getClientIP(ctx context.Context) (string, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("couldn't parse client IP address")
	}
	return p.Addr.String(), nil
}
