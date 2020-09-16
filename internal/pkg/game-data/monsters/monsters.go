package monsters

import (
	"encoding/csv"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/shn"
	bolt "go.etcd.io/bbolt"
	"os"
)

type MonsterData struct {
	Regens        MonsterRegenTable
	MobInfo       shn.MobInfo
	MobInfoServer shn.MobInfoServer
}

type MonsterRegenTable struct {
	groups map[string]RegenEntry
}

type RegenEntry struct {
	IsFamily                                             bool
	X, Y, Width, Height, RangeDegree                     uint32
	MobIndex                                             uint32
	MobNum                                               uint16
	KillNumber                                           uint16
	RespawnSeconds, RespawnSecondsMin, RespawnSecondsMax uint32
	RespawnDeltas                                        [9]uint32
}

// load SHN files
// load into MonsterData
// persist monsterData

// load all regen files
// store them into MonsterRegenTable
// persist to bolt DB



func LoadRegenData(shineFolder string, db *bolt.DB) error {

	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("monsters"))
		if err != nil {
			return err
		}
		return nil
	})

	// load monsters/regens/*.txt

	// serialize into

	if err != nil {
		return err
	}

	return nil
}

func loadTxtFile(filePath string) ([][]string, error) {
	var data [][]string
	txtFile, err := os.Open(filePath)
	if err != nil {
		return data, err
	}
	reader := csv.NewReader(txtFile)

	reader.Comma = '\t'
	reader.FieldsPerRecord = -1

	data, err = reader.ReadAll()
	if err != nil {
		return data, err
	}
	return data, err
}
