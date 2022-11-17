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

	e := echo.New()

	e.GET("/balance", api.GetBalance)           //получение баланса пользователя
	e.PUT("/replenish", api.ReplenishBalance)   //пополнение счета пользователя
	e.PUT("/reserve", api.Reserve)              //резервирование средств со счета пользователя
	e.PUT("/revenue", api.Revenue)              //признание выручки
	e.GET("/monthly_report", api.MonthlyReport) //формирование месячного отчета
	e.GET("/report", api.Report)                //получение отчета

	e.Logger.Fatal(e.Start(c.ServerAddr))

}
