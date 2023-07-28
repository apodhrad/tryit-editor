/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/apodhrad/tryit-editor/server"
	"github.com/apodhrad/tryit-editor/service"
	"github.com/spf13/cobra"
)

var configFile string

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Tryit Editor",
	Long:  "Start the Tryit Editor",
	RunE: func(cmd *cobra.Command, args []string) error {
		var services []service.Service
		var err error

		if configFile == "" {
			services = []service.Service{service.BUILTIN_SERVICE_HTML, service.BUILTIN_SERVICE_MARKDOWN}
		} else {
			services, err = service.LoadServices(configFile)
			if err != nil {
				return err
			}
		}

		ctx, err := server.Start(services)
		if err != nil {
			return err
		}

		<-ctx.Done()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	startCmd.Flags().StringVarP(&configFile, "config", "c", "", "yaml file with services")
}
