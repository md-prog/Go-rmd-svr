package server

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/redis.v5"

	"github.com/dgrijalva/jwt-go"
	"gitlab.com/chrislewispac/rmd-server/services"
	"net/http"
	"github.com/ginuerzh/weedo"
)

var (
	SigningKey = []byte("secret")
)

// Server Configuration
type Config struct {
	APPName string `default:"app name"`

	DB struct {
		Name     string
		User     string
		Password string
		Host     string
	}

	Redis struct {
		Addr     string
		Password string
	}

	Sparkpost struct {
		Key string
	}

	Intercom struct {
		Key string
	}

	BaseUrl string

	Env struct {
		Production  bool
		Development bool
		Staging     bool
		Port        string
	}

	Seaweedfs struct {
		MasterUrl   string
  	}
}

type Server struct {
	db       *sqlx.DB
	redis    *redis.Client
	email    *services.EmailSettings
	intercom *services.IntercomService
	fs	 *weedo.Client

	JWT echo.MiddlewareFunc
}

func NewServer(c *Config) (*Server, error) {
	// Connect and ping redis.
	rd := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       0,
	})

	// Connect and ping postgres.
	pgString := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable", c.DB.Host, c.DB.User, c.DB.Name, c.DB.Password)
	db, err := sqlx.Open("postgres", pgString)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)

	// Connect to Seaweedfs master server
	fs := weedo.NewClient(c.Seaweedfs.MasterUrl)

	return &Server{
		db:       db,
		redis:    rd,
		email:    services.InitEmailService(c.BaseUrl, c.Sparkpost.Key),
		intercom: services.InitIntercomService(c.Intercom.Key),
		fs:	  fs,

		JWT: middleware.JWTWithConfig(middleware.JWTConfig{
			SigningKey: SigningKey,
		}),
	}, nil
}

func (s *Server) Close() {
	s.db.Close()
	s.redis.Close()
}

// ValidateSession is a session token validation middleware intended to be
// called after middleware.JWT.
//
// For stored tokens (live session), it calls the next handler.
// For missing, expired, or otherwise not-stored tokens it return "401 - Unauthorized" error.
func (s *Server) ValidateSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get(middleware.DefaultJWTConfig.ContextKey)
		if token == nil {
			return echo.ErrUnauthorized
		}
		if stored, err := s.isTokenStored(token.(*jwt.Token).Raw); err != nil {
			handleErr(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		} else if !stored {
			return echo.ErrUnauthorized
		}

		return next(c)
	}
}

// tx executes a transaction function and commits when successful, or rolls back
// when an error *or* a panic occurs. Original recovered errors are preserved.
func (s *Server) tx(f func(tx *sqlx.Tx) error) (err error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err2 := recover(); err2 != nil {
			_ = tx.Rollback()
			panic(err2)
		}
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()

	}()
	return f(tx)
}
