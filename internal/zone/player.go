package zone

import (
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

const (
	maxProximityEntities = 1500
)

type player struct {
	baseEntity
	proximity   *playerProximity
	conn        *playerConnection
	view        *playerView
	stats       *playerStats
	state       *playerState
	inventories *playerInventories
	money       *playerMoney
	titles      *playerTitles
	quests      *playerQuests
	skills      *playerSkills
	targeting   *targeting
	prompt      *prompt

	ticks       *entityTicks
	persistence *playerPersistence

	// dangerZone: only to be used when loading or other situation!!
	dz *sync.RWMutex
}

type playerProximity struct {
	players map[uint16]*player
	npcs    map[uint16]*npc
	*sync.RWMutex
}

type playerConnection struct {
	lastHeartBeat time.Time
	close         chan<- bool
	outboundData  chan<- []byte
	*sync.RWMutex
}

type playerView struct {
	name       string
	class      uint8
	gender     uint8
	hairType   uint8
	hairColour uint8
	faceType   uint8
	*sync.RWMutex
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
	*sync.RWMutex
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
	justSpawned bool
	*sync.RWMutex
}

type playerMoney struct {
	coins       uint64
	fame        uint32
	wastedCoins uint64
	wastedFame  uint64
	*sync.RWMutex
}

type playerTitles struct {
	current struct {
		id      uint8
		element uint8
		mobID   uint16
	}
	titles []title
	*sync.RWMutex
}

type playerQuests struct {
	read       []quest
	done       []quest
	doing      []quest
	repeatable []quest
	*sync.RWMutex
}

type playerSkills struct {
	active  []skill
	passive []skill
	*sync.RWMutex
}

type prompt struct {
	action int
	*sync.RWMutex
}

type entityTicks struct {
	list []*time.Ticker
	*sync.RWMutex
}

type playerPersistence struct {
	char *persistence.Character
	*sync.RWMutex
}

type promptAction int

type playerStatPoints struct {
	str                  uint8
	end                  uint8
	dex                  uint8
	int                  uint8
	spr                  uint8
	redistributionPoints uint8
	*sync.RWMutex
}

type skill struct {
	id       uint16
	coolTime uint32
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

func (p *player) selectsNPC(n *npc) byte {
	var order byte
	p.targeting.Lock()
	p.targeting.selectingP = nil
	p.targeting.selectingN = n
	p.targeting.selectionOrder += 32
	order = p.targeting.selectionOrder
	p.targeting.Unlock()
	return order
}

func (p *player) selectsPlayer(ap *player) byte {
	var order byte
	p.targeting.Lock()
	p.targeting.selectingP = ap
	p.targeting.selectingN = nil

	p.targeting.selectionOrder += 32
	order = p.targeting.selectionOrder
	p.targeting.Unlock()

	ap.targeting.Lock()
	ap.targeting.selectedByP = append(ap.targeting.selectedByP, p)
	ap.targeting.Unlock()

	return order
}

func (p *player) adjacentPlayers() <-chan *player {
	ch := make(chan *player, maxProximityEntities)

	go func(p *player, send chan<- *player) {
		p.proximity.RLock()
		for _, pp := range p.proximity.players {
			send <- pp
		}
		p.proximity.RUnlock()
		close(send)
	}(p, ch)

	return ch
}

func (p *player) adjacentNpcs() <-chan *npc {
	ch := make(chan *npc, maxProximityEntities)

	go func(p *player, send chan<- *npc) {
		p.proximity.RLock()
		for _, pn := range p.proximity.npcs {
			send <- pn
		}
		p.proximity.RUnlock()
		close(send)
	}(p, ch)

	return ch
}

func (p *player) removeAdjacentPlayer(h uint16) {
	p.proximity.Lock()
	delete(p.proximity.players, h)
	p.proximity.Unlock()
}

func (p *player) selectedByPlayers() chan *player {
	ch := make(chan *player, maxProximityEntities)

	go func(p *player, send chan<- *player) {
		p.targeting.RLock()
		for _, ap := range p.targeting.selectedByP {
			send <- ap
		}
		p.targeting.RUnlock()
		close(send)
	}(p, ch)
	return ch
}

func (p *player) selectedByNPCs() chan *npc {
	ch := make(chan *npc, maxProximityEntities)

	go func(p *player, send chan<- *npc) {
		p.targeting.RLock()
		for _, n := range p.targeting.selectedByN {
			send <- n
		}
		p.targeting.RUnlock()
		close(send)
	}(p, ch)

	return ch
}

func (p *player) ncBatTargetInfoCmd() *structs.NcBatTargetInfoCmd {
	var nc = &structs.NcBatTargetInfoCmd{}

	nc.Handle = p.getHandle()

	p.stats.RLock()
	nc.TargetHP = p.stats.hp
	nc.TargetMaxHP = p.stats.maxHP
	nc.TargetSP = p.stats.sp
	nc.TargetMaxSP = p.stats.maxSP
	nc.TargetLP = p.stats.lp
	nc.TargetMaxLP = p.stats.maxLP
	p.stats.RUnlock()

	p.state.RLock()
	nc.TargetLevel = p.state.level
	p.state.RUnlock()

	return nc
}

func (p *player) load(name string) error {
	char, err := persistence.GetCharacterByName(name)

	if err != nil {
		return err
	}

	p.persistence = &playerPersistence{
		char:    &char,
		RWMutex: &sync.RWMutex{},
	}

	p.baseEntity.current.mapName = char.Location.MapName
	p.baseEntity.current.mapID = int(char.Location.MapID)
	p.baseEntity.current.x = char.Location.X
	p.baseEntity.current.y = char.Location.Y
	p.baseEntity.current.d = char.Location.D

	p.proximity = &playerProximity{
		players: make(map[uint16]*player),
		npcs:    make(map[uint16]*npc),
		RWMutex: &sync.RWMutex{},
	}

	p.ticks = &entityTicks{
		RWMutex: &sync.RWMutex{},
	}

	p.prompt = &prompt{
		RWMutex: &sync.RWMutex{},
	}

	p.targeting = &targeting{
		RWMutex: &sync.RWMutex{},
	}

	wg := &sync.WaitGroup{}
	wg.Add(7)

	errC := make(chan error, 7)

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
	go func(err chan<- error) {
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

	// check if you can iterate over the channel

	wg.Wait()

	return <-errC
}

func (p *player) itemData() error {
	// for this character, load all items in each respective box
	// each item loaded should be validated so that, best way is to iterate all items and for each item launch a routine that validates it and returns the valid item through a channel
	// we also forward the error channel in case there is an error
	var ivs = &playerInventories{
		RWMutex: &sync.RWMutex{},
	}

	eiBox, err := loadInventory(persistence.EquippedInventory, p)
	if err != nil {
		return err
	}

	biBox, err := loadInventory(persistence.BagInventory, p)
	if err != nil {
		return err
	}

	mhiBox, err := loadInventory(persistence.MiniHouseInventory, p)
	if err != nil {
		return err
	}

	ivs.equipped = eiBox
	ivs.inventory = biBox
	ivs.miniHouse = mhiBox

	p.dz.Lock()
	p.inventories = ivs
	p.dz.Unlock()

	return nil
}

func (p *player) viewData() {
	v := &playerView{ // todo: validation just in case, so we don't log a bad player that could potentially bin other player
		name:       p.persistence.char.Name,
		class:      p.persistence.char.Appearance.Class,
		gender:     p.persistence.char.Appearance.Gender,
		hairType:   p.persistence.char.Appearance.HairType,
		hairColour: p.persistence.char.Appearance.HairColor,
		faceType:   p.persistence.char.Appearance.FaceType,
		RWMutex:    &sync.RWMutex{},
	}
	p.dz.Lock()
	p.view = v
	p.dz.Unlock()
}

func (p *player) stateData() {
	s := &playerState{
		prevExp: 100,
		exp:     150,
		nextExp: 800,
		level:   p.persistence.char.Attributes.Level,
		// player state should also include buffs and debuffs in the future
		autoPickup:  0,
		polymorph:   65535,
		moverHandle: 0,
		moverSlot:   0,
		miniPet:     0,
		RWMutex:     &sync.RWMutex{},
	}
	p.dz.Lock()
	p.state = s
	p.dz.Unlock()
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
		RWMutex: &sync.RWMutex{},
	}
	p.dz.Lock()
	p.stats = s
	p.dz.Unlock()
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
		RWMutex: &sync.RWMutex{},
	}
	p.dz.Lock()
	p.titles = t
	p.dz.Unlock()
}

func (p *player) moneyData() {
	m := &playerMoney{
		coins:       100000,
		fame:        100000,
		wastedCoins: 0,
		wastedFame:  0,
		RWMutex:     &sync.RWMutex{},
	}
	p.dz.Lock()
	p.money = m
	p.dz.Unlock()
}

func (p *player) skillData() {
	// all learned skills stored in the database
	p.dz.Lock()
	p.skills = &playerSkills{
		RWMutex: &sync.RWMutex{},
	}
	p.dz.Unlock()
}

func (p *player) charParameterData() structs.CharParameterData {
	p.stats.RLock()
	p.state.RLock()
	nc := structs.CharParameterData{
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
	p.stats.RUnlock()
	p.state.RUnlock()
	return nc
}

func (p *player) equip(nc * structs.NcItemEquipReq) (itemSlotChange, error) {
	var (
		change itemSlotChange
		fromItem *item
		toItem *item
		slot = int(nc.Slot)
	)

	item := p.inventories.get(persistence.BagInventory, slot)
	if item == nil {
		p.persistence.RLock()
		characterName := p.persistence.char.Name
		p.persistence.RUnlock()
		return change, errors.Err{
			Code:    errors.ZoneItemSlotEquipNoItem,
			Details: errors.ErrDetails{
				"slot": nc.Slot,
				"handle": p.getHandle(),
				"characterName": characterName,
			},
		}
	}

	fromItem = item

	// slot that will be occupied
	equip := int(item.itemData.itemInfo.Equip)

	equippedItem := p.inventories.get(persistence.EquippedInventory, equip)

	if equippedItem != nil {
		toItem = equippedItem
	}

	opItem, err := item.pItem.MoveTo(persistence.EquippedInventory, equip)

	if err != nil {
		return change ,err
	}

	p.inventories.Lock()

	if toItem != nil {
		toItem.Lock()
		toItem.pItem = opItem
		toItem.Unlock()
		delete(p.inventories.equipped.items, equip)
		p.inventories.inventory.items[slot] = toItem
	}

	p.inventories.equipped.items[equip] = fromItem
	delete(p.inventories.inventory.items, slot)

	p.inventories.Unlock()

	change = itemSlotChange{
		gameFrom: uint16(persistence.BagInventory) << 10 | uint16(slot) & 1023,
		gameTo:   uint16(item.pItem.InventoryType << 10 | item.pItem.Slot & 1023),
		from:     itemSlot{
			slot:         slot,
			inventoryType: persistence.BagInventory,
			item: fromItem,
		},
		to:       itemSlot{
			slot:          equip,
			inventoryType: persistence.EquippedInventory,
			item:          toItem,
		},
	}
	//

	return change, nil
}

func (p *player) newItem(i *item) error {
	if i.pItem == nil {
		i.pItem = &persistence.Item{}
	}
	i.pItem.CharacterID = p.persistence.char.ID
	i.pItem.ShnID = i.itemData.itemInfo.ID
	i.pItem.ShnInxName = i.itemData.itemInfo.InxName
	i.pItem.Amount = i.amount
	i.pItem.Stackable = i.stackable

	i.pItem.Attributes = &persistence.ItemAttributes{}

	i.pItem.Attributes.StrengthBase = i.stats.strength.base
	i.pItem.Attributes.StrengthExtra = i.stats.strength.extra

	i.pItem.Attributes.DexterityBase = i.stats.dexterity.base
	i.pItem.Attributes.DexterityExtra = i.stats.dexterity.extra

	i.pItem.Attributes.IntelligenceBase = i.stats.intelligence.base
	i.pItem.Attributes.IntelligenceExtra = i.stats.intelligence.extra

	i.pItem.Attributes.EnduranceBase = i.stats.endurance.base
	i.pItem.Attributes.EnduranceExtra = i.stats.endurance.extra

	i.pItem.Attributes.SpiritBase = i.stats.spirit.base
	i.pItem.Attributes.SpiritExtra = i.stats.spirit.extra

	i.pItem.Attributes.PAttackBase = i.stats.physicalAttack.base
	i.pItem.Attributes.PAttackExtra = i.stats.physicalAttack.extra

	i.pItem.Attributes.MAttackBase = i.stats.magicalAttack.base
	i.pItem.Attributes.MAttackExtra = i.stats.magicalAttack.extra

	i.pItem.Attributes.MDefenseBase = i.stats.magicalDefense.base
	i.pItem.Attributes.MDefenseExtra = i.stats.magicalDefense.extra

	i.pItem.Attributes.PDefenseBase = i.stats.physicalDefense.base
	i.pItem.Attributes.PDefenseExtra = i.stats.physicalDefense.extra

	i.pItem.Attributes.AimBase = i.stats.aim.base
	i.pItem.Attributes.AimExtra = i.stats.aim.extra

	i.pItem.Attributes.EvasionBase = i.stats.evasion.base
	i.pItem.Attributes.EvasionExtra = i.stats.evasion.extra

	i.pItem.Attributes.MaxHPBase = i.stats.maxHP.base
	i.pItem.Attributes.MaxHPBase = i.stats.maxHP.extra

	err := i.pItem.Insert()

	if err != nil {
		return err
	}

	p.inventories.Lock()
	p.inventories.inventory.items[i.pItem.Slot] = i
	p.inventories.Unlock()

	return nil
}

func loadInventory(it persistence.InventoryType, p *player) (itemBox, error) {
	var box itemBox
	items, err := persistence.GetCharacterItems(int(p.persistence.char.ID), it)

	if err != nil {
		return box, err
	}

	box.box = int(it)
	box.items = make(map[int]*item)
	for _, item := range items {
		// load with goroutines and waitgroups
		box.items[item.Slot] = loadItem(item)
	}
	return box, nil
}

func lastHeartbeat(p *player) float64 {
	p.conn.RLock()
	lastHeartBeat := time.Since(p.conn.lastHeartBeat).Seconds()
	p.conn.RUnlock()
	return lastHeartBeat
}

func protoAvatarShapeInfo(pv *playerView) *structs.ProtoAvatarShapeInfo {
	pv.RLock()
	nc := &structs.ProtoAvatarShapeInfo{
		BF:        1 | pv.class<<2 | pv.gender<<7,
		HairType:  pv.hairType,
		HairColor: pv.hairColour,
		FaceShape: pv.faceType,
	}
	pv.RUnlock()
	return nc
}

func ncBriefInfoLoginCharacterCmd(p *player) structs.NcBriefInfoLoginCharacterCmd {

	var nc = structs.NcBriefInfoLoginCharacterCmd{
		Mode: 2,
	}

	nc.Handle = p.baseEntity.getHandle()

	p.baseEntity.RLock()
	nc.Coordinates = structs.ShineCoordType{
		XY: structs.ShineXYType{
			X: uint32(p.baseEntity.current.x),
			Y: uint32(p.baseEntity.current.y),
		},
		Direction: byte(p.baseEntity.current.d),
	}
	p.baseEntity.RUnlock()

	nc.Shape = *protoAvatarShapeInfo(p.view)

	p.view.RLock()
	nc.CharID = structs.Name5{
		Name: p.view.name,
	}
	p.view.RUnlock()

	p.state.RLock()
	nc.Class = p.view.class
	nc.Polymorph = p.state.polymorph
	nc.Level = p.state.level
	nc.MoverHandle = p.state.moverHandle
	nc.MoverSlot = p.state.moverSlot
	nc.UsingMinipet = p.state.miniPet
	p.state.RUnlock()

	return nc
}
