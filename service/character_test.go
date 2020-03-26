package service

import (
	"context"
	"encoding/hex"
	"github.com/shine-o/shine.engine.networking/structs"
	"reflect"
	"testing"
)

// test character data is valid
func TestValidateCharacterRequest(t *testing.T)  {
	// assert it returns true
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
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

		// fetch character
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
		if err == nil {
			t.Error("expected an error but got nil")
		}

		if reflect.TypeOf(err) != reflect.TypeOf(new(noSlot)) {
			t.Errorf("expected error of type %v but instead got %v", reflect.TypeOf(new(noSlot)).String() , reflect.TypeOf(err).String())
		}
	}
}

func TestInvalidSlotError(t *testing.T)  {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
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
		if err == nil {
			t.Error("expected an error but got nil")
		}

		if reflect.TypeOf(err) != reflect.TypeOf(new(invalidSlot)) {
			t.Errorf("expected error of type %v but instead got %v", reflect.TypeOf(new(invalidSlot)).String() , reflect.TypeOf(err).String())
		}
	}
}

func TestInvalidNameError(t *testing.T)  {
	// make sure function returns error "invalid character name"

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

		if reflect.TypeOf(err) != reflect.TypeOf(new(invalidName)) {
			t.Errorf("expected error of type %v but instead got %v", reflect.TypeOf(new(invalidName)).String() , reflect.TypeOf(err).String())
		}
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

		if reflect.TypeOf(err) != reflect.TypeOf(new(invalidClassGender)) {
			t.Errorf("expected error of type %v but instead got %v", reflect.TypeOf(new(invalidClassGender)).String() , reflect.TypeOf(err).String())
		}
	}
}