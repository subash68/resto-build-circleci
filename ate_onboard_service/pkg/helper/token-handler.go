package helper

import (
	"context"
	"crypto/sha1"
	"log"
	"time"

	"github.com/subash68/ate/ate_onboard_service/pkg/api/token"
	"google.golang.org/grpc"
)

func HashPassword(plainPassword string) string {
	h := sha1.New()

	h.Write([]byte(plainPassword))
	bs := h.Sum(nil)

	return string(bs)
}

func GenerateToken(id int64, fullname string, email string, userType uint8) string {
	con, err := grpc.Dial("ate_token_service:6060", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer con.Close()
	c := token.NewTokenServiceClient(con)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := c.Generate(ctx, &token.GenerateRequest{
		Id:       id,
		Fullname: fullname,
		Email:    email,
		UserType: int32(userType),
	})
	log.Printf("Response from token service %v ", response)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}
	log.Printf("Create result: <%v>\n\n", response.Token)

	return response.Token
}
