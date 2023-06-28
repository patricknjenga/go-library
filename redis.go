package library

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	*redis.Client `gorm:"-"`
	Address       string
	Port          string
}

func (r Redis) New() Redis {
	r.Client = redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%s", r.Address, r.Port)})
	return r
}

func (r Redis) GetSecret(key string) string {
	val, _ := r.Client.Get(context.Background(), key).Result()
	decoded, _ := base64.StdEncoding.DecodeString(val)
	return string(decoded)
}

func (r Redis) PublishError(err error) {
	log.Println(err)
	r.Client.Publish(context.Background(), "error", err)
}

func (r Redis) CheckError(f func() error) {
	err := f()
	if err != nil {
		r.PublishError(err)
	}
}
