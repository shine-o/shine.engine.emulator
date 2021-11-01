package zone

type handler struct {
	index uint16
	inUse map[uint16]bool
}

type handlerPetition struct {
	newHandle                 chan uint16
	queryHandle, deleteHandle uint16
	used                      chan bool
	err                       chan error
}

func (h *handler) handleWorker() {
	for {
		select {
		case hp := <-queryHandler:
			_, used := h.inUse[hp.queryHandle]
			if used {
				hp.used <- true
			}
			hp.used <- false
		case hp := <-newHandler:
			attempts := 1000
			for attempts != 0 {
				attempts--
				h.index++
				_, used := h.inUse[h.index]
				if used {
					continue
				}
				h.inUse[h.index] = true
				hp.newHandle <- h.index
				break
			}
			// return err
		case hp := <-removeHandler:
			delete(h.inUse, hp.deleteHandle)
		}
	}
}

func newHandle() (uint16, error) {
	hp := &handlerPetition{
		newHandle: make(chan uint16),
		err:       make(chan error),
	}

	newHandler <- hp

	select {
	case h := <-hp.newHandle:
		return h, nil
	case err := <-hp.err:
		return 0, err
	}
}

func removeHandle(h uint16) {
	hp := &handlerPetition{
		deleteHandle: h,
	}

	removeHandler <- hp
}

func handleExists(h uint16) bool {
	hp := &handlerPetition{
		queryHandle: h,
		used:        make(chan bool),
	}

	queryHandler <- hp

	return <-hp.used
}
