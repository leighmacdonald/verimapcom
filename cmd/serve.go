package cmd

import (
	"fmt"
	"github.com/leighmacdonald/verimapcom/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net"
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
		host, err := cmd.Flags().GetString("host")
		if err != nil {
			log.Fatal("Invalid host value")
		}
		port, err2 := cmd.Flags().GetUint16("port")
		if err2 != nil {
			log.Fatal("Invalid port value")
		}
		listenAddr := fmt.Sprintf("%s:%d", host, port)
		lis, err3 := net.Listen("tcp", listenAddr)
		if err3 != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := server.NewServer(server.Opts{
			Tls: false,
		})
		log.Infof("Listening on %s", listenAddr)
		if err := s.Serve(lis); err != nil {
			log.Errorf("Failed to serve: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
