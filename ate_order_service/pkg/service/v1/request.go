package v1

import (
	"context"
	"github.com/subash68/ate/ate_order_service/pkg/api/order"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strconv"
)

func (s *orderServiceServer) Request(ctx context.Context, req *order.RequestOrderRequest) (*order.RequestOrder, error) {

	if ctx == nil {
		return &order.RequestOrder{
			Api: apiVersion,
			Error: &order.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}
	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &order.RequestOrder{
			Api: apiVersion,
			Error: &order.ResponseStatus{
				Status:  false,
				Message: "unsupported API version",
			},
		}, nil
	}

	c, err := s.connect(ctx)

	if err != nil {
		return nil, err
	}

	// defer c.Close()

	log.Print(req.CartId)

	_, err = c.ExecContext(ctx, "UPDATE cart SET cartState = ? where id = ?", REQUESTED, req.CartId)

	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to Update Cart information"+err.Error())
	}

	cart := s.ReadCart(c, ctx, req.CartId)
	if cart == nil {
		return nil, status.Errorf(codes.Unknown, "failed to Collect Cart information"+err.Error())
	}

	tokens := s.collectTokensFromIDs(c, ctx, req.CartId)
	s.firebaseMsgSend("Cart number: "+strconv.Itoa(int(req.CartId))+" Had been Canceled", "New Notification", *tokens)

	return &order.RequestOrder{
		Api:  apiVersion,
		Cart: cart,
		Error: &order.ResponseStatus{
			Status:  false,
			Message: "",
		},
	}, nil
}
