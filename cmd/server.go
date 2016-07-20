package cmd

import (
	"github.com/namely/broadway/env"
	"github.com/namely/broadway/server"
	"github.com/namely/broadway/store"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts a Broadway server",
	Long:  "Starts a Broadway server",
	Run: func(cmd *cobra.Command, args []string) {
		s := server.New(store.New())
		s.Init()
		err := s.Run(env.ServerHost)
		if err != nil {
			panic(err)
		}
	},
}
