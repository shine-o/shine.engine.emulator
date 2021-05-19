package zone

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/persistence"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
	"reflect"
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
	*baseEntity
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
	sync.RWMutex
}

func (p *player) alreadyNearbyEntity(e entity) bool {
	p.baseEntity.proximity.RLock()
	_, exists := p.baseEntity.proximity.entities[e.getHandle()]
	p.baseEntity.proximity.RUnlock()
	return exists
}

func (p *player) newNearbyEntitiesTicker(zm *zoneMap) {
	log.Infof("[player_ticks] newNearbyEntitiesTicker for handle %v", p.getHandle())
	tick := time.NewTicker(200 * time.Millisecond)
	p.ticks.Lock()
	p.ticks.list = append(p.ticks.list, tick)
	p.ticks.Unlock()
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			addWithinRangeEntities(p, zm)
		}
	}
}

func (p *player) oldNearbyEntitiesTicker() {
	log.Infof("[player_ticks] oldNearbyEntitiesTicker for handle %v", p.getHandle())
	tick := time.NewTicker(200 * time.Millisecond)
	p.ticks.Lock()
	p.ticks.list = append(p.ticks.list, tick)
	p.ticks.Unlock()
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			removeOutOfRangeEntities(p)
		}
	}
}

func (p *player) getPacketData() interface{} {
	return ncBriefInfoLoginCharacterCmd(p)
}

func (p *player) notifyAboutNewEntity(e entity) {
	switch e.(type) {
	case *player:
		networking.Send(p.conn.outboundData, networking.NC_BRIEFINFO_LOGINCHARACTER_CMD, e.getPacketData())
		break
	case *npc:
		networking.Send(p.conn.outboundData, networking.NC_BRIEFINFO_REGENMOB_CMD, e.getPacketData())
		break
	default:
		log.Errorf("unknown entity type %v", reflect.TypeOf(e).String())
	}
}

func (p *player) notifyAboutRemovedEntity(e entity) {
	nc := &structs.NcBriefInfoDeleteHandleCmd{
		Handle: e.getHandle(),
	}
	networking.Send(p.conn.outboundData, networking.NC_BRIEFINFO_BRIEFINFODELETE_CMD, nc)
}

func (p *player) getNearbyEntities() <-chan entity {
	return getNearbyEntities(p.baseEntity.proximity)
}

func (p *player) removeNearbyEntity(e entity) {
	p.Lock()
	delete(p.baseEntity.proximity.entities, e.getHandle())
	p.Unlock()
}

func (p *player) addNearbyEntity(e entity) {
	h := e.getHandle()
	p.baseEntity.proximity.Lock()
	p.baseEntity.proximity.entities[h] = e
	p.baseEntity.proximity.Unlock()
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
	points              playerStatPoints
	str                 stat
	end                 stat
	dex                 stat
	int                 stat
	spr                 stat
	minPhysicalDamage   stat
	maxPhysicalDamage   stat
	minMagicalDamage    stat
	maxMagicalDamage    stat
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

type playerCondition int

const (
	camping playerCondition = iota
	normal
	vendor
	riding
)

type playerState struct {
	current     playerCondition
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

type playerSkills struct {
	active  []skill
	passive []skill
	sync.RWMutex
}

type prompt struct {
	action int
	sync.RWMutex
}

type entityTicks struct {
	list []*time.Ticker
	sync.RWMutex
}

type playerPersistence struct {
	char *persistence.Character
	sync.RWMutex
}

type promptAction int

type playerStatPoints struct {
	str                  uint8
	end                  uint8
	dex                  uint8
	int                  uint8
	spr                  uint8
	redistributionPoints uint8
	sync.RWMutex
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
		char: &char,
	}

	p.baseEntity.events = events{
		send: make(sendEvents),
		recv: make(recvEvents),
	}

	for _, ev := range playerEvents {
		c := make(chan event, 5)
		p.baseEntity.events.send[ev] = c
		p.baseEntity.events.recv[ev] = c
	}

	go p.eduPlayerEvents001()
	go p.eduPlayerEvents002()

	p.baseEntity.current.mapName = char.Location.MapName
	p.baseEntity.current.mapID = int(char.Location.MapID)
	p.baseEntity.current.x = char.Location.X
	p.baseEntity.current.y = char.Location.Y
	p.baseEntity.current.d = char.Location.D

	p.baseEntity.proximity = &entityProximity{
		entities: make(map[uint16]entity),
	}

	p.ticks = &entityTicks{}

	p.prompt = &prompt{}

	p.targeting = &targeting{}

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
	var ivs = &playerInventories{}

	eiBox, err := loadInventory(persistence.EquippedInventory, p)
	if err != nil {
		return err
	}

	biBox, err := loadInventory(persistence.BagInventory, p)
	if err != nil {
		return err
	}

	diBox, err := loadInventory(persistence.DepositInventory, p)
	if err != nil {
		return err
	}

	mhiBox, err := loadInventory(persistence.MiniHouseInventory, p)
	if err != nil {
		return err
	}

	riBox, err := loadInventory(persistence.RewardInventory, p)
	if err != nil {
		return err
	}

	ivs.equipped = eiBox
	ivs.inventory = biBox
	ivs.miniHouse = mhiBox
	ivs.deposit = diBox
	ivs.reward = riBox

	p.Lock()
	p.inventories = ivs
	p.Unlock()

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
	}
	p.Lock()
	p.view = v
	p.Unlock()
}

func (p *player) stateData() {
	s := &playerState{
		current: normal,
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
		justSpawned: false,
	}
	p.Lock()
	p.state = s
	p.Unlock()
}

func (p *player) statsData() {
	// given all:
	//  class base stats for current level, equippedID items, charged buffs, buffs/debuffs, assigned stat points
	// calculate base stats (class base stats for current level, assigned stat points) , and stats with gear on (equippedID items, charged buffs, buffs/debuffs)
	// given that equippedID
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
	p.skills = &playerSkills{}
	p.Unlock()
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

func canBeEquipped(equip int, class int) bool {
	if equip >= 1 && equip <= 29 {
		return true
	}

	switch data.ItemClassEnum(class) {
	case data.ItemBracelet:
	case data.ItemCosShield:
	case data.ItemClassBoot:
	case data.ItemClassShield:
	case data.ItemClassArmor:
	case data.ItemClassWeapon:
	case data.ItemClassAmulet:
		return true
	}

	return false
}

func (p *player) equip(slot int) (itemSlotChange, error) {
	var (
		change   itemSlotChange
		fromItem *item
		toItem   *item
	)

	fromItem, _ = p.inventories.get(persistence.BagInventory, slot)
	if fromItem == nil {
		p.persistence.RLock()
		characterName := p.persistence.char.Name
		p.persistence.RUnlock()
		return change, errors.Err{
			Code: errors.ZoneItemSlotEquipNoItem,
			Details: errors.ErrDetails{
				"slot":          slot,
				"handle":        p.getHandle(),
				"characterName": characterName,
			},
		}
	}

	// slot that will be occupied
	equip := int(fromItem.itemData.itemInfo.Equip)
	class := int(fromItem.itemData.itemInfo.Class)

	if !canBeEquipped(equip, class) {
		return change, errors.Err{
			Code: errors.ZoneItemEquipBadType,
			Details: errors.ErrDetails{
				"slot":   slot,
				"handle": p.getHandle(),
				"equip":  equip,
			},
		}
	}

	equippedItem, _ := p.inventories.get(persistence.EquippedInventory, equip)

	if equippedItem != nil {
		toItem = equippedItem
	}

	opItem, err := fromItem.pItem.MoveTo(persistence.EquippedInventory, equip)

	if err != nil {
		return change, err
	}

	p.inventories.Lock()
	if toItem != nil {
		toItem.Lock()
		toItem.pItem = opItem
		toItem.Unlock()
		delete(p.inventories.equipped.items, equip)
		p.inventories.inventory.items[slot] = toItem
	} else {
		delete(p.inventories.inventory.items, slot)
	}
	p.inventories.equipped.items[equip] = fromItem
	p.inventories.Unlock()

	change = itemSlotChange{
		gameFrom: uint16(persistence.BagInventory)<<10 | uint16(slot)&1023,
		gameTo:   uint16(fromItem.pItem.InventoryType<<10 | fromItem.pItem.Slot&1023),
		from: itemSlot{
			slot:          slot,
			inventoryType: persistence.BagInventory,
			item:          fromItem,
		},
		to: itemSlot{
			slot:          equip,
			inventoryType: persistence.EquippedInventory,
			item:          toItem,
		},
	}

	return change, nil
}

func (p *player) unEquip(from, to int) (itemSlotChange, error) {
	var (
		change   itemSlotChange
		fromItem *item
		toItem   *item
	)

	fromItem, _ = p.inventories.get(persistence.EquippedInventory, from)
	if fromItem == nil {
		p.persistence.RLock()
		characterName := p.persistence.char.Name
		p.persistence.RUnlock()
		return change, errors.Err{
			Code: errors.ZoneItemSlotEquipNoItem,
			Details: errors.ErrDetails{
				"equip":         from,
				"handle":        p.getHandle(),
				"characterName": characterName,
			},
		}
	}

	toItem, _ = p.inventories.get(persistence.BagInventory, to)
	if toItem != nil {
		p.persistence.RLock()
		characterName := p.persistence.char.Name
		p.persistence.RUnlock()
		return change, errors.Err{
			Code: errors.ZoneItemSlotInUse,
			Details: errors.ErrDetails{
				"equip":         from,
				"handle":        p.getHandle(),
				"characterName": characterName,
			},
		}
	}

	_, err := fromItem.pItem.MoveTo(persistence.BagInventory, to)

	if err != nil {
		return change, err
	}

	p.inventories.Lock()
	delete(p.inventories.equipped.items, from)
	p.inventories.inventory.items[to] = fromItem
	p.inventories.Unlock()

	change = itemSlotChange{
		gameFrom: uint16(persistence.EquippedInventory)<<10 | uint16(from)&1023,
		gameTo:   uint16(persistence.BagInventory)<<10 | uint16(to)&1023,
		from: itemSlot{
			slot:          from,
			inventoryType: persistence.BagInventory,
			item:          fromItem,
		},
		to: itemSlot{
			slot:          to,
			inventoryType: persistence.EquippedInventory,
			item:          toItem,
		},
	}

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

	return p.inventories.add(persistence.InventoryType(i.pItem.InventoryType), i.pItem.Slot, i)
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

// will return either the itemID or 65535
func equippedID(pi *playerInventories, equip data.ItemEquipEnum) uint16 {
	var id uint16 = 65535
	pi.RLock()
	item, ok := pi.equipped.items[int(equip)]
	if ok {
		id = item.itemData.itemInfo.ID
	}
	pi.RUnlock()
	return id
}

func justSpawned(p *player) bool {
	p.state.RLock()
	defer p.state.RUnlock()
	return p.state.justSpawned
}

func ncBriefInfoLoginCharacterCmd(p *player) *structs.NcBriefInfoLoginCharacterCmd {

	var nc = &structs.NcBriefInfoLoginCharacterCmd{
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
	nc.ShapeData = shapeData(p)

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

func shapeData(p *player) structs.NcBriefInfoLoginCharacterCmdShapeData {
	var (
		nc        structs.NcBriefInfoLoginCharacterCmdShapeData
		shapeData []byte
	)

	switch p.state.current {
	case vendor:
		//var inc = structs.CharBriefInfoBooth{
		//	Camp: structs.CharBriefInfoCamp{
		//		MiniHouse: 0,
		//		Dummy:     [10]byte{},
		//	},
		//	IsSelling: 1,
		//	SignBoard: structs.StreetBoothSignBoard{
		//		Text: "tutti frutti",
		//	},
		//}
		//d, err := structs.Pack(&inc)
		//if err != nil {
		//	log.Error(err)
		//	break
		//}
		//shapeData = d
		break
	case camping:
		//struct CHARBRIEFINFO_CAMP
		//type CharBriefInfoCamp struct {
		//	MiniHouse uint16
		//	Dummy     [10]byte //
		//}
		break
	case normal:
		inc := structs.CharBriefInfoNotCamp{
			Equip: structs.ProtoEquipment{
				EquHead:         equippedID(p.inventories, data.ItemEquipHat),
				EquMouth:        equippedID(p.inventories, data.ItemEquipMouth),
				EquRightHand:    equippedID(p.inventories, data.ItemEquipRightHand),
				EquBody:         equippedID(p.inventories, data.ItemEquipBody),
				EquLeftHand:     equippedID(p.inventories, data.ItemEquipLeftHand),
				EquPant:         equippedID(p.inventories, data.ItemEquipLeg),
				EquBoot:         equippedID(p.inventories, data.ItemEquipShoes),
				EquAccBoot:      equippedID(p.inventories, data.ItemEquipShoesAcc),
				EquAccPant:      equippedID(p.inventories, data.ItemEquipLegAcc),
				EquAccBody:      equippedID(p.inventories, data.ItemEquipBodyAcc),
				EquAccHeadA:     equippedID(p.inventories, data.ItemEquipHatAcc),
				EquMinimonR:     equippedID(p.inventories, data.ItemEquipMinimonR),
				EquEye:          equippedID(p.inventories, data.ItemEquipEye),
				EquAccLeftHand:  equippedID(p.inventories, data.ItemEquipLeftHandAcc),
				EquAccRightHand: equippedID(p.inventories, data.ItemEquipRightHandAcc),
				EquAccBack:      equippedID(p.inventories, data.ItemEquipBack),
				EquCosEff:       equippedID(p.inventories, data.ItemEquipCosEff),
				EquAccHip:       equippedID(p.inventories, data.ItemEquipTail),
				EquMinimon:      equippedID(p.inventories, data.ItemEquipMinimon),
				EquAccShield:    equippedID(p.inventories, data.ItemEquipShieldAcc),
				Upgrade: structs.EquipmentUpgrade{
					Gap: [2]uint8{0, 12},
					//Gap: 12,
					BF2: 1,
				},
			},
		}
		d, err := structs.Pack(&inc)
		if err != nil {
			log.Error(err)
			break
		}
		shapeData = d
		break
	case riding:
		////struct CHARBRIEFINFO_RIDE
		//type CharBriefInfoRide struct {
		//	Equip    ProtoEquipment
		//	RideInfo CharBriefInfoRideInfo
		//}
		break
	}

	for i, d := range shapeData {
		nc.Data[i] = d
	}

	return nc
}
