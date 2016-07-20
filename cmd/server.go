package cmd

import (
	"fmt"

	"github.com/namely/broadway/env"
	"github.com/namely/broadway/server"
	"github.com/namely/broadway/store"
	"gopkg.in/urfave/cli.v1"
)

var ServerCmd = func(c *cli.Context) error {
	fmt.Println("starting server...")
	s := server.New(store.New())
	s.Init()
	err := s.Run(env.ServerHost)
	if err != nil {
		panic(err)
	}
	return nil
}
