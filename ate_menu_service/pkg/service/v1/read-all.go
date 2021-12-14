package v1

import (
	"context"
	"log"

	"github.com/subash68/ate/ate_menu_service/pkg/api/menu"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/**
This is to Display list of items on the setting
*/
func (s *menuServiceServer) ReadAll(ctx context.Context, req *menu.ReadAllRequest) (*menu.ReadAllResponse, error) {

	if ctx == nil {
		return &menu.ReadAllResponse{
			Api: apiVersion,
			Error: &menu.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}

	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &menu.ReadAllResponse{
			Api: apiVersion,
			Error: &menu.ResponseStatus{
				Status:  false,
				Message: "unsupported API version",
			},
		}, nil
	}

	log.Println("Checked Api -> Next")

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}

	defer c.Close()

	// insert menu into database
	currentUser := ctx.Value("userId").(int64)

	log.Println("Current user is: ", currentUser)

	rows, err := c.QueryContext(ctx, "SELECT `id`, `name`, `incredients`, `status`, `menuOrder`, `isFeatured`, `position`, `price`, `discount`, `discountType`, `productImage` FROM products WHERE `userId` = ?", currentUser)

	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from categories -> "+err.Error())
	}
	// defer rows.Close()
	log.Println("Passed Error")
	data := []*menu.Menu{}

	for rows.Next() {
		td := new(menu.Menu)

		//TODO: few other items should also be added here
		if err := rows.Scan(&td.Id, &td.Name, &td.Incredients, &td.Status, &td.Order, &td.IsFeatured, &td.Position, &td.Price, &td.Discount, &td.DiscountType, &td.ProductImageUrl); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field values from menu row-> "+err.Error())
		}
		data = append(data, td)
	}

	_ = rows.Close()

	for _, v := range data {
		rows2, err := c.QueryContext(ctx, "SELECT `id`, `name` FROM addons WHERE `id` in (select `addonsId` from product_addons where `productId` = ?)", v.Id) // in (select `addonsId` from product_addons where `productId` = ?)", td.Id)

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
		v.Addons = addons
	}

	log.Println(data)

	// Returns: all menu details from server
	return &menu.ReadAllResponse{
		Api:   apiVersion,
		Menus: data,
		Error: &menu.ResponseStatus{
			Status:  false,
			Message: "",
		},
	}, nil
}
