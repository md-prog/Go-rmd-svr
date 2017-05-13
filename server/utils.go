package server

import (
	"time"

	"github.com/Jeffail/gabs"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	null "gopkg.in/guregu/null.v3"
	"strconv"
)

const (
	dateLayout string = "2014-06-26T00:00:00.000-04:00"
)
func MustParseJSON(s string) *gabs.Container {
	jsonParsed, err := gabs.ParseJSON([]byte(s))
	if err != nil {
		panic(err)
	}
	return jsonParsed
}

func userID(c echo.Context) string {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	return claims["id"].(string)
}

func generateToken(userID string, exp time.Duration, admin bool) *jwt.Token {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	if admin {
		claims["admin"] = true
	}
	claims["id"] = userID
	claims["exp"] = time.Now().Add(exp).Unix()

	return token
}

func parseNullableFloat(str string) null.Float {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return null.Float{}
	}
	return null.FloatFrom(f)
}
//GetJsonBody ...
func GetJsonBody(s string) *gabs.Container {
	jsonParsed, err := gabs.ParseJSON([]byte(s))
	if err != nil {
		handleErr(err)
	}
	return jsonParsed
}
