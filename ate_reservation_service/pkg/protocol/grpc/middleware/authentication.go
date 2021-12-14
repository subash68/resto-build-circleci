package middleware

import (
	"context"
	"log"
	"strings"

	"github.com/subash68/ate/ate_reservation_service/pkg/helper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

type user struct {
	id string
}

type ContextUserKey struct {
	message string
	status  bool
}

func EnsureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	// The keys within metadata.MD are normalized to lowercase.
	// See: https://godoc.org/google.golang.org/grpc/metadata#New
	state, userId, _ := valid(md["authorization"])

	if !state {
		log.Printf("Something is wrong with validation...")

		return handler(nil, req)
	}

	newCtx := context.WithValue(ctx, "userId", userId)
	return handler(newCtx, req)
}

// valid validates the authorization.
func valid(authorization []string) (bool, int64, error) {
	if len(authorization) < 1 {
		return false, 0, nil
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	response := helper.Validate(token)

	return response.Status, response.Id, nil
}
