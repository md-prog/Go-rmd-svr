package Controllers

import (
	"net/http"

	_ "github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

//GetUserByIDCtrl ...
func GetUserByIDCtrl(c echo.Context) error {
	return c.JSON(http.StatusOK, "user")
}
