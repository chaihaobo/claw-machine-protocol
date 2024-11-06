package main

import (
	"github.com/chaihaobo/claw-machine-protocol"
)

func main() {
	opt := &redis.Options{
		Network:  "tcp",
		Addr:     "114.55.105.195:6379",
		Password: "FEwG2zx7nPr2BAFC",
		DB:       3,
	}
	client := redis.NewClient(opt)
	server := protocol.NewServer(client)
	server.Start()
}
