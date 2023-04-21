package src

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"runtime"
)

type Error struct {
	err  string
	file string
	line int
	ok   bool
	pc   uintptr
}

type Redis struct {
	*redis.Client
	Address string
	Port    int
}

func (r Redis) New() Redis {
	r.Client = redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%d", r.Address, r.Port)})
	return r
}

func (r Redis) GetSecret(key string) string {
	val, _ := r.Client.Get(context.Background(), key).Result()
	decoded, _ := base64.StdEncoding.DecodeString(val)
	return string(decoded)
}

func (r Redis) PublishError(err error) {
	pc, file, line, ok := runtime.Caller(2)
	e, _ := json.Marshal(Error{pc: pc, file: file, line: line, ok: ok, err: err.Error()})
	log.Println(Error{pc: pc, file: file, line: line, ok: ok, err: err.Error()})
	r.Client.Publish(context.Background(), "error", e)
}
