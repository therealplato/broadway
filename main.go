package main

import (
	"fmt"
	"os"

	"github.com/namely/broadway/cmd"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "broadway"
	app.Usage = "deploy distributed container systems"
	app.Action = func(c *cli.Context) error {
		fmt.Println(`Use "broadway server"`)
		return nil
	}
	app.Flags = cmd.CommonFlags
	app.Commands = []cli.Command{
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "start broadway's HTTP API server",
			Action:  cmd.ServerCmd,
			Flags:   cmd.ServerCmdFlags,
		},
	}
	app.Run(os.Args)
}
