package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	Controllers "github.com/chrislewispac/rmd-server/controllers"
	Models "github.com/chrislewispac/rmd-server/models"
	Services "github.com/chrislewispac/rmd-server/services"
	"github.com/fvbock/endless"
	"github.com/jinzhu/configor"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	redis "gopkg.in/redis.v5"
	"os"
	"strings"
)

//Config Global Server Configuration
var Config = struct {
	APPName string `default:"app name"`

	DB struct {
		Name     string
		User     string
		Password string
		Host     string `default:"104.131.111.28"`
	}

	Redis struct {
		Addr string
	}

	Email struct {
		Address  string
		Password string
	}

	Env struct {
		Production  bool
		Development bool
		Port        string
	}
}{}

func main() {
	configor.Load(&Config, "config.yml", "config.production.yml")

	Services.Email.Address = Config.Email.Address
	Services.Email.Password = Config.Email.Password

	var err error
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	/* DATABASE INITIATIALIZATION AND OPENING */
	PgString := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable", Config.DB.Host, Config.DB.User, Config.DB.Name, Config.DB.Password)
	Models.DB, err = sqlx.Open("postgres", PgString)
	if err != nil {
		log.Fatal("Error: Could not establish a connection with the database")
	}
	err = Models.DB.Ping()
	if err != nil {
		log.Fatal("Error: Could not establish a connection with the database")
	}
	Models.DB.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)

	defer Models.DB.Close()

	Models.Client = redis.NewClient(&redis.Options{
		Addr:     Config.Redis.Addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	/* END DATABASE INITIATIALIZATION AND OPENING */

	/* BEGIN ROUTE HANDLERS */
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.POST("/login", Controllers.LoginCtrl)
	e.POST("/register", Controllers.RegisterCtrl)

	e.Static("/", "static")

	// Token restricted routes
	r := e.Group("/auth")

	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("secret"),
	}))

	//FACILITY ROUTES
	r.GET("/facilities", Controllers.GetUserFacilitiesCtrl)
	r.GET("/facility/:id", Controllers.GetFacilityByIDCtrl)
	r.DELETE("/facility/:id", Controllers.DeleteFacilityByIDCtrl)
	r.POST("/facility/update/:id", Controllers.UpdateFacilityByIDCtrl)
	r.POST("/facility", Controllers.PostFacilityCtrl)
	r.POST("/facilities/csv", Controllers.PostFacilitiesFromCsvCtrl)

	//CONTRACT ROUTES
	r.GET("/contracts", Controllers.GetUserContractsCtrl)
	r.GET("/contract/:id", Controllers.GetContractByIDCtrl)
	r.POST("/contract/update/:id", Controllers.UpdateContractByIDCtrl)
	r.POST("/contract", Controllers.PostContractCtrl)
	r.POST("/contracts/csv", Controllers.PostContractsFromCsvCtrl)

	//PERSONS ROUTES
	r.GET("/persons", Controllers.GetUserPersonsCtrl)
	r.GET("/person/:id", Controllers.GetPersonByIDCtrl)
	r.POST("/person/update/:id", Controllers.UpdatePersonByIDCtrl)
	r.POST("/person", Controllers.PostPersonCtrl)
	r.POST("/persons/csv", Controllers.PostPersonsFromCsvCtrl)

	//USER ROUTES
	r.GET("/users/:id", Controllers.GetUserByIDCtrl)
	r.POST("/logout", Controllers.LogoutCtrl)
	r.POST("/password_reset", Controllers.ForgotPasswordCtrl)

	/* END ROUTE HANDLERS */

	if Config.Env.Production {
		log.Println("Starting in production mode...")
		serverErr := endless.ListenAndServe("localhost:8000", e)
		if serverErr != nil {
			log.Println(serverErr)
		}
		log.Println("Server on 8000 stopped")

		os.Exit(0)
	} else {
		log.Println("Starting in development mode...")
		e.Logger.Fatal(e.Start(Config.Env.Port))
	}

}
