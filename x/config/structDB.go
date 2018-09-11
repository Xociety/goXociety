package config

import (
	"database/sql"

	"github.com/globalsign/mgo"
	"github.com/go-redis/redis"
	// _ "github.com/lib/pq"
)

// net
type ConnPostgres struct {
	DB *sql.DB
}

type ConnMongo struct {
	Session *mgo.Session
}
type ConnRedis struct {
	Client *redis.Client
}
type UserDB struct {
	UserID      int64
	Username    string
	Email       string
	Password    string
	Name        string
	Phone       string
	Gender      int
	Bio         string
	Credit      int
	PhotoURL    string
	LanguageID  int
	CountryCode string
	Timezone    int
	LastIP      string
	Updatetime  int
	Createtime  int
}
