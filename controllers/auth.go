package Controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	Models "github.com/chrislewispac/rmd-server/models"
	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

//LoginCtrl ...
func LoginCtrl(c echo.Context) (err error) {
	u := new(Models.User)

	if err = c.Bind(u); err != nil {
		return c.JSON(http.StatusInternalServerError, "There was nothing submitted")
	}

	var user Models.User
	err = Models.DB.QueryRowx(`SELECT * FROM users WHERE email=$1`, u.Email).StructScan(&user)
	if err != nil {
		handleErr(err)
	}

	password := []byte(u.Password)
	dbPassword := []byte(user.Password)
	err = bcrypt.CompareHashAndPassword(dbPassword, password)
	if err != nil {
		r := &Models.Response{
			Errors: "ERROR: Wrong Password.",
		}
		return c.JSON(http.StatusUnauthorized, r)
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as Models.Response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	err = Models.Client.Set(t, "valid", 0).Err()
	if err != nil {
		handleErr(err)
	}

	anon := struct {
		User Models.User `json:"user"`
	}{
		user,
	}

	data, _ := json.Marshal(anon)

	res := Models.NewResponse()
	res.Msg = "Successfully Logged In"
	res.Token = t
	res.Data = data

	return c.JSON(http.StatusOK, res)

}

//RegisterCtrl ...
func RegisterCtrl(c echo.Context) (err error) {
	u := new(Models.User)

	if err = c.Bind(u); err != nil {
		return c.JSON(http.StatusInternalServerError, "unable to parse data")
	}

	password := []byte(u.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	if err != nil {
		handleErr(err)
	}

	var user Models.User
	err = Models.DB.QueryRowx(`
		INSERT INTO users
			( email
			, password )
		VALUES ($1, $2)
		RETURNING *`, u.Email, hashedPassword).StructScan(&user)
	if err != nil {
		log.Warn(err)
		r := &Models.Response{
			Errors: fmt.Sprintf("ERROR: Account already exists with email %s", u.Email),
		}
		return c.JSON(http.StatusConflict, r)
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	claims["id"] = 16
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as Models.Response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	err = Models.Client.Set(t, "valid", 0).Err()
	if err != nil {
		handleErr(err)
	}

	res := Models.NewResponse()
	res.Msg = "Successfully Registered"
	res.Token = t

	return c.JSON(http.StatusOK, res)
}

//LogoutCtrl ...
func LogoutCtrl(c echo.Context) error {

	token := c.Get("user").(*jwt.Token)

	err := Models.Client.Del(token.Raw).Err()
	if err != nil {
		handleErr(err)
	}

	res := Models.NewResponse()
	res.Msg = "Logged Out and Token Revoked"
	res.Token = ""

	return c.JSON(http.StatusOK, res)
}

//LogoutCtrl ...
func ForgotPasswordCtrl(c echo.Context) error {

	body := GetJsonBody(c)

	fmt.Println(body)

	return c.JSON(http.StatusOK, "password reset email sent")
}
