package character

import (
	"encoding/hex"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/google/logger"
	"github.com/google/uuid"
	"github.com/shine-o/shine.engine.core/structs"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
)


var log *logger.Logger

func init() {
	log = logger.Init("character logger", true, false, ioutil.Discard)
	log.Info("character logger init()")
}

// ErrCharacter is used to handle known errors
type ErrCharacter struct {
	Code    int
	Message string
}

func (ec *ErrCharacter) Error() string {
	return ec.Message
}

// EquippedItems model for the database layer
type EquippedItems struct {
	tableName        struct{} `pg:"world.character_equipped_items"`
	ID               uint64
	CharacterID      uint64 //
	Character        *Character
	Head             uint16
	Face             uint16
	Body             uint16
	Pants            uint16
	Boots            uint16
	LeftHand         uint16
	RightHand        uint16
	LeftMiniPet      uint16
	RightMiniPet     uint16
	ApparelHead      uint16
	ApparelFace      uint16
	ApparelEye       uint16
	ApparelBody      uint16
	ApparelPants     uint16
	ApparelBoots     uint16
	ApparelLeftHand  uint16
	ApparelRightHand uint16
	ApparelBack      uint16
	ApparelTail      uint16
	ApparelAura      uint16
	ApparelShield    uint16
	DeletedAt        time.Time `pg:",soft_delete"`
}

// Character model for the database layer
type Character struct {
	tableName     struct{} `pg:"world.characters"`
	ID            uint64
	UserID        uint64 `pg:",notnull"`
	Name          string `pg:",notnull,unique"`
	Appearance    *Appearance
	Attributes    *Attributes
	Location      *Location
	Options       *ClientOptions
	Items         []Items
	EquippedItems *EquippedItems
	AdminLevel    uint8 `pg:",notnull,use_zero"`
	Slot          uint8 `pg:",notnull,use_zero"`
	IsDeleted     bool  `pg:",use_zero"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time `pg:",soft_delete"`
}

// Appearance model for the database layer
type Appearance struct {
	tableName   struct{} `pg:"world.character_appearance"`
	ID          uint64
	CharacterID uint64 //
	Character   *Character
	Class       uint8 `pg:",notnull"`
	Gender      uint8 `pg:",notnull,use_zero"`
	HairType    uint8 `pg:",notnull,use_zero"`
	HairColor   uint8 `pg:",notnull,use_zero"`
	FaceType    uint8 `pg:",notnull,use_zero"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time `pg:",soft_delete"`
}

// Attributes model for the database layer
type Attributes struct {
	tableName    struct{} `pg:"world.character_attributes"`
	ID           uint64
	CharacterID  uint64
	Character    *Character
	Level        uint8  `pg:",notnull"`
	Experience   uint64 `pg:",notnull,use_zero"`
	Fame         uint32 `pg:",notnull,use_zero"`
	Hp           uint32 `pg:",notnull"`
	Sp           uint32 `pg:",notnull"`
	Intelligence uint8  `pg:",notnull,use_zero"`
	Strength     uint8  `pg:",notnull,use_zero"`
	Dexterity    uint8  `pg:",notnull,use_zero"`
	Endurance    uint8  `pg:",notnull,use_zero"`
	Spirit       uint8  `pg:",notnull,use_zero"`
	Money        uint64 `pg:",notnull,use_zero"`
	KillPoints   uint32 `pg:",notnull,use_zero"`
	HpStones     uint16 `pg:",notnull"`
	SpStones     uint16 `pg:",notnull"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time `pg:",soft_delete"`
}

// Location model for the database layer
type Location struct {
	tableName   struct{} `pg:"world.character_location"`
	ID          uint64
	CharacterID uint64 //
	Character   *Character
	MapID       uint32 `pg:",notnull"`
	MapName     string `pg:",notnull"`
	X           uint32 `pg:",notnull"`
	Y           uint32 `pg:",notnull"`
	D           uint8  `pg:",notnull,use_zero"`
	IsKQ        bool   `pg:",notnull,use_zero"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time `pg:",soft_delete"`
}

type ClientOptions struct {
	tableName   struct{} `pg:"world.client_options"`
	ID          uint64
	CharacterID uint64 //
	Character   *Character
	GameOptions []byte `pg:",notnull"`
	Keymap      []byte `pg:",notnull"`
	Shortcuts   []byte `pg:",notnull"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// todo: move to inventory package
type Items struct {
	tableName struct{} `pg:"world.character_items"`
	//ID          uint64
	CharacterID uint64 `pg:",pk,use_zero"`
	Character   *Character
	// box 2 = reward inventory
	// box 3 = mini house furniture
	// box 8 = equipped items
	// box 9 = inventory, storage
	// box 12 = mini houses
	// box 13 = mini house accessories
	// box 14 = mini house tile all inventory
	// box 15 = premium actions inventory(dances)
	// box 16 = mini house mini game inventory
	InventoryType int    `pg:",pk,use_zero"`
	Slot          uint16 `pg:",pk,use_zero"`
	ShnID         uint16 `pg:",notnull"`
	Stackable     bool   `pg:",notnull,use_zero"`
	Amount        uint32
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time `pg:",soft_delete"`
}

const (
	startLevel 	 = 1
	startMapID   = 1
	startMapName   = "Rou"
)

// ErrInvalidSlot happens if the client tries to bypass client side verification
var ErrInvalidSlot = &ErrCharacter{
	Code:    0,
	Message: "invalid slot",
}

// ErrNameTaken name is reserved or in use
var ErrNameTaken = &ErrCharacter{
	Code:    1,
	Message: "name taken",
}

// ErrNoSlot happens if the client tries to bypass client side verification
var ErrNoSlot = &ErrCharacter{
	Code:    2,
	Message: "no slot available",
}

// ErrInvalidName happens if the client tries to bypass client side verification
var ErrInvalidName = &ErrCharacter{
	Code:    3,
	Message: "invalid name",
}

// ErrInvalidClassGender happens if the client tries to bypass client side verification
var ErrInvalidClassGender = &ErrCharacter{
	Code:    4,
	Message: "invalid class gender data",
}

// Validate checks data sent by the client is valid
func Validate(db *pg.DB, userID uint64, req structs.NcAvatarCreateReq) error {

	if req.SlotNum > 5 {
		return ErrInvalidSlot
	}

	name := req.Name.Name

	var charName string
	err := db.Model((*Character)(nil)).Column("name").Where("name = ?", name).Select(&charName)

	if err == nil {
		return ErrNameTaken
	}

	var chars []Character
	err = db.Model(&chars).Where("user_id = ?", userID).Select()

	if len(chars) == 6 {
		return ErrNoSlot
	}

	alphaNumeric := regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString
	if !alphaNumeric(name) {
		return ErrInvalidName
	}


	// todo: missing validation: default hair, color, face values
	// todo: missing validation: default starter class values ( mage, cleric, archer, fighter)
	isMale := (req.Shape.BF >> 7) & 1
	class := (req.Shape.BF >> 2) & 31

	if isMale > 1 || isMale < 0 {
		return ErrInvalidClassGender
	}

	if class < 1 || class > 27 {
		return ErrInvalidClassGender
	}

	return nil
}

// New creates character for the User with userID and returns data the client can understand
func New(db *pg.DB, userID uint64, req structs.NcAvatarCreateReq) (structs.AvatarInformation, error) {
	newTx, err := db.Begin()

	if err != nil {
		return structs.AvatarInformation{}, err
	}

	defer newTx.Close()

	name := strings.TrimRight(string(req.Name.Name[:]), "\x00")
	char := Character{
		UserID:     userID,
		AdminLevel: 0,
		Name:       name,
		Slot:       req.SlotNum,
	}

	_, err = newTx.Model(&char).Returning("*").Insert()

	if err != nil {
		return structs.AvatarInformation{}, &ErrCharacter{
			Code:    7,
			Message: fmt.Sprintf("%v:%v", err, newTx.Rollback()),
		}
	}

	char.
		initialAppearance(req.Shape).
		initialAttributes().
		initialLocation().
		initialClientOptions().
		initialEquippedItems()

	if _, err = newTx.Model(char.Appearance).Returning("*").Insert(); err != nil {
		return structs.AvatarInformation{}, &ErrCharacter{
			Code:    7,
			Message: fmt.Sprintf("%v:%v", err, newTx.Rollback()),
		}
	}

	if _, err = newTx.Model(char.Attributes).Returning("*").Insert(); err != nil {
		return structs.AvatarInformation{}, &ErrCharacter{
			Code:    7,
			Message: fmt.Sprintf("%v:%v", err, newTx.Rollback()),
		}
	}

	if _, err = newTx.Model(char.Location).Returning("*").Insert(); err != nil {
		return structs.AvatarInformation{}, &ErrCharacter{
			Code:    7,
			Message: fmt.Sprintf("%v:%v", err, newTx.Rollback()),
		}
	}

	if _, err = newTx.Model(char.Options).Returning("*").Insert(); err != nil {
		return structs.AvatarInformation{}, &ErrCharacter{
			Code:    7,
			Message: fmt.Sprintf("%v:%v", err, newTx.Rollback()),
		}
	}

	if _, err = newTx.Model(char.EquippedItems).Returning("*").Insert(); err != nil {
		return structs.AvatarInformation{}, &ErrCharacter{
			Code:    7,
			Message: fmt.Sprintf("%v:%v", err, newTx.Rollback()),
		}
	}
	return char.NcRepresentation(), newTx.Commit()
}

func Get(db *pg.DB, characterID uint64) (Character, error) {
	var c Character
	c.ID = characterID
	err := db.Model(&c).
		WherePK().
		Relation("Appearance").
		Relation("Attributes").
		Relation("Location").
		Select()
	return c, err
}

func GetByName(db *pg.DB, name string) (Character, error) {
	var c Character
	err := db.Model(&c).
		Relation("Appearance").
		Relation("Attributes").
		Relation("Location").
		Where("name = ?", name).
		Select() // query the world for a character with name
	return c, err
}

// Delete character for User with userID
// soft deletion is performed
func Delete(db *pg.DB, userID uint64, req structs.NcAvatarEraseReq) error {
	deleteTx, err := db.Begin()
	defer deleteTx.Close()
	if err != nil {
		return &ErrCharacter{
			Code:    1,
			Message: fmt.Sprintf("database error, could not start transaction: %v", err),
		}
	}
	var char Character
	err = deleteTx.Model(&char).Where("user_id = ?", userID).Where("slot = ?", req.Slot).Select()

	if err != nil {
		txErr := deleteTx.Rollback()
		return &ErrCharacter{
			Code:    2,
			Message: fmt.Sprintf("database error, character not found: %v, %v", err, txErr),
		}
	}

	name := fmt.Sprintf("%v@%v", char.Name, uuid.New().String())
	_, err = deleteTx.Model((*Character)(nil)).Set("name = ?", name).Where("user_id = ?", userID).Where("slot = ? ", req.Slot).Update()
	if err != nil {
		txErr := deleteTx.Rollback()
		return &ErrCharacter{
			Code:    3,
			Message: fmt.Sprintf("database error: %v, %v", err, txErr),
		}
	}

	if _, err = deleteTx.Model(&char).Where("user_id = ?", userID).Where("slot = ?", req.Slot).Delete(); err != nil {
		return &ErrCharacter{
			Code:    7,
			Message: fmt.Sprintf("%v:%v", err, deleteTx.Rollback()),
		}
	}

	if _, err = deleteTx.Model(char.Appearance).Where("character_id = ?", char.ID).Delete(); err != nil {
		return &ErrCharacter{
			Code:    7,
			Message: fmt.Sprintf("%v:%v", err, deleteTx.Rollback()),
		}
	}

	if _, err = deleteTx.Model(char.Attributes).Where("character_id = ?", char.ID).Delete(); err != nil {
		return &ErrCharacter{
			Code:    7,
			Message: fmt.Sprintf("%v:%v", err, deleteTx.Rollback()),
		}
	}

	if _, err = deleteTx.Model(char.Location).Where("character_id = ?", char.ID).Delete(); err != nil {
		return &ErrCharacter{
			Code:    7,
			Message: fmt.Sprintf("%v:%v", err, deleteTx.Rollback()),
		}
	}

	if _, err = deleteTx.Model(char.EquippedItems).Where("character_id = ?", char.ID).Delete(); err != nil {
		return &ErrCharacter{
			Code:    7,
			Message: fmt.Sprintf("%v:%v", err, deleteTx.Rollback()),
		}
	}

	return deleteTx.Commit()
}

func (c *Character) initialAppearance(shape structs.ProtoAvatarShapeInfo) *Character {
	isMale := (shape.BF >> 7) & 1
	class := (shape.BF >> 2) & 31

	c.Appearance = &Appearance{
		CharacterID: c.ID,
		Class:       class,
		Gender:      isMale,
		HairType:    shape.HairType,
		HairColor:   shape.HairColor,
		FaceType:    shape.FaceShape,
	}
	return c
}

func (c *Character) initialAttributes() *Character {
	c.Attributes = &Attributes{
		CharacterID:  c.ID,
		Level:        startLevel,
		Experience:   0,
		Fame:         0,
		Hp:           500,
		Sp:           500,
		Intelligence: 27,
		Strength:     27,
		Dexterity:    27,
		Endurance:    27,
		Spirit:       27,
		Money:        100,
		KillPoints:   0,
		HpStones:     15,
		SpStones:     15,
	}
	return c
}

func (c *Character) initialLocation() *Character {
	// loadDB
	c.Location = &Location{
		CharacterID: c.ID,
		MapID:       startMapID,
		MapName:     startMapName,
		X:           5323,
		Y:           4501,
		D:           90,
		IsKQ:        false,
	}
	return c
}

func (c *Character) initialClientOptions() *Character {
	// hardcoded defaults :)
	// game_options: 20000000010100010200010300010400000500010600000700010800010900010a00010b00010c00000d00000e00000f00001000011100001200001300001400011500001600001700001800011900011a00011b00011c00011d00001e00001f0001
	goData, _ := hex.DecodeString("20000000010100010200010300010400000500010600000700010800010900010a00010b00010c00000d00000e00000f00001000011100001200001300001400011500001600001700001800011900011a00011b00011c00011d00001e00001f0001")
	// keymap: 5f00000000790100001b02000043030000490400004b0500004c0600004607000048080000560900000d0a00114e0b0011470c0011500d0011570e0000de0f000058100000471100000012000000130000001400000015000052160011411700005718000053190000001a0000411b0000441c00005a1d0000201e0000261f000028200000252100002722000024230000542400005125000045260000f527000042280000502900004d2a0000552b00105a2c0000002d0000002e0000232f000031300000323100003332000034330000353400003635000037360000383700003938000030390000bd3a0000bb3b0010313c0010323d0010333e0010343f0010354000103641001037420010384300103944001030450010bd460010bb4700123148001232490012334a0012344b0012354c0012364d0012374e0012384f00123950001230510012bd520012bb530000005400000055000000560000005700000058000000590000005a0000005b0000005c0000005d0000005e000000
	kmData, _ := hex.DecodeString("5f00000000790100001b02000043030000490400004b0500004c0600004607000048080000560900000d0a00114e0b0011470c0011500d0011570e0000de0f000058100000471100000012000000130000001400000015000052160011411700005718000053190000001a0000411b0000441c00005a1d0000201e0000261f000028200000252100002722000024230000542400005125000045260000f527000042280000502900004d2a0000552b00105a2c0000002d0000002e0000232f000031300000323100003332000034330000353400003635000037360000383700003938000030390000bd3a0000bb3b0010313c0010323d0010333e0010343f0010354000103641001037420010384300103944001030450010bd460010bb4700123148001232490012334a0012344b0012354c0012364d0012374e0012384f00123950001230510012bd520012bb530000005400000055000000560000005700000058000000590000005a0000005b0000005c0000005d0000005e000000")
	// shortcuts: 040000040000000000010400010000000a0100ac0d00000b0100b10d0000
	scData, _ := hex.DecodeString("040000040000000000010400010000000a0100ac0d00000b0100b10d0000")

	c.Options = &ClientOptions{
		CharacterID: c.ID,
		GameOptions: goData, // hardcoded byte slice
		Keymap:      kmData, // hardcoded byte slice
		Shortcuts:   scData, // hardcoded byte slice
	}
	return c
}

func (c *Character) initialEquippedItems() *Character {
	c.EquippedItems = &EquippedItems{
		CharacterID:      c.ID,
		Head:             65535,
		Face:             65535,
		Body:             65535,
		Pants:            65535,
		Boots:            65535,
		LeftHand:         65535,
		RightHand:        65535,
		LeftMiniPet:      65535,
		RightMiniPet:     65535,
		ApparelHead:      65535,
		ApparelFace:      65535,
		ApparelEye:       65535,
		ApparelBody:      65535,
		ApparelPants:     65535,
		ApparelBoots:     65535,
		ApparelLeftHand:  65535,
		ApparelRightHand: 65535,
		ApparelBack:      65535,
		ApparelTail:      65535,
		ApparelAura:      65535,
		ApparelShield:    65535,
	}
	return c
}

func (c *Character) AllEquippedItems(db *pg.DB) *structs.NcCharClientItemCmd {
	return c.getItemsByInventory(db, 8)
}

func (c *Character) InventoryItems(db *pg.DB) *structs.NcCharClientItemCmd {
	return c.getItemsByInventory(db, 9)
}

func (c *Character) MiniHouseItems(db *pg.DB) *structs.NcCharClientItemCmd {
	return c.getItemsByInventory(db, 12)
}

func (c *Character) PremiumActionItems(db *pg.DB) *structs.NcCharClientItemCmd {
	return c.getItemsByInventory(db, 15)
}

// if not 65535, add item to the list
// todo: shn item processing
// get all items where character id and inventory type (8 for equipped items) match
func (c *Character) getItemsByInventory(db *pg.DB, inventoryType uint8) *structs.NcCharClientItemCmd {
	nc := &structs.NcCharClientItemCmd{
		NumOfItem: 0,
		Box:       inventoryType,
		Flag: structs.ProtoNcCharClientItemCmdFlag{
			BF0: 183,
		},
	}
	var items []Items
	err := db.Model(&items).Where("character_id = ?", c.ID).Where("inventory_type = ?", inventoryType).Select()

	switch inventoryType {
	case 8:
		nc.Flag.BF0 = 183
		break
	case 9:
		nc.Flag.BF0 = 165
		break
	case 12:
		nc.Flag.BF0 = 209
		break
	case 15:
		nc.Flag.BF0 = 243
		break
	}

	if err != nil {
		return nc
	}

	return nc
}

// NcRepresentation returns a struct that can be serialized into bytes and can be sent to the client
func (c *Character) NcRepresentation() structs.AvatarInformation {
	nc := structs.AvatarInformation{
		ChrRegNum: uint32(c.ID),
		Name: structs.Name5{
			Name: c.Name,
		},
		Level: uint16(c.Attributes.Level),
		Slot:  c.Slot,
		LoginMap: structs.Name3{
			Name: c.Location.MapName,
		},
		DelInfo: structs.ProtoAvatarDeleteInfo{},
		Shape:   c.Appearance.NcRepresentation(),
		Equip:   c.EquippedItems.NcRepresentation(),
		TutorialInfo: structs.ProtoTutorialInfo{ // x(
			TutorialState: 2,
			TutorialStep:  byte(0),
		},
	}
	return nc
}

// NcRepresentation returns a struct that can be serialized into bytes and can be sent to the client
func (cei *EquippedItems) NcRepresentation() structs.ProtoEquipment {
	return structs.ProtoEquipment{
		EquHead:         cei.Head,
		EquMouth:        cei.ApparelFace,
		EquRightHand:    cei.RightHand,
		EquBody:         cei.Body,
		EquLeftHand:     cei.LeftHand,
		EquPant:         cei.Pants,
		EquBoot:         cei.Boots,
		EquAccBoot:      cei.ApparelBoots,
		EquAccPant:      cei.ApparelPants,
		EquAccBody:      cei.ApparelBody,
		EquAccHeadA:     cei.ApparelHead,
		EquMinimonR:     cei.RightMiniPet,
		EquEye:          cei.Face,
		EquAccLeftHand:  cei.ApparelLeftHand,
		EquAccRightHand: cei.ApparelRightHand,
		EquAccBack:      cei.ApparelBack,
		EquCosEff:       cei.ApparelAura,
		EquAccHip:       cei.ApparelTail,
		EquMinimon:      cei.LeftMiniPet,
		EquAccShield:    cei.ApparelShield,
		Upgrade:         structs.EquipmentUpgrade{},
	}
}

// NcRepresentation returns a struct that can be serialized into bytes and can be sent to the client
func (ca *Appearance) NcRepresentation() structs.ProtoAvatarShapeInfo {
	return structs.ProtoAvatarShapeInfo{
		BF:        1 | ca.Class<<2 | ca.Gender<<7,
		HairType:  ca.HairType,
		HairColor: ca.HairColor,
		FaceShape: ca.FaceType,
	}
}
