package main

import (
	"github.com/raja-dettex/go-cache/cache"
	"github.com/raja-dettex/go-cache/server"
)

func main() {
	opts := server.ServerOpts{
		ListenToAddr: ":4000",
		IsLeader:     true,
	}
	cache := cache.New()
	server := server.NewServer(opts, cache)
	server.Start()

}
