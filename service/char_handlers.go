package service

import (
	"context"
	"github.com/shine-o/shine.engine.core/game/character"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
)

//NC_MAP_LOGIN_REQ
func NcMapLoginReq(ctx context.Context, np * networking.Parameters) {

	// now how die fukken we know which sector and which handle

	// here the player logs in, we take the map hes at and given the coordinates we know which sector

	// we have a function which given a map id and coordinates it returns us the sector channel which can take event data



	// todo: shn files checksum
	nc := structs.NcMapLoginReq{}

	err := structs.Unpack(np.Command.Base.Data, &nc)
	if err != nil {
		log.Error(err)
		return
	}
	charName := nc.CharData.CharID.Name
	var char character.Character
	err = db.Model(&char).
		Relation("Appearance").
		Relation("Attributes").
		Relation("Location").
		Where("name = ?", charName).
		Select() // query the world for a character with name

	if err != nil {
		log.Error(err)
		return
	}

	// todo: check if these packets should be sent sequentially
	NcCharClientBaseCmd(ctx, &char) // todo: check if race condition
	NcCharClientShapeCmd(ctx, char.Appearance)
	NcCharClientQuestDoingCmd(ctx, &char)
	NcCharClientQuestDoneCmd(ctx, &char)
	NcCharClientQuestReadCmd(ctx, &char)
	NcCharClientQuestRepeatCmd(ctx, &char)
	NcCharClientPassiveCmd(ctx, &char)
	NcCharClientSkillCmd(ctx, &char)
	NcCharClientItemCmd(ctx, char.AllEquippedItems(db))
	NcCharClientItemCmd(ctx, char.InventoryItems(db))
	NcCharClientItemCmd(ctx, char.MiniHouseItems(db))
	NcCharClientItemCmd(ctx, char.PremiumActionItems(db))
	NcCharClientCharTitleCmd(ctx, &char)
	NcCharClientGameCmd(ctx)
	NcCharClientChargedBuffCmd(ctx, &char)
	NcCharClientCoinInfoCmd(ctx, &char)
	NcQuestResetTimeClientCmd(ctx, &char)
	NcMapLoginAck(ctx, &char)
}

//NC_CHAR_CLIENT_BASE_CMD
func NcCharClientBaseCmd(ctx context.Context, char *character.Character) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4152,
		},
		NcStruct: &structs.NcCharClientBaseCmd{
			ChrRegNum:  uint32(char.ID),
			CharName:   structs.Name5{
				Name: char.Name,
			},
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
			Unk:        1,
			Fame:       char.Attributes.Fame,
			Cen:        54983635,
			LoginInfo: structs.NcCharBaseCmdLoginLocation{
				CurrentMap: structs.Name3{
					Name: char.Location.MapName,
				},
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
	pc.Send(ctx)
}

//NC_CHAR_CLIENT_SHAPE_CMD
func NcCharClientShapeCmd(ctx context.Context, ca *character.Appearance) {
	shapeInfo := ca.NcRepresentation()
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4153,
		},
		NcStruct: &shapeInfo,
	}
	pc.Send(ctx)
}

//NC_CHAR_CLIENT_ITEM_CMD
func NcCharClientItemCmd(ctx context.Context, cmd *structs.NcCharClientItemCmd) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4167,
		},
		NcStruct: cmd,
	}
	pc.Send(ctx)
}

//NC_MAP_LOGIN_ACK
func NcMapLoginAck(ctx context.Context, char *character.Character) {
	// handle ID

	// character complete Parameters (resultant stats from base + items )
	// todo: character stat calculation given the assigned stats, title, equipped items, abstate
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 6146,
		},
		NcStruct: &structs.NcMapLoginAck{
			Handle: 7000, // id of the entity inside this map
			Params: structs.CharParameterData{
				PrevExp: 91582354,
				NextExp: 103941051,
				Strength: structs.ShineCharStatVar{
					Base:   213,
					Change: 291,
				},
				Constitute: structs.ShineCharStatVar{
					Base:   213,
					Change: 291,
				},
				Dexterity: structs.ShineCharStatVar{
					Base:   213,
					Change: 291,
				},
				Intelligence: structs.ShineCharStatVar{
					Base:   213,
					Change: 291,
				},
				Wisdom: structs.ShineCharStatVar{
					Base:   213,
					Change: 291,
				},
				MentalPower: structs.ShineCharStatVar{
					Base:   213,
					Change: 291,
				},
				WCLow: structs.ShineCharStatVar{
					Base:   213,
					Change: 291,
				},
				WCHigh: structs.ShineCharStatVar{
					Base:   213,
					Change: 291,
				},
				AC: structs.ShineCharStatVar{
					Base:   213,
					Change: 291,
				},
				TH: structs.ShineCharStatVar{
					Base:   213,
					Change: 291,
				},
				TB: structs.ShineCharStatVar{
					Base:   213,
					Change: 291,
				},
				MALow: structs.ShineCharStatVar{
					Base:   213,
					Change: 291,
				},
				MAHigh: structs.ShineCharStatVar{
					Base:   213,
					Change: 291,
				},
				MR: structs.ShineCharStatVar{
					Base:   213,
					Change: 291,
				},
				MH: structs.ShineCharStatVar{
					Base:   213,
					Change: 291,
				},
				MB: structs.ShineCharStatVar{
					Base:   213,
					Change: 291,
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
	pc.Send(ctx)
	// Login coordinates
}

//NC_CHAR_CLIENT_QUEST_READ_CMD
//4302
func NcCharClientQuestReadCmd(ctx context.Context, char *character.Character) {
	// todo: quest logic
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4302,
		},
		NcStruct: &structs.NcCharClientQuestReadCmd{
			CharID:          uint32(char.ID),
			NumOfReadQuests: 0,
		},
	}
	pc.Send(ctx)
}

//NC_CHAR_CLIENT_QUEST_DOING_CMD
//4154
func NcCharClientQuestDoingCmd(ctx context.Context, char *character.Character) {
	// todo: quest logic
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4154,
		},
		NcStruct: &structs.NcCharClientQuestDoingCmd{
			CharID:          uint32(char.ID),
			NeedClear:       0,
			NumOfDoingQuest: 0,
		},
	}
	pc.Send(ctx)
}

//NC_CHAR_CLIENT_QUEST_DONE_CMD
//4155
func NcCharClientQuestDoneCmd(ctx context.Context, char *character.Character) {
	// todo: quest logic
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4155,
		},
		NcStruct: &structs.NcCharClientQuestDoneCmd{
			CharID:             uint32(char.ID),
			TotalDoneQuest:     0,
			TotalDoneQuestSize: 0,
			Count:              0,
			Index:              0,
		},
	}
	pc.Send(ctx)
}

//NC_CHAR_CLIENT_QUEST_REPEAT_CMD
//4311
func NcCharClientQuestRepeatCmd(ctx context.Context, char *character.Character) {
	// todo: quest logic
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4311,
		},
		NcStruct: &structs.NcCharClientQuestRepeatCmd{
			CharID: uint32(char.ID),
			Count:  0,
		},
	}
	pc.Send(ctx)
}

//NC_CHAR_CLIENT_PASSIVE_CMD
//4158
func NcCharClientPassiveCmd(ctx context.Context, char *character.Character) {
	// todo: skill logic
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4158,
		},
		NcStruct: &structs.NcCharClientPassiveCmd{
			Number: 0,
		},
	}
	pc.Send(ctx)
}

//NC_CHAR_CLIENT_SKILL_CMD
//4157
func NcCharClientSkillCmd(ctx context.Context, char *character.Character) {
	// todo: skill logic
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4157,
		},
		NcStruct: &structs.NcCharClientSkillCmd{
			RestEmpower: 0,
			PartMark: structs.PartMark{
				BF0: 3,
			},
			MaxNum: 0,
			Skills: structs.NcCharSkillClientCmd{
				ChrRegNum: uint32(char.ID),
				Number:    0,
			},
		},
	}
	pc.Send(ctx)
}

//NC_CHAR_CLIENT_CHARTITLE_CMD
//4169
func NcCharClientCharTitleCmd(ctx context.Context, char *character.Character) {
	// todo: character title logic
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4169,
		},
		NcStruct: &structs.NcClientCharTitleCmd{
			CurrentTitle:        0,
			CurrentTitleElement: 0,
			CurrentTitleMobID:   0,
			NumOfTitle:          0,
		},
	}
	pc.Send(ctx)
}

//NC_CHAR_CLIENT_GAME_CMD
//4168
func NcCharClientGameCmd(ctx context.Context) {
	// no idea what this is for
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4168,
		},
		NcStruct: &structs.NcCharClientGameCmd{
			Filler0: 65535,
			Filler1: 65535,
		},
	}
	pc.Send(ctx)
}

//NC_CHAR_CLIENT_CHARGEDBUFF_CMD
//4170
func NcCharClientChargedBuffCmd(ctx context.Context, char *character.Character) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4170,
		},
		NcStruct: &structs.NcCharClientChargedBuffCmd{
			Count: 0,
		},
	}
	pc.Send(ctx)
}

//NC_CHAR_CLIENT_COININFO_CMD
//4318
func NcCharClientCoinInfoCmd(ctx context.Context, char *character.Character) {
	// todo: money && fame logic
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4318,
		},
		NcStruct: &structs.NcCharClientCoinInfoCmd{
			Coin:          73560,
			ExchangedCoin: 100000,
		},
	}
	pc.Send(ctx)
}

//NC_QUEST_RESET_TIME_CLIENT_CMD
//17438
func NcQuestResetTimeClientCmd(ctx context.Context, char *character.Character) {
	// todo: quest logic
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 17438,
		},
		NcStruct: &structs.NcQuestResetTimeClientCmd{
			ResetYear:  1577862000,
			ResetMonth: 1585724400,
			ResetWeek:  1586761200,
			ResetDay:   1587279600,
		},
	}
	pc.Send(ctx)
}

//NC_MAP_LOGINCOMPLETE_CMD
//6147
func NcMapLoginCompleteCmd(ctx context.Context, pc *networking.Command) {}
