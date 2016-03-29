package main

import (
	"os"

	"github.com/namely/broadway/server"
	"github.com/namely/broadway/store"
)

func main() {
	s := server.New(store.New())
	s.Init()
	err := s.Run(os.Getenv("HOST"))
	if err != nil {
		panic(err)
	}

}
