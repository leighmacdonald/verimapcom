package cmd

import (
	"context"
	"github.com/leighmacdonald/verimapcom/core"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "geosync",
	Short: "geosync",
	Long:  `geosync`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("Shutdown uncleanly: %s", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringP("host", "H", "localhost", "Listen host")
	rootCmd.PersistentFlags().Uint16P("port", "p", 9090, "Listen port")
	rootCmd.PersistentFlags().StringP("dir", "d", "./", "Directory to watch for files")
	rootCmd.PersistentFlags().StringP("name", "n", "project_name", "Name of the project")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vm_uploader.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetDefault("migrate", true)
	viper.SetDefault("listen_http", ":8001")
	viper.SetDefault("redis", "localhost:6379")
	viper.SetDefault("dsn", "postgres:///verimapcom")
	viper.SetDefault("cms_host_internal", "https://cms.verimap.com")
	viper.SetDefault("static_path", "dist")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".config" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath("./")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Infof("Using config file: %s", viper.ConfigFileUsed())
	}
}
