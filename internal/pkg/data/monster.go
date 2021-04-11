package data

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type MonsterData struct {
	MapRegens     map[string]MonsterRegenTable
	MobInfo       ShineMobInfo
	MobInfoServer ShineMobInfoServer
}

type MonsterRegenTable struct {
	Groups map[string]RegenEntry
}

type RegenEntry struct {
	IsFamily                         bool
	X, Y, Width, Height, RangeDegree int
	Mobs                             []RegenEntryMob
}

type RegenEntryMob struct {
	Index                                                string
	Num                                                  uint16
	KillNumber                                           int
	RespawnSeconds, RespawnSecondsMin, RespawnSecondsMax int
	RespawnDeltas                                        [9]uint32
}

// load SHN files
// load into MonsterData
// persist monsterData
func LoadMonsterData(shineFolder string) (*MonsterData, error) {
	var monsterData = &MonsterData{}
	var mobInfo ShineMobInfo
	var mobInfoServer ShineMobInfoServer

	mobInfoPath, err := ValidPath(shineFolder + "/shn/" + "MobInfo.shn")
	if err != nil {
		return monsterData, err
	}

	err = Load(mobInfoPath, &mobInfo)

	if err != nil {
		return monsterData, err
	}

	mobInfoServerPath, err := ValidPath(shineFolder + "/shn/" + "MobInfoServer.shn")

	if err != nil {
		return monsterData, err
	}

	err = Load(mobInfoServerPath, &mobInfoServer)

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

	root, err := ValidPath(root)

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

				y, err = strconv.Atoi(row[5])

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
					e.Mobs = append(e.Mobs, RegenEntryMob{
						Index:             mobIndex,
						Num:               uint16(mobNum),
						KillNumber:        killNum,
						RespawnSeconds:    regSec,
						RespawnSecondsMin: regSecMin,
						RespawnSecondsMax: regSecMax,
						RespawnDeltas:     [9]uint32{},
					})
					mrt.Groups[groupIndex] = e
				}

			}
		}

		// remove incomplete entries
		for inx, g := range mrt.Groups {
			if len(g.Mobs) == 0 {
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
