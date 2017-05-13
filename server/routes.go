package server

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func (s *Server) InitRoutes(e *echo.Echo) {
	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.POST("/login", s.Login)
	e.POST("/register", s.Register)
	e.POST("/forgot_password", s.ForgotPassword)
	e.GET("/testing", s.GetFacilityByID)

	e.POST("/email_forwarding", s.ReceiveEmailCV)

	e.Static("/", "static")

	// Token restricted routes
	r := e.Group("/auth")

	r.Use(s.JWT, s.ValidateSession)

	//FACILITY ROUTES
	r.GET("/facilities", s.GetUserFacilities)
	r.GET("/facility/:id", s.GetFacilityByID)
	r.DELETE("/facility/:id", s.DeleteFacilityByID)
	r.POST("/facility/update/:id", s.UpdateFacilityByID)
	r.POST("/facility", s.PostFacility)
	r.POST("/facilities/csv", s.PostFacilitiesFromCsv)

	//CONTRACT ROUTES
	r.GET("/contracts", s.GetUserContracts)
	r.GET("/contract/:id", s.GetContractByID)
	r.DELETE("/contract/:id", s.DeleteContractByID)
	r.POST("/contract/update/:id", s.UpdateContractByID)
	r.POST("/contract", s.PostContract)
	r.POST("/contracts/csv", s.PostContractsFromCsv)

	//PERSONS ROUTES
	r.GET("/persons", s.GetUserPersons)
	r.GET("/person/:id", s.GetPersonByID)
	r.DELETE("/person/:id", s.DeletePersonByID)
	r.POST("/person/update/:id", s.UpdatePersonByID)
	r.POST("/person", s.PostPerson)
	r.POST("/persons/csv", s.PostPersonsFromCsv)

	//USER ROUTES
	r.POST("/logout", s.Logout)
	r.GET("/users/:token", s.GetUserByToken)
	r.POST("/reset_password", s.ResetPassword)

	//PROVIDER ROUTES
	r.GET("/providers", s.GetUserProviders)
	r.GET("/provider/:id", s.GetProviderByID)
	r.DELETE("/provider/:id", s.DeleteProviderByID)
	r.POST("/provider/update/:id", s.UpdateProviderByID)
	r.POST("/provider", s.PostProvider)
	r.POST("/providers/csv", s.PostProvidersFromCsv)

	//FILES ROUTES

	r.DELETE("/file/:fid", s.DeleteFile)
	r.POST("/file", s.PostFile)
	r.GET("/files", s.GetUserFileList)
	r.GET("/file/:fid", s.GetFile)

	return
}
