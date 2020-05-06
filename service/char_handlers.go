package service

import (
	"context"
	"github.com/shine-o/shine.engine.core/game/character"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
)

// NC_MAP_LOGIN_REQ
func ncMapLoginReq(ctx context.Context, np *networking.Parameters) {
	// todo: shn files checksum
	nc := structs.NcMapLoginReq{}

	err := structs.Unpack(np.Command.Base.Data, &nc)
	if err != nil {
		log.Error(err)
		return
	}

	// NC_MAP_LOGIN_REQ process
	// 		shnFileCheck event
	// 		character data event(shape, stats, quests, items, etc..)
	//
	// 		NC_MAP_LOGIN_ACK
	//			registerPlayerHandle event  (save handle id in session, for when client sends NC_MAP_LOGINCOMPLETE_CMD)
	//
	//			NC_MAP_LOGINCOMPLETE_CMD (continued in handler: ncMapLoginCompleteCmd)
	// 			playerAppearedEvent
	//			socialNotificationsEvent


	


	//// todo: check if these packets should be sent sequentially
	//
	//// todo: quest wrapper
	//ncCharClientBaseCmd(ctx, &char) // todo: check if race condition
	//ncCharClientShapeCmd(ctx, char.Appearance)
	//
	//// todo: quest wrapper
	//ncCharClientQuestDoingCmd(ctx, &char)
	//ncCharClientQuestDoneCmd(ctx, &char)
	//ncCharClientQuestReadCmd(ctx, &char)
	//ncCharClientQuestRepeatCmd(ctx, &char)
	//
	//// todo: skills wrapper
	//ncCharClientPassiveCmd(ctx, &char)
	//ncCharClientSkillCmd(ctx, &char)
	//
	//ncCharClientItemCmd(ctx, char.AllEquippedItems(db))
	//ncCharClientItemCmd(ctx, char.InventoryItems(db))
	//ncCharClientItemCmd(ctx, char.MiniHouseItems(db))
	//ncCharClientItemCmd(ctx, char.PremiumActionItems(db))
	//
	//ncCharClientCharTitleCmd(ctx, &char)
	//
	//ncCharClientGameCmd(ctx)
	//ncCharClientChargedBuffCmd(ctx, &char)
	//ncCharClientCoinInfoCmd(ctx, &char)
	//ncQuestResetTimeClientCmd(ctx, &char)
	//
	//// this should be sent after the client sends  the cmd packet acknowledging the map login is complete
	//var pae playerAppearedEvent
	//err = pae.process(np, &char)
	//
	//if err != nil {
	//	log.Error(err)
	//}
	//
	//ncMapLoginAck(pae.player, &char)

	// also send nearby players, mobs, mounts
	// NC_BRIEFINFO_CHARACTER_CMD
	// NC_BRIEFINFO_MOB_CMD
	// NC_BRIEFINFO_MOVER_CMD

	// register player to map

	//runPlayerAppearedEvent(np, char)

	// to this player send him info about nearby players, mobs, npc
}

//NC_CHAR_CLIENT_BASE_CMD
func ncCharClientBaseCmd(ctx context.Context, char *character.Character) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4152,
		},
		NcStruct: &structs.NcCharClientBaseCmd{
			ChrRegNum: uint32(char.ID),
			CharName: structs.Name5{
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
func ncCharClientShapeCmd(ctx context.Context, ca *character.Appearance) {
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
func ncCharClientItemCmd(ctx context.Context, cmd *structs.NcCharClientItemCmd) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4167,
		},
		NcStruct: cmd,
	}
	pc.Send(ctx)
}

//NC_MAP_LOGIN_ACK
func ncMapLoginAck(p *player, char *character.Character) {
	// handle ID

	// character complete Parameters (resultant stats from base + items )
	// todo: character stat calculation given the assigned stats, title, equipped items, abstate
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 6146,
		},
		NcStruct: &structs.NcMapLoginAck{
			Handle: p.handle, // id of the entity inside this map
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
				X: uint32(char.Location.X),
				Y: uint32(char.Location.Y),
			},
		},
	}
	pc.SendDirectly(p.conn.outboundData)
}

//NC_CHAR_CLIENT_QUEST_READ_CMD
//4302
func ncCharClientQuestReadCmd(ctx context.Context, char *character.Character) {
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
func ncCharClientQuestDoingCmd(ctx context.Context, char *character.Character) {
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
func ncCharClientQuestDoneCmd(ctx context.Context, char *character.Character) {
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
func ncCharClientQuestRepeatCmd(ctx context.Context, char *character.Character) {
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
func ncCharClientPassiveCmd(ctx context.Context, char *character.Character) {
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
func ncCharClientSkillCmd(ctx context.Context, char *character.Character) {
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
func ncCharClientCharTitleCmd(ctx context.Context, char *character.Character) {
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
func ncCharClientGameCmd(ctx context.Context) {
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
func ncCharClientChargedBuffCmd(ctx context.Context, char *character.Character) {
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
func ncCharClientCoinInfoCmd(ctx context.Context, char *character.Character) {
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
func ncQuestResetTimeClientCmd(ctx context.Context, char *character.Character) {
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

// NC_MAP_LOGINCOMPLETE_CMD
// 6147
func ncMapLoginCompleteCmd(ctx context.Context, pc *networking.Command) {}

//NC_CHAR_LOGOUTREADY_CMD
func ncCharLogoutReadyCmd(ctx context.Context, np *networking.Parameters) {
	np.NetVars.CloseConnection <- true
}
