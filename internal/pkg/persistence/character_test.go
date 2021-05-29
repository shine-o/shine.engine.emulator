package persistence

import (
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/google/logger"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/database"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
)

const (
	unexpectedErrorType    = "expected %v, got %v"
	unexpectedErrorCode    = "expected error code %v, got %v"
	nilError               = "expected an error but got nil"
	nilItem                = "no item"
	unexpectedEquippedItem = "id =%v, expected id =%v"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	log = logger.Init("test logger", true, false, ioutil.Discard)
	log.Info("test logger")

	InitDB(database.ConnectionParams{
		User:     "user",
		Password: "password",
		Host:     "127.0.0.1",
		Port:     "54320",
		Database: "shine",
		Schema:   "world",
	})

	err := database.CreateSchema(db, "world")
	if err != nil {
		log.Fatal(err)
	}

	cleanDB()

	os.Exit(m.Run())
}

// test character data is valid
func TestValidateCharacterRequest(t *testing.T) {
	// {"packetType":"small","length":27,"department":5,"command":"1","opCode":5121,"data":"0046696768747265726f6f0000000000000000000085060000","rawData":"1b01140046696768747265726f6f0000000000000000000085060000","friendlyName":"NC_AVATAR_CREATE_REQ"}
	if data, err := hex.DecodeString("0046696768747265726f6f0000000000000000000085060000"); err != nil {
		t.Error(err)
	} else {
		nc := structs.NcAvatarCreateReq{}
		err := structs.Unpack(data, &nc)
		if err != nil {
			t.Error(err)
		}
		err = ValidateCharacter(1, &nc)
		if err != nil {
			t.Error(err)
		}
	}
}

// foolishly assuming data is okay
func TestCreateCharacter(t *testing.T) {
	cleanDB()
	// {"packetType":"small","length":27,"department":5,"command":"1","opCode":5121,"data":"0046696768747265726f6f0000000000000000000085060000","rawData":"1b01140046696768747265726f6f0000000000000000000085060000","friendlyName":"NC_AVATAR_CREATE_REQ"}
	if data, err := hex.DecodeString("0046696768747265726f6f0000000000000000000085060000"); err != nil {
		t.Error(err)
	} else {
		nc := structs.NcAvatarCreateReq{}
		err := structs.Unpack(data, &nc)
		if err != nil {
			t.Error(err)
		}
		_, err = NewCharacter(1, &nc, false)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestDeleteCharacter(t *testing.T) {
	cleanDB()

	createDummyCharacters()

	// {"packetType":"small","length":3,"department":5,"command":"7","opCode":5127,"data":"03","rawData":"03071403","friendlyName":"NC_AVATAR_ERASE_REQ"}
	nc := structs.NcAvatarEraseReq{
		Slot: 0,
	}

	err := DeleteCharacter(1, int(nc.Slot))
	if err != nil {
		log.Error(err)
	}

	// try and fetch a character for this user id in the current slot.
	var deletedAt time.Time
	err = db.Model((*Character)(nil)).
		Column("deleted_at").
		Deleted().
		Where("user_id = ?", 1).
		Where("slot = ?", 0).
		Select(&deletedAt)

	if err != nil {
		t.Error(err)
	}
}

func TestCharacterNameInUseError(t *testing.T) {
	cleanDB()

	createDummyCharacters()

	name := fmt.Sprintf("mob%v", 1)
	c := structs.NcAvatarCreateReq{
		SlotNum: byte(0),
		Name: structs.Name5{
			Name: name,
		},
		Shape: structs.ProtoAvatarShapeInfo{
			BF:        133,
			HairType:  6,
			HairColor: 0,
			FaceShape: 0,
		},
	}
	err := ValidateCharacter(1, &c)
	if err == nil {
		log.Error(err)
	}

	errChar, ok := err.(errors.Err)
	if !ok {
		t.Errorf(unexpectedErrorType, "errors.Err", reflect.TypeOf(errChar))
	}
	if errChar.Code != errors.PersistenceCharNameTaken {
		t.Errorf(unexpectedErrorCode, errors.PersistenceCharNameTaken, errChar.Code)
	}
}

func TestNoSlotAvailableError(t *testing.T) {
	cleanDB()
	createDummyCharacters()

	// {"packetType":"small","length":27,"department":5,"command":"1","opCode":5121,"data":"0046696768747265726f6f0000000000000000000085060000","rawData":"1b01140046696768747265726f6f0000000000000000000085060000","friendlyName":"NC_AVATAR_CREATE_REQ"}
	if data, err := hex.DecodeString("0046696768747265726f6f0000000000000000000085060000"); err != nil {
		t.Error(err)
	} else {
		nc := structs.NcAvatarCreateReq{}
		err := structs.Unpack(data, &nc)
		if err != nil {
			t.Error(err)
		}

		err = ValidateCharacter(1, &nc)
		if err != nil {
			cErr, ok := err.(errors.Err)
			if !ok {
				t.Errorf(unexpectedErrorType, "errors.Err", reflect.TypeOf(cErr))
			}
			if cErr.Code != errors.PersistenceCharNoSlot {
				t.Errorf(unexpectedErrorCode, errors.PersistenceCharNoSlot, cErr.Code)
			}
			return
		}
		t.Error(nilError)
	}
}

func TestInvalidSlotError(t *testing.T) {
	// {"packetType":"small","length":27,"department":5,"command":"1","opCode":5121,"data":"0946696768747265726f6f0000000000000000000085060000","rawData":"1b01140046696768747265726f6f0000000000000000000085060000","friendlyName":"NC_AVATAR_CREATE_REQ"}
	if data, err := hex.DecodeString("0946696768747265726f6f0000000000000000000085060000"); err != nil {
		t.Error(err)
	} else {
		nc := structs.NcAvatarCreateReq{}
		err := structs.Unpack(data, &nc)
		if err != nil {
			t.Error(err)
		}

		err = ValidateCharacter(1, &nc)
		if err != nil {
			cErr, ok := err.(errors.Err)
			if !ok {
				t.Error("unexpected error type")
			}
			if cErr.Code != errors.PersistenceCharInvalidSlot {
				t.Errorf(unexpectedErrorCode, errors.PersistenceCharInvalidSlot, cErr.Code)
			}
			return
		}
		t.Error(nilError)
	}
}

func TestInvalidNameError(t *testing.T) {
	cleanDB()

	// {"packetType":"small","length":27,"department":5,"command":"1","opCode":5121,"data":"0046325B5D747265726F6F0000000000000000000085060000","rawData":"1b01140046325B5D747265726F6F0000000000000000000085060000","friendlyName":"NC_AVATAR_CREATE_REQ"}
	if data, err := hex.DecodeString("0046325B5D747265726F6F0000000000000000000085060000"); err != nil {
		t.Error(err)
	} else {
		nc := structs.NcAvatarCreateReq{}
		err := structs.Unpack(data, &nc)
		if err != nil {
			t.Error(err)
		}

		err = ValidateCharacter(1, &nc)
		if err != nil {
			cErr, ok := err.(errors.Err)
			if !ok {
				t.Error("unexpected error type")
			}
			if cErr.Code != errors.PersistenceCharInvalidName {
				t.Errorf(unexpectedErrorCode, errors.PersistenceCharInvalidName, cErr.Code)
			}
			return
		}
		t.Error(nilError)
	}
}

// KOREAN MONKEYS!
// test that correct gender and class are extracted using binary operators
func TestInvalidGenderClassBinaryOperation(t *testing.T) {
	cleanDB()

	// {"packetType":"small","length":27,"department":5,"command":"1","opCode":5121,"data":"0946696768747265726f6f0000000000000000000085060000","rawData":"1b01140046325B5D747265726F6F0000000000000000000085060000","friendlyName":"NC_AVATAR_CREATE_REQ"}
	if data, err := hex.DecodeString("0246696768747265726f6f00000000000000000000ff060000"); err != nil {
		t.Error(err)
	} else {
		nc := structs.NcAvatarCreateReq{}
		err := structs.Unpack(data, &nc)
		if err != nil {
			t.Error(err)
		}

		err = ValidateCharacter(1, &nc)
		if err != nil {
			errChar, ok := err.(errors.Err)
			if !ok {
				t.Errorf(unexpectedErrorType, "errors.Err", reflect.TypeOf(errChar))
			}
			if errChar.Code != errors.PersistenceCharInvalidClassGender {
				t.Errorf(unexpectedErrorCode, errors.PersistenceCharInvalidClassGender, errChar.Code)
			}
		} else {
			t.Error(nilError)
		}
	}
}

func TestNewCharacterDefaultItems(t *testing.T) {
	cleanDB()
	createDummyCharacters()
	// assert user has an inventory
	characters, err := getCharacters(false)
	if err != nil {
		t.Error(err)
	}

	// 31000	House_MushRoom	Mushroom House
	var miniHouseID uint16 = 31000
	for _, character := range characters {

		clauses := make(map[string]interface{})

		clauses[characterIDEquals] = character.ID
		clauses[itemShnIDEquals] = miniHouseID

		item, err := GetItemWhere(clauses, false)
		if err != nil {
			t.Error(err)
		}

		if item == nil {
			t.Error(nilItem)
		}
	}
}

func TestLoadNewCharacterMageEquippedItems(t *testing.T) {
	cleanDB()
	// should have 1 staff
	// 1750	ShortStaff	Short Staff
	var rightHand uint16 = 1750
	character := newCharacter("mage")

	if character.EquippedItems.RightHand != rightHand {
		t.Errorf(unexpectedEquippedItem, character.EquippedItems.RightHand, rightHand)
	}

	clauses := make(map[string]interface{})

	clauses[characterIDEquals] = character.ID
	clauses[itemShnIDEquals] = rightHand

	_, err := GetItemWhere(clauses, false)
	if err != nil {
		t.Error(err)
	}
}

func TestLoadNewCharacterFighterEquippedItems(t *testing.T) {
	cleanDB()
	// 250	ShortSword	Short Sword
	// bitField := 1 | 1 << 2 | 1 << 7

	var rightHand uint16 = 250
	character := newCharacter("fighter")

	if character.EquippedItems.RightHand != rightHand {
		t.Errorf(unexpectedEquippedItem, character.EquippedItems.RightHand, rightHand)
	}

	clauses := make(map[string]interface{})

	clauses[characterIDEquals] = character.ID
	clauses[itemShnIDEquals] = rightHand

	_, err := GetItemWhere(clauses, false)
	if err != nil {
		t.Error(err)
	}
}

func TestLoadNewCharacterArcherEquippedItems(t *testing.T) {
	cleanDB()
	// 1250	ShortBow	Short Bow
	// bitField := 1 | 11 << 2 | 1 << 7
	var rightHand uint16 = 1250
	character := newCharacter("archer")

	if character.EquippedItems.RightHand != rightHand {
		t.Errorf(unexpectedEquippedItem, character.EquippedItems.RightHand, rightHand)
	}

	clauses := make(map[string]interface{})

	clauses[characterIDEquals] = character.ID
	clauses[itemShnIDEquals] = rightHand

	_, err := GetItemWhere(clauses, false)
	if err != nil {
		t.Error(err)
	}
}

func TestLoadNewCharacterClericEquippedItems(t *testing.T) {
	cleanDB()
	// 750	ShortMace	Short Mace
	// bitField := 1 | 6 << 2 | 1 << 7
	var rightHand uint16 = 750
	character := newCharacter("cleric")

	if character.EquippedItems.RightHand != rightHand {
		t.Errorf(unexpectedEquippedItem, character.EquippedItems.RightHand, rightHand)
	}

	clauses := make(map[string]interface{})

	clauses[characterIDEquals] = character.ID
	clauses[itemShnIDEquals] = rightHand

	_, err := GetItemWhere(clauses, false)
	if err != nil {
		t.Error(err)
	}
}

func getCharacters(deleted bool) ([]*Character, error) {
	var chars []*Character

	query := db.Model(&chars)

	if !deleted {
		query.Where("character.deleted_at IS NULL")
	}

	err := query.
		Relation("Appearance").
		Relation("Attributes").
		Relation("Items").
		Relation("Options").
		Relation("Location").
		Select()

	return chars, err
}

func newCharacter(class string) *Character {
	var (
		bitField byte
		name     string
	)

	switch class {
	case "mage":
		bitField = byte(1 | 16<<2 | 1<<7)
		name = fmt.Sprintf("mage%v", 1)
		break
	case "fighter":
		bitField = byte(1 | 1<<2 | 1<<7)
		name = fmt.Sprintf("fighter%v", 1)
		break
	case "archer":
		bitField = byte(1 | 11<<2 | 1<<7)
		name = fmt.Sprintf("archer%v", 1)
		break
	case "cleric":
		bitField = byte(1 | 6<<2 | 1<<7)
		name = fmt.Sprintf("cleric%v", 1)
		break
	}

	c := structs.NcAvatarCreateReq{
		SlotNum: byte(0),
		Name: structs.Name5{
			Name: name,
		},
		Shape: structs.ProtoAvatarShapeInfo{
			BF:        bitField,
			HairType:  6,
			HairColor: 0,
			FaceShape: 0,
		},
	}

	char, err := NewCharacter(1, &c, false)
	if err != nil {
		log.Fatal(err)
	}

	return char
}

func cleanDB() {
	err := DeleteTables()
	if err != nil {
		log.Fatal(err)
	}
	err = CreateTables()
	if err != nil {
		log.Fatal(err)
	}
}

func createDummyCharacters() {
	for i := 0; i <= 5; i++ {
		name := fmt.Sprintf("mob%v", i+1)
		c := structs.NcAvatarCreateReq{
			SlotNum: byte(i),
			Name: structs.Name5{
				Name: name,
			},
			Shape: structs.ProtoAvatarShapeInfo{
				BF:        133,
				HairType:  6,
				HairColor: 0,
				FaceShape: 0,
			},
		}
		_, err := NewCharacter(1, &c, false)
		if err != nil {
			log.Fatal(err)
		}
	}
}
