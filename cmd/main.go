package main

import (
	"fmt"
	"github.com/bossncn/go-boilerplate/cmd/app"
	"github.com/bossncn/go-boilerplate/config"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		fmt.Println("Error Load config:", err)
		return
	}

	app.Run(cfg)
}
