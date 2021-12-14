package v1

import (
	"context"
	"log"

	"github.com/subash68/ate/ate_menu_service/pkg/api/menu"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *menuServiceServer) Read(ctx context.Context, req *menu.ReadRequest) (*menu.ReadResponse, error) {

	if ctx == nil {
		return &menu.ReadResponse{
			Api: apiVersion,
			Error: &menu.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}

	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &menu.ReadResponse{
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

	rows, err := c.QueryContext(ctx, "SELECT `id`, `name`, `incredients`, `status`, `menuOrder`, `isFeatured`, `position`, `price`, `discount`, `discountType`, `productImage` FROM products WHERE `id` = ?", req.Id)

	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from categories -> "+err.Error())
	}

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from categories -> "+err.Error())
		}

		return &menu.ReadResponse{
			Api: apiVersion,
			Error: &menu.ResponseStatus{
				Status:  true,
				Message: "menu not found",
			},
		}, nil
	}

	var data menu.Menu

	if err := rows.Scan(&data.Id, &data.Name, &data.Incredients, &data.Status, &data.Order, &data.IsFeatured, &data.Position, &data.Price, &data.Discount, &data.DiscountType, &data.ProductImageUrl); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from menu row-> "+err.Error())
	}

	_ = rows.Close()

	log.Println(data)

	log.Println("Start: ", data.Id)
	if ctx != nil {
		log.Println(c)
		rows2, err := c.QueryContext(ctx, "SELECT `id`, `name` FROM addons WHERE `id` in (select `addonsId` from product_addons where `productId` = ?)", data.Id) // in (select `addonsId` from product_addons where `productId` = ?)", td.Id)
		log.Println("Error: ", err)
		if err != nil {
			return nil, status.Error(codes.Unknown, "failed to select from categories -> "+err.Error())
		}
		log.Println("Error Passed")
		addons := []*menu.Addons{}
		for rows2.Next() {
			a := new(menu.Addons)
			if err := rows2.Scan(&a.Id, &a.Name); err != nil {
				return nil, status.Error(codes.Unknown, "failed to retrieve field values from menu row-> "+err.Error())
			}
			log.Println("Passed", a.Name)
			addons = append(addons, a)
		}
		data.Addons = addons
	}

	// {GRPC} Get reservation info here

	// {GRPC} Get all addons here
	/**
	Get all addons for select product from addons service
	*/
	// data.Addons =

	return &menu.ReadResponse{
		Api:  apiVersion,
		Menu: &data,
		Error: &menu.ResponseStatus{
			Status:  false,
			Message: "",
		},
	}, nil
}
