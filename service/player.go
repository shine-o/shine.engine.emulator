package service

import (
	"github.com/shine-o/shine.engine.core/game/character"
	"github.com/shine-o/shine.engine.core/structs"
)

type player struct {
	baseEntity
	conn     playerConnection
	view     playerView
	stats    playerStats
	state    playerState
	items    playerItems
	money    playerMoney
	titles   playerTitles
	quests   playerQuests
	skills   []skill
	passives []passive
}

type playerConnection struct {
	close        chan<- bool
	outboundData chan<- []byte
}

type playerView struct {
	name       string
	class      uint8
	gender     uint8
	hairType   uint8
	hairColour uint8
	faceType   uint8
}

type playerStats struct {
	prevExp             uint64
	nextExp             uint64
	str                 stat
	end                 stat
	dex                 stat
	int                 stat
	spr                 stat
	physicalDamage      stat
	magicalDamage       stat
	physicalDefense     stat
	magicalDefense      stat
	evasion             stat
	aim                 stat
	hp                  uint32
	sp                  uint32
	lp                  uint32
	hpStones            uint32
	spStones            uint32
	curseResistance     stat
	restraintResistance stat
	poisonResistance    stat
	rollbackResistance  stat
}

type playerState struct {
	level       uint8
	autoPickup  uint8
	polymorph   uint16
	moverHandle uint16
	moverSlot   uint8
	miniPet     uint8
}

type playerItems struct {
	equipped  itemBox
	inventory itemBox
	miniHouse itemBox
	reward    itemBox
	premium   itemBox
}

type playerMoney struct {
	coins       uint64
	fame        uint64
	wastedCoins uint64
	wastedFame  uint64
}

type playerTitles struct {
	current struct {
		id      uint8
		element uint8
		mobID   uint16
	}
	titles []title
}

type playerQuests struct {
	read       []quest
	done       []quest
	doing      []quest
	repeatable []quest
}

type skill struct {
	id       uint16
	coolTime uint32
}

type passive struct {
	id uint16
}

type quest struct {
	id             uint16
	status         uint8
	startTime      int32
	endTime        int32
	step           uint8
	completedCount int
	deadline       int32
}

type title struct {
	// fields should be different
	// MainChar::InitCharactorTitle(PROTO_NC_CHAR_CLIENT_CHARTITLE_CMD *)
	tType        uint8
	bitOperation uint8
}

type itemBox struct {
	box   uint8
	items []item
}

type item struct {
	slot uint16
	id   uint16
}

type stat struct {
	base     uint32
	withGear uint32
}

func (p *player) load(name string) error {
	char, err := character.GetByName(db, name)

	if err != nil {
		return err
	}

	view := make(chan playerView)
	state := make(chan playerState)
	stats := make(chan playerStats)
	items := make(chan playerItems)
	money := make(chan playerMoney)
	titles := make(chan playerTitles)
	quests := make(chan playerQuests)
	skills := make(chan []skill)
	passives := make(chan []passive)
	pErr := make(chan error)

	go p.viewData(view, &char, pErr)
	go p.stateData(state, &char, pErr)
	go p.statsData(stats, &char, pErr)
	go p.itemData(items, &char, pErr)
	go p.moneyData(money, &char, pErr)
	go p.titleData(titles, &char, pErr)
	go p.skillData(skills, &char, pErr)
	go p.passiveData(passives, &char, pErr)

	select {
	case v := <-view:
		p.view = v
	case s := <-state:
		p.state = s
	case s := <-stats:
		p.stats = s
	case i := <-items:
		p.items = i

	case m := <-money:
		p.money = m

	case t := <-titles:
		p.titles = t

	case q := <-quests:
		p.quests = q

	case s := <-skills:
		p.skills = s

	case pa := <-passives:
		p.passives = pa

	case err := <-pErr:
		return err
	}
	// for p launch routines to create player inner structs view, state, stats
	//
	return nil
}

func (p *player) viewData(view chan <- playerView, c *character.Character, err chan error) {
	v := playerView { // todo: validation just in case, so we don't log a bad player that could potentially bin other player
		name:       c.Name,
		class:      c.Appearance.Class,
		gender:     c.Appearance.Gender,
		hairType:   c.Appearance.HairType,
		hairColour: c.Appearance.HairColor,
		faceType:   c.Appearance.FaceType,
	}
	view <- v
}

func (p *player) stateData(state chan  <-playerState, c *character.Character, err chan  <-error) {
	s := playerState{
		level:       c.Attributes.Level,
		// player state should also include buffs and debuffs in the future
		autoPickup:  0,
		polymorph:   65535,
		moverHandle: 0,
		moverSlot:   0,
		miniPet:     0,
	}
	state <- s
}

func (p *player) statsData(stats chan  <-playerStats, c *character.Character, err chan  <-error) {
	// given all:
	//  class base stats for current level, equipped items, charged buffs, buffs/debuffs, assigned stat points
	// calculate base stats (class base stats for current level, assigned stat points) , and stats with gear on (equipped items, charged buffs, buffs/debuffs)
	s := playerStats{
		prevExp: 0,
		nextExp: 0,
		str: stat{
			base:     0,
			withGear: 0,
		},
		end: stat{
			base:     0,
			withGear: 0,
		},
		dex: stat{
			base:     0,
			withGear: 0,
		},
		int: stat{
			base:     0,
			withGear: 0,
		},
		spr: stat{
			base:     0,
			withGear: 0,
		},
		physicalDamage: stat{
			base:     0,
			withGear: 0,
		},
		magicalDamage: stat{
			base:     0,
			withGear: 0,
		},
		physicalDefense: stat{
			base:     0,
			withGear: 0,
		},
		magicalDefense: stat{
			base:     0,
			withGear: 0,
		},
		evasion: stat{
			base:     0,
			withGear: 0,
		},
		aim: stat{
			base:     0,
			withGear: 0,
		},
		hp:       100,
		sp:       100,
		lp:       0,
		hpStones: 0,
		spStones: 0,
		curseResistance: stat{
			base:     0,
			withGear: 0,
		},
		restraintResistance: stat{
			base:     0,
			withGear: 0,
		},
		poisonResistance: stat{
			base:     0,
			withGear: 0,
		},
		rollbackResistance: stat{
			base:     0,
			withGear: 0,
		},
	}

	stats <- s
}

func (p *player) itemData(items chan  <-playerItems, c *character.Character, err chan <- error) {

}

func (p *player) titleData(titles chan <- playerTitles, c *character.Character, err chan  <-error) {

}

func (p *player) moneyData(money chan  <-playerMoney, c *character.Character, err chan <- error) {

}

func (p *player) skillData(skills chan <- []skill, c *character.Character, err chan  <-error) {

}

func (p *player) passiveData(passives chan <- []passive, c *character.Character, err chan <- error) {

}

//ncCharClientBaseCmd(ctx, &char) 
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
func (p *player) ncLoginRepresentation() structs.NcBriefInfoLoginCharacterCmd {
	nc := structs.NcBriefInfoLoginCharacterCmd{
		Handle: p.getHandle(),
		CharID: structs.Name5{
			Name: p.view.name,
		},
		Coordinates: structs.ShineCoordType{
			XY: structs.ShineXYType{
				X: p.location.x,
				Y: p.location.y,
			},
			Direction: p.location.d,
		},
		Mode:  0,
		Class: p.view.class,
		Shape: structs.ProtoAvatarShapeInfo{
			BF:        1 | p.view.class<<2 | p.view.gender<<7,
			HairType:  0,
			HairColor: 0,
			FaceShape: 0,
		},
		ShapeData:       structs.NcBriefInfoLoginCharacterCmdShapeData{},
		Polymorph:       p.state.polymorph,
		Emoticon:        structs.StopEmoticonDescript{},
		CharTitle:       structs.CharTitleBriefInfo{},
		AbstateBit:      structs.AbstateBit{},
		MyGuild:         0,
		Type:            0,
		IsAcademyMember: 0,
		IsAutoPick:      0,
		Level:           p.state.level,
		Animation:       [32]byte{},
		MoverHandle:     p.state.moverHandle,
		MoverSlot:       p.state.moverSlot,
		KQTeamType:      0,
		UsingMinipet:    p.state.miniPet,
		Unk:             0,
	}
	return nc
}
