package character

import (
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

var log     *logger.Logger

func init() {
	log = logger.Init("NetworkingLogger", true, false, ioutil.Discard)
	log.Info("networking logger init()")
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
	Inventory     *Inventory
	EquippedItems *EquippedItems
	AdminLevel    uint8     `pg:",notnull,use_zero"`
	Slot          uint8     `pg:",notnull,use_zero"`
	IsDeleted     bool      `pg:",use_zero"`
	DeletedAt     time.Time `pg:",soft_delete"`
}

// Appearance model for the database layer
type Appearance struct {
	tableName   struct{} `pg:"world.character_appearance"`
	ID          uint64
	CharacterID uint64 //
	Character   *Character
	Class       uint8     `pg:",notnull"`
	Gender      uint8     `pg:",notnull,use_zero"`
	HairType    uint8     `pg:",notnull,use_zero"`
	HairColor   uint8     `pg:",notnull,use_zero"`
	FaceType    uint8     `pg:",notnull,use_zero"`
	DeletedAt   time.Time `pg:",soft_delete"`
}

// Attributes model for the database layer
type Attributes struct {
	tableName    struct{} `pg:"world.character_attributes"`
	ID           uint64
	CharacterID  uint64
	Character    *Character
	Level        uint16    `pg:",notnull"`
	Experience   uint64    `pg:",notnull,use_zero"`
	Fame         uint64    `pg:",notnull,use_zero"`
	Hp           uint32    `pg:",notnull"`
	Sp           uint32    `pg:",notnull"`
	Intelligence uint32    `pg:",notnull"`
	Strength     uint32    `pg:",notnull"`
	Dexterity    uint32    `pg:",notnull"`
	Endurance    uint32    `pg:",notnull"`
	Spirit       uint32    `pg:",notnull"`
	Money        uint64    `pg:",notnull,use_zero"`
	KillPoints   uint32    `pg:",notnull,use_zero"`
	HpStones     uint16    `pg:",notnull"`
	SpStones     uint16    `pg:",notnull"`
	DeletedAt    time.Time `pg:",soft_delete"`
}

// Location model for the database layer
type Location struct {
	tableName   struct{} `pg:"world.character_location"`
	ID          uint64
	CharacterID uint64 //
	Character   *Character
	MapName     string    `pg:",notnull"`
	X           uint32    `pg:",notnull"`
	Y           uint32    `pg:",notnull"`
	D           uint32    `pg:",notnull"`
	IsKQ        bool      `pg:",notnull,use_zero"`
	DeletedAt   time.Time `pg:",soft_delete"`
}

// Inventory model for the database layer
type Inventory struct {
	tableName   struct{} `pg:"world.character_inventory"`
	ID          uint64
	CharacterID uint64 //
	Character   *Character
	ShnID       uint16    `pg:",notnull"`
	Slot        uint16    `pg:",notnull,use_zero"`
	IsStack     bool      `pg:",notnull,use_zero"`
	IsStored    bool      `pg:",notnull,use_zero"`
	FromMonarch bool      `pg:",notnull,use_zero"`
	FromStore   bool      `pg:",notnull,use_zero"`
	DeletedAt   time.Time `pg:",soft_delete"`
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

var ErrInvalidSlot = &ErrCharacter{
	Code:    0,
	Message: "invalid slot",
}
var ErrNameTaken = &ErrCharacter{
	Code:    1,
	Message: "name taken",
}
var ErrNoSlot = &ErrCharacter{
	Code:    2,
	Message: "no slot available",
}
var ErrInvalidName = &ErrCharacter{
	Code:    3,
	Message: "invalid name",
}
var ErrInvalidClassGender = &ErrCharacter{
	Code:    4,
	Message: "invalid class gender data",
}

type ErrCharacter struct {
	Code    int
	Message string
}

func (ec *ErrCharacter) Error() string {
	return ec.Message
}

const (
	startLevel = 1
	startMap   = "Rou"
)

func Validate(db *pg.DB, userID uint64, req structs.NcAvatarCreateReq) error {

	if req.SlotNum > 5 {
		return ErrInvalidSlot
	}

	name := strings.TrimRight(string(req.Name.Name[:]), "\x00")

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

func New(db *pg.DB, userID uint64, req structs.NcAvatarCreateReq) (structs.AvatarInformation, error) {
	//err := Validate(db, userID, req)
	//if err != nil {
	//	switch err.Error() {
	//	case ErrInvalidName.Error():
	//		// name taken error
	//		break
	//	case ErrNameTaken.Error():
	//		// errcode 385
	//		break
	//	case ErrInvalidSlot.Error():
	//		break
	//	case ErrInvalidClassGender.Error():
	//		break
	//	case ErrNoSlot.Error():
	//		break
	//	}
	//	return structs.AvatarInformation{}, err
	//}
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
		newTx.Rollback()
		return structs.AvatarInformation{}, err
	}

	char.
		initialAppearance(req.Shape).
		initialAttributes().
		initialLocation().
		initialInventory().
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

	if _, err = newTx.Model(char.Inventory).Returning("*").Insert(); err != nil {
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
	return char.ncRepresentation(), newTx.Commit()
}

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

	if _, err = deleteTx.Model(char.Inventory).Where("character_id = ?", char.ID).Delete(); err != nil {
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
	c.Location = &Location{
		CharacterID: c.ID,
		MapName:     startMap,
		X:           3000,
		Y:           11666,
		D:           90,
		IsKQ:        false,
	}
	return c
}

func (c *Character) initialInventory() *Character {
	c.Inventory = &Inventory{
		CharacterID: c.ID,
		ShnID:       1,
		Slot:        0,
		IsStack:     false,
		IsStored:    false,
		FromMonarch: false,
		FromStore:   false,
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

func (c *Character) ncRepresentation() structs.AvatarInformation {
	var name [20]byte
	var mapName [12]byte
	copy(name[:], c.Name)
	copy(mapName[:], c.Location.MapName)

	nc := structs.AvatarInformation{
		ChrRegNum: uint32(c.ID),
		Name: structs.Name5{
			Name: name,
		},
		Level: c.Attributes.Level,
		Slot:  c.Slot,
		LoginMap: structs.Name3{
			Name: mapName,
		},
		DelInfo: structs.ProtoAvatarDeleteInfo{},
		Shape:   c.Appearance.ncRepresentation(),
		Equip:   c.EquippedItems.ncRepresentation(),
		TutorialInfo: structs.ProtoTutorialInfo{ // x(
			TutorialState: 2,
			TutorialStep:  byte(0),
		},
	}
	return nc
}

func (cei *EquippedItems) ncRepresentation() structs.ProtoEquipment {
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

func (ca *Appearance) ncRepresentation() structs.ProtoAvatarShapeInfo {
	return structs.ProtoAvatarShapeInfo{
		BF:        1 | ca.Class<<2 | ca.Gender<<7,
		HairType:  ca.HairType,
		HairColor: ca.HairColor,
		FaceShape: ca.FaceType,
	}
}
