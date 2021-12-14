package main

import (
	"fmt"
	"github.com/subash68/ate/ate_setting_service/pkg/cmd"
	"os"
)

func main() {

	if err := cmd.RunServer(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	/*
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Setting service")
			fmt.Println("Hitting endpoint")

			con, err := grpc.Dial("ate_token_service:6060", grpc.WithInsecure())
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}

			defer con.Close()
			c := v1.NewTokenServiceClient(con)
			ctx, cancel := context.WithTimeout(context.Background(),5*time.Second)
			defer cancel()

			res1, err := c.Validate(ctx, &v1.ValidateRequest{
				Token: "Some token to validate",
			})
			fmt.Println(res1)
			if err != nil {
				log.Fatalf("Create failed: %v", err)
			}
			log.Printf("Create result: <%v>\n\n", res1.Status)

		})

		http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "pong pong pong")
		})

		http.ListenAndServe(":8080", nil)*/
}
