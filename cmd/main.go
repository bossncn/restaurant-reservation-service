package main

import (
	"fmt"
	"github.com/bossncn/restaurant-reservation-service/cmd/app"
	"github.com/bossncn/restaurant-reservation-service/config"
)

// @title Restaurant Reservation Service
// @version 1.0
// @description Service for managing table reservations in a restaurant.
// @termsOfService http://swagger.io/terms/
//
// @license.name MIT
// @license.url https://github.com/bossncn/restaurant-reservation-service/blob/main/LICENSE
func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		fmt.Println("Error Load config:", err)
		return
	}

	app.Run(cfg)
}
