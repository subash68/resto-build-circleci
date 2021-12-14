package v1

import (
	"context"
	"log"

	"github.com/subash68/ate/ate_cart_service/pkg/api/cart"
)

func (s *cartServiceServer) DeleteCart(ctx context.Context, req *cart.DeleteCartRequest) (*cart.DeleteCartResponse, error) {
	if ctx == nil {
		return &cart.DeleteCartResponse{
			Status: &cart.ResponseStatus{IsSuccess: false, SuccessMessage: "", ErrorCode: "", ErrorMessage: "", ErrorDetail: ""},
		}, nil
	}

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
		"call sp_REMOVE_CART(?)",
		currentCart.Cart.Id)

	defer result.Close()

	return &cart.DeleteCartResponse{
		Status: &cart.ResponseStatus{
			IsSuccess:      false,
			SuccessMessage: "",
			ErrorCode:      "",
			ErrorMessage:   "",
			ErrorDetail:    "",
		},
	}, nil

}
