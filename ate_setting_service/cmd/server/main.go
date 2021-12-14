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
}
