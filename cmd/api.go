package cmd

import (
	"backend/internal/infra/api/router"
	"backend/wire"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "API",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		application := wire.Init(appConfig)
		router.StartServer(application)
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
