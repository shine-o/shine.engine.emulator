package mobs

import (
	"encoding/csv"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/shn"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var log *logger.Logger

func init() {
	log = logger.Init("maps logger", true, false, os.Stdout)
}

type MonsterData struct {
	MapRegens     map[string]MonsterRegenTable
	MobInfo       shn.ShineMobInfo
	MobInfoServer shn.ShineMobInfoServer
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
func LoadMonsterData(shineFolder string) (MonsterData, error) {
	var monsterData MonsterData
	var mobInfo shn.ShineMobInfo
	var mobInfoServer shn.ShineMobInfoServer

	mobInfoPath, err := utils.ValidPath(shineFolder + "/shn/client/" + "MobInfo.shn")
	if err != nil {
		return monsterData, err
	}

	err = shn.Load(mobInfoPath, &mobInfo)

	if err != nil {
		return monsterData, err
	}

	mobInfoServerPath, err := utils.ValidPath(shineFolder + "/shn/server/" + "MobInfoServer.shn")

	if err != nil {
		return monsterData, err
	}

	err = shn.Load(mobInfoServerPath, &mobInfoServer)

	if err != nil {
		return monsterData, err
	}

	monsterData.MobInfo = mobInfo
	monsterData.MobInfoServer = mobInfoServer

	mapRegens, err := loadRegenData(shineFolder)

	if err != nil {
		return monsterData, err
	}

	monsterData.MapRegens = mapRegens

	return monsterData, err
}

func loadRegenData(shineFolder string) (map[string]MonsterRegenTable, error) {

	mapRegens := make(map[string]MonsterRegenTable)

	// serialize into
	var files []string
	root := shineFolder + "/monsters/regens"

	root, err := utils.ValidPath(root)

	if err != nil {
		return mapRegens, err
	}

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	if err != nil {
		return mapRegens, err
	}

	for _, f := range files[1:] {

		data, err := loadTxtFile(f)

		if err != nil {
			return mapRegens, err
		}

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

				x, err := strconv.Atoi(strings.TrimSpace(row[4]))

				if err != nil {
					log.Error(err)
					continue
				}

				y, err = strconv.Atoi(row[ 5])

				if err != nil {
					log.Error(err)
					continue
				}

				width, err = strconv.Atoi(row[6])

				if err != nil {
					log.Error(err)
					continue
				}

				height, err = strconv.Atoi(row[7])

				if err != nil {
					log.Error(err)
					continue
				}

				rangeDegree, err = strconv.Atoi(row[8])

				if err != nil {
					log.Error(err)
					continue
				}

				mrt.Groups[groupIndex] = RegenEntry{
					IsFamily:    isFamily,
					X:           x,
					Y:           y,
					Width:       width,
					Height:      height,
					RangeDegree: rangeDegree,
				}
			}
		}

		// iterate only the table MobRegen last
		for _, row := range data[mobRegenStart:] {
			if row[0] == "#record" {
				var (
					groupIndex                   string
					mobIndex                     string
					mobNum                       int
					killNum                      int
					regSec, regSecMin, regSecMax int
				)

				groupIndex = row[2]

				mobIndex = row[3]

				mobNum, err := strconv.Atoi(row[4])

				if err != nil {
					log.Error(err)
					continue
				}

				killNum, err = strconv.Atoi(row[5])

				if err != nil {
					log.Error(err)
					continue
				}

				regSec, err = strconv.Atoi(row[6])

				if err != nil {
					log.Error(err)
					continue
				}

				regSecMin, err = strconv.Atoi(row[7])

				if err != nil {
					log.Error(err)
					continue
				}

				regSecMax, err = strconv.Atoi(row[8])

				if err != nil {
					log.Error(err)
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

		res := strings.Split(f, "/")

		if len(res) == 1 {
			// windows path
			res = strings.Split(f, "\\")
		}

		mapIndex := strings.Split(res[len(res)-1], ".")

		mapRegens[mapIndex[0]] = mrt
	}

	return mapRegens, err
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
