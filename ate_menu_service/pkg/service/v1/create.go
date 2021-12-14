package v1

import (
	"context"
	"log"

	"github.com/subash68/ate/ate_menu_service/pkg/api/menu"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *menuServiceServer) Create(ctx context.Context, req *menu.CreateRequest) (*menu.CreateResponse, error) {

	if ctx == nil {
		return &menu.CreateResponse{
			Api: apiVersion,
			Error: &menu.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}
	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &menu.CreateResponse{
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
	currentUser := ctx.Value("userId").(int64)
	log.Print(req.Menu)

	res, err := c.ExecContext(ctx, "INSERT INTO products(`name`, `incredients`, `categoryId`, `status`, `menuOrder`, `isFeatured`, `position`, `price`, `discount`, `discountType`, `userId`, `productImage`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", req.Menu.Name, req.Menu.Incredients, req.Menu.CategoryId, req.Menu.Status, req.Menu.Order, req.Menu.IsFeatured, req.Menu.Position, req.Menu.Price, req.Menu.Discount, req.Menu.DiscountType, currentUser, req.Menu.ProductImageUrl)

	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to insert menu information"+err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created menu "+err.Error())
	}

	//Create menu addons mapping here
	for i := 0; i < len(req.Menu.Addons); i++ {
		addonsResponse, err := c.ExecContext(ctx, "INSERT INTO product_addons (`productId`, `addonsId`) VALUES (?, ?)", id, req.Menu.Addons[i].Id)

		if err != nil {
			return nil, status.Errorf(codes.Unknown, "failed to insert menu information"+err.Error())
		}

		addonsId, err := addonsResponse.LastInsertId()
		if err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve id for created menu "+err.Error())
		}

		log.Print(addonsId)
	}

	return &menu.CreateResponse{
		Api: apiVersion,
		Id:  id,
		Error: &menu.ResponseStatus{
			Status:  false,
			Message: "menu created successfully",
		},
	}, nil
}
