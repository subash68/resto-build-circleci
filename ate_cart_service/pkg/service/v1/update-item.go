package v1

import (
	"context"
	"log"

	"github.com/subash68/ate/ate_cart_service/pkg/api/cart"
)

func (s *cartServiceServer) UpdateItem(ctx context.Context, req *cart.UpdateItemRequest) (*cart.UpdateItemResponse, error) {
	// Check if context is present
	log.Println("serivce: remove item to cart")
	if ctx == nil {
		return &cart.UpdateItemResponse{
			Status: &cart.ResponseStatus{
				IsSuccess:    false,
				ErrorCode:    "E101",
				ErrorMessage: "User authentication failed",
			},
		}, nil
	}

	// currentUser := ctx.Value("userId").(int64)
	c, err := s.connect(ctx)

	if err != nil {
		return nil, err
	}

	defer c.Close()

	currentCart, _ := s.ReadCart(ctx, &cart.ReadCartRequest{
		Api: "v1",
	})

	result := c.QueryRowContext(ctx,
		"call sp_UPDATE_CART_ITEM(?, ?, ?)",
		currentCart.Cart.GetId(),
		req.GetItemId(),
		req.GetItemCount())

	// defer result.Close()

	var updated int32
	result.Scan(&updated)

	// update cart price

	return &cart.UpdateItemResponse{
		Status: &cart.ResponseStatus{
			IsSuccess:      updated > 0,
			SuccessMessage: "",
			ErrorCode:      "",
			ErrorMessage:   "",
			ErrorDetail:    "",
		},
	}, nil
}
