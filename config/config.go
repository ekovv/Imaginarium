package config

import (
	"flag"
	"os"
)

type Config struct {
	Token  string
	DB     string
	TypeDB string
}

type F struct {
	token  *string
	db     *string
	typeDB *string
}

var f F

const addr = "localhost:8080"

func init() {
	f.token = flag.String("t", "a", "-t=token")
	f.db = flag.String("d", "", "-d=db")
	f.typeDB = flag.String("type", "postgresql", "-type=")
}

func New() (c Config) {
	flag.Parse()
	if envRunToken := os.Getenv("TOKEN"); envRunToken != "" {
		f.token = &envRunToken
	}
	if envDB := os.Getenv("DATABASE_DSN"); envDB != "" {
		f.db = &envDB
	}
	if envTypeDB := os.Getenv("DATABASE_DSN"); envTypeDB != "" {
		f.db = &envTypeDB
	}
	c.Token = *f.token
	c.DB = *f.db
	c.TypeDB = *f.typeDB
	return c

}
