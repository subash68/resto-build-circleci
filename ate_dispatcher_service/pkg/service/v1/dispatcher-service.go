package v1

import (
	"context"
	"database/sql"
	"log"

	"github.com/subash68/ate/ate_dispatcher_service/pkg/api/dispatcher"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiVersion = "v1"
)

type dispatcherServiceServer struct {
	db *sql.DB
	dispatcher.UnimplementedDispatcherServiceServer
}

func NewDispatcherServiceServer(db *sql.DB) dispatcher.DispatcherServiceServer {
	return &dispatcherServiceServer{db: db}
}

func (s *dispatcherServiceServer) checkAPI(api string) error {
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements api version %s but asked for %s", apiVersion, api)
		}
	}

	return nil
}

func (s *dispatcherServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "Failed to connect to database -> "+err.Error())
	}

	return c, nil
}

// Get all pending orders
func (s *dispatcherServiceServer) PendingOrders(ctx context.Context, req *dispatcher.PendingOrdersRequest) (*dispatcher.PendingOrdersResponse, error) {

	if ctx == nil {
		return &dispatcher.PendingOrdersResponse{
			// Api: apiVersion,
			// Error: &dispatcher.ResponseStatus{
			// 	Status:  true,
			// 	Message: "user authentication failed",
			// },
		}, nil
	}
	// Check api version
	// if err := s.checkAPI(req.Api); err != nil {
	// 	return &dispatcher.CreateResponse{
	// 		Api: apiVersion,
	// 		Error: &dispatcher.ResponseStatus{
	// 			Status:  false,
	// 			Message: "unsupported API version",
	// 		},
	// 	}, nil
	// }

	// c, err := s.connect(ctx)

	// if err != nil {
	// 	return nil, err
	// }

	// defer c.Close()

	// // insert dispatcher into database
	// currentUser := ctx.Value("userId").(int64)
	// log.Print(req.dispatcher)

	// res, err := c.ExecContext(ctx, "INSERT INTO products(`name`, `incredients`, `categoryId`, `status`, `dispatcherOrder`, `isFeatured`, `position`, `price`, `discount`, `discountType`, `userId`, `productImage`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", req.dispatcher.Name, req.dispatcher.Incredients, req.dispatcher.CategoryId, req.dispatcher.Status, req.dispatcher.Order, req.dispatcher.IsFeatured, req.dispatcher.Position, req.dispatcher.Price, req.dispatcher.Discount, req.dispatcher.DiscountType, currentUser, req.dispatcher.ProductImageUrl)

	// if err != nil {
	// 	return nil, status.Errorf(codes.Unknown, "failed to insert dispatcher information"+err.Error())
	// }

	// id, err := res.LastInsertId()
	// if err != nil {
	// 	return nil, status.Error(codes.Unknown, "failed to retrieve id for created dispatcher "+err.Error())
	// }

	//Create dispatcher addons mapping here
	// for i := 0; i < len(req.dispatcher.Addons); i++ {
	// 	addonsResponse, err := c.ExecContext(ctx, "INSERT INTO product_addons (`productId`, `addonsId`) VALUES (?, ?)", id, req.dispatcher.Addons[i].Id)

	// 	if err != nil {
	// 		return nil, status.Errorf(codes.Unknown, "failed to insert dispatcher information"+err.Error())
	// 	}

	// 	addonsId, err := addonsResponse.LastInsertId()
	// 	if err != nil {
	// 		return nil, status.Error(codes.Unknown, "failed to retrieve id for created dispatcher "+err.Error())
	// 	}

	// 	log.Print(addonsId)
	// }

	log.Print("HITTING DISPATCHER SERVICE")
	return &dispatcher.PendingOrdersResponse{
		// Api: apiVersion,
		// Id:  id,
		// Error: &dispatcher.ResponseStatus{
		// 	Status:  false,
		// 	Message: "dispatcher created successfully",
		// },
	}, nil
}
