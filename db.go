package main

import (
	"github.com/garyburd/redigo/redis"
)

func RedisConnect() redis.Conn {
	c, err := redis.Dial("tcp", ":6379")
	failOnError(err, "Fail to connect db")
	return c
}
