package cmd

import (
	"context"
	"flag"
	"github.com/leighmacdonald/verimapcom/web"
	"github.com/leighmacdonald/verimapcom/web/store"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/cobra"
)

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
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

		flag.Parse()
		if flagMigrate {
			if err := store.Migrate(os.Getenv("DATABASE_URL")); err != nil {
				if err.Error() != "no change" {
					log.Fatalf("Could not do migrations: %v", err)
				}
				log.Infof("No migration performed")
			}
		}

		w := web.New(ctx)
		if err := w.Setup(); err != nil {
			log.Fatalf("Could not run setup: %v", err)
		}
		defer w.Close()

		opts := web.DefaultHTTPOpts()
		opts.Handler = w.Handler

		srv := web.NewHTTPServer(opts)

		// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
		if err := srv.ListenAndServe(); err != nil {
			log.Errorf("Shutdown unclean: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(webCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// webCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// webCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
