package zone

import (
	"context"

	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
)

// NC_ITEM_RELOC_REQ
func ncItemRelocReq(ctx context.Context, np *networking.Parameters) {
	var e itemIsMovedEvent

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		log.Error(errors.Err{
			Code: errors.ZoneNoSessionAvailable,
		})
		return
	}

	e = itemIsMovedEvent{
		nc:      &structs.NcitemRelocateReq{},
		session: session,
	}

	err := structs.Unpack(np.Command.Base.Data, e.nc)
	if err != nil {
		log.Error(err)
		return
	}

	zm, ok := maps.list[session.mapID]
	if !ok {
		log.Error(errors.Err{
			Code: errors.ZoneMapNotFound,
			Details: errors.ErrDetails{
				"session": session,
			},
		})
		return
	}

	zm.events.send[itemIsMoved] <- &e
}

// NC_ITEM_EQUIP_REQ
func ncItemEquipReq(ctx context.Context, np *networking.Parameters) {
	var e itemEquipEvent

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	e = itemEquipEvent{
		nc:      &structs.NcItemEquipReq{},
		session: session,
	}

	err := structs.Unpack(np.Command.Base.Data, e.nc)
	if err != nil {
		log.Error(err)
		return
	}

	zm, ok := maps.list[session.mapID]

	if !ok {
		log.Error(errors.Err{
			Code: errors.ZoneMapNotFound,
			Details: errors.ErrDetails{
				"session": session,
			},
		})
		return
	}

	zm.events.send[itemEquip] <- &e
}

// NC_ITEM_UNEQUIP_REQ
func ncItemUnEquipReq(ctx context.Context, np *networking.Parameters) {
	var e itemUnEquipEvent

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	e = itemUnEquipEvent{
		nc:      &structs.NcItemUnequipReq{},
		session: session,
	}

	err := structs.Unpack(np.Command.Base.Data, e.nc)
	if err != nil {
		log.Error(err)
		return
	}

	zm, ok := maps.list[session.mapID]

	if !ok {
		log.Error(errors.Err{
			Code: errors.ZoneMapNotFound,
			Details: errors.ErrDetails{
				"session": session,
			},
		})
		return
	}

	zm.events.send[itemUnEquip] <- &e
}
