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
	app.Usage = "deploy container topographies"
	app.Commands = []cli.Command{
		{
			Name:    "server",
			Aliases: []string{"a"},
			Usage:   "add a task to the list",
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
