package main

import (
	"github.com/namely/broadway/env"
	"github.com/namely/broadway/server"
	"github.com/namely/broadway/store/etcdstore"
)

func main() {
	s := server.New(etcdstore.New())
	s.Init()
	err := s.Run(env.ServerHost)
	if err != nil {
		panic(err)
	}

}
