package v1

import (
	"context"
	"database/sql"

	"github.com/subash68/ate/ate_location_service/pkg/api/location"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiVersion = "v1"
)

type locationServiceServer struct {
	db *sql.DB
	location.UnimplementedLocationServiceServer
}

func NewLocationServiceServer(db *sql.DB) location.LocationServiceServer {
	return &locationServiceServer{db: db}
}

//TODO: This should be used where required
func (s *locationServiceServer) checkAPI(api string) error {
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements api version %s but asked for %s", apiVersion, api)
		}
	}

	return nil
}

func (s *locationServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "Failed to connect to database -> "+err.Error())
	}

	return c, nil
}
