package Models

import (
	"github.com/jmoiron/sqlx"
	redis "gopkg.in/redis.v5"
)

//DB databse connection PostgreSQL
var DB *sqlx.DB

//Client redis client for session management
var Client *redis.Client
