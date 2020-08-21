package service

import "time"

func playerLogout(cancel, conclude <-chan event, m * zoneMap, p * player) {
	t := time.NewTicker(15 * time.Second)

	finish := func() {
		t.Stop()
		select {
		case p.conn.close <- true:
			m.send[playerDisappeared] <- &playerDisappearedEvent{}
			return
		default:
			log.Error("unexpected error occurred while closing connection")
			return
		}
	}

	for {
		select {
		case _, ok := <-cancel:
			if !ok {
				log.Error("failed to receive event")
			}
			return
		case _, ok := <-conclude:
			if !ok {
				log.Error("failed to receive event")
				return
			}
			finish()
			return
		case <- t.C:
			finish()
			return
		}
	}
}

