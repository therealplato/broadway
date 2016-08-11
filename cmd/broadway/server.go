package broadway

import (
	"fmt"
	"os"

	"gopkg.in/urfave/cli.v1"

	"github.com/namely/broadway/pkg/cfg"
	"github.com/namely/broadway/pkg/deployment"
	"github.com/namely/broadway/pkg/server"
	"github.com/namely/broadway/pkg/store/etcdstore"
)

// ServerCmd is executed by cli on `broadway server`
var ServerCmd = func(c *cli.Context) error {
	etcdstore.Setup(cfg.GlobalCfg)  // configure etcd before using
	deployment.Setup(cfg.GlobalCfg) // configure kubernetes deployments before using
	fmt.Printf("starting server with config...\n%+v", cfg.GlobalCfg)
	s := server.New(cfg.GlobalCfg, etcdstore.New())
	s.Init()
	if err := s.Run(cfg.GlobalCfg.ServerHost); err != nil {
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
		Destination: &cfg.GlobalCfg.ServerHost,
	},
	cli.StringFlag{
		Name:        "auth-token",
		Usage:       "a global bearer token required for http api requests", // but not GET/POST command/
		EnvVar:      "BROADWAY_AUTH_TOKEN",
		Destination: &cfg.GlobalCfg.AuthBearerToken,
	},
	// slack-token is sent from Slack to POST command/ and can be found on Slack's Custom Command configuration page.
	// broadway denies the request if the received token doesn't match this config value
	cli.StringFlag{
		Name:        "slack-token",
		Usage:       "the expected Slack custom command token",
		EnvVar:      "SLACK_VERIFICATION_TOKEN",
		Destination: &cfg.GlobalCfg.SlackToken,
	},
	cli.StringFlag{
		Name:        "slack-webhook",
		Usage:       "slack.com webhook URL where broadway sends notifications",
		EnvVar:      "SLACK_WEBHOOK",
		Destination: &cfg.GlobalCfg.SlackWebhook,
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "broadway"
	app.Usage = "deploy distributed container systems"
	app.Action = func(c *cli.Context) error {
		fmt.Println(`Use "broadway server"`)
		return nil
	}
	app.Flags = CommonFlags
	app.Commands = []cli.Command{
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "start broadway's HTTP API server",
			Action:  ServerCmd,
			Flags:   ServerCmdFlags,
		},
	}
	app.Run(os.Args)
}
