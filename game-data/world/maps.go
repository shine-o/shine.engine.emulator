package world

import (
	"encoding/csv"
	"fmt"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.core/game-data/blocks"
	"github.com/shine-o/shine.engine.core/game-data/shn"
	"github.com/shine-o/shine.engine.core/game-data/utils"
	"github.com/shine-o/shine.engine.core/structs"
	bolt "go.etcd.io/bbolt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
)

var log *logger.Logger

func init() {
	log = logger.Init("maps logger", true, false, ioutil.Discard)
	log.Info("maps logger init()")
}

type Map struct {
	ID             int `struct:"int32"`
	MapAttributeID int `struct:"int32"`
}

type MapAttributes struct {
	ID               int    `struct:"int32"`
	MapInfoIndex     string `struct:"[12]byte"`
	ShineMapName     string `struct:"[32]byte"`
	SectorX          int    `struct:"int32"`
	SectorY          int    `struct:"int32"`
	IdleDuration     int    `struct:"int32"`
	ScriptFile       string `struct:"[64]byte"`
	PKPoints         bool   `struct:"byte"`
	HideName         bool   `struct:"byte"`
	TeleportTo       bool   `struct:"byte"`
	TeleportFrom     bool   `struct:"byte"`
	RegenMapIndex    string `struct:"[12]byte"`
	RegenXA          int    `struct:"int32"`
	RegenYA          int    `struct:"int32"`
	RegenXB          int    `struct:"int32"`
	RegenYB          int    `struct:"int32"`
	RegenXC          int    `struct:"int32"`
	RegenYC          int    `struct:"int32"`
	RegenXD          int    `struct:"int32"`
	RegenYD          int    `struct:"int32"`
	RegenSpots       int    `struct:"int8"`
	CanTrade         bool   `struct:"byte"`
	CanRest          bool   `struct:"byte"`
	UseItem          bool   `struct:"byte"`
	CastSkill        bool   `struct:"byte"`
	UseChat          bool   `struct:"byte"`
	UseShout         bool   `struct:"byte"`
	OpenBooth        bool   `struct:"byte"`
	Produce          bool   `struct:"byte"`
	UseMount         bool   `struct:"byte"`
	UseStones        bool   `struct:"byte"`
	PartyType        int    `struct:"int8"`
	MobExpPenalty    int    `struct:"int32"`
	PlayerExpPenalty int    `struct:"int32"`
}

type MapData struct {
	ID         int `struct:"int32"`
	Attributes MapAttributes
	Info       shn.MapInfo
	SHBD       blocks.SHBD
}

type mapData struct {
	attributes  map[int]MapAttributes
	mapInfo     *shn.ShineMapInfo
	shineFolder string
}

func LoadMapData(shineFolder string, db *bolt.DB) error {
	var attributes map[int]MapAttributes
	var mapInfo shn.ShineMapInfo

	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("maps"))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	mapFiles := []string{"NormalMaps.txt", "DungeonMaps.txt", "KingdomQuestMaps.txt", "GuildTournamentMaps.txt"}

	mapAttributesPath, err := utils.ValidPath(shineFolder + "/world/" + "MapAttributes.txt")

	attributes, err = loadMapAttributes(mapAttributesPath)
	if err != nil {
		return err
	}

	mapInfoPath, err := utils.ValidPath(shineFolder + "/shn/client/" + "MapInfo.shn")
	if err != nil {
		return err
	}

	err = shn.Load(mapInfoPath, &mapInfo)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, file := range mapFiles {
		var maps []Map
		mapsPath, err := utils.ValidPath(shineFolder + "/world/" + file)
		maps, err = loadMaps(mapsPath)
		if err != nil {
			return err
		}
		for _, m := range maps {
			wg.Add(1)
			md := mapData{
				shineFolder: shineFolder,
				attributes:  attributes,
				mapInfo:     &mapInfo,
			}
			go md.persistMap(&wg, m, db)
		}
	}
	wg.Wait()
	return nil
}

func (md *mapData) persistMap(wg *sync.WaitGroup, m Map, db *bolt.DB) {
	defer wg.Done()
	var lm MapData
	if attr, ok := md.attributes[m.MapAttributeID]; ok {
		lm.Attributes = attr
	} else {
		log.Errorf("unkown map attribute entry with ID %v", m.MapAttributeID)
		return
	}

	for _, row := range md.mapInfo.Rows {
		if row.MapName.Name == lm.Attributes.MapInfoIndex {
			lm.Info = row
		}
	}

	if lm.Info == (shn.MapInfo{}) {
		log.Errorf("no MapInfo.shn entry found for normal map entry with ID %v, ignoring map", lm.Attributes.ID)
		return
	}

	// load shbd
	var s *blocks.SHBD

	shbdPath, err := utils.ValidPath(md.shineFolder + "/blocks/" + lm.Info.MapName.Name + ".shbd")
	if err != nil {
		log.Errorf("shbd file found for normal map entry with ID %v, ignoring map %v", lm.Attributes.ID, err)
	}

	s, err = blocks.LoadSHBDFile(shbdPath)
	if err != nil {
		log.Errorf("failed to load shbd file for map entry with ID %v, ignoring map %v", lm.Attributes.ID, err)
	}

	lm.SHBD = *s
	lm.ID = m.ID
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("maps"))
		data, err := structs.Pack(&lm)
		if err != nil {
			return err
		}

		err = b.Put([]byte(fmt.Sprintf("%v", m.ID)), data)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Errorf("failed to persist map with id %v %v", lm.Attributes.ID, err)
	}
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

func loadMaps(filePath string) ([]Map, error) {
	var maps []Map
	data, err := loadTxtFile(filePath)
	if err != nil {
		return maps, err
	}
	for i, datum := range data {
		if i < 2 {
			continue
		}

		id, err := strconv.Atoi(datum[0])
		if err != nil {
			return maps, err
		}
		mapAttributeID, err := strconv.Atoi(datum[1])

		if err != nil {
			return maps, err
		}
		maps = append(maps, Map{
			ID:             id,
			MapAttributeID: mapAttributeID,
		})
	}
	return maps, nil
}

func loadMapAttributes(filePath string) (map[int]MapAttributes, error) {
	var attributes = make(map[int]MapAttributes, 0)

	data, err := loadTxtFile(filePath)

	if err != nil {
		return attributes, err
	}

	for i, datum := range data {

		if i < 2 {
			continue
		}

		if len(datum) != 34 {
			return attributes, fmt.Errorf("unexpected number of columns %v", len(datum))
		}

		ma, err := mapAttributeFields(datum)
		if err != nil {
			return attributes, err
		}
		attributes[ma.ID] = ma
	}
	return attributes, nil
}

func mapAttributeFields(fields []string) (MapAttributes, error) {
	var ma MapAttributes

	//todo:  a better way to iterate over these fields
	id, err := strconv.Atoi(fields[0])
	if err != nil {
		return ma, err
	}

	sectorX, err := strconv.Atoi(fields[3])
	if err != nil {
		return ma, err
	}

	sectorY, err := strconv.Atoi(fields[4])
	if err != nil {
		return ma, err
	}

	idleDuration, err := strconv.Atoi(fields[5])
	if err != nil {
		return ma, err
	}

	pKPoints, err := strconv.ParseBool(fields[7])
	if err != nil {
		return ma, err
	}

	hideName, err := strconv.ParseBool(fields[8])
	if err != nil {
		return ma, err
	}

	teleportTo, err := strconv.ParseBool(fields[9])
	if err != nil {
		return ma, err
	}

	teleportFrom, err := strconv.ParseBool(fields[10])
	if err != nil {
		return ma, err
	}

	regenXA, err := strconv.Atoi(fields[12])
	if err != nil {
		return ma, err
	}

	regenYA, err := strconv.Atoi(fields[13])
	if err != nil {
		return ma, err
	}

	regenXB, err := strconv.Atoi(fields[14])
	if err != nil {
		return ma, err
	}

	regenYB, err := strconv.Atoi(fields[15])
	if err != nil {
		return ma, err
	}

	regenXC, err := strconv.Atoi(fields[16])
	if err != nil {
		return ma, err
	}

	regenYC, err := strconv.Atoi(fields[17])
	if err != nil {
		return ma, err
	}

	regenXD, err := strconv.Atoi(fields[18])
	if err != nil {
		return ma, err
	}

	regenYD, err := strconv.Atoi(fields[19])
	if err != nil {
		return ma, err
	}

	regenSpots, err := strconv.Atoi(fields[20])
	if err != nil {
		return ma, err
	}

	canTrade, err := strconv.ParseBool(fields[21])
	if err != nil {
		return ma, err
	}
	canRest, err := strconv.ParseBool(fields[22])
	if err != nil {
		return ma, err
	}
	useItem, err := strconv.ParseBool(fields[23])
	if err != nil {
		return ma, err
	}
	castSkill, err := strconv.ParseBool(fields[24])
	if err != nil {
		return ma, err
	}
	useChat, err := strconv.ParseBool(fields[25])
	if err != nil {
		return ma, err
	}
	useShout, err := strconv.ParseBool(fields[26])
	if err != nil {
		return ma, err
	}
	openBooth, err := strconv.ParseBool(fields[27])
	if err != nil {
		return ma, err
	}
	produce, err := strconv.ParseBool(fields[28])
	if err != nil {
		return ma, err
	}
	useMount, err := strconv.ParseBool(fields[29])
	if err != nil {
		return ma, err
	}
	useStones, err := strconv.ParseBool(fields[30])
	if err != nil {
		return ma, err
	}
	partyType, err := strconv.Atoi(fields[31])
	if err != nil {
		return ma, err
	}

	mobExpPenalty, err := strconv.Atoi(fields[32])
	if err != nil {
		return ma, err
	}

	playerExpPenalty, err := strconv.Atoi(fields[33])
	if err != nil {
		return ma, err
	}

	ma = MapAttributes{
		ID:               id,
		MapInfoIndex:     fields[1],
		ShineMapName:     fields[2],
		SectorX:          sectorX,
		SectorY:          sectorY,
		IdleDuration:     idleDuration,
		ScriptFile:       fields[6],
		PKPoints:         pKPoints,
		HideName:         hideName,
		TeleportTo:       teleportTo,
		TeleportFrom:     teleportFrom,
		RegenMapIndex:    fields[11],
		RegenXA:          regenXA,
		RegenYA:          regenYA,
		RegenXB:          regenXB,
		RegenYB:          regenYB,
		RegenXC:          regenXC,
		RegenYC:          regenYC,
		RegenXD:          regenXD,
		RegenYD:          regenYD,
		RegenSpots:       regenSpots,
		CanTrade:         canTrade,
		CanRest:          canRest,
		UseItem:          useItem,
		CastSkill:        castSkill,
		UseChat:          useChat,
		UseShout:         useShout,
		OpenBooth:        openBooth,
		Produce:          produce,
		UseMount:         useMount,
		UseStones:        useStones,
		PartyType:        partyType,
		MobExpPenalty:    mobExpPenalty,
		PlayerExpPenalty: playerExpPenalty,
	}
	return ma, nil
}

