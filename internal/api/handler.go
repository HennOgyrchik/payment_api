package api

import (
	"encoding/json"
	_ "fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"turbo-carnival/internal/postgresql"
)

func GetBalance(c echo.Context) error {
	var user postgresql.User

	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}

	err = postgresql.GetBalance(&user)
	if err != nil {

		return c.String(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, user)

}

func ReplenishBalance(c echo.Context) error {
	user := struct {
		Id, Count uint
	}{}
	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}

	err = postgresql.Replenish(&user)
	if err != nil {
		return c.String(http.StatusNotFound, "User not found")
	}
	return c.String(http.StatusOK, "Ok")
}

func Reserve(c echo.Context) error {
	transaction := struct {
		UserID, ServiceID, OrderID, Cost uint
	}{}
	err := json.NewDecoder(c.Request().Body).Decode(&transaction)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}
	err = postgresql.WriteTransaction(transaction)
	if err != nil {
		return c.String(http.StatusNoContent, "Insufficient funds")
	}
	return c.String(http.StatusOK, "Ok")
}
