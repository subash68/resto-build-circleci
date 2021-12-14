package v1

import (
	"context"
	"database/sql"

	"github.com/subash68/ate/ate_menu_service/pkg/api/menu"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiVersion = "v1"
)

//TODO: get db instance here from service configuration
type menuServiceServer struct {
	db *sql.DB
	menu.UnimplementedMenuServiceServer
}

func NewMenuServiceServer(db *sql.DB) menu.MenuServiceServer {
	return &menuServiceServer{db: db}
}

func (s *menuServiceServer) checkAPI(api string) error {
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements api version %s but asked for %s", apiVersion, api)
		}
	}

	return nil
}

func (s *menuServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "Failed to connect to database -> "+err.Error())
	}

	return c, nil
}
