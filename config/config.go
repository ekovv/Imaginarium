package config

import (
	"flag"
	"os"
)

type Config struct {
	Token string
}

type F struct {
	token *string
}

var f F

const addr = "localhost:8080"

func init() {
	f.token = flag.String("t", "a", "-t=token")
}

func New() (c Config) {
	flag.Parse()
	if envRunToken := os.Getenv("TOKEN"); envRunToken != "" {
		f.token = &envRunToken
	}
	c.Token = *f.token
	return c

}
