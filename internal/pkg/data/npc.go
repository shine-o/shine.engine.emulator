package data

import (
	"strconv"
	"strings"
)

type NpcData struct {
	MapNPCs map[string][]*ShineNPC
	VendorNPCs map[string]*VendorItems
}

type VendorItems struct {
	Tabs map[string]VendorTab
}

type VendorTab struct {
	Pages []VendorGrid
}

type VendorGrid [8][6]string

type ShineNPC struct {
	MobIndex string
	MapIndex string
	X, Y, D  int
	NPCMenu  bool
	Role     string
	RoleArg  string // if Gate then ShinePortal is not nil
	*ShinePortal
}

type ShinePortal struct {
	RoleArg        string // use it later to assign ShinePortal to ShineNPC
	ServerMapIndex string
	ClientMapIndex string
	X, Y, D        int // landing coordinates
	FromLevel      int
	Party          bool
}

func LoadNPCData(shineFolder string) (*NpcData, error) {
	var (
		npc = &NpcData{
			MapNPCs: make(map[string][]*ShineNPC),
		}
	)

	err := loadMapNPCs(shineFolder, npc)

	if err != nil {
		return npc, err
	}

	// load vendor npcs

	return npc, nil
}

func loadMapNPCs(shineFolder string, npc *NpcData) error {
	var (
		npcs    []*ShineNPC
		portals []*ShinePortal
	)

	filesPath, err := ValidPath(shineFolder + "/world/" + "NPC.txt")

	if err != nil {
		return err
	}

	data, err := loadTxtFile(filesPath)
	if err != nil {
		return  err
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
					return err
				}

				y, err := strconv.Atoi(row[6])
				if err != nil {
					return err
				}

				d, err := strconv.Atoi(row[7])
				if err != nil {
					return err
				}

				npcMenu, err := strconv.ParseBool(row[8])
				if err != nil {
					return err
				}

				role := row[9]

				roleArg := row[10]

				sn := &ShineNPC{
					MobIndex: monsterIndex,
					MapIndex: mapIndex,
					X:        x,
					Y:        y,
					D:        d,
					NPCMenu:  npcMenu,
					Role:     role,
					RoleArg:  roleArg,
				}

				npcs = append(npcs, sn)

			} else {
				roleArg := row[10]
				serverMapIndex := row[11]
				clientMapIndex := row[12]

				x, err := strconv.Atoi(row[13])
				if err != nil {
					return err
				}

				y, err := strconv.Atoi(row[14])
				if err != nil {
					return err
				}

				fromLevel, err := strconv.Atoi(row[15])
				if err != nil {
					return err
				}

				party, err := strconv.ParseBool(row[16])
				if err != nil {
					return err
				}

				portals = append(portals, &ShinePortal{
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
					n.ShinePortal = p
				}
			}
		}
		npc.MapNPCs[n.MapIndex] = append(npc.MapNPCs[n.MapIndex], n)
	}
	return nil
}
