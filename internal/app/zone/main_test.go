package zone

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/database"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/shn"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/persistence"
	"io/ioutil"
	"os"
	"testing"
)

type testParams struct {
	db * pg.DB
	itemInfo * shn.ShineItemInfo
	itemInfoServer * shn.ShineItemInfoServer
}

var tp = testParams{}

func TestMain(m *testing.M) {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	log = logger.Init("test logger", true, false, ioutil.Discard)
	log.Info("test logger")
	db := database.Connection(ctx, database.ConnectionParams{
		User:     "user",
		Password: "password",
		Host:     "127.0.0.1",
		Port:     "54320",
		Database: "shine",
		Schema:   "world",
	})
	err := database.CreateSchema(db, "world")
	if err != nil {
		log.Fatal(err)
	}
	persistence.CleanDB(db)
	tp.db = db
	os.Exit(m.Run())
}
