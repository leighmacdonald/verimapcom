package cmd

import (
	"context"
	"github.com/leighmacdonald/verimapcom/core"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFormatter(&log.TextFormatter{ForceColors: true})
		log.SetLevel(log.DebugLevel)
		ctx := context.Background()
		app, err := core.New(ctx)
		if err != nil {
			log.Fatalf("Could not start service: %v", err)
		}
		if err := app.ListenAndServe(); err != nil {
			log.Errorf("Failed to cleanly shutdown: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
