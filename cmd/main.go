package main

import (
	"github.com/labstack/echo/v4"
	"os"
	"turbo-carnival/internal/api"
	"turbo-carnival/internal/config"
)

func main() {
	err := config.InitConfig(os.Args)
	if err != nil {
		panic(err)
		return
	}
	c := config.GetConfig()

	////////////////////////////////////////////////////////////

	e := echo.New()

	// Route => handler

	e.GET("/balance", api.GetBalance)
	e.PUT("/replenish", api.ReplenishBalance)
	e.PUT("/reserve", api.Reserve)
	e.PUT("/revenue", api.Revenue)

	// Start the Echo server
	e.Logger.Fatal(e.Start(c.ServerAddr))

}
