package zone

import "github.com/shine-o/shine.engine.emulator/pkg/structs"

func playerDisappearedLogic(zm *zoneMap, ev *playerDisappearedEvent) {
	zm.entities.players.Lock() // TODO: check if its necessary
	defer zm.entities.players.Unlock()

	for _, p := range zm.entities.players.active {
		if p.handle == ev.handle {
			continue
		}
		go ncMapLogoutCmd(p, &structs.NcMapLogoutCmd{
			Handle: ev.handle,
		})
	}
}
