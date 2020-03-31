package service

import (
	"context"
	"encoding/hex"
	"fmt"
	networking "github.com/shine-o/shine.engine.networking"
	"github.com/shine-o/shine.engine.networking/structs"
	"testing"
	"time"
)

// test character data is valid
func TestValidateCharacterRequest(t *testing.T)  {
	// assert it returns true

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s := &session{
		ID:       "bcd1fde6-f9d0-451d-a4b6-4992bd6207e1",
		WorldID:  "1",
		UserName: "admin",
		UserID: 1,
	}
	ctx = context.WithValue(ctx, networking.ShineSession, s)

	// {"packetType":"small","length":27,"department":5,"command":"1","opCode":5121,"data":"0046696768747265726f6f0000000000000000000085060000","rawData":"1b01140046696768747265726f6f0000000000000000000085060000","friendlyName":"NC_AVATAR_CREATE_REQ"}
	if data, err := hex.DecodeString("0046696768747265726f6f0000000000000000000085060000"); err != nil {
		t.Error(err)
	} else {
		nc := structs.NcAvatarCreateReq{}
		err := nc.Unpack(data)
		if err != nil {
			t.Error(err)
		}
		err = validateCharacter(ctx, nc)
		if err != nil {
			t.Error(err)
		}
	}
}

func cleanDB()  {
	purge(worldDB)
	createSchema(worldDB)
}

// foolishly assuming data is okay
func TestCreateCharacter(t *testing.T) {
	defer cleanDB()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s := &session{
		ID:       "bcd1fde6-f9d0-451d-a4b6-4992bd6207e1",
		WorldID:  "1",
		UserName: "admin",
		UserID: 1,
	}
	ctx = context.WithValue(ctx, networking.ShineSession, s)

	// {"packetType":"small","length":27,"department":5,"command":"1","opCode":5121,"data":"0046696768747265726f6f0000000000000000000085060000","rawData":"1b01140046696768747265726f6f0000000000000000000085060000","friendlyName":"NC_AVATAR_CREATE_REQ"}
	if data, err := hex.DecodeString("0046696768747265726f6f0000000000000000000085060000"); err != nil {
		t.Error(err)
	} else {
		nc := structs.NcAvatarCreateReq{}
		err := nc.Unpack(data)
		if err != nil {
			t.Error(err)
		}

		_, err = newCharacter(ctx, nc)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestDeleteCharacter(t *testing.T) {
	defer cleanDB()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s := &session{
		ID:       "bcd1fde6-f9d0-451d-a4b6-4992bd6207e1",
		WorldID:  "1",
		UserName: "admin",
		UserID: 1,
	}
	ctx = context.WithValue(ctx, networking.ShineSession, s)

	createDummyCharacters(ctx)

	// {"packetType":"small","length":3,"department":5,"command":"7","opCode":5127,"data":"03","rawData":"03071403","friendlyName":"NC_AVATAR_ERASE_REQ"}
	nc := structs.NcAvatarEraseReq{
		Slot: 0,
	}

	err := deleteCharacter(ctx, nc)

	if err != nil {
		log.Error(err)
	}

	// try and fetch a character for this user id in the current slot.
	var deletedAt time.Time
	err = worldDB.Model((*Character)(nil)).
		Column("deleted_at").
		Where("user_id = ?", 1).
		Where("slot = ?", 0).
		Select(&deletedAt)

	if err != nil {
		t.Error(err)
	}

}

func TestCharacterNameInUseError(t *testing.T)  {
	
}

func createDummyCharacters(ctx context.Context) {

	for i := 0; i <= 5; i++ {
		name := fmt.Sprintf("mob%v", i+1)
		c :=  structs.NcAvatarCreateReq{
			SlotNum: byte(i),
			Name:    structs.Name5{},
			Shape:   structs.ProtoAvatarShapeInfo{
				BF:        133,
				HairType:  6,
				HairColor: 0,
				FaceShape: 0,
			},
		}
		copy(c.Name.Name[:], name)
		_, err := newCharacter(ctx, c)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func TestNoSlotAvailableError(t *testing.T)  {
	defer cleanDB()

	// fill database with characters for dummy account
	// dummy data prep goes here or in the main function, data is only used for this test function
	// account session needs to be in place
	// make sure function returns error of type "no slots available"
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	s := &session{
		ID:       "bcd1fde6-f9d0-451d-a4b6-4992bd6207e1",
		WorldID:  "1",
		UserName: "admin",
		UserID: 1,
	}
	ctx = context.WithValue(ctx, networking.ShineSession, s)

	createDummyCharacters(ctx)

	// {"packetType":"small","length":27,"department":5,"command":"1","opCode":5121,"data":"0046696768747265726f6f0000000000000000000085060000","rawData":"1b01140046696768747265726f6f0000000000000000000085060000","friendlyName":"NC_AVATAR_CREATE_REQ"}
	if data, err := hex.DecodeString("0046696768747265726f6f0000000000000000000085060000"); err != nil {
		t.Error(err)
	} else {
		nc := structs.NcAvatarCreateReq{}
		err := nc.Unpack(data)
		if err != nil {
			t.Error(err)
		}

		err = validateCharacter(ctx, nc)
		if err != nil {
			if err.Error() != "no slot available" {
				t.Errorf("unexpected error message %v", err.Error())
			}
			return
		}
		t.Error("expected an error but got nil")
	}
}

func TestInvalidSlotError(t *testing.T)  {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s := &session{
		ID:       "bcd1fde6-f9d0-451d-a4b6-4992bd6207e1",
		WorldID:  "1",
		UserName: "admin",
		UserID: 1,
	}
	ctx = context.WithValue(ctx, networking.ShineSession, s)

	// {"packetType":"small","length":27,"department":5,"command":"1","opCode":5121,"data":"0946696768747265726f6f0000000000000000000085060000","rawData":"1b01140046696768747265726f6f0000000000000000000085060000","friendlyName":"NC_AVATAR_CREATE_REQ"}
	if data, err := hex.DecodeString("0946696768747265726f6f0000000000000000000085060000"); err != nil {
		t.Error(err)
	} else {
		nc := structs.NcAvatarCreateReq{}
		err := nc.Unpack(data)
		if err != nil {
			t.Error(err)
		}

		err = validateCharacter(ctx, nc)
		if err != nil {
			if err.Error() != "invalid slot" {
				t.Errorf("unexpected error message %v", err.Error())
			}
			return
		}
		t.Error("expected an error but got nil")
	}
}

func TestInvalidNameError(t *testing.T)  {
	defer cleanDB()
	// make sure function returns error "invalid character name"

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	s := &session{
		ID:       "bcd1fde6-f9d0-451d-a4b6-4992bd6207e1",
		WorldID:  "1",
		UserName: "admin",
		UserID: 1,
	}
	ctx = context.WithValue(ctx, networking.ShineSession, s)

	// {"packetType":"small","length":27,"department":5,"command":"1","opCode":5121,"data":"0046325B5D747265726F6F0000000000000000000085060000","rawData":"1b01140046325B5D747265726F6F0000000000000000000085060000","friendlyName":"NC_AVATAR_CREATE_REQ"}
	if data, err := hex.DecodeString("0046325B5D747265726F6F0000000000000000000085060000"); err != nil {
		t.Error(err)
	} else {
		nc := structs.NcAvatarCreateReq{}
		err := nc.Unpack(data)
		if err != nil {
			t.Error(err)
		}

		err = validateCharacter(ctx, nc)
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
func TestInvalidGenderClassBinaryOperation(t *testing.T)  {
	defer cleanDB()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	s := &session{
		ID:       "bcd1fde6-f9d0-451d-a4b6-4992bd6207e1",
		WorldID:  "1",
		UserName: "admin",
		UserID: 1,
	}
	ctx = context.WithValue(ctx, networking.ShineSession, s)


	// {"packetType":"small","length":27,"department":5,"command":"1","opCode":5121,"data":"0946696768747265726f6f0000000000000000000085060000","rawData":"1b01140046325B5D747265726F6F0000000000000000000085060000","friendlyName":"NC_AVATAR_CREATE_REQ"}
	if data, err := hex.DecodeString("0246696768747265726f6f00000000000000000000ff060000"); err != nil {
		t.Error(err)
	} else {
		nc := structs.NcAvatarCreateReq{}
		err := nc.Unpack(data)
		if err != nil {
			t.Error(err)
		}

		err = validateCharacter(ctx, nc)
		if err != nil {
			if err.Error() != "invalid Class Gender data" {
				t.Errorf("unexpected error message: %v", err.Error())
			}
			return
		}
		t.Error("expected an error but got nil")
	}
}