package zone

import (
	"context"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/database"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/persistence"
	"io/ioutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	log = logger.Init("test logger", true, false, ioutil.Discard)
	log.Info("test logger")

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

	loadGameData("../../../files")

	os.Exit(m.Run())
}
