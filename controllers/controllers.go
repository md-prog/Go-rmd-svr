package Controllers

import (
	"net/http"

	Models "github.com/chrislewispac/rmd-server/models"
	_ "github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

func TempCtrl(c echo.Context) error {

	//token := c.Get("user").(*jwt.Token)

	r := &Models.Response{
		Text: "This is a temp API handler",
	}

	return c.JSON(http.StatusOK, r)
}

//TempPostCtrl ...
func TempPostCtrl(c echo.Context) error {

	r := &Models.Response{
		Text: "Temp post Ctrl",
	}

	return c.JSON(http.StatusOK, r)
}
