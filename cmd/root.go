package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "broadway",
	Short: "Controls Broadway",
	Long:  "Controls Broadway",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`Use "broadway server"`)
		os.Exit(-1)
	},
}

// func init() {
// RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
// RootCmd.PersistentFlags().StringVarP(&projectBase, "projectbase", "b", "", "base project directory eg. github.com/spf13/")
// RootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "Author name for copyright attribution")
// RootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "Name of license for the project (can provide `licensetext` in config)")
// RootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
// viper.BindPFlag("author", RootCmd.PersistentFlags().Lookup("author"))
// viper.BindPFlag("projectbase", RootCmd.PersistentFlags().Lookup("projectbase"))
// viper.BindPFlag("useViper", RootCmd.PersistentFlags().Lookup("viper"))
// viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
// viper.SetDefault("license", "apache")
// }
