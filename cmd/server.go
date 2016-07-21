package cmd

import (
	"fmt"

	"github.com/namely/broadway/cfg"
	"github.com/namely/broadway/server"
	"github.com/namely/broadway/store"
	"gopkg.in/urfave/cli.v1"
)

// ServerCmd is executed by cli on `broadway server`
var ServerCmd = func(c *cli.Context) error {
	fmt.Println("starting server...")
	s := server.New(store.New())
	s.Init()
	// err := s.Run(env.ServerHost)
	err := s.Run(ServerCfg.ServerHost)
	if err != nil {
		panic(err)
	}
	return nil
}

// ServerCfg is the configuration object for the server
var ServerCfg cfg.BroadwayServer

// ServerCmdFlags declare what flags can be passed to the `server` subcommand
var ServerCmdFlags = []cli.Flag{
	cli.StringFlag{
		Name:        "host, server-host",
		Value:       "0.0.0.0:3000",
		Usage:       "listen address for broadway http api server",
		EnvVar:      "HOST",
		Destination: &ServerCfg.ServerHost,
	},
}
