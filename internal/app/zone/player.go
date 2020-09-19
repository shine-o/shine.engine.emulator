package zone

import (
	"github.com/go-pg/pg/v9"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game/character"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
	"sync"
	"time"
)

type player struct {
	baseEntity
	players  map[uint16]*player
	monsters map[uint16]*monster
	char     *character.Character
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
	tickers  []*time.Ticker
	sync.RWMutex
}

func (p *player) getHandle() uint16 {
	p.RLock()
	h := p.handle
	p.RUnlock()
	return h
}

func lastHeartbeat(p *player) float64 {
	p.RLock()
	lastHeartBeat := time.Since(p.conn.lastHeartBeat).Seconds()
	p.RUnlock()
	return lastHeartBeat
}

func (p *player) adjacentPlayers() <-chan *player {
	p.RLock()
	ch := make(chan *player, len(p.players))
	p.RUnlock()

	go func(send chan<- *player) {
		p.RLock()
		for _, ap := range p.players {
			send <- ap
		}
		p.RUnlock()
		close(send)
	}(ch)

	return ch
}

func (p *player) removeAdjacentPlayer(h uint16) {
	p.Lock()
	delete(p.players, h)
	p.Unlock()
}

func (p *player) adjacentMonsters() <-chan *monster {
	p.RLock()
	ch := make(chan *monster, len(p.monsters))
	p.RUnlock()

	go func(send chan<- *monster) {
		p.RLock()
		for _, ap := range p.monsters {
			send <- ap
		}
		p.RUnlock()
		close(send)
	}(ch)

	return ch
}

type playerConnection struct {
	lastHeartBeat time.Time
	close         chan<- bool
	outboundData  chan<- []byte
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
	points            playerStatPoints
	str               stat
	end               stat
	dex               stat
	int               stat
	spr               stat
	minPhysicalDamage stat
	maxPhysicalDamage stat

	minMagicalDamage stat
	maxMagicalDamage stat

	physicalDefense     stat
	magicalDefense      stat
	evasion             stat
	aim                 stat
	hp                  uint32
	sp                  uint32
	lp                  uint32
	maxHpStones         uint32
	maxSpStones         uint32
	hpStones            uint16
	spStones            uint16
	curseResistance     stat
	restraintResistance stat
	poisonResistance    stat
	rollbackResistance  stat
}

type playerStatPoints struct {
	str                  uint8
	end                  uint8
	dex                  uint8
	int                  uint8
	spr                  uint8
	redistributionPoints uint8
}

type playerState struct {
	prevExp     uint64
	exp         uint64
	nextExp     uint64
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
	fame        uint32
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
	base       uint32
	withExtras uint32
}

func (p *player) load(name string, worldDB *pg.DB) error {

	char, err := character.GetByName(worldDB, name)

	if err != nil {
		return err
	}
	p.Lock()
	defer p.Unlock()

	p.char = &char

	p.location.mapName = char.Location.MapName
	p.location.mapID = int(char.Location.MapID)
	p.location.x = char.Location.X
	p.location.y = char.Location.Y
	p.location.d = char.Location.D

	p.players = make(map[uint16]*player)
	p.monsters = make(map[uint16]*monster)

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

	done := 0
	for {
		select {
		case v := <-view:
			p.view = v
			view = nil
			done++
		case st := <-state:
			p.state = st
			state = nil
			done++
		case s := <-stats:
			p.stats = s
			stats = nil
			done++
		case i := <-items:
			p.items = i
			items = nil
			done++
		case m := <-money:
			p.money = m
			money = nil
			done++
		case t := <-titles:
			p.titles = t
			titles = nil
			done++
		case q := <-quests:
			p.quests = q
			quests = nil
			done++
		case sk := <-skills:
			p.skills = sk
			skills = nil
			done++
		case pa := <-passives:
			p.passives = pa
			passives = nil
			done++
		case err := <-pErr:
			return err
		}

		if done == 8 { // risky, checking if channels are nil is safer
			return nil
		}
	}
	// for p launch routines to create player inner structs view, state, stats
}

func (p *player) viewData(view chan<- playerView, c *character.Character, err chan error) {
	v := playerView{ // todo: validation just in case, so we don't log a bad player that could potentially bin other player
		name:       c.Name,
		class:      c.Appearance.Class,
		gender:     c.Appearance.Gender,
		hairType:   c.Appearance.HairType,
		hairColour: c.Appearance.HairColor,
		faceType:   c.Appearance.FaceType,
	}
	view <- v
}

func (p *player) stateData(state chan<- playerState, c *character.Character, err chan<- error) {
	s := playerState{
		prevExp: 100,
		exp:     150,
		nextExp: 800,
		level:   c.Attributes.Level,
		// player state should also include buffs and debuffs in the future
		autoPickup:  0,
		polymorph:   65535,
		moverHandle: 0,
		moverSlot:   0,
		miniPet:     0,
	}
	state <- s
}

func (p *player) statsData(stats chan<- playerStats, c *character.Character, err chan<- error) {
	// given all:
	//  class base stats for current level, equipped items, charged buffs, buffs/debuffs, assigned stat points
	// calculate base stats (class base stats for current level, assigned stat points) , and stats with gear on (equipped items, charged buffs, buffs/debuffs)
	// given that equipped
	s := playerStats{
		str: stat{
			base:       0,
			withExtras: 0,
		},
		end: stat{
			base:       150,
			withExtras: 150,
		},
		dex: stat{
			base:       0,
			withExtras: 0,
		},
		int: stat{
			base:       0,
			withExtras: 0,
		},
		spr: stat{
			base:       0,
			withExtras: 0,
		},
		minPhysicalDamage: stat{
			base:       0,
			withExtras: 0,
		},
		maxPhysicalDamage: stat{
			base:       0,
			withExtras: 0,
		},
		minMagicalDamage: stat{
			base:       0,
			withExtras: 0,
		},
		maxMagicalDamage: stat{
			base:       0,
			withExtras: 0,
		},
		physicalDefense: stat{
			base:       0,
			withExtras: 0,
		},
		magicalDefense: stat{
			base:       0,
			withExtras: 0,
		},
		evasion: stat{
			base:       0,
			withExtras: 0,
		},
		aim: stat{
			base:       0,
			withExtras: 0,
		},
		hp:       1000,
		sp:       1000,
		lp:       0,
		hpStones: 15,
		spStones: 15,
		curseResistance: stat{
			base:       0,
			withExtras: 0,
		},
		restraintResistance: stat{
			base:       0,
			withExtras: 0,
		},
		poisonResistance: stat{
			base:       0,
			withExtras: 0,
		},
		rollbackResistance: stat{
			base:       0,
			withExtras: 0,
		},
	}

	stats <- s
}

func (p *player) itemData(items chan<- playerItems, c *character.Character, err chan<- error) {
	// for this character, load all items in each respective box
	// each item loaded should be validated so that, best way is to iterate all items and for each item launch a routine that validates it and returns the valid item through a channel
	// we also forward the error channel in case there is an error
	i := playerItems{
		equipped: itemBox{
			box: 8,
		},
		inventory: itemBox{
			box: 9,
		},
		miniHouse: itemBox{
			box: 12,
		},
		premium: itemBox{
			box: 15,
		},
	}
	items <- i
}

func (p *player) titleData(titles chan<- playerTitles, c *character.Character, err chan<- error) {
	// bit operation for titles u.u
	t := playerTitles{
		current: struct {
			id      uint8
			element uint8
			mobID   uint16
		}{
			id:      0,
			element: 0,
			mobID:   0,
		},
	}
	titles <- t
}

func (p *player) moneyData(money chan<- playerMoney, c *character.Character, err chan<- error) {
	m := playerMoney{
		coins:       100000,
		fame:        100000,
		wastedCoins: 0,
		wastedFame:  0,
	}
	money <- m
}

func (p *player) skillData(skills chan<- []skill, c *character.Character, err chan<- error) {
	// all learned skills stored in the database
	var s []skill
	skills <- s
}

func (p *player) passiveData(passives chan<- []passive, c *character.Character, err chan<- error) {
	var pa []passive
	passives <- pa
}

func (p *player) ncBriefInfoLoginCharacterCmd() structs.NcBriefInfoLoginCharacterCmd {

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
		Mode:            0,
		Class:           p.view.class,
		Shape:           *p.view.protoAvatarShapeInfo(),
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
		Unk:             1,
	}
	return nc
}
func (p *player) charParameterData() structs.CharParameterData {
	return structs.CharParameterData{
		PrevExp: p.state.prevExp,
		NextExp: p.state.nextExp,
		Strength: structs.ShineCharStatVar{
			Base:   p.stats.str.base,
			Change: p.stats.str.withExtras,
		},
		Constitute: structs.ShineCharStatVar{
			Base:   p.stats.end.base,
			Change: p.stats.end.withExtras,
		},
		Dexterity: structs.ShineCharStatVar{
			Base:   p.stats.dex.base,
			Change: p.stats.dex.withExtras,
		},
		Intelligence: structs.ShineCharStatVar{
			Base:   p.stats.int.base,
			Change: p.stats.int.withExtras,
		},
		Wisdom: structs.ShineCharStatVar{
			Base:   p.stats.spr.base,
			Change: p.stats.spr.withExtras,
		},
		MentalPower: structs.ShineCharStatVar{
			Base:   0,
			Change: 0,
		},
		WCLow: structs.ShineCharStatVar{
			Base:   p.stats.minPhysicalDamage.base,
			Change: p.stats.minPhysicalDamage.withExtras,
		},
		WCHigh: structs.ShineCharStatVar{
			Base:   p.stats.maxPhysicalDamage.base,
			Change: p.stats.maxPhysicalDamage.withExtras,
		},
		AC: structs.ShineCharStatVar{
			Base:   p.stats.physicalDefense.base,
			Change: p.stats.physicalDefense.withExtras,
		},
		TH: structs.ShineCharStatVar{
			Base:   p.stats.aim.base,
			Change: p.stats.aim.withExtras,
		},
		TB: structs.ShineCharStatVar{
			Base:   p.stats.evasion.base,
			Change: p.stats.evasion.withExtras,
		},
		MALow: structs.ShineCharStatVar{
			Base:   p.stats.minMagicalDamage.base,
			Change: p.stats.minMagicalDamage.withExtras,
		},
		MAHigh: structs.ShineCharStatVar{
			Base:   p.stats.maxMagicalDamage.base,
			Change: p.stats.maxMagicalDamage.withExtras,
		},
		MR: structs.ShineCharStatVar{
			Base:   p.stats.magicalDefense.base,
			Change: p.stats.magicalDefense.withExtras,
		},
		MH: structs.ShineCharStatVar{
			Base:   500, // ?
			Change: 500, // ?
		},
		MB: structs.ShineCharStatVar{
			Base:   500, // ?
			Change: 500, // ?
		},
		MaxHP:      p.stats.hp,
		MaxSP:      p.stats.sp,
		MaxLP:      p.stats.lp,
		MaxAP:      0, // ¿?
		MaxHPStone: p.stats.maxHpStones,
		MaxSPStone: p.stats.maxSpStones,
		PwrStone: structs.CharParameterDataPwrStone{ // ¿?
			Flag:      0,
			EPPPhysic: 0,
			EPMagic:   0,
			MaxStone:  0,
		},
		GrdStone: structs.CharParameterDataPwrStone{ // ??
			Flag:      0,
			EPPPhysic: 0,
			EPMagic:   0,
			MaxStone:  0,
		},
		PainRes: structs.ShineCharStatVar{
			Base:   p.stats.poisonResistance.base,
			Change: p.stats.poisonResistance.withExtras,
		},
		RestraintRes: structs.ShineCharStatVar{
			Base:   p.stats.restraintResistance.base,
			Change: p.stats.restraintResistance.withExtras,
		},
		CurseRes: structs.ShineCharStatVar{
			Base:   p.stats.curseResistance.base,
			Change: p.stats.curseResistance.withExtras,
		},
		ShockRes: structs.ShineCharStatVar{
			Base:   p.stats.rollbackResistance.base,
			Change: p.stats.rollbackResistance.withExtras,
		},
	}
}

func (pv *playerView) protoAvatarShapeInfo() *structs.ProtoAvatarShapeInfo {
	return &structs.ProtoAvatarShapeInfo{
		BF:        1 | pv.class<<2 | pv.gender<<7,
		HairType:  pv.hairType,
		HairColor: pv.hairColour,
		FaceShape: pv.faceType,
	}
}

func (pi *playerItems) ncCharClientItemCmd() []structs.NcCharClientItemCmd {
	var ncs []structs.NcCharClientItemCmd
	// for now empty, later on process each box type item
	ncs = []structs.NcCharClientItemCmd{
		{
			NumOfItem: 0,
			Box:       pi.equipped.box,
			Flag: structs.ProtoNcCharClientItemCmdFlag{
				BF0: 0,
			},
		},
		{
			NumOfItem: 0,
			Box:       pi.inventory.box,
			Flag: structs.ProtoNcCharClientItemCmdFlag{
				BF0: 0,
			},
		},
		{
			NumOfItem: 0,
			Box:       pi.miniHouse.box,
			Flag: structs.ProtoNcCharClientItemCmdFlag{
				BF0: 0,
			},
		},
		{
			NumOfItem: 0,
			Box:       pi.premium.box,
			Flag: structs.ProtoNcCharClientItemCmdFlag{
				BF0: 0,
			},
		},
	}
	return ncs
}
