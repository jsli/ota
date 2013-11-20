package models

import (
	"github.com/gosexy/redis"
)

const (
	HOST = "localhost"
	PORT = 6379
)

type RedisDal struct {
	*redis.Client
}

func NewRedisDal() (*RedisDal, error) {
	client := redis.New()
	err := client.Connect(HOST, PORT)
//	client.HS
	if err != nil {
		return nil, err
	}

	return &RedisDal{client}, nil
}

func (rdal *RedisDal) Close() {
	rdal.Quit()
}
