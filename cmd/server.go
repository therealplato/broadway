package cmd

import (
	"fmt"

	"github.com/namely/broadway/server"
	"github.com/namely/broadway/store"
	"gopkg.in/urfave/cli.v1"
)

// ServerCmd is executed by cli on `broadway server`
var ServerCmd = func(c *cli.Context) error {
	fmt.Println("starting server...")
	s := server.New(store.New(), cfg.CommonCfg, cfg.ServerCfg)
	s.Init()
	// err := s.Run(env.ServerHost)
	err := s.Run(cfg.ServerCfg.ServerHost)
	if err != nil {
		panic(err)
	}
	return nil
}

// ServerCmdFlags declares what flags can be passed to the `server` subcommand
var ServerCmdFlags = []cli.Flag{
	cli.StringFlag{
		Name:        "host, server-host",
		Value:       "0.0.0.0:3000",
		Usage:       "listen address for broadway http api server",
		EnvVar:      "HOST",
		Destination: &cfg.ServerCfg.ServerHost,
	},
	cli.StringFlag{
		Name:        "auth-token",
		Usage:       "a global bearer token required for http api requests", // but not GET/POST command/
		EnvVar:      "BROADWAY_AUTH_TOKEN",
		Destination: &cfg.ServerCfg.AuthBearerToken,
	},
	// slack-token is sent from Slack to POST command/ and can be found on Slack's Custom Command configuration page.
	// broadway denies the request if the received token doesn't match this config value
	cli.StringFlag{
		Name:        "slack-token",
		Usage:       "the expected Slack custom command token",
		EnvVar:      "SLACK_VERIFICATION_TOKEN",
		Destination: &cfg.ServerCfg.SlackToken,
	},
	cli.StringFlag{
		Name:        "slack-webhook",
		Usage:       "slack.com webhook URL where broadway sends notifications",
		EnvVar:      "SLACK_WEBHOOK",
		Destination: &cfg.ServerCfg.SlackWebhook,
	},
	cli.StringFlag{
		Name:        "manifest-dir",
		Usage:       "path to a folder containing broadway manifests",
		Value:       "./manifests",
		EnvVar:      "BROADWAY_MANIFESTS_PATH",
		Destination: &cfg.ServerCfg.ManifestsPath,
	},
	cli.StringFlag{
		Name:        "playbook-dir",
		Usage:       "path to a folder containing broadway playbooks",
		Value:       "./playbooks",
		EnvVar:      "BROADWAY_PLAYBOOKS_PATH",
		Destination: &cfg.ServerCfg.PlaybooksPath,
	},
}
