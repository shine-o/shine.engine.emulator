package world

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/utils"
	"strconv"
	"strings"
)

//#Table	ShineNPC
//#ColumnType			STRING[33]	STRING[20]	DWRD	DWRD	WORD	BYTE	INDEX	INDEX
//#ColumnName			MobName	Map	Coord-X	Coord-Y	Direct	NPCMenu	Role	RoleArg0
//
//#Table	LinkTable
//#ColumnType		Index	String[33]	String[33]	DWRD	DWRD	WORD	WORD	WORD	BYTE
//#ColumnName		argument	MapServer	MapClient	Coord-X	Coord-Y	Direct	LevelFrom	LevelTo	Party


type NPC struct {
	Data map[string][]*ShineNPC
}

type ShineNPC struct {
	MobIndex string
	MapIndex string
	X,Y,D int
	NPCMenu bool
	Role string
	RoleArg string // if Gate then LinkTable is not nil
	*LinkTable
}

type LinkTable struct {
	RoleArg string // use it later to assign LinkTable to ShineNPC
	ServerMapIndex string
	ClientMapIndex string
	X,Y int
	FromLevel int
	Party bool
}

func LoadNPCData(shineFolder string) (*NPC, error) {
	var (
		npc = &NPC{
			Data: make(map[string][]*ShineNPC),
		}
		npcs []*ShineNPC
		portals []*LinkTable
	)

	filePath, err := utils.ValidPath(shineFolder + "/world/" + "NPC.txt")

	if err != nil {
		return npc, err
	}

	data, err := loadTxtFile(filePath)
	if err != nil {
		return npc, err
	}

	for _, row := range data {

		if row[0] == ";" || strings.Contains(row[0], ";") {
			continue
		}

		if row[1] == "#recordin" || strings.Contains(row[1], "#recordin") {
			if len(row) <= 2 {
				continue
			}
			if row[2] == "ShineNPC" {

				monsterIndex := row[3]

				mapIndex := row[4]

				x, err := strconv.Atoi(row[5])
				if err != nil {
					return npc, err
				}

				y, err := strconv.Atoi(row[6])
				if err != nil {
					return npc, err
				}

				d, err := strconv.Atoi(row[7])
				if err != nil {
					return npc, err
				}

				npcMenu, err := strconv.ParseBool(row[8])
				if err != nil {
					return npc, err
				}

				role := row[9]

				roleArg := row[10]

				sn := &ShineNPC{
					MobIndex:  monsterIndex,
					MapIndex:  mapIndex,
					X:         x,
					Y:         y,
					D:         d,
					NPCMenu:   npcMenu,
					Role:      role,
					RoleArg:   roleArg,
				}

				//npc.Data[sn.MapIndex] = sn
				npcs = append(npcs, sn)

			} else {
				roleArg := row[10]
				serverMapIndex := row[11]
				clientMapIndex := row[12]

				x, err := strconv.Atoi(row[13])
				if err != nil {
					return npc, err
				}

				y, err := strconv.Atoi(row[14])
				if err != nil {
					return npc, err
				}

				fromLevel, err := strconv.Atoi(row[15])
				if err != nil {
					return npc, err
				}

				party, err := strconv.ParseBool(row[16])
				if err != nil {
					return npc, err
				}

				portals = append(portals, &LinkTable{
					RoleArg:        roleArg,
					ServerMapIndex: serverMapIndex,
					ClientMapIndex: clientMapIndex,
					X:              x,
					Y:              y,
					FromLevel:      fromLevel,
					Party:          party,
				})
			}
		}
	}

	for i, _ := range npcs {
		n := npcs[i]
		if n.Role == "IDGate" || n.Role == "Gate" || n.Role == "RandomGate" {
			for j, _ := range portals {
				p := portals[j]
				if p.RoleArg == n.RoleArg {
					n.LinkTable = p
				}
			}
		}
		npc.Data[n.MapIndex] = append(npc.Data[n.MapIndex], n)
	}

	return npc, nil
}