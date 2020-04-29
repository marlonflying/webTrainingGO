package main

import (
	redisStore "gopkg.in/boj/redistore.v1"
)

const (
	connHost = "localhost"
	connPort = "8080"
)

var store *redisStore.RediStore
var err error
