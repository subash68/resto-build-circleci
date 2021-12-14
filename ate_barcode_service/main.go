package main

import (
	"fmt"
	"os"

	"github.com/subash68/ate/ate_category_service/pkg/cmd"
)

func main() {
	if err := cmd.RunServer(); err != nil {
		fmt.Fprint(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
