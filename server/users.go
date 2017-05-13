package server

import (
	"net/http"

	"github.com/labstack/echo"
)

//GetUserByIDCtrl ...
func (s *Server) GetUserByToken(c echo.Context) error {
	return c.JSON(http.StatusOK, userID(c))
}
