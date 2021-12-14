package v1

import (
	"context"
	"log"

	"github.com/subash68/ate/ate_cart_service/pkg/api/cart"
)

func (s *cartServiceServer) ReadItems(ctx context.Context, cartId int32) []*cart.CartItem {
	c, err := s.connect(ctx)

	if err != nil {
		return nil
	}

	defer c.Close()

	result, _ := c.QueryContext(ctx,
		"call sp_GET_CART_ITEMS(?)", cartId)

	defer result.Close()

	// result should return item inserted
	Items := []*cart.CartItem{}
	for result.Next() {

		var id int32
		var itemId int32
		var itemCount int32
		var itemPrice float32

		scanerr := result.Scan(&id, &itemId, &itemCount, &itemPrice)
		if scanerr != nil {
			log.Print(scanerr)
		}

		Items = append(Items, &cart.CartItem{
			Id: id, ItemId: itemId, ItemCount: itemCount,
		})
		log.Println(Items)
	}

	return Items

}

// Get all items added in cart
func (s *cartServiceServer) ReadCart(ctx context.Context, req *cart.ReadCartRequest) (*cart.ReadCartResponse, error) {

	if ctx == nil {
		return &cart.ReadCartResponse{
			Cart:   &cart.Cart{},
			Status: &cart.ResponseStatus{IsSuccess: false, SuccessMessage: "", ErrorCode: "", ErrorMessage: "", ErrorDetail: ""},
		}, nil
	}

	currentUser := ctx.Value("userId").(int64)
	c, err := s.connect(ctx)

	if err != nil {
		return nil, err
	}

	defer c.Close()

	row := c.QueryRowContext(ctx,
		"call sp_READ_CART(?)", currentUser)

	// defer result.Close()

	var id int32
	var orderType string
	var address string
	var coupon string
	var hasCoupon bool
	var shipCost float32
	var cartTotal float32
	var instructions string
	row.Scan(
		&id,
		&orderType,
		&address,
		&coupon,
		&hasCoupon,
		&shipCost,
		&cartTotal,
		&instructions,
	)

	return &cart.ReadCartResponse{
		Cart: &cart.Cart{
			Id:           id,
			Type:         orderType,
			Address:      address,
			Coupon:       coupon,
			HasCoupon:    hasCoupon,
			ShippingCost: shipCost,
			Total:        cartTotal,
			Instructions: instructions,
			Items:        s.ReadItems(ctx, id),
		},
		Status: &cart.ResponseStatus{IsSuccess: false, SuccessMessage: "", ErrorCode: "", ErrorMessage: "", ErrorDetail: ""},
	}, nil
}
