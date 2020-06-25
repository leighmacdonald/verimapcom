package cmd

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/leighmacdonald/verimapcom/client"
	"github.com/leighmacdonald/verimapcom/gs"
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
		app := client.New(client.Opts{
			ListenAddr: fmt.Sprintf("%s:%d", host, port),
			RootDir:    dir,
		})
		if err := app.Connect(); err != nil {
			log.Fatalf("Could not connect to server: %s", err)
		}
		if err := app.Ping(); err != nil {
			log.Fatalf("Failed to contact server: %s", err.Error())
		}
		var project gs.Project
		projectConfig := path.Join(dir, "project.yml")
		if gs.Exists(projectConfig) {
			b, err := ioutil.ReadFile(projectConfig)
			if err != nil {
				log.Fatalf("Error reading config file %s: %v", projectConfig, err)
			}
			if err := yaml.Unmarshal(b, &project); err != nil {
				log.Fatalf("Failed to decode config file %s: %v", projectConfig, err)
			}
		} else {
			project.ProjectID = 0
			project.Path = dir
			project.Name = name
		}
		if err := app.OpenProject(&project); err != nil {
			log.Fatalf(err.Error())
		}
		b, err := yaml.Marshal(&project)
		if err != nil {
			log.Fatalf("Could not encode project config: %s", err)
		}
		if err := ioutil.WriteFile(projectConfig, b, 0766); err != nil {
			log.Fatalf("Could not write project config: %s", err)
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
	rootCmd.PersistentFlags().Uint16P("port", "p", 9999, "Listen port")
	rootCmd.PersistentFlags().StringP("dir", "d", "./", "Directory to watch for files")
	rootCmd.PersistentFlags().StringP("name", "n", "project_name", "Name of the project")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vm_uploader.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
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

		// Search config in home directory with name ".vm_uploader" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".vm_uploader")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
