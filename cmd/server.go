package cmd

import (
	"fmt"

	"github.com/namely/broadway/cfg"
	"github.com/namely/broadway/server"
	"github.com/namely/broadway/store"
	"gopkg.in/urfave/cli.v1"
)

// ServerCfg is the configuration object for the server
var ServerCfg cfg.BroadwayServer

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

// ServerCmdFlags declare what flags can be passed to the `server` subcommand
var ServerCmdFlags = []cli.Flag{
	cli.StringFlag{
		Name:        "host, server-host",
		Value:       "0.0.0.0:3000",
		Usage:       "listen address for broadway http api server",
		EnvVar:      "HOST",
		Destination: &ServerCfg.ServerHost,
	},
	cli.StringFlag{
		Name:        "auth, auth-token",
		Usage:       "a global bearer token required for http api requests", // but not GET/POST command/
		EnvVar:      "BROADWAY_AUTH_TOKEN",
		Destination: &ServerCfg.AuthBearerToken,
	},
	// slack-token is sent from Slack to POST command/ and can be found on Slack's
	// Custom Command configuration page.
	// broadway denies the request if it doesn't match this config value
	cli.StringFlag{
		Name:        "slack, slack-token",
		Usage:       "the expected Slack custom command token",
		EnvVar:      "BROADWAY_AUTH_TOKEN",
		Destination: &ServerCfg.AuthBearerToken,
	},
}
