package library

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"runtime"

	"github.com/redis/go-redis/v9"
)

type Error struct {
	Err   string
	Trace string
}
type Redis struct {
	*redis.Client
	Address string
	Port    string
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
	var b = make([]byte, 512)
	runtime.Stack(b, false)
	var e = Error{Err: err.Error(), Trace: string(b)}
	log.Println(e)
	r.Client.Publish(context.Background(), "error", e)
}

func (r Redis) CheckError(f func() error) {
	err := f()
	if err != nil {
		r.PublishError(err)
	}
}
