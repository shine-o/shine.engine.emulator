package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/go-pg/migrations/v8"
	"github.com/google/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/shine-o/shine.engine.emulator/internal/pkg/database"
)

var (
	cfgFile string
	log     *logger.Logger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "migrations",
	Short: "apply migrations to the database",
	Long: `This program runs command on the db. Supported commands are:
  - init - creates version info table in the database
  - up - runs all available migrations.
  - up [target] - runs available migrations up to the target one.
  - down - reverts last migration.
  - reset - reverts all migrations.
  - version - prints current db version.
  - set_version [version] - sets db version without running migrations.
`,
	Run:   Migrate,
}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Migrate(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	db := database.Connection(ctx, database.ConnectionParams{
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
		Host:     viper.GetString("database.host"),
		Port:     viper.GetString("database.port"),
		Database: viper.GetString("database.db_name"),
	})

	oldVersion, newVersion, err := migrations.Run(db, args...)
	if err != nil {
		exitf(err.Error())
	}
	if newVersion != oldVersion {
		fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		fmt.Printf("version is %d\n", oldVersion)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/migrations.yml)")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigFile(cfgFile)

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func errorf(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s+"\n", args...)
}

func exitf(s string, args ...interface{}) {
	errorf(s, args...)
	os.Exit(1)
}
