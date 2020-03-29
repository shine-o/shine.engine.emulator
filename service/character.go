package service

import (
	"context"
	"errors"
	"github.com/shine-o/shine.engine.networking/structs"
)

type noSession error
type noSlot error
type invalidSlot error
type invalidName error
type invalidClassGender error

var errNoSession = errors.New("no login session was found").(noSession)
var errNoSlot = errors.New("no slot available").(noSlot)
var errInvalidSlot = errors.New("invalid slot").(invalidSlot)
var errInvalidName = errors.New("invalid name").(invalidName)
var errInvalidClassGender = errors.New("invalid Class Gender data").(invalidClassGender)

const (
	startLevel = 1
	startMap = "Rou"
)

func newCharacter(ctx context.Context, req structs.NcAvatarCreateReq) (structs.AvatarInformation, error) {
	select {
	case <- ctx.Done():
		return structs.AvatarInformation{}, errCC
	default:
		//
		avatar := structs.AvatarInformation{}

		// new Character model

		copy(avatar.Name.Name[:], "works")
		avatar.Slot = 0
		avatar.Level = 1
		copy(avatar.LoginMap.Name[:], "Rou")

		avatar.Shape = req.Shape

		avatar.ChrRegNum = 1007

		avatar.TutorialInfo.TutorialState = 2
		avatar.TutorialInfo.TutorialStep = 0

		ei := structs.ProtoEquipment{
			EquHead:         65535,
			EquMouth:        65535,
			EquRightHand:    65535,
			EquBody:         11,
			EquLeftHand:     65535,
			EquPant:         10,
			EquBoot:         8,
			EquAccBoot:      35421,
			EquAccPant:      65535,
			EquAccHeadA:     65535,
			EquMinimonR:     65535,
			EquEye:          65535,
			EquAccLeftHand:  65535,
			EquAccRightHand: 65535,
			EquAccBack:      65535,
			EquCosEff:       65535,
			EquAccHip:       65535,
			EquMinimon:      65535,
			EquAccShield:    65535,
		}

		avatar.Equip = ei

		return avatar, nil
	}
}

func validateCharacter(ctx context.Context, req structs.NcAvatarCreateReq) error {
	// fetch session

	session, ok := ctx.Value("session").(session)

	if !ok {
		return errNoSession
	}

	// query all characters for this account
	// if no slot is available, return errNoSlot
	//req.SlotNum
	var chars []*Character
	worldDB.Model(chars).Where("user_id = ?" , session.UserID)

	// if name contains special characters, return errInvalidName (name should only contain alphanumeric characters)
	// if a character with the same name already exists, return errNameTaken
	// req.Name

	//
	//req.Shape.BF


	return nil
}