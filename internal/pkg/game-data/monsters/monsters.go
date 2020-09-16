package monsters

import (
	"encoding/csv"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/shn"
	bolt "go.etcd.io/bbolt"
	"os"
	"strconv"
)

type MonsterData struct {
	Regens        MonsterRegenTable
	MobInfo       shn.MobInfo
	MobInfoServer shn.MobInfoServer
}

type MonsterRegenTable struct {
	Groups map[string]RegenEntry
}

type RegenEntry struct {
	IsFamily                                             bool
	X, Y, Width, Height, RangeDegree                     int
	MobIndex                                             string
	MobNum                                               uint16
	KillNumber                                           int
	RespawnSeconds, RespawnSecondsMin, RespawnSecondsMax int
	RespawnDeltas                                        [9]uint32
}

// load SHN files
// load into MonsterData
// persist monsterData

// load all regen files
// store them into MonsterRegenTable
// persist to bolt DB

// LoadSHNData

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

	// if #table == MobRegenGroup
	// 	for each record we make a RegenEntry and add it to MonsterRegenTable using GroupIndex as key

	// load all file names inside  monsters/regens

	// load each of them and create a MonsterRegenTable for it

	mrt := MonsterRegenTable{
		Groups: make(map[string]RegenEntry),
	}

	var mobRegenStart int

	// iterate only the table MobRegenGroup first
	for i, row := range data {

		if row[1] == "MobRegen" {
			mobRegenStart = i
			break
		}

		if row[0] == "#record" {
			var (
				groupIndex                       string
				isFamily                         = false
				x, y, width, height, rangeDegree int
			)

			groupIndex = row[2]

			if row[3] == "Y" {
				isFamily = true
			}

			x, err := strconv.Atoi(row[4])

			if err != nil {
				logger.Error(err)
				continue
			}

			y, err = strconv.Atoi(row[5])

			if err != nil {
				logger.Error(err)
				continue
			}

			width, err = strconv.Atoi(row[6])

			if err != nil {
				logger.Error(err)
				continue
			}

			height, err = strconv.Atoi(row[7])

			if err != nil {
				logger.Error(err)
				continue
			}

			rangeDegree, err = strconv.Atoi(row[8])

			if err != nil {
				logger.Error(err)
				continue
			}

			mrt.Groups[groupIndex] = monsters.RegenEntry{
				IsFamily:    isFamily,
				X:           x,
				Y:           y,
				Width:       width,
				Height:      height,
				RangeDegree: rangeDegree,
			}
		}
	}

	for _, row := range data[mobRegenStart:] {
		if row[0] == "#record" {
			var (
				groupIndex                string
				mobIndex                     string
				mobNum                       int
				killNum                      int
				regSec, regSecMin, regSecMax int
			)

			groupIndex = row[2]

			mobIndex = row[3]

			mobNum, err := strconv.Atoi(row[4])

			if err != nil {
				logger.Error(err)
				continue
			}

			killNum, err = strconv.Atoi(row[5])

			if err != nil {
				logger.Error(err)
				continue
			}

			regSec, err = strconv.Atoi(row[6])

			if err != nil {
				logger.Error(err)
				continue
			}

			regSecMin, err = strconv.Atoi(row[7])

			if err != nil {
				logger.Error(err)
				continue
			}

			regSecMax, err = strconv.Atoi(row[8])

			if err != nil {
				logger.Error(err)
				continue
			}

			e, ok := mrt.Groups[groupIndex]

			if ok {
				e.MobIndex = mobIndex
				e.MobNum = uint16(mobNum)
				e.KillNumber = killNum
				e.RespawnSeconds = regSec
				e.RespawnSecondsMin = regSecMin
				e.RespawnSecondsMax = regSecMax
				mrt.Groups[groupIndex] = e
			}

		}
	}

	// remove incomplete entries
	for inx, g := range mrt.Groups {
		if g.MobIndex == "" {
			delete(mrt.Groups, inx)
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func LoadTxtFile(filePath string) ([][]string, error) {
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
