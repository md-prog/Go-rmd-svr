package Controllers

import (
	"github.com/Jeffail/gabs"
	"github.com/labstack/echo"
	"io/ioutil"
)

func GetJsonBody(c echo.Context) *gabs.Container {
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		handleErr(err)
	}
	jsonParsed, _ := gabs.ParseJSON([]byte(b))
	return jsonParsed
}
