package main

import (
	"github.com/labstack/echo/v4"
	"turbo-carnival/internal/api"
)

func main() {

	e := echo.New()

	// Route => handler

	e.GET("/balance", api.GetBalance)
	e.GET("/replenish", api.ReplenishBalance)
	e.GET("/reserve", api.Reserve)

	// Start the Echo server
	e.Logger.Fatal(e.Start(":1010"))

}
