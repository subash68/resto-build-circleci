package v1

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/subash68/ate/ate_menu_service/pkg/api/menu"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *menuServiceServer) Update(ctx context.Context, req *menu.UpdateRequest) (*menu.UpdateResponse, error) {

	if ctx == nil {
		return &menu.UpdateResponse{
			Api: apiVersion,
			Error: &menu.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}
	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &menu.UpdateResponse{
			Api: apiVersion,
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

	// insert menu into database

	modifiedAt := time.Now().UTC()

	res, err := c.ExecContext(ctx, "UPDATE products SET `name` = ?, `incredients` = ?, `status` = ?, `menuOrder` = ?, `modifiedAt` = ?, `productImage` = ? WHERE `id` = ?",
		req.Menu.Name, req.Menu.Incredients, req.Menu.Status, req.Menu.Order, modifiedAt, req.Menu.ProductImageUrl, req.Menu.Id)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to update menu information"+err.Error())
	}

	_, err = c.ExecContext(ctx, "DELETE FROM product_addons WHERE `productId`=?", req.Menu.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete menu -> "+err.Error())
	}

	for i := 0; i < len(req.Menu.Addons); i++ {
		addonsResponse, err := c.ExecContext(ctx, "INSERT INTO product_addons (`productId`, `addonsId`) VALUES (?, ?)", req.Menu.Id, req.Menu.Addons[i].Id)

		if err != nil {
			return nil, status.Errorf(codes.Unknown, "failed to insert menu information"+err.Error())
		}

		addonsId, err := addonsResponse.LastInsertId()
		if err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve id for created menu "+err.Error())
		}

		log.Print(addonsId)
	}

	r, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for update rows count "+err.Error())
	}

	if r == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("menu with ID='%d' is not found",
			req.Menu.Id))
	}

	return &menu.UpdateResponse{
		Api:     apiVersion,
		Updated: r,
		Error: &menu.ResponseStatus{
			Status:  false,
			Message: "menu updated successfully",
		},
	}, nil
}
