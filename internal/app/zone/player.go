package zone

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/persistence"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
	"sync"
	"time"
)

const (
	portal promptAction = iota
	partnerSummon
	deleteItem
	sellItem
	equipItem
	addFriend
	joinGuild
	joinParty
	joinExpedition
	duel
	enhanceItem
)

// todo: use pointers for high use variables (state, skills, passives, targetting, etc..), locking the whole entity is gonna be costly for latency
type player struct {
	baseEntity
	//todo: group into another struct that can be locked
	players     map[uint16]*player
	monsters    map[uint16]*monster
	npcs        map[uint16]*npc
	conn        *playerConnection
	view        *playerView
	stats       *playerStats
	state       *playerState
	inventories *playerInventories
	money       *playerMoney
	titles      *playerTitles
	quests      *playerQuests
	skills      []skill
	passives    []passive
	*targeting
	*prompt
	justSpawned bool
	tickers     []*time.Ticker
	char        *persistence.Character
	sync.RWMutex

}

type promptAction int

type prompt struct {
	action int
	sync.RWMutex
}

func (p *player) selectsNPC(n *npc) byte {
	var order byte
	p.Lock()
	p.targeting.selectingP = nil
	p.targeting.selectingM = nil
	p.targeting.selectingN = n
	p.targeting.selectionOrder += 32
	order = p.targeting.selectionOrder
	p.Unlock()
	//ep.Lock()
	//ep.targeting.selectedByP = append(ep.targeting.selectedByP, p)
	//ep.Unlock()
	return order
}

func (p *player) selectsPlayer(ap *player) byte {
	var order byte
	p.Lock()
	p.targeting.selectingP = ap
	p.targeting.selectingM = nil
	p.targeting.selectingN = nil

	p.targeting.selectionOrder += 32
	order = p.targeting.selectionOrder
	p.Unlock()

	ap.Lock()
	ap.targeting.selectedByP = append(ap.targeting.selectedByP, p)
	ap.Unlock()

	return order
}

func (p *player) selectsMonster(m *monster) byte {
	var order byte
	p.Lock()
	p.targeting.selectingP = nil
	p.targeting.selectingM = m
	p.targeting.selectingN = nil
	p.targeting.selectionOrder += 32
	order = p.targeting.selectionOrder
	p.Unlock()
	//ep.Lock()
	//ep.targeting.selectedByP = append(ep.targeting.selectedByP, p)
	//ep.Unlock()
	return order
}

func (p *player) getHandle() uint16 {
	p.RLock()
	h := p.handle
	p.RUnlock()
	return h
}

func (p *player) spawned() bool {
	p.RLock()
	defer p.RUnlock()

	return p.justSpawned
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

func (p *player) selectedByMonsters() chan *monster {
	p.RLock()
	ch := make(chan *monster, len(p.targeting.selectedByM))
	p.RUnlock()

	go func(send chan<- *monster) {
		p.RLock()
		for _, m := range p.targeting.selectedByM {
			send <- m
		}
		p.RUnlock()
		close(send)
	}(ch)
	return ch
}

func (p *player) selectedByPlayers() chan *player {
	p.RLock()
	ch := make(chan *player, len(p.targeting.selectedByP))
	p.RUnlock()

	go func(send chan<- *player) {
		p.RLock()
		for _, ap := range p.targeting.selectedByP {
			send <- ap
		}
		p.RUnlock()
		close(send)
	}(ch)
	return ch
}

func (p *player) selectedByNPCs() chan *npc {
	p.RLock()
	ch := make(chan *npc, len(p.targeting.selectedByN))
	p.RUnlock()

	go func(send chan<- *npc) {
		p.RLock()
		for _, n := range p.targeting.selectedByN {
			send <- n
		}
		p.RUnlock()
		close(send)
	}(ch)
	return ch
}

func (p *player) ncBatTargetInfoCmd() *structs.NcBatTargetInfoCmd {
	var nc structs.NcBatTargetInfoCmd
	p.RLock()
	nc = structs.NcBatTargetInfoCmd{
		Order:         0,
		Handle:        p.handle,
		TargetHP:      p.stats.hp,
		TargetMaxHP:   p.stats.maxHP,
		TargetSP:      p.stats.sp,
		TargetMaxSP:   p.stats.maxSP,
		TargetLP:      p.stats.lp,
		TargetMaxLP:   p.stats.maxLP,
		TargetLevel:   p.state.level,
		HpChangeOrder: 0,
	}
	p.RUnlock()
	return &nc
}

func lastHeartbeat(p *player) float64 {
	p.RLock()
	lastHeartBeat := time.Since(p.conn.lastHeartBeat).Seconds()
	p.RUnlock()
	return lastHeartBeat
}

type playerConnection struct {
	lastHeartBeat time.Time
	close         chan<- bool
	outboundData  chan<- []byte
	sync.RWMutex

}

type playerView struct {
	name       string
	class      uint8
	gender     uint8
	hairType   uint8
	hairColour uint8
	faceType   uint8
	sync.RWMutex
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
	maxHP               uint32
	maxSP               uint32
	maxLP               uint32
	maxHpStones         uint32
	maxSpStones         uint32
	hpStones            uint16
	spStones            uint16
	curseResistance     stat
	restraintResistance stat
	poisonResistance    stat
	rollbackResistance  stat
	sync.RWMutex
}

type playerStatPoints struct {
	str                  uint8
	end                  uint8
	dex                  uint8
	int                  uint8
	spr                  uint8
	redistributionPoints uint8
	sync.RWMutex
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
	sync.RWMutex
}

type playerMoney struct {
	coins       uint64
	fame        uint32
	wastedCoins uint64
	wastedFame  uint64
	sync.RWMutex
}

type playerTitles struct {
	current struct {
		id      uint8
		element uint8
		mobID   uint16
	}
	titles []title
	sync.RWMutex
}

type playerQuests struct {
	read       []quest
	done       []quest
	doing      []quest
	repeatable []quest
	sync.RWMutex
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

type stat struct {
	base       uint32
	withExtras uint32
}

type itemSlotChange struct {
	from int
	to   int
}

func (p *player) load(name string) error {
	char, err := persistence.GetCharacterByName(name)

	if err != nil {
		return err
	}

	p.char = &char

	p.current.mapName = char.Location.MapName
	p.current.mapID = int(char.Location.MapID)
	p.current.x = char.Location.X
	p.current.y = char.Location.Y
	p.current.d = char.Location.D

	p.players = make(map[uint16]*player)
	p.monsters = make(map[uint16]*monster)
	p.npcs = make(map[uint16]*npc)

	p.prompt = &prompt{}

	wg := &sync.WaitGroup{}
	wg.Add(8)

	errC := make(chan error, 8)

	go func() {
		defer wg.Done()
		p.viewData()
	}()
	go func() {
		defer wg.Done()
		p.stateData()
	}()
	go func() {
		defer wg.Done()
		p.statsData()
	}()
	go func(err chan <- error) {
		defer wg.Done()
		err <- p.itemData()
	}(errC)
	go func() {
		defer wg.Done()
		p.moneyData()
	}()
	go func() {
		defer wg.Done()
		p.titleData()
	}()
	go func() {
		defer wg.Done()
		p.skillData()
	}()
	go func() {
		defer wg.Done()
		p.passiveData()
	}()

	wg.Wait()

	return <- errC
}

func (p *player) viewData() {
	v := &playerView{ // todo: validation just in case, so we don't log a bad player that could potentially bin other player
		name:       p.char.Name,
		class:      p.char.Appearance.Class,
		gender:     p.char.Appearance.Gender,
		hairType:   p.char.Appearance.HairType,
		hairColour: p.char.Appearance.HairColor,
		faceType:   p.char.Appearance.FaceType,
	}
	p.Lock()
	p.view = v
	p.Unlock()
}

func (p *player) stateData() {
	s := &playerState{
		prevExp: 100,
		exp:     150,
		nextExp: 800,
		level:   p.char.Attributes.Level,
		// player state should also include buffs and debuffs in the future
		autoPickup:  0,
		polymorph:   65535,
		moverHandle: 0,
		moverSlot:   0,
		miniPet:     0,
	}
	p.Lock()
	p.state = s
	p.Unlock()
}

func (p *player) statsData() {
	// given all:
	//  class base stats for current level, equipped items, charged buffs, buffs/debuffs, assigned stat points
	// calculate base stats (class base stats for current level, assigned stat points) , and stats with gear on (equipped items, charged buffs, buffs/debuffs)
	// given that equipped
	s := &playerStats{
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
		// todo: remove magick :(
		hp:       1000,
		sp:       1000,
		maxHP:    1000,
		maxSP:    1000,
		lp:       4294967295,
		maxLP:    4294967295,
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
	p.Lock()
	p.stats = s
	p.Unlock()
}

func (p *player) titleData() {
	// bit operation for titles u.u
	t := &playerTitles{
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
	p.Lock()
	p.titles = t
	p.Unlock()
}

func (p *player) moneyData() {
	m := &playerMoney{
		coins:       100000,
		fame:        100000,
		wastedCoins: 0,
		wastedFame:  0,
	}
	p.Lock()
	p.money = m
	p.Unlock()
}

func (p *player) skillData() {
	// all learned skills stored in the database
	p.Lock()
	p.skills = []skill{}
	p.Unlock()
}

func (p *player) passiveData() {
	p.Lock()
	p.passives = []passive{}
	p.Unlock()
}

func (p *player) ncBriefInfoLoginCharacterCmd() structs.NcBriefInfoLoginCharacterCmd {

	nc := structs.NcBriefInfoLoginCharacterCmd{
		Handle: p.getHandle(),
		CharID: structs.Name5{
			Name: p.view.name,
		},
		Coordinates: structs.ShineCoordType{
			XY: structs.ShineXYType{
				X: uint32(p.current.x),
				Y: uint32(p.current.y),
			},
			Direction: byte(p.current.d),
		},
		Mode:            2,
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
		MaxHP:      p.stats.maxHP,
		MaxSP:      p.stats.maxSP,
		MaxLP:      p.stats.maxLP,
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

func (p *player) equip(i *item, slot data.ItemEquipEnum) (itemSlotChange, error) {
	slotChange := itemSlotChange{
		from: i.pItem.Slot,
		to:   0,
	}

	uItem, err := i.pItem.MoveTo(persistence.EquippedInventory, int(slot))

	if err != nil {
		return itemSlotChange{}, errors.Err{
			Code: errors.ZoneItemEquipFailed,
			Details: errors.ErrDetails{
				"err":     err,
				"pHandle": p.handle,
			},
		}
	}

	slotChange.to = int(slot)

	p.inventories.Lock()
	i.pItem = uItem
	p.inventories.Unlock()

	return slotChange, nil
}

func (p *player) newItem(i *item) error {
	i.pItem = &persistence.Item{}
	i.pItem.CharacterID = p.char.ID
	i.pItem.ShnID = i.itemData.itemInfo.ID
	i.pItem.ShnInxName = i.itemData.itemInfo.InxName
	i.pItem.Amount = i.amount
	i.pItem.Stackable = i.stackable

	i.pItem.Attributes = &persistence.ItemAttributes{}

	i.pItem.Attributes.StrengthBase = i.stats.strength.base
	i.pItem.Attributes.StrengthExtra = i.stats.strength.extra

	err := i.pItem.Insert()

	if err != nil {
		return err
	}

	p.inventories.Lock()
	p.inventories.inventory.items[i.pItem.Slot] = i
	p.inventories.Unlock()

	return nil
}

func (pi *playerInventories) ncCharClientItemCmd() []structs.NcCharClientItemCmd {
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
