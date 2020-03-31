package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/shine-o/shine.engine.networking"
	"github.com/shine-o/shine.engine.networking/structs"
	"regexp"
	"strings"
	"time"
)

// User model for schema: accounts
type Character struct {
	tableName     struct{} `pg:"world.characters"`
	ID            uint64
	UserID        uint64 `pg:",notnull"`
	Name          string `pg:",notnull,unique"`
	Appearance    *CharacterAppearance
	Attributes    *CharacterAttributes
	Location      *CharacterLocation
	Inventory     *CharacterInventory
	EquippedItems *CharacterEquippedItems
	AdminLevel    uint8     `pg:",notnull,use_zero"`
	Slot          uint8     `pg:",notnull,use_zero"`
	IsDeleted     bool      `pg:",use_zero"`
	DeletedAt     time.Time `pg:",soft_delete"`
}
// todo: checkout how these hooks really work
var _ pg.BeforeDeleteHook = (*Character)(nil)
func (c * Character) BeforeDelete(ctx context.Context)  (context.Context, error) {
	c.Name = fmt.Sprintf("%v@%v", c.Name, "12346")
	_, err := worldDB.Model(c).Set("name = ?name").Where("id = ?id").Update()
	if err != nil {
		log.Error(err)
		return ctx, err
	}
	return ctx, nil
}

type CharacterAppearance struct {
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

type CharacterAttributes struct {
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

type CharacterLocation struct {
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

type CharacterInventory struct {
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

type CharacterEquippedItems struct {
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

var errNoSession = errors.New("no login session was found")
var errNoSlot = errors.New("no slot available")
var errInvalidSlot = errors.New("invalid slot")
var errInvalidName = errors.New("invalid name")
var errInvalidRequest = errors.New("invalid packet data was sent")
var errNameTaken = errors.New("name already in use")
var errInvalidClassGender = errors.New("invalid Class Gender data")

type errCharacter struct {
	Code int
	Message string
}

func (ec *errCharacter) Error() string  {
	return ec.Message
}


const (
	startLevel = 1
	startMap   = "Rou"
)

func newCharacter(ctx context.Context, req structs.NcAvatarCreateReq) (structs.AvatarInformation, error) {
	select {
	case <-ctx.Done():
		return structs.AvatarInformation{}, errCC
	default:
		err := validateCharacter(ctx, req)
		if err != nil {

			switch err.Error() {
			case errInvalidName.Error():
				// name taken error
				break
			case errNameTaken.Error():
				// errcode 385
				break
			case errInvalidSlot.Error():
				break
			case errInvalidClassGender.Error():
				break
			case errNoSlot.Error():
				break
			}

			return structs.AvatarInformation{}, err
		}

		newCharacterTx, err := worldDB.Begin()

		if err != nil {
			return structs.AvatarInformation{}, err
		}

		defer newCharacterTx.Close()

		wsi := ctx.Value(networking.ShineSession)
		ws := wsi.(*session)

		name := strings.TrimRight(string(req.Name.Name[:]), "\x00")
		char := Character{
			UserID:     ws.UserID,
			AdminLevel: 0,
			Name:       name,
			Slot:       req.SlotNum,
		}

		_, err = newCharacterTx.Model(&char).Returning("*").Insert()

		if err != nil {
			newCharacterTx.Rollback()
			return structs.AvatarInformation{}, err
		}

		char.
			initialAppearance(req.Shape).
			initialAttributes().
			initialLocation().
			initialInventory().
			initialEquippedItems()

		if _, err = newCharacterTx.Model(char.Appearance).Returning("*").Insert(); err != nil {
			err := newCharacterTx.Rollback()
			return structs.AvatarInformation{}, err
		}

		if _, err = newCharacterTx.Model(char.Attributes).Returning("*").Insert(); err != nil {
			err := newCharacterTx.Rollback()
			return structs.AvatarInformation{}, err
		}

		if _, err = newCharacterTx.Model(char.Location).Returning("*").Insert(); err != nil {
			err := newCharacterTx.Rollback()
			return structs.AvatarInformation{}, err
		}

		if _, err = newCharacterTx.Model(char.Inventory).Returning("*").Insert(); err != nil {
			err := newCharacterTx.Rollback()
			return structs.AvatarInformation{}, err
		}

		if _, err = newCharacterTx.Model(char.EquippedItems).Returning("*").Insert(); err != nil {
			err := newCharacterTx.Rollback()
			return structs.AvatarInformation{}, err
		}

		err = newCharacterTx.Commit()

		if err != nil {
			return structs.AvatarInformation{}, err
		}

		return char.ncRepresentation(), nil
	}
}

func deleteCharacter(ctx context.Context, req structs.NcAvatarEraseReq) error {
	wsi := ctx.Value(networking.ShineSession)
	ws := wsi.(*session)
	// ws.
	deleteCharTx, err := worldDB.Begin()
	defer deleteCharTx.Close()
	if err != nil {
		return &errCharacter{
			Code:    1,
			Message: fmt.Sprintf("database error, could not start transaction: %v", err),
		}
	}
	var char Character
	err = deleteCharTx.Model(&char).Where("user_id = ?", ws.UserID).Where("slot = ?", req.Slot).Select()

	if err != nil {
		return &errCharacter{
			Code:    2,
			Message: fmt.Sprintf("database error, character not found: %v", err),
		}
	}

	_, err = deleteCharTx.Model(&char).Where("user_id = ?", ws.UserID).Where("slot = ?", req.Slot).Delete()
	if err != nil {
		return &errCharacter{
			Code:    3,
			Message: fmt.Sprintf("database error, failed to delete row: %v", err),
		}
	}

	name := fmt.Sprintf("%v@%v", char.Name, "12346")
	res, err :=  deleteCharTx.Model((*Character)(nil)).Set("name = ?", name).Where("user_id = ?", ws.UserID).Update()

	log.Info(res)

	_, err = deleteCharTx.Model(char.Appearance).Where("character_id = ?", char.ID).Delete()
	if err != nil {
		return &errCharacter{
			Code:    3,
			Message: fmt.Sprintf("database error, failed to delete row: %v", err),
		}
	}

	_, err = deleteCharTx.Model(char.Attributes).Where("character_id = ?", char.ID).Delete()
	if err != nil {
		return &errCharacter{
			Code:    3,
			Message: fmt.Sprintf("database error, failed to delete row: %v", err),
		}
	}

	_, err = deleteCharTx.Model(char.Location).Where("character_id = ?", char.ID).Delete()
	if err != nil {
		return &errCharacter{
			Code:    3,
			Message: fmt.Sprintf("database error, failed to delete row: %v", err),
		}
	}

	_, err = deleteCharTx.Model(char.Inventory).Where("character_id = ?", char.ID).Delete()
	if err != nil {
		return &errCharacter{
			Code:    3,
			Message: fmt.Sprintf("database error, failed to delete row: %v", err),
		}
	}

	_, err = deleteCharTx.Model(char.EquippedItems).Where("character_id = ?", char.ID).Delete()

	if err != nil {
		return &errCharacter{
			Code:    3,
			Message: fmt.Sprintf("database error, failed to delete row: %v", err),
		}
	}

	err = deleteCharTx.Commit()

	if err != nil {
		return &errCharacter{
			Code:    3,
			Message: fmt.Sprintf("database error, could not commit transaction: %v", err),
		}
	}

	return nil
}



func validateCharacter(ctx context.Context, req structs.NcAvatarCreateReq) error {
	// fetch session
	wsi := ctx.Value(networking.ShineSession)
	ws := wsi.(*session)

	if req.SlotNum > 5 {
		return errInvalidSlot
	}

	name := strings.TrimRight(string(req.Name.Name[:]), "\x00")

	var charName string
	err := worldDB.Model((*Character)(nil)).Column("name").Where("name = ?", name).Select(&charName)

	if err == nil {
		log.Error(charName)
		return errInvalidName
	}

	if charName == name {
		return errNameTaken
	}

	var chars []Character
	err = worldDB.Model(&chars).Where("user_id = ?", ws.UserID).Select()

	if len(chars) == 6 {
		return errNoSlot
	}

	alphaNumeric := regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString
	if !alphaNumeric(name) {
		return errInvalidName
	}

	isMale := (req.Shape.BF >> 7) & 1
	class := (req.Shape.BF >> 2) & 31

	if isMale > 1 || isMale < 0 {
		return errInvalidClassGender
	}

	if class < 1 || class > 27 {
		return errInvalidClassGender
	}

	return nil
}

func (c *Character) initialAppearance(shape structs.ProtoAvatarShapeInfo) *Character {
	isMale := (shape.BF >> 7) & 1
	class := (shape.BF >> 2) & 31

	c.Appearance = &CharacterAppearance{
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
	c.Attributes = &CharacterAttributes{
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
	c.Location = &CharacterLocation{
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
	c.Inventory = &CharacterInventory{
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
	c.EquippedItems = &CharacterEquippedItems{
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

func (cei *CharacterEquippedItems) ncRepresentation() structs.ProtoEquipment {
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

func (ca *CharacterAppearance) ncRepresentation() structs.ProtoAvatarShapeInfo {
	return structs.ProtoAvatarShapeInfo{
		BF:        1 | ca.Class<<2 | ca.Gender<<7,
		HairType:  ca.HairType,
		HairColor: ca.HairColor,
		FaceShape: ca.FaceType,
	}
}
