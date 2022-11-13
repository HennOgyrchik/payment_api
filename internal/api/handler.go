package api

import (
	"encoding/json"
	_ "fmt"
	"github.com/labstack/echo/v4"
	"turbo-carnival/internal/postgresql"
)

func GetBalance(c echo.Context) error {
	user := struct {
		Id   int
		Cash string
	}{}
	err := json.NewDecoder(c.Request().Body).Decode(&user)
	if err != nil {
		return err
	}

	user.Cash, err = postgresql.GetBalance(user.Id)

	return c.JSON(200, user)

}
