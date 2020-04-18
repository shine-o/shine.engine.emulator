package service

import (
	"context"
	"github.com/shine-o/shine.engine.core/game/character"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
)

//NC_MAP_LOGIN_REQ
func mapLoginReq(ctx context.Context, pc *networking.Command) {
	// todo: shn files checksum
	nc := structs.NcMapLoginReq{}

	err := structs.Unpack(pc.Base.Data, &nc)
	if err != nil {
		return
	}
	charName := nc.CharData.CharID.String()
	var char character.Character
	err = db.Model(&char).
		Relation("Appearance").
		Relation("Attributes").
		Relation("Location").
		Where("name = ?", charName).
		Select() // query the world for a character with name

	if err != nil {
		return
	}

	// todo: check if these packets should be sent sequentially
	go charClientBaseCmd(ctx, &char) // todo: check if race condition
	go charClientShapeCmd(ctx, char.Appearance)

	go clientItemCmd(ctx, char.AllEquippedItems(db))
	go clientItemCmd(ctx, char.InventoryItems(db))
	go clientItemCmd(ctx, char.MiniHouseItems(db))
	go clientItemCmd(ctx, char.PremiumActionItems(db))

	go mapLoginAck(ctx, &char)
	go loginCompleteCmd(ctx)

}

//NC_CHAR_CLIENT_BASE_CMD
func charClientBaseCmd(ctx context.Context, char *character.Character) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4152,
		},
		NcStruct: structs.NcCharClientBaseCmd{
			ChrRegNum:  uint32(char.ID),
			CharName:   structs.NewName5(char.Name),
			Slot:       char.Slot,
			Level:      char.Attributes.Level,
			Experience: char.Attributes.Experience,
			PwrStone:   0,
			GrdStone:   0,
			HPStone:    char.Attributes.HpStones,
			SPStone:    char.Attributes.SpStones,
			CurHP:      char.Attributes.Hp,
			CurSP:      char.Attributes.Sp,
			CurLP:      0,
			Unk:        0,
			Fame:       char.Attributes.Fame,
			Cen:        0,
			LoginInfo: structs.NcCharBaseCmdLoginLocation{
				CurrentMap: structs.NewName3(char.Name),
				CurrentCoord: structs.ShineCoordType{
					XY: structs.ShineXYType{
						X: char.Location.X,
						Y: char.Location.Y,
					},
					Direction: char.Location.D,
				},
			},
			Stats: structs.CharStats{
				Strength:          char.Attributes.Strength,
				Constitute:        char.Attributes.Endurance,
				Dexterity:         char.Attributes.Dexterity,
				Intelligence:      char.Attributes.Intelligence,
				MentalPower:       char.Attributes.Spirit,
				RedistributePoint: 0,
			},
			IdleTime:   0,
			PkCount:    char.Attributes.KillPoints,
			PrisonMin:  0,
			AdminLevel: char.AdminLevel,
			Flag: structs.NcCharBaseCmdFlag{
				Val: 0,
			},
		},
	}
	go pc.Send(ctx)
}

//NC_CHAR_CLIENT_SHAPE_CMD
func charClientShapeCmd(ctx context.Context, ca *character.Appearance) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4152,
		},
		NcStruct: ca.NcRepresentation(),
	}
	go pc.Send(ctx)
}

//NC_CHAR_CLIENT_ITEM_CMD
func clientItemCmd(ctx context.Context, cmd structs.NcCharClientItemCmd) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4167,
		},
		NcStruct: cmd,
	}
	go pc.Send(ctx)
}

//NC_MAP_LOGIN_ACK
func mapLoginAck(ctx context.Context, char *character.Character) {
	// handle ID

	// character complete Parameters (resultant stats from base + items )
	// todo: character stat calculation given the assigned stats, title, equipped items, abstate
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 6146,
		},
		NcStruct: structs.NcMapLoginAck{
			Handle: 0, // id of the entity inside this map
			Params: structs.CharParameterData{
				PrevExp: 0,
				NextExp: 0,
				Strength: structs.ShineCharStatVar{
					Base:   0,
					Change: 0,
				},
				Constitute: structs.ShineCharStatVar{
					Base:   0,
					Change: 0,
				},
				Dexterity: structs.ShineCharStatVar{
					Base:   0,
					Change: 0,
				},
				Intelligence: structs.ShineCharStatVar{
					Base:   0,
					Change: 0,
				},
				Wisdom: structs.ShineCharStatVar{
					Base:   0,
					Change: 0,
				},
				MentalPower: structs.ShineCharStatVar{
					Base:   0,
					Change: 0,
				},
				WCLow: structs.ShineCharStatVar{
					Base:   0,
					Change: 0,
				},
				WCHigh: structs.ShineCharStatVar{
					Base:   0,
					Change: 0,
				},
				AC: structs.ShineCharStatVar{
					Base:   0,
					Change: 0,
				},
				TH: structs.ShineCharStatVar{
					Base:   0,
					Change: 0,
				},
				TB: structs.ShineCharStatVar{
					Base:   0,
					Change: 0,
				},
				MALow: structs.ShineCharStatVar{
					Base:   0,
					Change: 0,
				},
				MAHigh: structs.ShineCharStatVar{
					Base:   0,
					Change: 0,
				},
				MR: structs.ShineCharStatVar{
					Base:   0,
					Change: 0,
				},
				MH: structs.ShineCharStatVar{
					Base:   0,
					Change: 0,
				},
				MB: structs.ShineCharStatVar{
					Base:   0,
					Change: 0,
				},
				MaxHP:      1000,
				MaxSP:      1000,
				MaxLP:      0,
				MaxAP:      0,
				MaxHPStone: 50,
				MaxSPStone: 50,
				PwrStone: structs.CharParameterDataPwrStone{
					Flag:      0,
					EPPPhysic: 0,
					EPMagic:   0,
					MaxStone:  0,
				},
				GrdStone: structs.CharParameterDataPwrStone{
					Flag:      0,
					EPPPhysic: 0,
					EPMagic:   0,
					MaxStone:  0,
				},
				PainRes: structs.ShineCharStatVar{
					Base:   0,
					Change: 0,
				},
				RestraintRes: structs.ShineCharStatVar{
					Base:   0,
					Change: 0,
				},
				CurseRes: structs.ShineCharStatVar{
					Base:   0,
					Change: 0,
				},
				ShockRes: structs.ShineCharStatVar{
					Base:   0,
					Change: 0,
				},
			},
			LoginCoord: structs.ShineXYType{
				X: char.Location.X,
				Y: char.Location.Y,
			},
		},
	}
	go pc.Send(ctx)
	// Login coordinates
}

//NC_MAP_LOGINCOMPLETE_CMD
//6147
func loginCompleteCmd(ctx context.Context) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 6147,
		},
	}
	go pc.Send(ctx)
}
