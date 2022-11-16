package api

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
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

		return c.String(http.StatusBadRequest, "Bad request")
	}

	return c.JSON(http.StatusOK, struct {
		Cash uint
	}{user.Cash})

}

func ReplenishBalance(c echo.Context) error {
	var user postgresql.User
	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}

	err = postgresql.Replenish(&user)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}
	return c.String(http.StatusOK, "Ok")
}

func Reserve(c echo.Context) error {
	var user postgresql.User
	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}
	err = postgresql.WriteTransaction(&user)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}
	return c.String(http.StatusOK, "Ok")
}

func Revenue(c echo.Context) error {
	var user postgresql.User
	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}

	err = postgresql.RecognizeRevenue(&user)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}

	return c.String(http.StatusOK, "Ok")
}

func Task1(c echo.Context) error {
	val := struct{ Date string }{}

	err := json.NewDecoder(c.Request().Body).Decode(&val)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}

	layout := "2006-01"

	t, err := time.Parse(layout, val.Date)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}

	err = postgresql.Task1(t)

	return err
}
