package api

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"turbo-carnival/internal/postgresql"
)

func GetBalance(c echo.Context) error {
	var user postgresql.User

	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if (err != nil) || (user.Id == 0) {
		return c.String(http.StatusBadRequest, "Bad request")
	}
	//fmt.Println(user)
	err = postgresql.GetBalance(&user)
	if err != nil {

		return c.String(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, user)

}

func ReplenishBalance(c echo.Context) error {
	var user postgresql.User
	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if (err != nil) || (user.Id == 0) {
		return c.String(http.StatusBadRequest, "Bad request")
	}

	err = postgresql.Replenish(&user)
	if err != nil {
		return c.String(http.StatusNotFound, "User not found")
	}
	return c.String(http.StatusOK, "Ok")
}

func Reserve(c echo.Context) error {
	var user postgresql.User
	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if (err != nil) || (user.Id == 0) {
		return c.String(http.StatusBadRequest, "Bad request")
	}
	err = postgresql.WriteTransaction(&user)
	if err != nil {
		return c.String(http.StatusNoContent, "Insufficient funds")
	}
	return c.String(http.StatusOK, "Ok")
}

func Revenue(c echo.Context) error {
	var user postgresql.User
	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if (err != nil) || (user.Id == 0) {
		return c.String(http.StatusBadRequest, "Bad request")
	}

	err = postgresql.RecognizeRevenue(&user)
	if err != nil {
		return c.String(http.StatusNoContent, "Insufficient funds") //изменить текст
	}

	return c.String(http.StatusOK, "Ok")
}
