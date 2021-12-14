package helper

import (
	"context"
	"log"
	"time"

	"github.com/subash68/ate/ate_category_service/pkg/api/token"
	"google.golang.org/grpc"
)

// Validate create a client here to validate token from token service
func Validate(tokenString string) *token.ValidateResponse {

	con, err := grpc.Dial("ate_token_service:6060", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer con.Close()
	c := token.NewTokenServiceClient(con)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := c.Validate(ctx, &token.ValidateRequest{
		Token: tokenString,
	})
	log.Println(response)
	if err != nil {
		log.Fatalf("Validation failed failed: %v", err)
	}
	log.Printf("Validation result: <%v>\n\n", response.Status)

	return response
}
