package persistence

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"regexp"
	"time"

	"github.com/go-pg/pg/v10"

	"github.com/google/logger"
	"github.com/google/uuid"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
)

const (
	genericIDEquals = "id = ?"

	characterSlotEquals   = "slot = ?"
	characterUserIDEquals = "user_id = ?"
	characterIDEquals     = "character_id = ?"
	characterNameEquals   = "name = ?"

	itemIDEquals            = "item_id = ?"
	itemShnIDEquals         = "shn_id = ?"
	itemInventoryTypeEquals = "inventory_type = ?"
	itemSlotEquals          = "slot = ?"
)

var log *logger.Logger

func init() {
	log = logger.Init("persistence logger", true, false, ioutil.Discard)
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
	UserID        uint64         `pg:",notnull"`
	Name          string         `pg:",notnull,unique"`
	Appearance    *Appearance    `pg:"rel:belongs-to"`
	Attributes    *Attributes    `pg:"rel:belongs-to"`
	Location      *Location      `pg:"rel:belongs-to"`
	Options       *ClientOptions `pg:"rel:belongs-to"`
	Items         []*Item        `pg:"rel:has-many"`
	EquippedItems *EquippedItems `pg:"rel:belongs-to"`
	AdminLevel    uint8          `pg:",notnull,use_zero"`
	Slot          uint8          `pg:",notnull,use_zero"`
	IsDeleted     bool           `pg:",use_zero"`
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
	X           int    `pg:",notnull"`
	Y           int    `pg:",notnull"`
	D           int    `pg:",notnull,use_zero"`
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
	DeletedAt   time.Time `pg:",soft_delete"`
}

const (
	// move to global config
	startLevel   = 1
	startMapID   = 1
	startMapName = "Rou"
)

// ValidateCharacter checks data sent by the client is valid
func ValidateCharacter(userID uint64, req *structs.NcAvatarCreateReq) error {
	if req.SlotNum > 5 {
		return errors.Err{
			Code: errors.PersistenceCharInvalidSlot,
			Details: errors.Details{
				"userID":   userID,
				"slotNum":  req.SlotNum,
				"charName": req.Name.Name,
			},
		}
	}

	name := req.Name.Name

	var charName string
	err := db.Model((*Character)(nil)).Column("name").Where(characterNameEquals, name).Select(&charName)

	if err == nil {
		// return ErrNameTaken
		return errors.Err{
			Code: errors.PersistenceCharNameTaken,
			Details: errors.Details{
				"userID":   userID,
				"charName": req.Name.Name,
			},
		}
	}

	var chars []Character
	err = db.Model(&chars).Where(characterUserIDEquals, userID).Select()

	if len(chars) == 6 {
		return errors.Err{
			Code: errors.PersistenceCharNoSlot,
			Details: errors.Details{
				"userID": userID,
			},
		}
	}

	alphaNumeric := regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString
	if !alphaNumeric(name) {
		return errors.Err{
			Code: errors.PersistenceCharInvalidName,
			Details: errors.Details{
				"userID":   userID,
				"charName": req.Name.Name,
			},
		}
	}

	// todo: missing validation: default hair, color, face values
	// todo: missing validation: default starter class values ( mage, cleric, archer, fighter)
	isMale := (req.Shape.BF >> 7) & 1
	class := (req.Shape.BF >> 2) & 31

	if isMale > 1 || isMale < 0 {
		return errors.Err{
			Code: errors.PersistenceCharInvalidClassGender,
			Details: errors.Details{
				"userID":   userID,
				"charName": req.Name.Name,
				"bfValue":  req.Shape.BF,
				"isMale":   isMale,
				"class":    class,
			},
		}
	}

	if class < 1 || class > 27 {
		return errors.Err{
			Code: errors.PersistenceCharInvalidClassGender,
			Details: errors.Details{
				"userID":   userID,
				"charName": req.Name.Name,
				"bfValue":  req.Shape.BF,
				"isMale":   isMale,
				"class":    class,
			},
		}
	}

	return nil
}

// NewCharacter creates character for the User with userID and returns data the client can understand
func NewCharacter(userID uint64, req *structs.NcAvatarCreateReq, initialItems bool) (*Character, error) {
	char := &Character{
		UserID:     userID,
		AdminLevel: 0,
		Name:       req.Name.Name,
		Slot:       req.SlotNum,
	}
	tx, err := db.Begin()
	if err != nil {
		return char, err
	}

	defer closeTx(tx)

	_, err = tx.Model(char).Returning("*").Insert()

	if err != nil {
		return char, persistenceError(err, tx)
	}

	char.initialAppearance(req.Shape)
	char.initialAttributes()
	char.initialLocation()
	char.initialClientOptions()
	char.initialEquippedItems()

	_, err = tx.Model(char.Appearance).Returning("*").Insert()
	if err != nil {
		return char, persistenceError(err, tx)
	}

	_, err = tx.Model(char.Attributes).Returning("*").Insert()
	if err != nil {
		return char, persistenceError(err, tx)
	}

	_, err = tx.Model(char.Location).Returning("*").Insert()
	if err != nil {
		return char, persistenceError(err, tx)
	}

	_, err = tx.Model(char.Options).Returning("*").Insert()
	if err != nil {
		return char, persistenceError(err, tx)
	}

	_, err = tx.Model(char.EquippedItems).Returning("*").Insert()
	if err != nil {
		return char, persistenceError(err, tx)
	}

	if initialItems {
		err = tx.Commit()
		if err != nil {
			return char, persistenceError(err, tx)
		}

		char.initialItems()

		itx, err := db.Begin()
		if err != nil {
			return char, err
		}

		defer closeTx(itx)

		for _, item := range char.Items {
			err := item.Insert()
			if err != nil {
				return char, persistenceError(err, tx)
			}
		}
		return char, itx.Commit()
	}

	return char, tx.Commit()
}

func persistenceError(err error, tx *pg.Tx) error {
	return errors.Err{
		Code: errors.PersistenceErrDB,
		Details: errors.Details{
			"err":   err,
			"txErr": tx.Rollback(),
		},
	}
}

func GetCharacter(characterID uint64) (Character, error) {
	var c Character
	c.ID = characterID
	err := db.Model(&c).
		WherePK().
		Relation("Appearance").
		Relation("Attributes").
		Relation("Items").
		Relation("Options").
		Relation("Location").
		Select()
	return c, err
}

func GetCharacterByName(name string) (Character, error) {
	var c Character
	err := db.Model(&c).
		Where("name = ?", name).
		Select() // query the world for a character with name

	err = db.Model(&c).
		WherePK().
		Relation("Appearance").
		Relation("Attributes").
		Relation("Items").
		Relation("Options").
		Relation("Location").
		Select() // query the world for a character with name

	return c, err
}

func GetCharacterBySlot(slot byte, userID uint64) (Character, error) {
	var c Character
	err := db.Model(&c).
		Relation("Appearance").
		Relation("Attributes").
		Relation("Items").
		Relation("Options").
		Relation("Location").
		Where(characterUserIDEquals, userID).
		Where(characterSlotEquals, slot).Select()
	return c, err
}

func UserCharacters(id uint64) ([]*Character, error) {
	var chars []*Character

	err := db.Model(&chars).
		Relation("Appearance").
		Relation("Attributes").
		Relation("Location").
		Relation("EquippedItems").
		Where("user_id = ?", id).
		Select()

	return chars, err
}

func UpdateCharacter(c *Character) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer closeTx(tx)

	_, err = tx.Model(c).
		WherePK().Returning("*").Update()

	_, err = tx.Model(c.Appearance).
		WherePK().Returning("*").Update()

	_, err = tx.Model(c.Attributes).
		WherePK().Returning("*").Update()

	_, err = tx.Model(c.Location).
		WherePK().Returning("*").Update()

	_, err = tx.Model(c.Options).
		WherePK().Returning("*").Update()

	if err != nil {
		return tx.Rollback()
	}

	return tx.Commit()
}

// todo: switch to method
func UpdateLocation(c *Character) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer closeTx(tx)

	_, err = tx.Model(c.Location).
		WherePK().Returning("*").Update()

	if err != nil {
		return errors.Err{
			Code: errors.PersistenceErrDB,
			Details: errors.Details{
				"err":   err,
				"txErr": tx.Rollback(),
			},
		}
	}

	return tx.Commit()
}

// DeleteCharacter todo: switch to method
func DeleteCharacter(userID uint64, slot int) error {
	tx, err := db.Begin()
	if err != nil {
		return persistenceError(err, tx)
	}

	defer closeTx(tx)

	var char Character
	err = tx.Model(&char).Where(characterUserIDEquals, userID).Where(characterSlotEquals, slot).Select()

	if err != nil {
		return persistenceError(err, tx)
	}

	name := fmt.Sprintf("%v@%v", char.Name, uuid.New().String())
	_, err = tx.Model((*Character)(nil)).Set(characterNameEquals, name).Where(characterUserIDEquals, userID).Where(characterSlotEquals, slot).Update()
	if err != nil {
		return persistenceError(err, tx)
	}

	_, err = tx.Model(&char).Where(characterUserIDEquals, userID).Where(characterSlotEquals, slot).Delete()
	if err != nil {
		return persistenceError(err, tx)
	}

	_, err = tx.Model(char.Appearance).Where(characterIDEquals, char.ID).Delete()
	if err != nil {
		return persistenceError(err, tx)
	}

	_, err = tx.Model(char.Attributes).Where(characterIDEquals, char.ID).Delete()
	if err != nil {
		return persistenceError(err, tx)
	}

	_, err = tx.Model(char.Location).Where(characterIDEquals, char.ID).Delete()
	if err != nil {
		return persistenceError(err, tx)
	}

	_, err = tx.Model(char.EquippedItems).Where(characterIDEquals, char.ID).Delete()
	if err != nil {
		return persistenceError(err, tx)
	}

	return tx.Commit()
}

func (c *Character) initialAppearance(shape structs.ProtoAvatarShapeInfo) {
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
}

func (c *Character) initialAttributes() {
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
}

func (c *Character) initialLocation() {
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
}

func (c *Character) initialClientOptions() {
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
}

// TODO: remove this, as I will be using the equipped inventory to send this data
func (c *Character) initialEquippedItems() {
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
}

func (c *Character) initialItems() {
	var (
		shnID  uint16 = 0
		shnInx        = ""
		slot   int
	)

	switch c.Appearance.Class {
	case 1: // fighter
		shnID = 250
		shnInx = "ShortSword"
		slot = 12
		break

	case 6: // cleric
		shnID = 750
		shnInx = "ShortMace"
		slot = 12
		break

	case 11: // archer
		shnID = 1250
		shnInx = "ShortBow"
		slot = 10
		break

	case 16: // mage
		shnID = 1750
		shnInx = "ShortStaff"
		slot = 12
		break
	}

	item := &Item{
		InventoryType: int(EquippedInventory),
		Slot:          slot,
		CharacterID:   c.ID,
		Character:     c,
		ShnID:         shnID,
		ShnInxName:    shnInx,
		Stackable:     false,
		Amount:        1,
	}
	c.Items = append(c.Items, item)
}
