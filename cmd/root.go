package cmd

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/leighmacdonald/verimapcom/client"
	"github.com/leighmacdonald/verimapcom/core"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path"
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
		host, err := cmd.Flags().GetString("host")
		if err != nil {
			log.Fatal("Invalid host value")
		}
		port, err2 := cmd.Flags().GetUint16("port")
		if err2 != nil {
			log.Fatal("Invalid port value")
		}
		dir, err3 := cmd.Flags().GetString("dir")
		if err3 != nil {
			log.Fatal("Invalid dir specified, doesn't exist")
		}
		name, err4 := cmd.Flags().GetString("name")
		if err4 != nil {
			log.Fatal("Invalid name specified, doesn't exist")
		}
		caCertPath := viper.GetString("ssl_ca")
		app := client.New(client.Opts{
			ListenAddr: fmt.Sprintf("%s:%d", host, port),
			RootDir:    dir,
			CaCert:     caCertPath,
		})
		if err := app.Connect(); err != nil {
			log.Fatalf("Could not connect to server: %s", err)
		}
		var missionConfig client.MissionConfig
		projectConfig := path.Join(dir, "missionConfig.yaml")
		if core.Exists(projectConfig) {
			b, err := ioutil.ReadFile(projectConfig)
			if err != nil {
				log.Fatalf("Error reading config file %s: %v", projectConfig, err)
			}
			if err := yaml.Unmarshal(b, &missionConfig); err != nil {
				log.Fatalf("Failed to decode config file %s: %v", projectConfig, err)
			}
		} else {
			missionConfig.MissionID = 0
			missionConfig.Name = name
		}
		if missionConfig.MissionID <= 0 {
			mid, err := app.CreateMission(missionConfig.Name)
			if err != nil {
				if err == core.ErrDuplicate {
					log.Fatalf("Duplicate mission name")
				}
				log.Fatalf("Failed to create mission: %v", err)
			}
			missionConfig.MissionID = mid
		}
		if err := app.OpenMission(missionConfig.MissionID); err != nil {
			log.Fatalf("Failed to open mission: %v", err)
		}
		b, err2 := yaml.Marshal(&missionConfig)
		if err2 != nil {
			log.Fatalf("Could not encode missionConfig config: %s", err2)
		}
		if err := ioutil.WriteFile(projectConfig, b, 0766); err != nil {
			log.Fatalf("Could not write missionConfig config: %s", err)
		}
		if err := app.Start(); err != nil {
			log.Errorf("Failed to shutdown cleanly: %v", err)
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
