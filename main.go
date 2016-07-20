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
	app.Usage = "Deploy distributed container systems"
	app.Commands = []cli.Command{
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "Start Broadway's HTTP API server",
			Action:  cmd.ServerCmd,
		},
	}
	app.Action = func(c *cli.Context) error {
		fmt.Println(`Use "broadway server"`)
		return nil
	}

	app.Run(os.Args)

	// if err := cmd.RootCmd.Execute(); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(-1)
	// }
}
