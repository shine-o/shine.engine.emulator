package service

import (
	"context"
	"github.com/shine-o/shine.engine.networking/structs"
)

func newCharacter(ctx context.Context, req structs.NcAvatarCreateReq) (structs.AvatarInformation, error) {
	select {
	case <- ctx.Done():
		return structs.AvatarInformation{}, errCC
	default:
		//

		avatar := structs.AvatarInformation{}

		copy(avatar.Name.Name[:], "works")
		avatar.Slot = 0
		avatar.Level = 1
		copy(avatar.LoginMap.Name[:], "Rou")
		avatar.DelInfo.Year = byte(43)

		avatar.Shape = req.Shape

		avatar.ChrRegNum = 1007

		avatar.TutorialInfo.TutorialState = 2
		avatar.TutorialInfo.TutorialStep = 0

		ei := structs.ProtoEquipment{
			EquHead:         -1,
			EquMouth:        -1,
			EquRightHand:    -1,
			EquBody:         -1,
			EquLeftHand:     -1,
			EquPant:         -1,
			EquBoot:         -1,
			EquAccBoot:      -1,
			EquAccPant:      -1,
			EquAccHeadA:     -1,
			EquMinimonR:     -1,
			EquEye:          -1,
			EquAccLeftHand:  -1,
			EquAccRightHand: -1,
			EquAccBack:      -1,
			EquCosEff:       -1,
			EquAccHip:       -1,
			EquMinimon:      -1,
			EquAccShield:    -1,
		}
		avatar.Equip = ei

		return avatar, nil
	}
}