package character

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.core/database"
	"github.com/shine-o/shine.engine.core/structs"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var db *pg.DB

func TestMain(m *testing.M) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	log = logger.Init("test logger", true, false, ioutil.Discard)
	log.Info("test logger")
	db = database.Connection(ctx, database.ConnectionParams{
		User:     "user",
		Password: "password",
		Host:     "postgres",
		//Host:     "192.168.1.238",
		Port:     "5432",
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

func cleanDB() {
	err := DeleteTables(db)
	if err != nil {
		log.Error(err)
	}
	err = CreateTables(db)
	if err != nil {
		log.Error(err)
	}
}

func createDummyCharacters() {
	for i := 0; i <= 5; i++ {
		name := fmt.Sprintf("mob%v", i+1)
		c := structs.NcAvatarCreateReq{
			SlotNum: byte(i),
			Name:    structs.Name5{
				Name: name,
			},
			Shape: structs.ProtoAvatarShapeInfo{
				BF:        133,
				HairType:  6,
				HairColor: 0,
				FaceShape: 0,
			},
		}
		_, err := New(db, 1, &c)
		if err != nil {
			log.Fatal(err)
		}
	}
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
		err = Validate(db, 1, &nc)
		if err != nil {
			t.Error(err)
		}
	}
}

// foolishly assuming data is okay
func TestCreateCharacter(t *testing.T) {
	defer cleanDB()
	// {"packetType":"small","length":27,"department":5,"command":"1","opCode":5121,"data":"0046696768747265726f6f0000000000000000000085060000","rawData":"1b01140046696768747265726f6f0000000000000000000085060000","friendlyName":"NC_AVATAR_CREATE_REQ"}
	if data, err := hex.DecodeString("0046696768747265726f6f0000000000000000000085060000"); err != nil {
		t.Error(err)
	} else {
		nc := structs.NcAvatarCreateReq{}
		err := structs.Unpack(data, &nc)
		if err != nil {
			t.Error(err)
		}
		_, err = New(db, 1, &nc)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestDeleteCharacter(t *testing.T) {
	defer cleanDB()

	createDummyCharacters()

	// {"packetType":"small","length":3,"department":5,"command":"7","opCode":5127,"data":"03","rawData":"03071403","friendlyName":"NC_AVATAR_ERASE_REQ"}
	nc := structs.NcAvatarEraseReq{
		Slot: 0,
	}

	err := Delete(db, 1, &nc)

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
	defer cleanDB()

	createDummyCharacters()
	
	name := fmt.Sprintf("mob%v", 1)
	c := structs.NcAvatarCreateReq{
		SlotNum: byte(0),
		Name:    structs.Name5{
			Name: name,
		},
		Shape: structs.ProtoAvatarShapeInfo{
			BF:        133,
			HairType:  6,
			HairColor: 0,
			FaceShape: 0,
		},
	}
	err := Validate(db, 1, &c)
	if err == nil {
		log.Error(err)
	}

	errChar, ok := err.(*ErrCharacter)

	if !ok {
		log.Error(err)
		t.Error("expected error of type ErrCharacter")
		return
	}

	if errChar.Code != 1 {
		log.Error(err)
		t.Errorf("expected errorCharacter with code %v, instead got %v", 1, errChar.Code)
	}
}

func TestNoSlotAvailableError(t *testing.T) {
	defer cleanDB()
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

		err = Validate(db, 1, &nc)
		if err != nil {
			if err.Error() != "no slot available" {
				t.Errorf("unexpected error message %v", err.Error())
			}
			return
		}
		t.Error("expected an error but got nil")
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

		err = Validate(db, 1, &nc)
		if err != nil {
			if err.Error() != "invalid slot" {
				t.Errorf("unexpected error message %v", err.Error())
			}
			return
		}
		t.Error("expected an error but got nil")
	}
}

func TestInvalidNameError(t *testing.T) {
	defer cleanDB()

	// {"packetType":"small","length":27,"department":5,"command":"1","opCode":5121,"data":"0046325B5D747265726F6F0000000000000000000085060000","rawData":"1b01140046325B5D747265726F6F0000000000000000000085060000","friendlyName":"NC_AVATAR_CREATE_REQ"}
	if data, err := hex.DecodeString("0046325B5D747265726F6F0000000000000000000085060000"); err != nil {
		t.Error(err)
	} else {
		nc := structs.NcAvatarCreateReq{}
		err := structs.Unpack(data, &nc)
		if err != nil {
			t.Error(err)
		}

		err = Validate(db, 1, &nc)
		if err != nil {
			if err.Error() != "invalid name" {
				t.Errorf("unexpected error message: %v", err.Error())
			}
			return
		}
		t.Error("expected an error but got nil")
	}
}

// KOREAN MONKEYS!
// test that correct gender and class are extracted using binary operators
func TestInvalidGenderClassBinaryOperation(t *testing.T) {
	defer cleanDB()

	// {"packetType":"small","length":27,"department":5,"command":"1","opCode":5121,"data":"0946696768747265726f6f0000000000000000000085060000","rawData":"1b01140046325B5D747265726F6F0000000000000000000085060000","friendlyName":"NC_AVATAR_CREATE_REQ"}
	if data, err := hex.DecodeString("0246696768747265726f6f00000000000000000000ff060000"); err != nil {
		t.Error(err)
	} else {
		nc := structs.NcAvatarCreateReq{}
		err := structs.Unpack(data, &nc)
		if err != nil {
			t.Error(err)
		}

		err = Validate(db, 1, &nc)
		if err != nil {
			errChar, ok := err.(*ErrCharacter)
			if !ok {
				t.Error("expected an error of type ErrCharacter")
			}
			if errChar.Code != 4 {
				t.Errorf("expected errorCharacter with code %v, instead got %v", 4, errChar.Code)
			}
		} else {
			t.Error("expected an error but got nil")
		}
	}
}
