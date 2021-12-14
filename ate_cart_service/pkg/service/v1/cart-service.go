package v1

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/subash68/ate/ate_cart_service/pkg/api/cart"
	"github.com/subash68/ate/ate_cart_service/pkg/api/menu"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	apiVersion = "v1"
)

type cartServiceServer struct {
	db *sql.DB
	cart.UnimplementedCartServiceServer
}

func NewCartServiceServer(db *sql.DB) cart.CartServiceServer {
	return &cartServiceServer{db: db}
}

//TODO: This should be used where required
func (s *cartServiceServer) checkAPI(api string) error {
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements api version %s but asked for %s", apiVersion, api)
		}
	}

	return nil
}

func (s *cartServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "Failed to connect to database -> "+err.Error())
	}

	return c, nil
}

func (s *cartServiceServer) GetPrice(ctx context.Context, itemId int32) *menu.CalculatePriceResponse {
	// GRPC calls are routed to 6060 port
	con, err := grpc.Dial("ate_menu_service:6060", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer con.Close()

	md, _ := metadata.FromIncomingContext(ctx)

	log.Println("token from out going request")
	log.Println(md["authorization"])

	c := menu.NewMenuServiceClient(con)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", md["authorization"][0])

	response, err := c.CalculatePrice(ctx, &menu.CalculatePriceRequest{
		Api:     apiVersion,
		Product: itemId,
	})
	log.Println(response)
	if err != nil {
		log.Fatalf("Validated failed : %v", err)
	}
	log.Printf("validation result : %v \n\n", response.DiscountedPrice)

	// update datebase here

	return response

}
