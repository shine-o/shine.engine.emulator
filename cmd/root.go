/*
Copyright © 2020 Marius Bora

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/google/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"io/ioutil"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var log *logger.Logger

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "login",
	Short: "Serve the login services needed and offer CLI support for account operations",
	Long:  `The purpose of the login service is to handle packets related to user account login and server selection.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	log = logger.Init("InitLogger", true, true, ioutil.Discard)
	log.Info("root init()")
	cobra.OnInitialize(initConfig)
	err := doc.GenMarkdownTree(rootCmd, "docs")
	if err != nil {
		log.Fatal(err)
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.login.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

		// Search config in home directory with name ".shine.engine.login" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".login")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//fmt.Println("Using config file:", viper.ConfigFileUsed())
		log.Infof("Using config file: %v", viper.ConfigFileUsed())
	}

	viper.SetDefault("serve.port", 9010)
	viper.SetDefault("crypt.xorKey", "0759694a941194858c8805cba09ecd583a365b1a6a16febddf9402f82196c8e99ef7bfbdcfcdb27a009f4022fc11f90c2e12fba7740a7d78401e2ca02d06cba8b97eefde49ea4e13161680f43dc29ad486d7942417f4d665bd3fdbe4e10f50f6ec7a9a0c273d2466d322689c9a520be0f9a50b25da80490dfd3e77d156a8b7f40f9be80f5247f56f832022db0f0bb14385c1cba40b0219dff08becdb6c6d66ad45be89147e2f8910b89360d860def6fe6e9bca06c1759533cfc0b2e0cca5ce12f6e5b5b426c5b2184f2a5d261b654df545c98414dc7c124b189cc724e73c64ffd63a2cee8c8149396cb7dcbd94e232f7dd0afc020164ec4c940ab156f5c9a934de0f3827bc81300f7b3825fee83e29ba5543bf6b9f1f8a4952187f8af888245c4fe1a830878e501f2fd10cb4fd0abcdc1285e252ee4a5838abffc63db960640ab450d54089179ad585cfec0d7e817fe3c3040122ec27ccfa3e21a654c8de00b6df279ff625340785bfa7a5a5e0830c3d5d2040af60a36456f305c41c7d3798c3e85a6e5885a49a6b6af4a37b619b09401e604b32d951a4fef95d4e4afb4ad47c330233d59dce5baa5a7cd8f805fa1f2b8c725750ae6c1989ca01fcfc299b61126863654626c45b50aa2bbeef9a790223752c2013fdd95a7623f10bb5b859f99f7ae606e9a53ab450bf165898b39a6e36ee8deb")
	viper.SetDefault("crypt.xorLimit", 350)
	viper.SetDefault("protocol.nc-data", "defaults/protocol-commands.yml")

	requiredParams := []string{
		"database.postgres.host",
		"database.postgres.port",
		"database.postgres.db_user",
		"database.postgres.db_password",
	}
	for _, rp := range requiredParams {
		if !viper.IsSet(rp) {
			log.Fatalf("missing required parameter %v", rp)
		}
	}
}
