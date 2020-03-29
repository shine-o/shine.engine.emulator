package service

import (
	"context"
	"encoding/hex"
	networking "github.com/shine-o/shine.engine.networking"
	"github.com/shine-o/shine.engine.networking/structs"
	"testing"
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

// foolishly assuming data is okay
func TestCreateCharacter(t *testing.T) {
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

}

func TestCharacterNameInUseError(t *testing.T)  {
	
}

func TestNoSlotAvailableError(t *testing.T)  {
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
				t.Errorf("unexpected error message %v", err.Error())
			}
			return
		}
		t.Error("expected an error but got nil")
	}
}

// KOREAN MONKEYS!
// test that correct gender and class are extracted using binary operators
func TestInvalidGenderClassBinaryOperation(t *testing.T)  {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
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
		if err == nil {
			t.Error("expected an error but got nil")
		}

		err = validateCharacter(ctx, nc)
		if err != nil {
			if err.Error() != "invalid Class Gender data" {
				t.Errorf("unexpected error message %v", err.Error())
			}
			return
		}
		t.Error("expected an error but got nil")
	}
}