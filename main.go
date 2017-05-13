package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/fvbock/endless"
	"github.com/jinzhu/configor"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"

	"gitlab.com/chrislewispac/rmd-server/server"
)

type CustomBinder struct{}

func (cb *CustomBinder) Bind(i interface{}, c echo.Context) (err error) {
	// You may use default binder
	db := new(echo.DefaultBinder)
	if err = db.Bind(i, c); err != echo.ErrUnsupportedMediaType {
		return
	}

	return
}

func main() {
	Config := server.Config{}
	if err := configor.Load(&Config, "config.yml"); err != nil {
		log.Fatal("Error loading config file:", err)
	}

	s, err := server.NewServer(&Config)
	if err != nil {
		log.Fatal("failed to create server:", err)
	}

	defer s.Close()

	e := echo.New()
	e.Binder = &CustomBinder{}
	s.InitRoutes(e)

	if Config.Env.Production {
		log.Println("Starting in production mode...")
		serverErr := endless.ListenAndServe("localhost:8000", e)
		if serverErr != nil {
			log.Println(serverErr)
		}
		log.Println("Server on 8000 stopped")

		os.Exit(0)
	} else if Config.Env.Staging {
		log.Println("Starting in staging mode...")
		serverErr := endless.ListenAndServe("localhost:3000", e)
		if serverErr != nil {
			log.Println(serverErr)
		}
		log.Println("Server on 3000 stopped")

		os.Exit(0)
	} else {
		log.Println("Starting in development mode...")
		e.Logger.Fatal(e.Start(Config.Env.Port))
	}

}
