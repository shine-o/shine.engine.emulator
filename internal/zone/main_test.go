package zone

import (
	"context"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/database"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/persistence"
	"io/ioutil"
	"os"
	"sync"
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

	//loadGameData("../../../files")
	loadTestData("../../files")

	os.Exit(m.Run())
}

func loadTestData(filesPath string)  {
	var (
		wg sync.WaitGroup
	)

	wg.Add(1)
	//go func(path string) {
	//	defer wg.Done()
	//	md, err := data.LoadMonsterData(path)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	monsterData = md
	//}(filesPath)
	//
	//go func(path string) {
	//	defer wg.Done()
	//	md, err := data.LoadMapData(path)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	mapData = md
	//}(filesPath)
	//

	//go func(path string) {
	//	defer wg.Done()
	//	nd, err := data.LoadNPCData(path)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	npcData = nd
	//}(filesPath)

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