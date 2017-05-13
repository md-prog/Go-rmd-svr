package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	Models "gitlab.com/chrislewispac/rmd-server/models"
	"golang.org/x/crypto/bcrypt"
	null "gopkg.in/guregu/null.v3"
	"gopkg.in/redis.v5"
)

//Login ...
func (s *Server) Login(c echo.Context) (err error) {
	u := new(Models.User)
	var user Models.User

	if err = c.Bind(u); err != nil {
		res := createAuthErrorResponse(user, "No user was found")
		return c.JSON(http.StatusOK, res)
	}

	u.Email = strings.ToLower(u.Email)

	err = s.db.QueryRowx(Models.GetUserByEmail, u.Email).StructScan(&user)
	if err != nil {
		handleErr(err)
		res := createAuthErrorResponse(user, "User Not Found")
		return c.JSON(http.StatusOK, res)
	}

	password := []byte(u.Password)
	dbPassword := []byte(user.Password)
	err = bcrypt.CompareHashAndPassword(dbPassword, password)
	if err != nil {
		handleErr(err)
		res := createAuthErrorResponse(user, "Wrong Password")
		return c.JSON(http.StatusOK, res)
	}
	// Generate encoded token and send it as Models.Response.
	t, err := generateToken(user.ID, 72*time.Hour, true).SignedString(SigningKey)
	if err != nil {
		handleErr(err)
		createAuthErrorResponse(user, "unable to create token")
		return err
	}

	err = s.storeToken(t, time.Hour*72)
	if err != nil {
		handleErr(err)
		createAuthErrorResponse(user, "unable to create redis token")
		handleErr(err)
	}
	user.IntercomHash = null.StringFrom(s.intercom.CalculateHash(user.ID))
	user.Password = ""
	user.Token = null.StringFrom(t)
	res := createAuthSuccessResponse(user, "Successfully Logged In")
	return c.JSON(http.StatusOK, res)
}

func (s *Server) storeToken(t string, d time.Duration) error {
	return s.redis.Set(t, "valid", d).Err()
}

func (s *Server) isTokenStored(t string) (bool, error) {
	r, err := s.redis.Get(t).Result()
	if err == redis.Nil {
		return false, nil
	}
	return r == "valid", err
}

//Register ...
func (s *Server) Register(c echo.Context) (err error) {
	u := new(Models.User)

	if err = c.Bind(u); err != nil {
		handleErr(err)
		return c.JSON(http.StatusInternalServerError, "unable to parse data")
	}

	u.Email = strings.ToLower(u.Email)

	hashedPassword, err := s.hashPass(u.Password)
	if err != nil {
		handleErr(err)
	}

	var user Models.User
	if err := s.insertUser(u.Email, hashedPassword, &user); err != nil {
		handleErr(err)
		res := createAuthErrorResponse(user, fmt.Sprintf("Account already exists with email %s", u.Email))
		return c.JSON(http.StatusOK, res)
	}

	// Generate encoded token and send it as Models.Response.
	t, err := generateToken(user.ID, 72*time.Hour, true).SignedString(SigningKey)
	if err != nil {
		handleErr(err)
		return err
	}

	err = s.storeToken(t, time.Hour*72)
	if err != nil {
		handleErr(err)
	}

	user.Password = ""
	user.Token = null.StringFrom(t)

	res := createAuthSuccessResponse(user, "Successfully Registered")
	return c.JSON(http.StatusOK, res)
}

func (s *Server) hashPass(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (s *Server) insertUser(email string, hashedPass []byte, user *Models.User) error {
	return s.db.QueryRowx(Models.CreateUser, email, hashedPass).StructScan(user)
}

//Logout TODO
func (s *Server) Logout(c echo.Context) error {
	var user Models.User
	if c.Get("user") == nil {
		res := createAuthErrorResponse(user, "Invalid Token")
		return c.JSON(http.StatusOK, res)
	}
	token := c.Get("user").(*jwt.Token)
	numDel, _ := s.redis.Del(token.Raw).Result()
	if numDel < 1 {
		res := createAuthErrorResponse(user, "Token not in storage")
		return c.JSON(http.StatusOK, res)
	}
	user.Password = ""
	user.Token = null.StringFrom("")
	res := createAuthSuccessResponse(user, "Logged Out")
	return c.JSON(http.StatusOK, res)
}

//ForgotPassword TODO
func (s *Server) ForgotPassword(c echo.Context) error {
	var user Models.User

	body := struct {
		Email string `json:"email" form:"email"`
	}{}
	if err := c.Bind(&body); err != nil {
		handleErr(err)
		res := createAuthErrorResponse(user, errMsgBadRequest)
		return c.JSON(http.StatusOK, res)
	}
	body.Email = strings.ToLower(body.Email)

	err := s.db.QueryRowx(`SELECT * FROM users WHERE email=$1`, body.Email).StructScan(&user)
	if err != nil {
		handleErr(err)
		res := createAuthErrorResponse(user, "That email is not associated with a user.")
		return c.JSON(http.StatusOK, res)
	}

	//generate timed token
	t, err := generateToken(user.ID, 30*time.Minute, false).SignedString(SigningKey)
	if err != nil {
		handleErr(err)
		res := createAuthErrorResponse(user, errMsgSomethingWrong)
		return c.JSON(http.StatusOK, res)
	}

	//add token to redis
	err = s.storeToken(t, time.Minute*30)
	if err != nil {
		handleErr(err)
		res := createAuthErrorResponse(user, errMsgSomethingWrong)
		return c.JSON(http.StatusOK, res)
	}

	//send reset email
	go s.email.SendResetPasswordEmail(user.Email, t)

	res := createAuthSuccessResponse(user, "Please check your email for a link to reset your password.")
	return c.JSON(http.StatusOK, res)
}

//ResetPassword TODO
func (s *Server) ResetPassword(c echo.Context) error {
	var user Models.User

	body := struct {
		Password string `json:"password"`
	}{}
	if err := c.Bind(&body); err != nil {
		handleErr(err)
		res := createAuthErrorResponse(user, errMsgBadRequest)
		return c.JSON(http.StatusOK, res)

	}

	userID := userID(c)

	password := []byte(body.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		handleErr(err)
		res := createAuthErrorResponse(user, errMsgSomethingWrong)
		return c.JSON(http.StatusOK, res)
	}

	err = s.db.QueryRowx(Models.UpdateUser, userID, hashedPassword).StructScan(&user)
	if err != nil {
		handleErr(err)
		res := createAuthErrorResponse(user, fmt.Sprintf("Unable to Update Password Please report this error"))
		return c.JSON(http.StatusOK, res)
	}
	res := createAuthSuccessResponse(user, "Your Password Was Successfully Changed")
	return c.JSON(http.StatusOK, res)
}

func createAuthErrorResponse(user Models.User, errMsg string) *Models.Res {
	user.Password = ""
	user.ID = ""
	anon := struct {
		User Models.User `json:"user"`
	}{
		user,
	}
	data, _ := json.Marshal(anon)
	res := Models.NewResponse()
	res.Msg = "There was an Error"
	res.Data = data
	res.Error = null.StringFrom(errMsg)
	return res
}

func createAuthSuccessResponse(user Models.User, successMsg string) *Models.Res {
	user.Password = ""
	anon := struct {
		User Models.User `json:"user"`
	}{
		user,
	}
	data, _ := json.Marshal(anon)
	res := Models.NewResponse()
	res.Msg = successMsg
	res.Data = data
	return res
}
