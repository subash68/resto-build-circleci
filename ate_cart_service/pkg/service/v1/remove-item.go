package v1

import (
	"context"
	"log"

	"github.com/subash68/ate/ate_cart_service/pkg/api/cart"
)

func (s *cartServiceServer) RemoveItem(ctx context.Context, req *cart.RemoveItemRequest) (*cart.RemoveItemResponse, error) {
	// Check if context is present
	log.Println("serivce: remove item to cart")
	if ctx == nil {
		return &cart.RemoveItemResponse{
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

	log.Print(currentCart)

	result, _ := c.QueryContext(ctx,
		"call sp_REMOVE_CART_ITEM(?, ?)",
		currentCart.Cart.Id,
		req.ItemId)

	defer result.Close()

	return &cart.RemoveItemResponse{
		Status: &cart.ResponseStatus{
			IsSuccess:      false,
			SuccessMessage: "",
			ErrorCode:      "",
			ErrorMessage:   "",
			ErrorDetail:    "",
		},
	}, nil
}
