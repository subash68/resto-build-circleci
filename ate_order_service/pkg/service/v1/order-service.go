package v1

import (
	"context"
	"database/sql"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"github.com/subash68/ate/ate_order_service/pkg/api/order"
	"io/ioutil"
	"log"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiVersion = "v1"
)

type orderServiceServer struct {
	db *sql.DB
	order.UnimplementedOrderServiceServer
}

func NewOrderServiceServer(db *sql.DB) order.OrderServiceServer {
	return &orderServiceServer{db: db}
}

func (s *orderServiceServer) checkAPI(api string) error {
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements api version %s but asked for %s", apiVersion, api)
		}
	}

	return nil
}

func (s *orderServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "Failed to connect to database -> "+err.Error())
	}

	return c, nil
}

func LoadConfig() *firebase.Config {
	file, err := os.Open("./v1/firebase-config.json")

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	byteValue, _ := ioutil.ReadAll(file)

	var fireBaseConfig firebase.Config

	err = json.Unmarshal(byteValue, &fireBaseConfig)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}
	return &fireBaseConfig
}

func (s *orderServiceServer) collectTokensFromIDs(c *sql.Conn, ctx context.Context, id int64) *[]string {
	var tokens []string

	token := ""
	// Get Token of the client
	rows, err := c.QueryContext(ctx, "select us.notificationToken from cart c join users us on us.id = c.userId where id = ?;", id)
	if err != nil {
		return nil
	}

	if err := rows.Scan(token); err != nil {
		return nil
	} else {
		tokens = append(tokens, token)
	}
	_ = rows.Close()
	token = ""
	// Get Token of Restaurant
	rows2, err := c.QueryContext(ctx, "select us.notificationToken from cart_items c join products p on p.id = c.itemId join users us on us.id = p.userId where cartId = ?;", id)
	if err != nil {
		return nil
	}
	for rows2.Next() {
		token = ""
		if err := rows2.Scan(&token); err == nil {
			tokens = append(tokens, token)
		}
		token = ""
	}
	_ = rows2.Close()
	token = ""

	// Get Token of Delivery Guy

	return &tokens
}

func (s *orderServiceServer) firebaseMsgSend(msg, title string, tokens []string) {
	// Still need to configure the Configuration File of Firebase to allow access to the Firebase account
	config := LoadConfig()
	app, err := firebase.NewApp(context.Background(), config)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}
	registrationTokens := tokens
	message := &messaging.MulticastMessage{
		Data: map[string]string{
			"title": title,
			"msg":   msg,
		},
		Tokens: registrationTokens,
	}

	br, err := client.SendMulticast(context.Background(), message)
	if err != nil {
		log.Fatalln(err)
	}

	// See the BatchResponse reference documentation
	// for the contents of response.
	fmt.Printf("%d messages were sent successfully\n", br.SuccessCount)
}

func (s *orderServiceServer) ReadCart(c *sql.Conn, ctx context.Context, id int64) *order.Cart {

	var cart order.Cart
	var comments []*order.CartComment
	var items []*order.CartItem

	rows, err := c.QueryContext(ctx, "select id, deliveryAddress, instructions, coupon, hasCoupon, shipCost, cartTotal, cartState from cart where id = ?;", id)
	if err != nil {
		return nil
	}

	if err := rows.Scan(cart.Id, cart.Address, cart.Instructions, cart.Coupon, cart.HasCoupon, cart.ShippingCost, cart.Total, cart.State); err != nil {
		return nil
	}

	_ = rows.Close()

	rows2, err := c.QueryContext(ctx, "select id, userId, comment from cart_comment where cartId = ?", id)
	if err != nil {
		cart.Comments = nil
	} else {
		for rows2.Next() {
			comment := new(order.CartComment)
			if err := rows2.Scan(&comment.Id, &comment.UserId, &comment.Comment); err == nil {
				comments = append(comments, comment)
			}
		}
		cart.Comments = comments
	}

	_ = rows2.Close()

	rows3, err := c.QueryContext(ctx, "select id, itemId, itemCount, itemPrice from cart_items where cartId = ?", id)
	if err != nil {
		cart.Items = nil
	} else {
		for rows3.Next() {
			item := new(order.CartItem)
			if err := rows3.Scan(&item.Id, &item.ItemId, &item.ItemCount, &item.ItemPrice); err == nil {
				items = append(items, item)
			}
		}
		cart.Items = items
	}

	_ = rows3.Close()

	return &cart
}

// Request(context.Context, *RequestOrderRequest) (*RequestOrder, error)
//	// InPreparation
//	Preparation(context.Context, *InPreparationRequest) (*InPreparation, error)
//	// EnRoute
//	Route(context.Context, *EnRouteRequest) (*EnRoute, error)
//	// Delivered
//	Deliver(context.Context, *DeliveredRequest) (*Delivered, error)
//	// Disputed
//	Dispute(context.Context, *DisputedRequest) (*Disputed, error)
//	// Refunded
//	Refund(context.Context, *RefundedRequest) (*Refunded, error)
//	// Canceled
//	Cancel(context.Context, *CanceledRequest) (*Canceled, error)
//	// Finished
//	Finish(context.Context, *FinishedRequest) (*Finished, error)

//// Get all pending orders
//func (s *orderServiceServer) PendingOrders(ctx context.Context, req *order.PendingOrdersRequest) (*order.PendingOrdersResponse, error) {
//
//	if ctx == nil {
//		return &order.PendingOrdersResponse{
//			// Api: apiVersion,
//			// Error: &order.ResponseStatus{
//			// 	Status:  true,
//			// 	Message: "user authentication failed",
//			// },
//		}, nil
//	}
//	// Check api version
//	// if err := s.checkAPI(req.Api); err != nil {
//	// 	return &order.CreateResponse{
//	// 		Api: apiVersion,
//	// 		Error: &order.ResponseStatus{
//	// 			Status:  false,
//	// 			Message: "unsupported API version",
//	// 		},
//	// 	}, nil
//	// }
//
//	// c, err := s.connect(ctx)
//
//	// if err != nil {
//	// 	return nil, err
//	// }
//
//	// defer c.Close()
//
//	// // insert order into database
//	// currentUser := ctx.Value("userId").(int64)
//	// log.Print(req.order)
//
//	// res, err := c.ExecContext(ctx, "INSERT INTO products(`name`, `incredients`, `categoryId`, `status`, `orderOrder`, `isFeatured`, `position`, `price`, `discount`, `discountType`, `userId`, `productImage`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", req.order.Name, req.order.Incredients, req.order.CategoryId, req.order.Status, req.order.Order, req.order.IsFeatured, req.order.Position, req.order.Price, req.order.Discount, req.order.DiscountType, currentUser, req.order.ProductImageUrl)
//
//	// if err != nil {
//	// 	return nil, status.Errorf(codes.Unknown, "failed to insert order information"+err.Error())
//	// }
//
//	// id, err := res.LastInsertId()
//	// if err != nil {
//	// 	return nil, status.Error(codes.Unknown, "failed to retrieve id for created order "+err.Error())
//	// }
//
//	//Create order addons mapping here
//	// for i := 0; i < len(req.order.Addons); i++ {
//	// 	addonsResponse, err := c.ExecContext(ctx, "INSERT INTO product_addons (`productId`, `addonsId`) VALUES (?, ?)", id, req.order.Addons[i].Id)
//
//	// 	if err != nil {
//	// 		return nil, status.Errorf(codes.Unknown, "failed to insert order information"+err.Error())
//	// 	}
//
//	// 	addonsId, err := addonsResponse.LastInsertId()
//	// 	if err != nil {
//	// 		return nil, status.Error(codes.Unknown, "failed to retrieve id for created order "+err.Error())
//	// 	}
//
//	// 	log.Print(addonsId)
//	// }
//
//	log.Print("HITTING DISPATCHER SERVICE")
//	return &order.PendingOrdersResponse{
//		// Api: apiVersion,
//		// Id:  id,
//		// Error: &order.ResponseStatus{
//		// 	Status:  false,
//		// 	Message: "order created successfully",
//		// },
//	}, nil
//}
