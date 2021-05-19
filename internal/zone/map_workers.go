package zone

import (
	"fmt"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
	"reflect"
	"strings"
)

const (
	runSpeed             = 120
	walkSpeed            = 60
	playerHeartbeatLimit = 10
)

func (zm *zoneMap) mapHandles() {
	log.Infof("[map_worker] mapHandles worker for map %v", zm.data.Info.MapName)
	for {
		select {
		case <-zm.events.recv[playerHandleMaintenance]:
			go playerHandleMaintenanceLogic(zm)
		case e := <-zm.events.recv[playerHandle]:
			go playerHandleLogic(e, zm)
		}
	}
}

func (zm *zoneMap) playerActivity() {
	log.Infof("[map_worker] playerActivity worker for map %v", zm.data.Info.MapName)
	for {
		select {
		case e := <-zm.events.recv[playerAppeared]:
			go playerAppearedLogic(e, zm)
		case e := <-zm.events.recv[playerDisappeared]:
			go playerDisappearedLogic(e, zm)
		case e := <-zm.events.recv[playerWalks]:
			go playerWalksLogic(e, zm)
		case e := <-zm.events.recv[playerRuns]:
			go playerRunsLogic(e, zm)
		case e := <-zm.events.recv[playerStopped]:
			go playerStoppedLogic(e, zm)
		case e := <-zm.events.recv[playerJumped]:
			go playerJumpedLogic(e, zm)
		case e := <-zm.events.recv[unknownHandle]:
			go unknownHandleLogic(e, zm)
		case e := <-zm.events.recv[playerSelectsEntity]:
			go playerSelectsEntityLogic(zm, e)
		case e := <-zm.events.recv[playerUnselectsEntity]:
			go playerUnselectsEntityLogic(zm, e)
		case e := <-zm.events.recv[itemIsMoved]:
			go itemIsMovedLogic(e, zm)
		case e := <-zm.events.recv[itemEquip]:
			go itemEquipLogic(e, zm)
		case e := <-zm.events.recv[itemUnEquip]:
			go itemUnEquipLogic(e, zm)
		}
	}
}

func (zm *zoneMap) npcInteractions() {
	log.Infof("[map_worker] npcInteractions worker for map %v", zm.data.Info.MapName)
	for {
		select {
		case e := <-zm.events.recv[playerClicksOnNpc]:
			go playerClicksOnNpcLogic(zm, e)
		case e := <-zm.events.recv[playerPromptReply]:
			go playerPromptReplyLogic(zm, e)
		}
	}
}

func (zm *zoneMap) monsterActivity() {
	log.Infof("[map_worker] monsterActivity worker for map %v", zm.data.Info.MapName)
	for {
		select {
		case e := <-zm.events.recv[monsterAppeared]:
			log.Info(e)
		case e := <-zm.events.recv[monsterDisappeared]:
			log.Info(e)
		case e := <-zm.events.recv[monsterWalks]:
			go npcWalksLogic(zm, e)
		case e := <-zm.events.recv[monsterRuns]:
			go npcRunsLogic(zm, e)
		}
	}
}

func playerClicksOnNpcLogic(zm *zoneMap, e event) {
	ev, ok := e.(*playerClicksOnNpcEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&playerClicksOnNpcEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	p, err := zm.entities.getPlayer(ev.handle)
	if err != nil {
		log.Error(err)
		return
	}

	// find npc with handle in ev.nc.Handle
	// send id of mob
	var (
		nc1 *structs.NcServerMenuReq
		nc2 *structs.NcActNpcMenuOpenReq
	)

	for n := range zm.entities.allNpc() {
		if n.getHandle() == ev.nc.NpcHandle {
			nc2 = &structs.NcActNpcMenuOpenReq{
				MobID: n.data.mobInfo.ID,
			}

			p.targeting.Lock()
			p.targeting.selectingN = n
			p.targeting.Unlock()

			if n.nType == npcPortal {
				var md *data.Map

				for i, m := range mapData.Maps {
					if m.Info.MapName.Name == n.data.npcData.ShinePortal.ServerMapIndex {
						md = mapData.Maps[i]
						break
					}
				}

				var mapName string
				if md != nil {
					mapName = md.Info.Name
				} else {
					mapName = "MAP IS OFFLINE"
				}

				mapName = strings.Replace(mapName, "#", " ", -1)

				title := fmt.Sprintf("Do you want to move to %v", mapName)

				nc1 = &structs.NcServerMenuReq{
					Title:     title,
					Priority:  0,
					NpcHandle: n.getHandle(),
					NpcPosition: structs.ShineXYType{
						X: uint32(n.current.x),
						Y: uint32(n.current.y),
					},
					LimitRange: 350,
					MenuNumber: 2,
					Menu: []structs.ServerMenu{
						{
							Reply:   1,
							Content: "Yes.",
						},
						{
							Reply:   0,
							Content: "No.",
						},
					},
				}

				networking.Send(p.conn.outboundData, networking.NC_MENU_SERVERMENU_REQ, nc1)
				if md == nil {
					return
				}

				p.Lock()
				p.next = location{
					mapID:   md.ID,
					mapName: md.Info.MapName.Name,
					x:       n.data.npcData.ShinePortal.X,
					y:       n.data.npcData.ShinePortal.Y,
				}
				p.Unlock()
				return
			}

			networking.Send(p.conn.outboundData, networking.NC_ACT_NPCMENUOPEN_REQ, nc2)
			break
		}
	}
}

func playerPromptReplyLogic(zm *zoneMap, e event) {
	log.Info(e)
	ev, ok := e.(*playerPromptReplyEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&playerPromptReplyEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	if ev.nc.Reply == 0 {
		return
	}

	p, err := zm.entities.getPlayer(ev.s.handle)
	if err != nil {
		log.Error(err)
		return
	}

	if p.targeting.selectingN == nil {
		log.Warning("prompt cannot be answered, player is no longer selecting an NPC")
		return
	}

	// for now, its only about portals
	if p.targeting.selectingN.nType == npcPortal && portalMatchesLocation(p.targeting.selectingN.data.npcData.ShinePortal, p.baseEntity.next) {
		// move player to map
		nzm, ok := maps.list[p.next.mapID]
		if !ok {
			log.Error(errors.Err{
				Code: errors.ZoneMapNotFound,
				Details: errors.ErrDetails{
					"mapID": p.next.mapID,
				},
			})
			return
		}

		cme := changeMapEvent{
			p:    p,
			s:    ev.s,
			prev: zm,
			next: nzm,
		}

		zoneEvents[changeMap] <- &cme
	}
}

func playerSelectsEntityLogic(zm *zoneMap, e event) {
	ev, ok := e.(*playerSelectsEntityEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&playerSelectsEntityEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	vp, err := zm.entities.getPlayer(ev.handle)
	if err != nil {
		log.Error(err)
		return
	}

	p, n := findFirstEntity(zm, ev.nc.TargetHandle)

	// set timeout in case of nonexistent handle
	// or use bool channels, and in the default case check if all three are false and return if so
	var nc *structs.NcBatTargetInfoCmd

	var notP, notN bool
	for {
		select {
		case ap := <-p:
			if ap == nil {
				notP = true
				break
			}

			nc = ap.ncBatTargetInfoCmd()

			order := vp.selectsPlayer(ap)

			nc.Order = order
			networking.Send(vp.conn.outboundData, networking.NC_BAT_TARGETINFO_CMD, nc)

			if ap.targeting.selectingP != nil {
				nextNc := ap.targeting.selectingP.ncBatTargetInfoCmd()
				nextNc.Order = order + 1
				networking.Send(vp.conn.outboundData, networking.NC_BAT_TARGETINFO_CMD, nextNc)
			}

			if ap.targeting.selectingN != nil {
				nextNc := ncBatTargetInfoCmd(ap.targeting.selectingN)
				nextNc.Order = order + 1
				networking.Send(vp.conn.outboundData, networking.NC_BAT_TARGETINFO_CMD, nextNc)
			}

			for p := range vp.selectedByPlayers() {
				nextNc := *nc
				nextNc.Order++
				networking.Send(p.conn.outboundData, networking.NC_BAT_TARGETINFO_CMD, nextNc)
			}

			return
		case an := <-n:
			if an == nil {
				notN = true
				break
			}
			order := vp.selectsNPC(an)

			nc = ncBatTargetInfoCmd(an)

			nc.Order = order

			networking.Send(vp.conn.outboundData, networking.NC_BAT_TARGETINFO_CMD, nc)

			for p := range vp.selectedByPlayers() {
				nextNc := *nc
				nextNc.Order++
				networking.Send(p.conn.outboundData, networking.NC_BAT_TARGETINFO_CMD, nextNc)
			}
			return
		default:
			if notP && notN {
				return
			}
		}
	}
}

func playerUnselectsEntityLogic(zm *zoneMap, e event) {
	ev, ok := e.(*playerUnselectsEntityEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&playerUnselectsEntityEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	vp, err := zm.entities.getPlayer(ev.handle)
	if err != nil {
		log.Error(err)
		return
	}

	var order byte
	vp.targeting.Lock()
	order = vp.targeting.selectionOrder
	vp.targeting.selectingP = nil
	vp.targeting.selectingN = nil
	vp.targeting.Unlock()

	nc := structs.NcBatTargetInfoCmd{
		Order:         order + 1,
		Handle:        65535,
		TargetHP:      0,
		TargetMaxHP:   0,
		TargetSP:      0,
		TargetMaxSP:   0,
		TargetLP:      0,
		TargetMaxLP:   0,
		TargetLevel:   0,
		HpChangeOrder: 0,
	}

	vp.targeting.RLock()
	for _, p := range vp.targeting.selectedByP {
		go networking.Send(p.conn.outboundData, networking.NC_BAT_TARGETINFO_CMD, &nc)
	}
	vp.targeting.RUnlock()
}

func itemEquipLogic(e event, zm *zoneMap) {
	var (
		ev  *itemEquipEvent
		ev1 *eduEquipItemEvent
	)

	ev, ok := e.(*itemEquipEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&itemEquipEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	p, err := zm.entities.getPlayer(ev.session.handle)
	if err != nil {
		log.Error(err)
		return
	}

	ev1 = &eduEquipItemEvent{
		slot: int(ev.nc.Slot),
		err:  make(chan error),
	}

	p.events.send[eduEquipItem] <- ev1

	select {
	case err := <-ev1.err:
		if err != nil {
			networking.Send(p.conn.outboundData, networking.NC_ITEM_EQUIP_ACK, &structs.NcItemEquipAck{
				Code: ItemEquipFailed,
			})
			return
		}
	}

	nc1, _, err := ncItemEquipChangeCmd(ev1.change)
	if err != nil {
		networking.Send(p.conn.outboundData, networking.NC_ITEM_EQUIP_ACK, &structs.NcItemEquipAck{
			Code: ItemEquipFailed,
		})
		log.Error(err)
		return
	}

	_, nc2, err := ncItemCellChangeCmd(ev1.change)
	if err != nil {
		networking.Send(p.conn.outboundData, networking.NC_ITEM_EQUIP_ACK, &structs.NcItemEquipAck{
			Code: ItemEquipFailed,
		})
		log.Error(err)
		return
	}

	networking.Send(p.conn.outboundData, networking.NC_ITEM_EQUIPCHANGE_CMD, &nc1)

	networking.Send(p.conn.outboundData, networking.NC_ITEM_EQUIP_ACK, &structs.NcItemEquipAck{
		Code: ItemEquipSuccess,
	})

	networking.Send(p.conn.outboundData, networking.NC_ITEM_CELLCHANGE_CMD, nc2)

	var ph = p.getHandle()

	switch ev1.change.from.item.itemData.itemInfo.Class {
	case data.ItemClassWeapon, data.ItemClassShield:
		nc3 := &structs.NcBriefInfoChangeWeaponCmd{
			UpgradeInfo: structs.NcBriefInfoChangeUpgradeCmd{
				Handle: ph,
				Item:   ev1.change.from.item.itemData.itemInfo.ID,
				Slot:   byte(ev1.change.from.item.itemData.itemInfo.Equip),
			},
		}
		mapEpicenterBroadcast(zm, p, networking.NC_BRIEFINFO_CHANGEWEAPON_CMD, nc3)
		break
	case data.ItemClassArmor, data.ItemClassBoot, data.ItemClassAmulet, data.ItemBracelet:
		nc3 := &structs.NcBriefInfoChangeUpgradeCmd{
			Handle: ph,
			Item:   ev1.change.from.item.itemData.itemInfo.ID,
			Slot:   byte(ev1.change.from.item.itemData.itemInfo.Equip),
		}
		mapEpicenterBroadcast(zm, p, networking.NC_BRIEFINFO_CHANGEUPGRADE_CMD, nc3)
		break
	}
}

func itemUnEquipLogic(e event, zm *zoneMap) {
	var (
		ev  *itemUnEquipEvent
		ev1 *eduUnEquipItemEvent
	)

	ev, ok := e.(*itemUnEquipEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&itemUnEquipEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	p, err := zm.entities.getPlayer(ev.session.handle)
	if err != nil {
		log.Error(err)
		return
	}

	ev1 = &eduUnEquipItemEvent{
		from: int(ev.nc.SlotEquip),
		to:   int(ev.nc.SlotInven),
		err:  make(chan error),
	}

	p.events.send[eduUnEquipItem] <- ev1

	select {
	case err := <-ev1.err:
		if err != nil {
			networking.Send(p.conn.outboundData, networking.NC_ITEM_EQUIP_ACK, &structs.NcItemEquipAck{
				Code: ItemEquipFailed,
			})
			log.Error(err)
			return
		}
	}

	_, nc1, err := ncItemEquipChangeCmd(ev1.change)
	if err != nil {
		networking.Send(p.conn.outboundData, networking.NC_ITEM_EQUIP_ACK, &structs.NcItemEquipAck{
			Code: ItemEquipFailed,
		})
		log.Error(err)
		return
	}

	nc2, _, err := ncItemCellChangeCmd(ev1.change)
	if err != nil {
		networking.Send(p.conn.outboundData, networking.NC_ITEM_EQUIP_ACK, &structs.NcItemEquipAck{
			Code: ItemEquipFailed,
		})
		log.Error(err)
		return
	}

	networking.Send(p.conn.outboundData, networking.NC_ITEM_EQUIPCHANGE_CMD, &nc1)

	networking.Send(p.conn.outboundData, networking.NC_ITEM_EQUIP_ACK, &structs.NcItemEquipAck{
		Code: ItemEquipSuccess,
	})

	networking.Send(p.conn.outboundData, networking.NC_ITEM_CELLCHANGE_CMD, nc2)

	//NC_BRIEFINFO_UNEQUIP_CMD
	nc3 := &structs.NcBriefInfoUnEquipCmd{
		Handle: p.getHandle(),
		Slot:   ev.nc.SlotEquip,
	}

	for ap := range zm.entities.allPlayers() {
		go func(p1, p2 *player, nc *structs.NcBriefInfoUnEquipCmd) {
			if p1.getHandle() != p2.getHandle() {
				if withinRange(p2, p1) {
					go networking.Send(p2.conn.outboundData, networking.NC_BRIEFINFO_UNEQUIP_CMD, nc)
				}
			}
		}(p, ap, nc3)
	}
}

func itemIsMovedLogic(e event, zm *zoneMap) {
	ev, ok := e.(*itemIsMovedEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&itemIsMovedEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	p, err := zm.entities.getPlayer(ev.session.handle)
	if err != nil {
		log.Error(err)
		return
	}

	change, err := p.inventories.moveItem(ev.nc.From.Inventory, ev.nc.To.Inventory)
	if err != nil {
		log.Error(err)
		//TODO: check error type, send custom NC_ITEM_RELOC_ACK with error code
		networking.Send(p.conn.outboundData, networking.NC_ITEM_RELOC_ACK, &structs.NcItemRelocateAck{
			Code: ItemSlotChangeCommon,
		})
		return
	}

	cc1, cc2, err := ncItemCellChangeCmd(change)
	if err != nil {
		log.Error(err)
		//TODO: check error type, send custom NC_ITEM_RELOC_ACK with error code
		networking.Send(p.conn.outboundData, networking.NC_ITEM_RELOC_ACK, &structs.NcItemRelocateAck{
			Code: ItemSlotChangeCommon,
		})
		return
	}

	networking.Send(p.conn.outboundData, networking.NC_ITEM_RELOC_ACK, &structs.NcItemRelocateAck{
		Code: ItemSlotChangeOk,
	})

	networking.Send(p.conn.outboundData, networking.NC_ITEM_CELLCHANGE_CMD, cc1)
	networking.Send(p.conn.outboundData, networking.NC_ITEM_CELLCHANGE_CMD, cc2)
}

func npcWalksLogic(zm *zoneMap, e event) {
	ev, ok := e.(*npcWalksEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&npcWalksEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}
	for ap := range zm.entities.allPlayers() {
		go func(p *player, n *npc) {
			if withinRange(p, n) {
				networking.Send(p.conn.outboundData, networking.NC_ACT_SOMEONEMOVEWALK_CMD, ev.nc)
			}
		}(ap, ev.n)
	}
}

func npcRunsLogic(zm *zoneMap, e event) {
	ev, ok := e.(*npcRunsEvent)

	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&npcRunsEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	for ap := range zm.entities.allPlayers() {
		go func(p *player, n *npc) {
			if withinRange(p, n) {
				go networking.Send(p.conn.outboundData, networking.NC_ACT_SOMEONEMOVERUN_CMD, ev.nc)
			}
		}(ap, ev.n)
	}
}

func playerHandleMaintenanceLogic(zm *zoneMap) {
	for ap := range zm.entities.allPlayers() {
		go func(p *player) {

			if justSpawned(p) {
				return
			}

			if lastHeartbeat(p) < playerHeartbeatLimit {
				return
			}

			h := p.getHandle()

			pde := &playerDisappearedEvent{
				handle: h,
			}

			p.ticks.Lock()
			for _, t := range p.ticks.list {
				t.Stop()
			}
			p.ticks.Unlock()

			select {
			case zm.events.send[playerDisappeared] <- pde:
				break
			default:
				log.Error("failed to stop heartbeat")
				break
			}

			zm.entities.removePlayer(h)

			removeHandle(h)

			p.conn.close <- true

		}(ap)
	}
}

func playerHandleLogic(e event, zm *zoneMap) {
	ev, ok := e.(*playerHandleEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&playerHandleEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	handle, err := newHandle()

	if err != nil {
		ev.err <- err
		return
	}

	ev.player.baseEntity.Lock()
	ev.player.baseEntity.handle = handle
	ev.player.baseEntity.Unlock()

	zm.addEntity(ev.player)

	ev.session.handle = handle
	ev.session.mapID = ev.player.current.mapID

	ev.done <- true
}

func playerAppearedLogic(e event, zm *zoneMap) {
	ev, ok := e.(*playerAppearedEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerAppearedEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	p, err := zm.entities.getPlayer(ev.handle)
	if err != nil {
		log.Error(err)
		return
	}

	npcs := ncBriefInfoMobCmd(zm)
	networking.Send(p.conn.outboundData, networking.NC_BRIEFINFO_MOB_CMD, &npcs)

	go p.heartbeat()
	go p.persistPosition()

	go p.newNearbyEntitiesTicker(zm)
	go p.oldNearbyEntitiesTicker()

	//go adjacentMonstersInform(p, zm)
}

// send data about equippedID or unequipped slots
// this is used on changing maps, for example when you leave the KQ arena and the hammer is automatically removed when you leave the map
func equippedItems(p1 *player) {
	var (
		nc1 []structs.NcBriefInfoUnEquipCmd
		nc2 []structs.NcBriefInfoChangeWeaponCmd
		nc3 []structs.NcBriefInfoChangeUpgradeCmd
		nc4 []structs.NcBriefInfoChangeDecorateCmd
	)

	p1h := p1.getHandle()
	p1.inventories.RLock()
	for i := 1; i <= 29; i++ {
		eItem, ok := p1.inventories.equipped.items[i]
		if !ok {
			nc1 = append(nc1, structs.NcBriefInfoUnEquipCmd{
				Handle: p1h,
				Slot:   byte(i),
			})
			continue
		}

		eItem.RLock()
		switch eItem.itemData.itemInfo.Class {
		case data.ItemClassArmor:
		case data.ItemClassAmulet:
		case data.ItemBracelet:
		case data.ItemClassShield:
		case data.ItemClassBoot:
			nc3 = append(nc3, structs.NcBriefInfoChangeUpgradeCmd{
				Handle: p1h,
				Item:   eItem.itemData.itemInfo.ID,
				//Slot:   byte(eItem.itemData.itemInfo.Equip),
			})
			break
		case data.ItemClassWeapon:
			nc2 = append(nc2, structs.NcBriefInfoChangeWeaponCmd{
				UpgradeInfo: structs.NcBriefInfoChangeUpgradeCmd{
					Handle: p1h,
					Item:   eItem.itemData.itemInfo.ID,
					//Slot:   byte(eItem.itemData.itemInfo.Equip),
				},
				CurrentMobID:     65535,
				CurrentKillLevel: 255,
			})
			break
		case data.ItemCosWeapon:
		case data.ItemCosShield:
		case data.ItemClassDecoration:
			nc4 = append(nc4, structs.NcBriefInfoChangeDecorateCmd{
				Handle: p1h,
				Item:   eItem.itemData.itemInfo.ID,
				//Slot:   byte(eItem.itemData.itemInfo.Equip),
			})
			break
		}
		eItem.RUnlock()

	}
	p1.inventories.RUnlock()

	for _, nc := range nc1 {
		networking.Send(p1.conn.outboundData, networking.NC_BRIEFINFO_UNEQUIP_CMD, &nc)
	}
	for _, nc := range nc2 {
		networking.Send(p1.conn.outboundData, networking.NC_BRIEFINFO_CHANGEWEAPON_CMD, &nc)
	}
	for _, nc := range nc3 {
		networking.Send(p1.conn.outboundData, networking.NC_BRIEFINFO_CHANGEUPGRADE_CMD, &nc)
	}
	for _, nc := range nc4 {
		networking.Send(p1.conn.outboundData, networking.NC_BRIEFINFO_CHANGEDECORATE_CMD, &nc)
	}
}

func playerDisappearedLogic(e event, zm *zoneMap) {
	ev, ok := e.(*playerDisappearedEvent)

	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerDisappearedEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	for ap := range zm.entities.allPlayers() {
		go func(p2 *player) {
			if p2.getHandle() == ev.handle {
				return
			}
			nc := &structs.NcBriefInfoDeleteHandleCmd{
				Handle: ev.handle,
			}
			networking.Send(p2.conn.outboundData, networking.NC_BRIEFINFO_BRIEFINFODELETE_CMD, nc)
		}(ap)
	}

}

func playerWalksLogic(e event, zm *zoneMap) {
	// player has a fifo queue for the last 30 movements
	// for every movement
	//		verify collision
	//			if fails, return to previous movement
	// 		verify speed ( default 30 for unmounted/unbuffed player)
	//			if fails return to position 1 in queue
	//		broadcast to players within range
	var (
		ev  *playerWalksEvent
		ev1 *eduPositionEvent
	)
	ev, ok := e.(*playerWalksEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerAppearedEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	p1, err := zm.entities.getPlayer(ev.handle)
	if err != nil {
		log.Error(err)
		return
	}

	if !ok {
		log.Error("player not found during playerStoppedLogic")
		return
	}

	igX := int(ev.nc.To.X)
	igY := int(ev.nc.To.Y)

	ev1 = &eduPositionEvent{
		x:   igX,
		y:   igY,
		zm:  zm,
		err: make(chan error),
	}

	p1.events.send[eduPosition] <- ev1

	select {
	case err := <-ev1.err:
		if err != nil {
			log.Error(err)
			return
		}
	}

	nc := &structs.NcActSomeoneMoveWalkCmd{
		Handle: ev.handle,
		From:   ev.nc.From,
		To:     ev.nc.To,
		Speed:  walkSpeed,
	}

	mapEpicenterBroadcast(zm, p1, networking.NC_ACT_SOMEONEMOVEWALK_CMD, nc)
}

func playerRunsLogic(e event, zm *zoneMap) {
	// player has a fifo queue for the last 30 movements
	// for every movement
	//		verify collision
	//			if fails, return to previous movement
	// 		verify speed ( default 30 for unmounted/unbuffed player)
	//			if fails return to position 1 in queue
	//		broadcast to players within range
	// 		new to movements array
	var (
		ev  *playerRunsEvent
		ev1 *eduPositionEvent
	)

	ev, ok := e.(*playerRunsEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerAppearedEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	p1, err := zm.entities.getPlayer(ev.handle)
	if err != nil {
		log.Error(err)
		return
	}

	// check if distance is too high
	// check if distance matches the applied mount or speed buffs

	igX := int(ev.nc.To.X)
	igY := int(ev.nc.To.Y)

	ev1 = &eduPositionEvent{
		x:   igX,
		y:   igY,
		zm:  zm,
		err: make(chan error),
	}

	p1.events.send[eduPosition] <- ev1

	select {
	case err := <-ev1.err:
		if err != nil {
			log.Error(err)
			return
		}
	}

	nc := &structs.NcActSomeoneMoveRunCmd{
		Handle: ev.handle,
		From:   ev.nc.From,
		To:     ev.nc.To,
		Speed:  runSpeed,
	}

	mapEpicenterBroadcast(zm, p1, networking.NC_ACT_SOMEONEMOVERUN_CMD, nc)
}

func playerStoppedLogic(e event, zm *zoneMap) {
	// movements triggered by keys inmediately send a STOP packet to the server
	// movements triggered by mouse do not send a STOP packet
	// for every stop
	//		verify collision
	//		broadcast to players within range
	var (
		ev  *playerStoppedEvent
		ev1 *eduPositionEvent
	)
	ev, ok := e.(*playerStoppedEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerStoppedEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	p1, err := zm.entities.getPlayer(ev.handle)
	if err != nil {
		log.Error(err)
		return
	}

	igX := int(ev.nc.Location.X)
	igY := int(ev.nc.Location.Y)

	ev1 = &eduPositionEvent{
		x:   igX,
		y:   igY,
		zm:  zm,
		err: make(chan error),
	}

	p1.events.send[eduPosition] <- ev1

	select {
	case err := <-ev1.err:
		if err != nil {
			log.Error(err)
			return
		}
	}

	nc := &structs.NcActSomeoneStopCmd{
		Handle:   ev.handle,
		Location: ev.nc.Location,
	}

	mapEpicenterBroadcast(zm, p1, networking.NC_ACT_SOMEONESTOP_CMD, nc)
}

func playerJumpedLogic(e event, zm *zoneMap) {
	ev, ok := e.(*playerJumpedEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerJumpedEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	p1, err := zm.entities.getPlayer(ev.handle)
	if err != nil {
		log.Error(err)
		return
	}

	nc := &structs.NcActSomeoneJumpCmd{
		Handle: ev.handle,
	}

	mapEpicenterBroadcast(zm, p1, networking.NC_ACT_SOMEEONEJUMP_CMD, nc)
}

func unknownHandleLogic(e event, zm *zoneMap) {
	ev, ok := e.(*unknownHandleEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&unknownHandleEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	// handle is persisted across maps X_x
	if ev.handle != ev.nc.AffectedHandle {
		log.Errorf("mismatched handles %v %v", ev.handle, ev.nc.AffectedHandle)
		return
	}

	p1, err := zm.entities.getPlayer(ev.handle)
	if err != nil {
		log.Error(err)
		return
	}

	//TODO: could also be NPC, Item on the ground or a Monster
	p2, err := zm.entities.getPlayer(ev.handle)
	if err != nil {
		log.Error(err)
		return
	}

	if withinRange(p1, p2) {
		nc1 := ncBriefInfoLoginCharacterCmd(p1)
		nc2 := ncBriefInfoLoginCharacterCmd(p2)
		go networking.Send(p2.conn.outboundData, networking.NC_BRIEFINFO_LOGINCHARACTER_CMD, &nc1)
		go networking.Send(p1.conn.outboundData, networking.NC_BRIEFINFO_LOGINCHARACTER_CMD, &nc2)
	}

	n := zm.entities.getNpc(ev.nc.ForeignHandle)

	if n == nil {
		return
	}

	if withinRange(p1, n) {
		nc := ncBriefInfoRegenMobCmd(n)
		go networking.Send(p1.conn.outboundData, networking.NC_BRIEFINFO_REGENMOB_CMD, &nc)
	}

}

func findFirstEntity(zm *zoneMap, handle uint16) (chan *player, chan *npc) {
	p := make(chan *player, 1)
	n := make(chan *npc, 1)

	go func(p chan<- *player, zm *zoneMap, targetHandle uint16) {
		for ap := range zm.entities.allPlayers() {
			if ap.getHandle() == targetHandle {
				p <- ap
				return
			}
		}
		p <- nil
	}(p, zm, handle)

	go func(n chan<- *npc, zm *zoneMap, targetHandle uint16) {
		for an := range zm.entities.allNpc() {
			if an.getHandle() == targetHandle {
				n <- an
				return
			}
		}
		n <- nil
	}(n, zm, handle)
	return p, n
}

func portalMatchesLocation(portal *data.ShinePortal, next location) bool {
	var (
		md *data.Map
	)

	for i, m := range mapData.Maps {
		if m.MapInfoIndex == portal.ClientMapIndex {
			md = mapData.Maps[i]
			break
		}
	}

	if md == nil {
		return false
	}

	if md.ID == next.mapID {
		return true
	}

	return false
}

// broadcast data to all players on map given an epicenter player
func mapEpicenterBroadcast(zm *zoneMap, epicenter *player, code networking.OperationCode, nc interface{}) {
	for p := range zm.entities.allPlayers() {
		go func(p1, p2 *player, nc interface{}) {
			if p1.getHandle() != p2.getHandle() {
				if withinRange(p2, p1) {
					networking.Send(p2.conn.outboundData, code, nc)
				}
			}
		}(epicenter, p, nc)
	}
}

func ncBriefInfoMobCmd(zm *zoneMap) structs.NcBriefInfoMobCmd {
	var npcs structs.NcBriefInfoMobCmd
	for n := range zm.entities.allNpc() {
		if n.baseEntity.eType == isNPC {
			info := ncBriefInfoRegenMobCmd(n)

			if n.nType == npcPortal {
				info.FlagState = 1
				info.FlagData.Data = n.data.npcData.ClientMapIndex
			}

			npcs.Mobs = append(npcs.Mobs, *info)
		}
	}

	npcs.MobNum = byte(len(npcs.Mobs))
	return npcs
}
