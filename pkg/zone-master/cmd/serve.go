// Package cmd various CLI commands related to the login service
package cmd

import (
	"fmt"
	"github.com/google/logger"
	"github.com/mitchellh/go-homedir"
	zm "github.com/shine-o/shine.engine.emulator/internal/app/zone-master"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Listen for zone connections to the master",
	Long:  `The purpose of the zone master service is to coordinate the registered zones.`,
	Run:  zm.Start,
}

func init() {
	log = logger.Init("InitLogger", true, false, ioutil.Discard)
	log.Info("root init()")
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/zone-master.yml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(serveCmd)
	log.Info("serve init()")
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
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".master" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName("zone-master")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}