package v1

import (
	"context"
	"github.com/subash68/ate/ate_order_service/pkg/api/order"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strconv"
)

func (s *orderServiceServer) Dispute(ctx context.Context, req *order.DisputedRequest) (*order.Disputed, error) {
	if ctx == nil {
		return &order.Disputed{
			Api: apiVersion,
			Error: &order.ResponseStatus{
				Status:  true,
				Message: "user authentication failed",
			},
		}, nil
	}
	// Check api version
	if err := s.checkAPI(req.Api); err != nil {
		return &order.Disputed{
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

	currentUser := ctx.Value("userId").(int64)

	_, err = c.ExecContext(ctx, "UPDATE cart SET cartState = ? where id = ?", DISPUTED, req.CartId)

	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to Update Cart information"+err.Error())
	}

	_, err = c.ExecContext(ctx, "UPDATE cart_comment SET comment = ?, cartId = ?, userId = ? where id = ?", req.Comment, req.CartId, currentUser)

	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to Update Cart information"+err.Error())
	}

	cart := s.ReadCart(c, ctx, req.CartId)
	if cart == nil {
		return nil, status.Errorf(codes.Unknown, "failed to Collect Cart information"+err.Error())
	}

	tokens := s.collectTokensFromIDs(c, ctx, req.CartId)
	s.firebaseMsgSend("Cart number: "+strconv.Itoa(int(req.CartId))+" Had been Canceled", "New Notification", *tokens)

	return &order.Disputed{
		Api:  apiVersion,
		Cart: cart,
		Error: &order.ResponseStatus{
			Status:  false,
			Message: "",
		},
	}, nil
}
