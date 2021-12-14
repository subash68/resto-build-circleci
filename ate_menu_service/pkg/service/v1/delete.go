package v1

import (
	"context"
	"fmt"

	"github.com/subash68/ate/ate_menu_service/pkg/api/menu"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *menuServiceServer) Delete(ctx context.Context, req *menu.DeleteRequest) (*menu.DeleteResponse, error) {

	if ctx == nil {
		return &menu.DeleteResponse{
			Api: apiVersion,
			Error: &menu.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}

	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &menu.DeleteResponse{
			Error: &menu.ResponseStatus{
				Status:  false,
				Message: "unsupported API version",
			},
		}, nil
	}

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	// delete ToDo
	res, err := c.ExecContext(ctx, "DELETE FROM products WHERE `id`=?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete menu -> "+err.Error())
	}

	_, err = c.ExecContext(ctx, "DELETE FROM product_addons WHERE `productId`=?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete menu -> "+err.Error())
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("menu with ID='%d' is not found",
			req.Id))
	}

	return &menu.DeleteResponse{
		Api:     apiVersion,
		Deleted: int32(rows),
		Error: &menu.ResponseStatus{
			Status:  false,
			Message: "",
		},
	}, nil
}
