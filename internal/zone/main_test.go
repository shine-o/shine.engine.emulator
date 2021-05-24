package zone

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/google/logger"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/database"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/persistence"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	log = logger.Init("test logger", true, false, ioutil.Discard)

	persistence.InitDB(database.ConnectionParams{
		User:     "user",
		Password: "password",
		Host:     "127.0.0.1",
		Port:     "54320",
		Database: "shine",
		Schema:   "world",
	})

	err := database.CreateSchema(persistence.DB(), "world")
	if err != nil {
		log.Fatal(err)
	}

	persistence.CleanDB()

	initConfig()

	filesPath, err := filepath.Abs("../../files")
	if err != nil {
		log.Fatal(err)
	}

	loadTestData(filesPath)

	os.Exit(m.Run())
}

func loadTestData(filesPath string) {
	var wg sync.WaitGroup

	wg.Add(4)

	go func(path string) {
		defer wg.Done()
		md, err := data.LoadMonsterData(path)
		if err != nil {
			log.Fatal(err)
		}
		monsterData = md
	}(filesPath)

	go func(path string) {
		defer wg.Done()
		md, err := data.LoadMapData(path)
		if err != nil {
			log.Fatal(err)
		}
		mapData = md
	}(filesPath)

	go func(path string) {
		defer wg.Done()
		nd, err := data.LoadNPCData(path)
		if err != nil {
			log.Fatal(err)
		}
		npcData = nd
	}(filesPath)

	go func(path string) {
		defer wg.Done()
		id, err := data.LoadItemData(path)
		if err != nil {
			log.Fatal(err)
		}
		itemsData = id
	}(filesPath)

	wg.Wait()
}

func initConfig() {
	// Search config in home directory with name ".zone" (without extension).
	viper.AddConfigPath("../../configs")
	viper.SetConfigName("zone")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Using config file:", viper.ConfigFileUsed())
}
